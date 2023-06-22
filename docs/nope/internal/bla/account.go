// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package bla

import (
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

// AccountRequest defines the options available for creating a NATS account.
type AccountRequest struct {
	// Name of the account
	Name string
}

// Account represents a NATS account.
type Account struct {
	// ID of the account, which for NATS is the public key of the account
	// and the subject of the account's JWT.
	ID string
	// NKey of the account.
	// The NKey (or "seed") can be converted into the account public
	// and private keys. The public key must match the account ID.
	NKey []byte
	// JWT of the account.
	// The JWT contains the account claims (i.e. name, config, limits, etc.)
	// and is signed using an operator NKey.
	// In a way, the JWT *is* the account definition, which we send to
	// the NATS server to create an account.
	JWT string
}

// SyncAccount takes an [AccountRequest] and an optional [Account], and
// uses the provided [nats.Conn] to create or update the account.
//
// The account JWT is signed using the operator NKey.
//
// The [Account] is optional.
// In case it is not provided, a new NATS account is created, with a new key pair.
func SyncAccount(
	nc *nats.Conn,
	operatorNKey []byte,
	account *Account,
	req AccountRequest,
) (*Account, error) {
	// If account is nil, we need to create a new account
	if account == nil {
		account, err := newAccount(operatorNKey, req)
		if err != nil {
			return nil, fmt.Errorf("creating account: %w", err)
		}
		// Request the new account (send it to the server)
		if _, err := nc.Request(
			"$SYS.REQ.CLAIMS.UPDATE",
			[]byte(account.JWT),
			time.Second,
		); err != nil {
			return nil, fmt.Errorf("requesting new account: %w", err)
		}
		return account, nil
	}

	// Else, an account was previously created and we need to sync it.
	accountKeyPair, err := nkeys.FromSeed(account.NKey)
	if err != nil {
		return nil, fmt.Errorf("parsing account NKey: %w", err)
	}
	accountPublicKey, err := accountKeyPair.PublicKey()
	if err != nil {
		return nil, fmt.Errorf("getting account public key: %w", err)
	}
	// Validate that the account ID and NKey match
	if account.ID != accountPublicKey {
		return nil, errors.New("account ID does not match public key from account NKey")
	}

	// Create an account JWT based on the request
	accountJWT, err := generateAccountJWT(operatorNKey, accountPublicKey, req)
	if err != nil {
		return nil, fmt.Errorf("creating account JWT: %w", err)
	}
	account.JWT = accountJWT

	// Lookup the account on the server.
	// If the account does not exist, the observed behaviour is that an empty JWT
	// is returned, without an error.
	// And re-creating the account using the same ID (public key) seems to work
	// just fine, so no need to handle an account that does not exist.
	lookupMessage, err := nc.Request(
		fmt.Sprintf("$SYS.REQ.ACCOUNT.%s.CLAIMS.LOOKUP", account.ID),
		nil,
		time.Second,
	)
	if err != nil {
		return nil, fmt.Errorf("looking up account with ID %s: %w", account.ID, err)
	}

	// Compare the server-provided account JWT with the new one we just created.
	// If no difference, there is nothing to update.
	if string(lookupMessage.Data) == account.JWT {
		return account, nil
	}

	// If there is a different, we need to update the account
	if _, err := nc.Request(
		fmt.Sprintf("$SYS.REQ.ACCOUNT.%s.CLAIMS.UPDATE", account.ID),
		[]byte(account.JWT),
		time.Second,
	); err != nil {
		return nil, fmt.Errorf("updating account with ID %s: %w", account.ID, err)
	}

	return account, nil
}

// DeleteAccount sends a request to the NATS server using the provided connection
// to delete the account ID given.
//
// The request is a generic JWT signed with the operator NKey.
func DeleteAccount(nc *nats.Conn, operatorNKey []byte, id string) error {
	operatorKeyPair, err := nkeys.FromSeed(operatorNKey)
	if err != nil {
		return fmt.Errorf("getting operator key pair from nkey: %w", err)
	}
	operatorPublicKey, err := operatorKeyPair.PublicKey()
	if err != nil {
		return fmt.Errorf("getting oeprator public key from key pair: %w", err)
	}

	// Create generic claim with list of accounts to delete
	claim := jwt.NewGenericClaims(operatorPublicKey)
	claim.Data["accounts"] = []string{id}
	jwt, err := claim.Encode(operatorKeyPair)
	if err != nil {
		return fmt.Errorf("encoding generic account delete claim: %w", err)
	}

	deleteMsg, err := nc.Request("$SYS.REQ.CLAIMS.DELETE", []byte(jwt), time.Second)
	if err != nil {
		return fmt.Errorf("requesting account deletion: %w", err)
	}
	fmt.Println("DELETE MESSAGE:\n", string(deleteMsg.Data))

	return nil
}

// newAccount creates a whole new account and returns the key pair and JWT.
// It is used when a brand new account is requested, and it does not
// communicate with the NATS server (only returns the new [Account]).
func newAccount(operatorNKey []byte, req AccountRequest) (*Account, error) {
	accountKeyPair, err := nkeys.CreateAccount()
	if err != nil {
		return nil, fmt.Errorf("creating account nkeys: %w", err)
	}

	accountPublicKey, err := accountKeyPair.PublicKey()
	if err != nil {
		return nil, fmt.Errorf("getting account public key: %w", err)
	}

	accountNKey, err := accountKeyPair.Seed()
	if err != nil {
		return nil, fmt.Errorf("getting account seed from key pair: %w", err)
	}

	accountJWT, err := generateAccountJWT(operatorNKey, accountPublicKey, req)
	if err != nil {
		return nil, fmt.Errorf("creating account JWT: %w", err)
	}
	return &Account{
		ID:   accountPublicKey,
		NKey: accountNKey,
		JWT:  accountJWT,
	}, nil
}

// generateAccountJWT creates a new account JWT with the given subject (account ID / public key)
// using the given [AccountRequest] to populate the claims, and finally signing
// it with the operator NKey.
//
// It does not talk with the server, but it does validate the claims locally.
func generateAccountJWT(operatorNKey []byte, subject string, req AccountRequest) (string, error) {
	operatorKeyPair, err := nkeys.FromSeed(operatorNKey)
	if err != nil {
		return "", fmt.Errorf("getting operator key pair from nkey: %w", err)
	}

	accountClaims := jwt.NewAccountClaims(subject)
	accountClaims.Name = req.Name
	accountClaims.Limits.JetStreamLimits.Consumer = -1
	accountClaims.Limits.JetStreamLimits.DiskMaxStreamBytes = -1
	accountClaims.Limits.JetStreamLimits.DiskStorage = -1
	accountClaims.Limits.JetStreamLimits.MaxAckPending = -1
	accountClaims.Limits.JetStreamLimits.MemoryMaxStreamBytes = -1
	accountClaims.Limits.JetStreamLimits.MemoryStorage = -1
	accountClaims.Limits.JetStreamLimits.Streams = -1

	// Validate account claims
	vr := jwt.ValidationResults{}
	accountClaims.Validate(&vr)
	if vr.IsBlocking(true) {
		var vErr error
		for _, iss := range vr.Issues {
			vErr = errors.Join(vErr, iss)
		}
		return "", fmt.Errorf("invalid account claims: %w", err)
	}

	accountJWT, err := accountClaims.Encode(operatorKeyPair)
	if err != nil {
		return "", fmt.Errorf("encoding account claims: %w", err)
	}

	return accountJWT, nil
}
