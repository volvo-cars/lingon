// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package bla_test

import (
	"fmt"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/volvo-cars/nope/internal/bla"
	"github.com/volvo-cars/nope/internal/natsutil"
)

func TestCreateUser(t *testing.T) {
	ts := natsutil.StartTestServer(t)

	// Connect as sys user
	sysUserConn, err := nats.Connect(
		ts.NS.ClientURL(),
		bla.UserJWTOption(ts.Auth.SysUserJWT, ts.Auth.SysUserKeyPair),
	)
	if err != nil {
		t.Fatal("connecting to nats: ", err)
	}

	// Create account
	account, err := bla.SyncAccount(
		sysUserConn,
		ts.Auth.OperatorNKey,
		nil,
		bla.AccountRequest{
			Name: "MY_ACCOUNT",
		},
	)
	if err != nil {
		t.Fatal("creating account: ", err)
	}

	// Create user and test authentication
	userResult, err := bla.SyncUser(
		account.NKey,
		nil,
		bla.UserRequest{
			Name: "my_user",
		},
	)
	if err != nil {
		t.Fatal("creating user: ", err)
	}
	if err := checkUserAuth(ts.NS.ClientURL(), userResult.JWT, userResult.NKey); err != nil {
		t.Fatal("authenticating user: ", err)
	}
}

func checkUserAuth(url, userJWT string, userNKey []byte) error {
	userKeyPair, err := nkeys.FromSeed(userNKey)
	if err != nil {
		return fmt.Errorf("getting user key pair from nkey: %w", err)
	}
	unc, err := nats.Connect(
		url,
		bla.UserJWTOption(userJWT, userKeyPair),
	)
	if err != nil {
		return fmt.Errorf("connecting to NATS: %w", err)
	}
	js, err := unc.JetStream()
	if err != nil {
		return fmt.Errorf("getting jetstream context: %w", err)
	}
	if _, err := js.AccountInfo(); err != nil {
		return fmt.Errorf("getting account info: %w", err)
	}
	return nil
}
