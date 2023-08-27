package bla

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

type AccountActor struct {
	// OperatorNKey (seed) is used to sign the account JWTs.
	OperatorNKey []byte `json:"operator_nkey"`

	// SysAccountConn is a NATS connection to the SYS account.
	SysAccountConn *nats.Conn

	// ActorAccountConn is a NATS connection to the actor account.
	ActorAccountConn *nats.Conn
	// ActorAccountPubKey is the public key of the actor account.
	ActorAccountPubKey string

	subs []*nats.Subscription
}

func (ac *AccountActor) Close() {
	for _, sub := range ac.subs {
		sub.Unsubscribe()
	}
	ac.SysAccountConn.Close()
	ac.ActorAccountConn.Close()
}

func (ac *AccountActor) requestAccount(claims *jwt.AccountClaims) (string, error) {
	kp, err := nkeys.FromSeed(ac.OperatorNKey)
	if err != nil {
		return "", fmt.Errorf("getting operator key pair from nkey: %w", err)
	}
	accountJWT, err := claims.Encode(kp)
	if err != nil {
		return "", fmt.Errorf("encoding account claims: %w", err)
	}
	if _, err := ac.SysAccountConn.Request(
		"$SYS.REQ.CLAIMS.UPDATE",
		[]byte(accountJWT),
		time.Second,
	); err != nil {
		return "", fmt.Errorf("requesting new account: %w", err)
	}
	return accountJWT, nil
}

// On startup:
// 1. Bootstrap operator, sys account + user
// 2. Connect to sys and create "actor" account + user
//

type CreateAccountMsg struct {
	Name string `json:"name"`
}

type CreateAccountReply struct {
	ID string `json:"id"`
}

type CreateUserMsg struct {
	Name string `json:"name"`
}

type CreateUserReply struct {
	User
}

func RegisterAccountActor(ctx context.Context, actor *AccountActor) error {
	name := "account"
	js, err := actor.ActorAccountConn.JetStream()
	if err != nil {
		return fmt.Errorf("jetstream: %w", err)
	}
	accountsBucket, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "accounts",
	})
	if err != nil {
		return fmt.Errorf("create key value: %w", err)
	}

	{
		action := "create"
		sub, err := actor.ActorAccountConn.QueueSubscribe(
			fmt.Sprintf("actor.%s.%s", name, action),
			name,
			func(msg *nats.Msg) {
				var data CreateAccountMsg
				if err := json.Unmarshal(msg.Data, &data); err != nil {
					slog.Error("unmarshal", "error", err)
					return
				}
				keyPair, err := nkeys.CreateAccount()
				if err != nil {
					slog.Error("create account", "error", err)
					return
				}
				pubKey, err := keyPair.PublicKey()
				if err != nil {
					slog.Error("get public key", "error", err)
					return
				}
				seed, err := keyPair.Seed()
				if err != nil {
					slog.Error("get seed", "error", err)
					return
				}

				claims := jwt.NewAccountClaims(pubKey)
				claims.Name = data.Name
				claims.Limits.JetStreamLimits.Consumer = -1
				claims.Limits.JetStreamLimits.DiskMaxStreamBytes = -1
				claims.Limits.JetStreamLimits.DiskStorage = -1
				claims.Limits.JetStreamLimits.MaxAckPending = -1
				claims.Limits.JetStreamLimits.MemoryMaxStreamBytes = -1
				claims.Limits.JetStreamLimits.MemoryStorage = -1
				claims.Limits.JetStreamLimits.Streams = -1
				claims.Imports.Add(&jwt.Import{
					Type: jwt.Service,
					Name: "all-actors",
					// Account is the public key of the account which exported the service.
					Account: actor.ActorAccountPubKey,
					// Subject is the exported account's subject.
					Subject: jwt.Subject("actor.*.*." + pubKey),
					// LocalSubject is the subject local to this account.
					LocalSubject: jwt.RenamingSubject("actor.*.*"),
				})
				if err := validateClaims(claims); err != nil {
					slog.Error("validate claims", "error", err)
					return
				}
				jwt, err := actor.requestAccount(claims)
				if err != nil {
					slog.Error("request account", "error", err)
					return
				}

				acc := Account{
					ID:   pubKey,
					Name: data.Name,
					NKey: seed,
					JWT:  jwt,
				}
				accB, err := json.Marshal(acc)
				if err != nil {
					slog.Error("marshal", "error", err)
					return
				}

				rev, err := accountsBucket.Put(pubKey, accB)
				if err != nil {
					slog.Error("put", "error", err)
					return
				}
				slog.Info("put account", "key", pubKey, "rev", rev)

				reply := CreateAccountReply{
					ID: pubKey,
				}
				replyB, err := json.Marshal(reply)
				if err != nil {
					slog.Error("marshal", "error", err)
					return
				}
				if err := msg.Respond(replyB); err != nil {
					slog.Error("respond", "error", err)
					return
				}
			},
		)
		if err != nil {
			return fmt.Errorf("subscribing: %w", err)
		}
		actor.subs = append(actor.subs, sub)
	}
	{
		action := "create_user"
		sub, err := actor.ActorAccountConn.QueueSubscribe(
			fmt.Sprintf("actor.%s.%s.*", name, action),
			name,
			func(msg *nats.Msg) {
				// Get account public key from subject
				subjects := strings.Split(msg.Subject, ".")
				if len(subjects) != 4 {
					slog.Error("invalid subject", "subject", msg.Subject)
					return
				}
				accountPubKey := subjects[3]
				slog.Info("create user", "subject", msg.Subject, "pubkey", accountPubKey)
				// Get account key pair
				kve, err := accountsBucket.Get(accountPubKey)
				if err != nil {
					slog.Error("get account key pair from bucket", "error", err)
					return
				}
				var acc Account
				if err := json.Unmarshal(kve.Value(), &acc); err != nil {
					slog.Error("unmarshal", "error", err)
					return
				}
				accountKeyPair, err := nkeys.FromSeed(acc.NKey)
				if err != nil {
					slog.Error("get account key pair", "error", err)
					return
				}

				// msg.Subject
				var data CreateUserMsg
				if err := json.Unmarshal(msg.Data, &data); err != nil {
					slog.Error("unmarshal", "error", err)
					return
				}
				keyPair, err := nkeys.CreateUser()
				if err != nil {
					slog.Error("create user", "error", err)
					return
				}
				pubKey, err := keyPair.PublicKey()
				if err != nil {
					slog.Error("get public key", "error", err)
					return
				}
				claims := jwt.NewUserClaims(pubKey)
				claims.Name = data.Name
				claims.IssuerAccount = accountPubKey
				if err := validateClaims(claims); err != nil {
					slog.Error("validate claims", "error", err)
					return
				}
				jwt, err := claims.Encode(accountKeyPair)
				if err != nil {
					slog.Error("encode claims", "error", err)
					return
				}
				seed, err := keyPair.Seed()
				if err != nil {
					slog.Error("get seed", "error", err)
					return
				}
				user := User{
					ID:   pubKey,
					Name: data.Name,
					NKey: seed,
					JWT:  jwt,
				}
				userB, err := json.Marshal(user)
				if err != nil {
					slog.Error("marshal", "error", err)
					return
				}
				if err := msg.Respond(userB); err != nil {
					slog.Error("respond", "error", err)
					return
				}
			},
		)
		if err != nil {
			return fmt.Errorf("subscribing: %w", err)
		}
		actor.subs = append(actor.subs, sub)
	}

	return nil
}

