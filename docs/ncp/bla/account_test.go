package bla

import (
	"context"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

func TestAccount(t *testing.T) {
	ctx := context.Background()
	ts := StartTestServer(t)
	sysUserConn, err := ts.SysUserConn()
	if err != nil {
		t.Fatal("getting sys user connection: ", err)
	}
	actorUserConn, err := ts.ActorUserConn()
	if err != nil {
		t.Fatal("getting actor user connection: ", err)
	}
	actor := AccountActor{
		OperatorNKey:       ts.Auth.OperatorNKey,
		SysAccountConn:     sysUserConn,
		ActorAccountConn:   actorUserConn,
		ActorAccountPubKey: ts.Auth.ActorAccountPublicKey,
	}
	if err := RegisterAccountActor(ctx, &actor); err != nil {
		t.Fatal("starting account actor: ", err)
	}
	defer actor.Close()

	reply, err := CreateAccount(actorUserConn, CreateAccountMsg{
		Name: "test",
	})
	if err != nil {
		t.Fatal("creating account: ", err)
	}
	userReply, err := CreateUserForAccount(actorUserConn, CreateUserMsg{
		Name: "test",
	}, reply.ID)
	if err != nil {
		t.Fatal("creating user: ", err)
	}
	// Test the connection for the new user
	keyPair, err := nkeys.FromSeed(userReply.NKey)
	if err != nil {
		t.Fatal("getting key pair from seed: ", err)
	}
	nc, err := nats.Connect(
		ts.NS.ClientURL(),
		UserJWTOption(userReply.JWT, keyPair),
	)
	if err != nil {
		t.Fatal("connecting with user JWT: ", err)
	}
	defer nc.Close()
	if nc.Status() != nats.CONNECTED {
		t.Fatal("not connected")
	}

}
