// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package bla_test

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	"github.com/volvo-cars/nope/internal/bla"
	"github.com/volvo-cars/nope/internal/natsutil"
)

func TestStream(t *testing.T) {
	ts := natsutil.StartTestServer(t)
	// In order to test streams we need to create an account and a user.
	sysUserConn, err := nats.Connect(
		ts.NS.ClientURL(),
		bla.UserJWTOption(ts.Auth.SysUserJWT, ts.Auth.SysUserKeyPair),
	)
	if err != nil {
		t.Fatal("connecting to nats: ", err)
	}
	streamAccount, err := bla.SyncAccount(sysUserConn, ts.Auth.OperatorNKey, nil, bla.AccountRequest{
		Name: "TEST_STREAM_ACCOUNT",
	})
	tu.AssertNoError(t, err, "creating account")
	// Create the Stream Account User
	streamAccountUser, err := bla.SyncUser(streamAccount.NKey, nil, bla.UserRequest{
		Name: "TEST_STREAM_USER",
	})
	tu.AssertNoError(t, err, "creating user")
	// Now we can create a stream by logging in as the stream account user.
	keyPair, err := nkeys.FromSeed(streamAccountUser.NKey)
	tu.AssertNoError(t, err, "getting stream account user keypair")
	nc, err := nats.Connect(
		ts.NS.ClientURL(),
		bla.UserJWTOption(streamAccountUser.JWT, keyPair),
	)
	tu.AssertNoError(t, err, "connecting to nats as stream account user")

	var stream *bla.Stream
	t.Run("create stream", func(t *testing.T) {
		var err error
		stream, err = bla.SyncStream(nc, nil, bla.StreamRequest{
			Name:     "MY_STREAM",
			Subjects: []string{"foo", "bar"},
		})
		tu.AssertNoError(t, err, "creating stream")
		t.Log("stream: ", stream)
	})
	t.Run("existing stream with no change", func(t *testing.T) {
		var err error
		oldStream := *stream
		stream, err = bla.SyncStream(nc, stream, bla.StreamRequest{
			Name:     stream.Name,
			Subjects: []string{"foo", "bar"},
		})
		tu.AssertNoError(t, err, "updating stream")
		tu.AssertEqual[*bla.Stream](t, &oldStream, stream)
	})
	t.Run("existing stream with change", func(t *testing.T) {
		var err error
		oldStream := *stream
		stream, err = bla.SyncStream(nc, stream, bla.StreamRequest{
			Name:     stream.Name,
			Subjects: []string{"something", "else"},
		})
		tu.AssertNoError(t, err, "updating stream")
		tu.AssertNotEqual[*nats.StreamConfig](t, &oldStream.Config, &stream.Config)
	})

	t.Run("existing stream deleted on nats", func(t *testing.T) {
		err := bla.DeleteStream(nc, stream.Name)
		tu.AssertNoError(t, err, "deleting stream")

		_, err = bla.SyncStream(nc, stream, bla.StreamRequest{
			Name:     stream.Name,
			Subjects: []string{"foo", "bar"},
		})
		tu.AssertNoError(t, err, "creating stream")
	})

	var consumer *bla.Consumer
	t.Run("create consumer", func(t *testing.T) {
		var err error
		consumer, err = bla.SyncConsumer(nc, nil, bla.ConsumerRequest{
			Stream: stream.Name,
			Name:   "MY_CONSUMER",
		})
		tu.AssertNoError(t, err, "creating consumer")
	})
	t.Run("update consumer", func(t *testing.T) {
		consumer, err = bla.SyncConsumer(nc, consumer, bla.ConsumerRequest{
			Stream: stream.Name,
			Name:   consumer.Name,
		})
		tu.AssertNoError(t, err, "creating consumer")
	})
	t.Run("update deleted consumer", func(t *testing.T) {
		delErr := bla.DeleteConsumer(nc, stream.Name, consumer.Name)
		tu.AssertNoError(t, delErr, "deleting consumer")
		var err error
		consumer, err = bla.SyncConsumer(nc, consumer, bla.ConsumerRequest{
			Stream: stream.Name,
			Name:   consumer.Name,
		})
		tu.AssertNoError(t, err, "creating consumer")
	})
}
