package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"time"

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
	natsJWT, ok := os.LookupEnv("NATS_JWT")
	if !ok {
		cErr = errors.Join(cErr, errors.New("NATS_JWT not set"))
	}
	natsNKey, ok := os.LookupEnv("NATS_NKEY")
	if !ok {
		cErr = errors.Join(cErr, errors.New("NATS_NKEY not set"))
	}
	if cErr != nil {
		slog.Error("required environment variables not set", "error", cErr)
		os.Exit(1)
	}
	nc, err := nats.Connect(natsURL, nats.UserJWTAndSeed(natsJWT, natsNKey))
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

	cons, err := js.Consumer(ctx, stream, consumer)
	if err != nil {
		slog.Error("getting consumer", "error", err)
		os.Exit(1)
	}

	slog.Info("listening for messages", "stream", stream, "consumer", consumer)

	count := 0
	for {
		msg, err := cons.Next()
		if err != nil {
			// Ignore timeout
			if errors.Is(err, nats.ErrTimeout) {
				slog.Info("timeout getting next message")
				continue
			}
			slog.Error("getting next message", "error", err)
			os.Exit(1)
		}
		count++
		slog.Info("SLEEPING 10 SECONDS")
		time.Sleep(10 * time.Second)
		var event schema.Event
		if err := proto.Unmarshal(msg.Data(), &event); err != nil {
			slog.Info(
				"terminating message unmarshalling protobuf message failed",
				"msg_subject",
				msg.Subject,
			)
			if err := msg.Term(); err != nil {
				slog.Error("terminating message", "error", err)
			}
			continue
		}
		// IMPORTANT: This terminates every 5th message intentionally!!
		if count%5 == 0 {
			slog.Info("NACK message by 5th message", "count", count)
			if err := msg.NakWithDelay(time.Minute * 5); err != nil {
				slog.Error("nacking message", "error", err)
			}
			// slog.Info("Term message by 5th message", "count", count)
			// if err := msg.Term(); err != nil {
			// 	slog.Error("terminating message", "error", err)
			// }
			continue
		}
		slog.Info("acknowledging message", "event", event.String())
		if err := msg.Ack(); err != nil {
			slog.Error("acknowledging message", "error", err)
		}
	}
}
