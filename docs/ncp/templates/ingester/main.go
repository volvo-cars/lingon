package main

import (
	"context"
	"errors"
	"flag"
	"os"

	"ingester/schema"

	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
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
		natsURL = nats.DefaultURL
		// cErr = errors.Join(cErr, errors.New("NATS_URL not set"))
	}
	// natsCreds, ok := os.LookupEnv("NATS_CREDS")
	// if !ok {
	// 	cErr = errors.Join(cErr, errors.New("NATS_CREDS not set"))
	// }
	if cErr != nil {
		slog.Error("required environment variables not set", "error", cErr)
		os.Exit(1)
	}
	// nc, err := nats.Connect(natsURL, nats.UserCredentials(natsCreds))
	nc, err := nats.Connect(natsURL)
	if err != nil {
		slog.Error("connecting to NATS", "error", err)
		os.Exit(1)
	}
	defer nc.Close()

	ctx := context.TODO()

	js, err := jetstream.New(nc)
	if err != nil {
		slog.Error("getting JetStream instance", "error", err)
		os.Exit(1)
	}

	cons, err := js.CreateOrUpdateConsumer(ctx, stream, jetstream.ConsumerConfig{
		Name:          consumer,
		Durable:       consumer,
		DeliverPolicy: jetstream.DeliverAllPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		slog.Error("getting consumer", "error", err)
		os.Exit(1)
	}

	iter, _ := cons.Messages()
	defer iter.Stop()
	for {
		msg, err := iter.Next()
		if err != nil {
			slog.Error("getting next message", "error", err)
			os.Exit(1)
		}
		var event schema.Event
		if err := proto.Unmarshal(msg.Data(), &event); err != nil {
			slog.Info("terminating message unmarshalling protobuf message failed", "msg_subject", msg.Subject)
			if err := msg.Term(); err != nil {
				slog.Error("terminating message", "error", err)
			}
			continue
		}
		slog.Info("acknowledging message", "event", event.String())
		if err := msg.Ack(); err != nil {
			slog.Error("acknowledging message", "error", err)
		}
	}
}
