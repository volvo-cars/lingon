package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"time"

	"ingestion/schema"

	"github.com/nats-io/nats.go"
	"golang.org/x/exp/slog"
	"google.golang.org/protobuf/proto"
)

func main() {
	var (
		stream   string
		consumer string
	)
	flag.StringVar(&stream, "stream", "", "name of the stream to subscribe to")
	flag.StringVar(&consumer, "consumer", "", "name of the consumer to subscribe to")
	flag.Parse()

	var cErr error
	// Check flags
	if stream == "" {
		cErr = errors.Join(cErr, errors.New("stream argument is required"))
	}
	if consumer == "" {
		cErr = errors.Join(cErr, errors.New("consumer argument is required"))
	}

	// Check environment variables
	natsURL, ok := os.LookupEnv("NATS_URL")
	if !ok {
		cErr = errors.Join(cErr, errors.New("NATS_URL not set"))
	}
	natsCreds, ok := os.LookupEnv("NATS_CREDS")
	if !ok {
		cErr = errors.Join(cErr, errors.New("NATS_CREDS not set"))
	}
	if cErr != nil {
		slog.Error("required environment variables not set", "error", cErr)
		os.Exit(1)
	}
	nc, err := nats.Connect(natsURL, nats.UserCredentials(natsCreds))
	if err != nil {
		slog.Error("connecting to NATS", "error", err)
		os.Exit(1)
	}
	defer nc.Close()

	ctx := context.TODO()
	js, err := nc.JetStream(nats.Context(ctx))
	if err != nil {
		slog.Error("getting JetStream context", "error", err)
		os.Exit(1)
	}

	sub, err := js.PullSubscribe("", consumer, nats.Bind(stream, consumer), nats.Context(ctx))
	if err != nil {
		slog.Error("subscribing", "error", err)
		os.Exit(1)
	}
	defer sub.Unsubscribe()
	slog.Info("subscribed", "stream", stream, "consumer", consumer)

	for {
		// For some reason, using context.Context doesn't work here.
		// It will timeout with "context deadline exceeded".
		msgs, err := sub.Fetch(
			10,
			nats.MaxWait(time.Hour),
			// nats.Context(ctx),
		)
		if err != nil {
			slog.Error("fetching messages", "error", err)
			os.Exit(1)
		}
		for _, msg := range msgs {
			var event schema.Event
			if err := proto.Unmarshal(msg.Data, &event); err != nil {
				slog.Info("terminating message unmarshalling protobuf message failed", "msg_subject", msg.Subject)
				if err := msg.Term(); err != nil {
					slog.Error("terminating message", "error", err)
				}
			}
			slog.Info("acknowledging message", "event", event.String())
			if err := msg.Ack(); err != nil {
				slog.Error("acknowledging message", "error", err)
			}
		}
	}
}
