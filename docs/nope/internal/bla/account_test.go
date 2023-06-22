package bla_test

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/volvo-cars/nope/internal/bla"
	"github.com/volvo-cars/nope/internal/natsutil"
)

func TestAccountSync(t *testing.T) {
	ts := natsutil.StartTestServer(t)

	sysUserConn, err := nats.Connect(
		ts.NS.ClientURL(),
		bla.UserJWTOption(ts.Auth.SysUserJWT, ts.Auth.SysUserKeyPair),
	)
	if err != nil {
		t.Fatal("connecting to nats: ", err)
	}
	var account *bla.Account
	t.Run("create account", func(t *testing.T) {
		account, err = bla.SyncAccount(
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
	})

	t.Run("existing account with no change", func(t *testing.T) {
		_, err := bla.SyncAccount(
			sysUserConn,
			ts.Auth.OperatorNKey,
			account,
			bla.AccountRequest{
				Name: "MY_ACCOUNT",
			},
		)
		if err != nil {
			t.Fatal("updating account: ", err)
		}
	})
	t.Run("existing account with changed name", func(t *testing.T) {
		_, err := bla.SyncAccount(
			sysUserConn,
			ts.Auth.OperatorNKey,
			account,
			bla.AccountRequest{
				Name: "SOMETHING_ELSE",
			},
		)
		if err != nil {
			t.Fatal("updating account: ", err)
		}
	})
	t.Run("existing account deleted on nats", func(t *testing.T) {
		// First delete the account on the nats server
		if err := bla.DeleteAccount(sysUserConn, ts.Auth.OperatorNKey, account.ID); err != nil {
			t.Fatal("deleting account: ", err)
		}
		// sysUserConn.Request()
		_, err := bla.SyncAccount(
			sysUserConn,
			ts.Auth.OperatorNKey,
			account,
			bla.AccountRequest{
				Name: "MY_ACCOUNT",
			},
		)
		if err != nil {
			t.Fatal("updating account: ", err)
		}
	})
}
