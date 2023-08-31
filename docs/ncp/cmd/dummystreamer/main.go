package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
)

func main() {
	var (
		name    string
		version string
	)
	flag.StringVar(&name, "name", "", "name of the schema")
	flag.StringVar(&version, "version", "", "version of the schema")
	flag.Parse()
	if name == "" {
		slog.Error("name argument is required")
		os.Exit(1)
	}
	if version == "" {
		slog.Error("version argument is required")
		os.Exit(1)
	}
	natsCreds, ok := os.LookupEnv("NATS_CREDS")
	if !ok {
		slog.Error("NATS_CREDS not set")
		os.Exit(1)
	}

	ctx := context.TODO()
	subject := fmt.Sprintf("ingest.%s.%s", name, strings.ReplaceAll(version, ".", "_"))
	nc, err := nats.Connect(nats.DefaultURL, nats.UserCredentials(natsCreds))
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
		ack, err := js.Publish(ctx, subject, b)
		if err != nil {
			slog.Error("publishing", "error", err)
			os.Exit(1)
		}
		slog.Info("published", "ack", ack, "subject", subject)
	}
}
