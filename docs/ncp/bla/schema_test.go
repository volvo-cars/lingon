package bla_test

import (
	"context"
	"ncp/bla"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

func TestSchema(t *testing.T) {
	ctx := context.TODO()
	ts := bla.StartTestServer(t)
	sysUserConn, err := ts.SysUserConn()
	if err != nil {
		t.Fatal("getting sys user connection: ", err)
	}
	actorUserConn, err := ts.ActorUserConn()
	if err != nil {
		t.Fatal("getting actor user connection: ", err)
	}
	// Create account for schema
	aa := bla.AccountActor{
		OperatorNKey:       ts.Auth.OperatorNKey,
		SysAccountConn:     sysUserConn,
		ActorAccountConn:   actorUserConn,
		ActorAccountPubKey: ts.Auth.ActorAccountPublicKey,
	}
	if err := bla.RegisterAccountActor(ctx, &aa); err != nil {
		t.Fatal("starting account actor: ", err)
	}
	defer aa.Close()

	actor := bla.SchemaActor{
		Conn: actorUserConn,
	}
	if err := bla.RegisterSchemaActor(ctx, &actor); err != nil {
		t.Fatal("starting schema actor: ", err)
	}
	defer actor.Close()

	accountReply, err := bla.CreateAccount(actorUserConn, bla.CreateAccountMsg{
		Name: "test",
	})
	if err != nil {
		t.Fatal("creating account: ", err)
	}
	userReply, err := bla.CreateUserForAccount(actorUserConn, bla.CreateUserMsg{
		Name: "test",
	}, accountReply.ID)
	if err != nil {
		t.Fatal("creating user: ", err)
	}

	keyPair, err := nkeys.FromSeed(userReply.NKey)
	if err != nil {
		t.Fatal("getting key pair from seed: ", err)
	}
	nc, err := nats.Connect(
		ts.NS.ClientURL(),
		bla.UserJWTOption(userReply.JWT, keyPair),
	)
	if err != nil {
		t.Fatal("connecting with user JWT: ", err)
	}

	reply, err := bla.PublishSchema(nc, bla.PublishSchemaMsg{
		Name:    "test",
		Version: "1.0.0",
		Schema:  bla.EventTxtar,
	})
	if err != nil {
		t.Fatal("publishing schema: ", err)
	}
	t.Log("schema published", "reply", reply)

	// TODO: how to create user and pass to ingester service?!?!

}
