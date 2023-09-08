package bla_test

import (
	"context"
	"fmt"
	"ncp/bla"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

func TestDocker(t *testing.T) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatal("creating docker client: ", err)
	}
	defer cli.Close()

	conts, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		t.Fatal("listing containers: ", err)
	}
	fmt.Println(conts)
}

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

	accountReply, err := bla.SendAccountCreateMsg(actorUserConn, bla.AccountCreateMsg{
		Name: "test",
	})
	if err != nil {
		t.Fatal("creating account: ", err)
	}
	userReply, err := bla.SendUserCreateForAccountMsg(
		actorUserConn,
		accountReply.ID,
		bla.UserCreateMsg{
			Name: "test",
		},
	)
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

	reply, err := bla.SendSchemaPublishMsg(nc, bla.SchemaPublishMsg{
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
