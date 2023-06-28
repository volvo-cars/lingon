package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"streamer/schema"

	"github.com/nats-io/nats.go"
	"golang.org/x/exp/slog"
	"google.golang.org/protobuf/proto"
)

func main() {
	var subject string
	flag.StringVar(&subject, "subject", "", "subject to publish to")
	flag.Parse()

	var cErr error
	if subject == "" {
		cErr = errors.Join(cErr, errors.New("subject argument is required"))
	}
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

	js, err := nc.JetStream()
	if err != nil {
		slog.Error("getting JetStream context", "error", err)
		os.Exit(1)
	}

	event := schema.Event{
		Msg: &schema.Message{
			AuthorId: "dummystreamer",
			Title:    "Dummy streamer",
			Content:  "Dummy streamer content",
		},
	}
	ticker := time.NewTicker(2 * time.Second)
	eventID := 0
	for _ = range ticker.C {
		for i := 0; i < 50; i++ {
			slog.Info("publishing event", "event", eventID)
			event.Msg.Id = fmt.Sprintf("event-%d", eventID)
			msg, err := proto.Marshal(&event)
			if err != nil {
				slog.Error("marshalling event", "error", err)
			}
			if _, err := js.PublishAsync(subject, msg); err != nil {
				slog.Error("publishing event", "error", err)
			}
			eventID++
		}
	}
}
