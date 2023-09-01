package bla

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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

	accountsBucket nats.KeyValue
	usersBucket    nats.KeyValue

	subs []*nats.Subscription
}

func (ac *AccountActor) Close() {
	for _, sub := range ac.subs {
		sub.Unsubscribe()
	}
	ac.SysAccountConn.Close()
	ac.ActorAccountConn.Close()
}

func (ac *AccountActor) requestAccount(
	claims *jwt.AccountClaims,
) (string, error) {
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

type AccountListMsg struct {
}

type AccountListReply struct {
	Accounts []Account `json:"accounts"`
}

type AccountGetMsg struct {
	ID string `json:"id"`
}

type AccountGetReply struct {
	Account Account `json:"account"`
}

type AccountCreateMsg struct {
	Name string `json:"name"`
}

type AccountCreateReply struct {
	ID string `json:"id"`
}

type UserCreateMsg struct {
	Name string `json:"name"`
}

type UserCreateReply struct {
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
	actor.accountsBucket = accountsBucket
	usersBucket, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "users",
	})
	if err != nil {
		return fmt.Errorf("create key value: %w", err)
	}
	actor.usersBucket = usersBucket

	{
		action := "account_list"
		sub, err := actor.ActorAccountConn.QueueSubscribe(
			fmt.Sprintf(ActionSubjectSubscribe, name, action),
			name,
			func(msg *nats.Msg) {
				reply, err := actor.accountList(ctx, msg)
				if err != nil {
					slog.Error("action", "action", action, "error", err)
					// reply = &CreateAccountReply{
					// 	Error: err.Error(),
					// }
					// TODO: handle errors
					if err := msg.Respond(nil); err != nil {
						slog.Error("respond", "error", err)
					}
				}
				bReply, err := json.Marshal(reply)
				if err != nil {
					slog.Error("marshal", "error", err)
					return
				}
				if err := msg.Respond(bReply); err != nil {
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
		action := "account_get"
		sub, err := actor.ActorAccountConn.QueueSubscribe(
			fmt.Sprintf(ActionSubjectSubscribe, name, action),
			name,
			func(msg *nats.Msg) {
				reply, err := actor.accountGet(ctx, msg)
				if err != nil {
					slog.Error("action", "action", action, "error", err)
					// reply = &CreateAccountReply{
					// 	Error: err.Error(),
					// }
					// TODO: handle errors
					if err := msg.Respond(nil); err != nil {
						slog.Error("respond", "error", err)
					}
				}
				bReply, err := json.Marshal(reply)
				if err != nil {
					slog.Error("marshal", "error", err)
					return
				}
				if err := msg.Respond(bReply); err != nil {
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
		action := "account_create"
		sub, err := actor.ActorAccountConn.QueueSubscribe(
			fmt.Sprintf(ActionSubjectSubscribe, name, action),
			name,
			func(msg *nats.Msg) {
				reply, err := actor.accountCreate(ctx, msg)
				if err != nil {
					slog.Error("action", "error", err)
					// reply = &CreateAccountReply{
					// 	Error: err.Error(),
					// }
					// TODO: handle errors
					if err := msg.Respond(nil); err != nil {
						slog.Error("respond", "error", err)
					}
				}
				bReply, err := json.Marshal(reply)
				if err != nil {
					slog.Error("marshal", "error", err)
					return
				}
				if err := msg.Respond(bReply); err != nil {
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
		action := "user_create"
		sub, err := actor.ActorAccountConn.QueueSubscribe(
			fmt.Sprintf(ActionSubjectSubscribeForAccount, name, action),
			name,
			func(msg *nats.Msg) {
				reply, err := actor.userCreate(ctx, msg)
				if err != nil {
					slog.Error("action", "error", err)
					// TODO: handle errors
					if err := msg.Respond(nil); err != nil {
						slog.Error("respond", "error", err)
					}
				}
				bReply, err := json.Marshal(reply)
				if err != nil {
					slog.Error("marshal", "error", err)
					return
				}
				if err := msg.Respond(bReply); err != nil {
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

func (aa *AccountActor) accountList(
	ctx context.Context,
	msg *nats.Msg,
) (*AccountListReply, error) {
	// TODO: unmarshal AccountListMsg, when needed
	keys, err := aa.accountsBucket.Keys()
	if err != nil {
		if errors.Is(err, nats.ErrNoKeysFound) {
			return &AccountListReply{
				Accounts: []Account{},
			}, nil
		}
		return nil, fmt.Errorf("get keys: %w", err)
	}

	accounts := make([]Account, len(keys))
	for i, key := range keys {
		kve, err := aa.accountsBucket.Get(key)
		if err != nil {
			return nil, fmt.Errorf("get account: %w", err)
		}
		var account Account
		if err := json.Unmarshal(kve.Value(), &account); err != nil {
			return nil, fmt.Errorf("unmarshal account: %w", err)
		}
		accounts[i] = account
	}
	return &AccountListReply{
		Accounts: accounts,
	}, nil
}

func (aa *AccountActor) accountGet(
	ctx context.Context,
	msg *nats.Msg,
) (*AccountGetReply, error) {
	// TODO: unmarshal AccountListMsg, when needed
	var data AccountGetMsg
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	kve, err := aa.accountsBucket.Get(data.ID)
	if err != nil {
		if errors.Is(err, nats.ErrKeyNotFound) {
			return &AccountGetReply{
				Account: Account{},
			}, nil
		}
		return nil, fmt.Errorf("get key: %w", err)
	}

	var account Account
	if err := json.Unmarshal(kve.Value(), &account); err != nil {
		return nil, fmt.Errorf("unmarshal account: %w", err)
	}
	return &AccountGetReply{
		Account: account,
	}, nil
}

func (aa *AccountActor) accountCreate(
	ctx context.Context,
	msg *nats.Msg,
) (*AccountCreateReply, error) {
	var data AccountCreateMsg
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	keyPair, err := nkeys.CreateAccount()
	if err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}
	pubKey, err := keyPair.PublicKey()
	if err != nil {
		return nil, fmt.Errorf("get public key: %w", err)
	}
	seed, err := keyPair.Seed()
	if err != nil {
		return nil, fmt.Errorf("get seed: %w", err)
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
		Account: aa.ActorAccountPubKey,
		// Subject is the exported account's subject.
		Subject: jwt.Subject(fmt.Sprintf(ActionImportSubject, pubKey)),
		// LocalSubject is the subject local to this account.
		LocalSubject: jwt.RenamingSubject(ActionImportLocalSubject),
	})
	// Export the Jetstream API for this account, which we will import into
	// the actor account, making this account's Jetstream API available to
	// connections from the actor account.
	claims.Exports.Add(&jwt.Export{
		Type:    jwt.Service,
		Name:    "js-api",
		Subject: jwt.Subject("$JS.API.>"),
	})
	// claims.Imports.Add(&jwt.Import{
	// 	Type: jwt.Service,
	// 	Name: "ingest",
	// 	// Account is the public key of the account which exported the service.
	// 	Account: aa.ActorAccountPubKey,
	// 	// Subject is the exported account's subject.
	// 	// Subject: ingest.<account-id>.<schema-name>.<schema-version>
	// 	Subject: jwt.Subject(fmt.Sprintf("ingest.%s.*.*", pubKey)),
	// 	// LocalSubject is the subject local to this account.
	// 	// Subject: ingest.<schema-name>.<schema-version>
	// 	LocalSubject: jwt.RenamingSubject("ingest.*.*"),
	// })
	if err := validateClaims(claims); err != nil {
		return nil, fmt.Errorf("validate claims: %w", err)
	}
	accJWT, err := aa.requestAccount(claims)
	if err != nil {
		return nil, fmt.Errorf("request account: %w", err)
	}

	acc := Account{
		ID:   pubKey,
		Name: data.Name,
		NKey: seed,
		JWT:  accJWT,
	}
	accB, err := json.Marshal(acc)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	// Use the account ID as the key.
	// TODO: can we have two accounts with the same name?
	rev, err := aa.accountsBucket.Put(acc.ID, accB)
	if err != nil {
		return nil, fmt.Errorf("put account: %w", err)
	}
	slog.Info("put account", "key", acc.ID, "rev", rev)

	// TODO: update the actor account JWT
	aac, err := claimsForAccount(aa.SysAccountConn, aa.ActorAccountPubKey)
	if err != nil {
		return nil, fmt.Errorf("getting actor account claims: %w", err)
	}
	aac.Imports.Add(&jwt.Import{
		Type:         jwt.Service,
		Name:         "js-api",
		Account:      pubKey,
		Subject:      jwt.Subject("$JS.API.>"),
		LocalSubject: jwt.RenamingSubject(fmt.Sprintf("$JS.API.%s.>", pubKey)),
	})
	if err := validateClaims(aac); err != nil {
		return nil, fmt.Errorf("validate updated actor account claims: %w", err)
	}
	actorJWT, err := aa.requestAccount(aac)
	if err != nil {
		return nil, fmt.Errorf("request updated actor account: %w", err)
	}
	slog.Info("UPDATED ACTOR ACCOUNT", "jwt", actorJWT)

	return &AccountCreateReply{
		ID: pubKey,
	}, nil
}

func (aa *AccountActor) userCreate(
	ctx context.Context,
	msg *nats.Msg,
) (*UserCreateReply, error) {
	accountPubKey, err := AccountFromSubject(msg.Subject)
	if err != nil {
		return nil, fmt.Errorf("getting account from subject: %w", err)
	}
	slog.Info("create user", "subject", msg.Subject, "pubkey", accountPubKey)
	// Get account key pair
	kve, err := aa.accountsBucket.Get(accountPubKey)
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}
	var acc Account
	if err := json.Unmarshal(kve.Value(), &acc); err != nil {
		return nil, fmt.Errorf("unmarshal account: %w", err)
	}
	accountKeyPair, err := nkeys.FromSeed(acc.NKey)
	if err != nil {
		return nil, fmt.Errorf("get account key pair: %w", err)
	}

	// msg.Subject
	var data UserCreateMsg
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	keyPair, err := nkeys.CreateUser()
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	pubKey, err := keyPair.PublicKey()
	if err != nil {
		return nil, fmt.Errorf("get public key: %w", err)
	}
	claims := jwt.NewUserClaims(pubKey)
	claims.Name = data.Name
	claims.IssuerAccount = accountPubKey
	if err := validateClaims(claims); err != nil {
		return nil, fmt.Errorf("validate claims: %w", err)
	}
	jwt, err := claims.Encode(accountKeyPair)
	if err != nil {
		return nil, fmt.Errorf("encode claims: %w", err)
	}
	seed, err := keyPair.Seed()
	if err != nil {
		return nil, fmt.Errorf("get seed: %w", err)
	}
	user := User{
		ID:   pubKey,
		Name: data.Name,
		NKey: seed,
		JWT:  jwt,
	}
	userB, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	if _, err := aa.usersBucket.Put(user.ID, userB); err != nil {
		return nil, fmt.Errorf("put user: %w", err)
	}
	return &UserCreateReply{
		User: user,
	}, nil
}

func SendAccountListMsg(
	nc *nats.Conn,
	msg AccountListMsg,
) (*AccountListReply, error) {
	msgB, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	replyB, err := nc.Request(
		fmt.Sprintf(ActionSubjectSend, "account", "account_list"),
		msgB,
		time.Second,
	)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	var reply AccountListReply
	if err := json.Unmarshal(replyB.Data, &reply); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &reply, nil
}

func SendAccountGetMsg(
	nc *nats.Conn,
	msg AccountGetMsg,
) (*AccountGetReply, error) {
	msgB, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	replyB, err := nc.Request(
		fmt.Sprintf(ActionSubjectSend, "account", "account_get"),
		msgB,
		time.Second,
	)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	var reply AccountGetReply
	if err := json.Unmarshal(replyB.Data, &reply); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &reply, nil
}

func SendAccountCreateMsg(
	nc *nats.Conn,
	msg AccountCreateMsg,
) (*AccountCreateReply, error) {
	msgB, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	replyB, err := nc.Request(
		fmt.Sprintf(ActionSubjectSend, "account", "account_create"),
		msgB,
		time.Second,
	)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	var reply AccountCreateReply
	if err := json.Unmarshal(replyB.Data, &reply); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &reply, nil
}

func SendUserCreateMsg(
	nc *nats.Conn,
	msg UserCreateMsg,
) (*UserCreateReply, error) {
	userMsgB, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	userB, err := nc.Request(
		fmt.Sprintf(ActionSubjectSend, "account", "user_create"),
		userMsgB,
		time.Second*5,
	)
	if err != nil {
		return nil, fmt.Errorf("publishing user create: %w", err)
	}
	var reply UserCreateReply
	if err := json.Unmarshal(userB.Data, &reply); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &reply, nil
}

// SendUserCreateForAccountMsg creates a user for a target account.
// In this case, the nats connection has to be for the actor account.
//
// This is the scenario when the target account may not have any users yet,
// and we want to create an initial user for that account.
func SendUserCreateForAccountMsg(
	nc *nats.Conn,
	targetAccountID string,
	msg UserCreateMsg,
) (*UserCreateReply, error) {
	userMsgB, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	userB, err := nc.Request(
		fmt.Sprintf(
			ActionSubjectSendForAccount,
			"account",
			"user_create",
			targetAccountID,
		),
		userMsgB,
		time.Second*5,
	)
	if err != nil {
		return nil, fmt.Errorf("publishing user create: %w", err)
	}
	var reply UserCreateReply
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
