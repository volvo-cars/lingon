package jsm

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"
)

func TestJSM(t *testing.T) {
	nc, err := nats.Connect(nats.DefaultURL, nats.UserCredentials("../../test.creds"))
	// nc, err := nats.Connect(nats.DefaultURL, nats.UserCredentials("../../test.creds"))
	if err != nil {
		t.Fatal("connecting to nats: ", err)
	}
	defer nc.Close()
	mgr, err := jsm.New(
		nc,
		// jsm.WithAPIPrefix("$JS.API.ABKE7EQ5PSVCM7LCJJQ4PHQ5RWFOP4G57YIJ7FYD5L5BKRV6ONSY44DO"),
		// jsm.WithAPIPrefix("$JS.API.AAVVR2R4XOBQWELPJL5A7WMLCBSRVJVI4KT7BIHVDA35RZDGQMZCST54"),
	)
	if err != nil {
		t.Fatal("creating jsm manager: ", err)
	}
	// ai, err := mgr.JetStreamAccountInfo()
	// if err != nil {
	// 	t.Fatal("getting jetstream account info: ", err)
	// }
	// fmt.Println("CONSUMERS: ", ai.Consumers)
	streams, err := mgr.StreamNames(&jsm.StreamNamesFilter{})
	if err != nil {
		t.Fatal("getting stream names: ", err)
	}

	// stream, err := mgr.LoadStream("event:1_0_7")
	stream, err := mgr.LoadStream("event")
	if err != nil {
		t.Fatal("loading stream: ", err)
	}
	consumers, err := stream.ConsumerNames()
	if err != nil {
		t.Fatal("getting consumer names: ", err)
	}
	si, err := stream.LatestInformation()
	if err != nil {
		t.Fatal("getting stream latest information: ", err)
	}
	// si.State.Consumers
	fmt.Println("STREAMS: ", streams)
	fmt.Printf("CONSUMERS FOR STREAM \"%s\": %+v\n", stream.Name(), consumers)
	fmt.Println("SUBJECTS: ", si.Config.Subjects)
	fmt.Printf("%+v\n", si.State)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	cons, err := stream.LoadConsumer("1_0_0")
	if err != nil {
		t.Fatal("loading consumer: ", err)
	}
	ci, err := cons.State()
	if err != nil {
		t.Fatal("getting consumer latest state: ", err)
	}
	b, err := json.MarshalIndent(ci, "", "  ")
	if err != nil {
		t.Fatal("marshaling consumer state: ", err)
	}
	fmt.Println(string(b))
	// fmt.Printf("CONSUMER: %+v\n", ci.Delivered)
	// fmt.Printf("CONSUMER: %+v\n", ci)
	// msg, _ := stream.ReadMessage(1)
	// // msg.
	// // si.State.Lost
	// fmt.Println("LOST: ", si.State.Lost)

	//
	// IMPORTANT: how to get last message for a subject and check that a message
	// actually exists
	//
	// msg, err := stream.ReadLastMessageForSubject(fmt.Sprintf("ingest.%s.%s", "event", "1_0_0"))
	// if err != nil {
	// 	if !jsm.IsNatsError(err, uint16(nats.JSErrCodeMessageNotFound)) {
	// 		t.Fatal("reading last message: ", err)
	// 	}
	// 	t.Log("MESSAGE NOT FOUND")
	// }
	// fmt.Println(msg)
}