func CreateAccount(nc *nats.Conn, msg CreateAccountMsg) (*CreateAccountReply, error) {
	msgB, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	replyB, err := nc.Request("actor.account.create", msgB, time.Second)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	var reply CreateAccountReply
	if err := json.Unmarshal(replyB.Data, &reply); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &reply, nil
}

func CreateUser(nc *nats.Conn, msg CreateUserMsg) (*CreateUserReply, error) {
	userMsgB, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	userB, err := nc.Request(
		"actor.account.create_user",
		userMsgB,
		time.Second*5,
	)
	if err != nil {
		return nil, fmt.Errorf("publishing user create: %w", err)
	}
	var reply CreateUserReply
	if err := json.Unmarshal(userB.Data, &reply); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &reply, nil
}

// CreateUserForAccount creates a user for a target account.
// In this case, the nats connection has to be for the actor account.
//
// This is the scenario when the target account may not have any users yet,
// and we want to create an initial user for that account.
func CreateUserForAccount(
	nc *nats.Conn,
	msg CreateUserMsg,
	targetAccountID string,
) (*CreateUserReply, error) {
	userMsgB, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	userB, err := nc.Request(
		"actor.account.create_user."+targetAccountID,
		userMsgB,
		time.Second*5,
	)
	if err != nil {
		return nil, fmt.Errorf("publishing user create: %w", err)
	}
	var reply CreateUserReply
	if err := json.Unmarshal(userB.Data, &reply); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &reply, nil
}

// Account represents a NATS account.
type Account struct {
	// ID of the account, which for NATS is the public key of the account
	// and the subject of the account's JWT.
	ID string `json:"id"`
	// Name is the user-friendly name of the account.
	Name string `json:"name"`
	// NKey of the account.
	// The NKey (or "seed") can be converted into the account public
	// and private keys. The public key must match the account ID.
	NKey []byte `json:"nkey"`
	// JWT of the account.
	// The JWT contains the account claims (i.e. name, config, limits, etc.)
	// and is signed using an operator NKey.
	// In a way, the JWT *is* the account definition, which we send to
	// the NATS server to create an account.
	JWT string `json:"jwt"`
}

// User represents a NATS user.
type User struct {
	// ID of the user, which for NATS is the public key.
	ID string `json:"id"`
	// Name is the user-friendly name of the user.
	Name string `json:"name"`
	// NKey of the user.
	// The NKey (or "seed") can be converted into the user public
	// and private keys. The public key must match the user ID.
	NKey []byte `json:"nkey"`
	// JWT of the user.
	// The JWT contains the user claims (i.e. name, config, limits, etc.)
	// and is signed using an account NKey.
	JWT string `json:"jwt"`
}
