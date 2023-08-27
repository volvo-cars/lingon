package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
)

func main() {
	ctx := context.TODO()
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		slog.Error("connecting to NATS", "error", err)
		os.Exit(1)
	}
	defer nc.Close()
	js, err := jetstream.New(nc)
	if err != nil {
		slog.Error("connecting to JetStream", "error", err)
		os.Exit(1)
	}
	event := Event{
		Msg: &Message{
			Title:   "Some Title",
			Content: "Contents here...",
		},
	}
	b, err := proto.Marshal(&event)
	if err != nil {
		slog.Error("marshaling event", "error", err)
		os.Exit(1)
	}
	for i := 0; i < 10; i++ {
		ack, err := js.Publish(ctx, "ingest.test", b)
		if err != nil {
			slog.Error("publishing", "error", err)
			os.Exit(1)
		}
		slog.Info("published", "ack", ack)
	}
}
