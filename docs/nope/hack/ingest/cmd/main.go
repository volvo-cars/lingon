package main

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/volvo-cars/nope/hack/ingest/templates"

	"github.com/volvo-cars/nope/hack/ingest"

	"github.com/nats-io/nats.go"
	"golang.org/x/exp/slog"
)

const (
	bucket = "schema"
)

func main() {

	var cErr error
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

	js, err := nc.JetStream()
	if err != nil {
		slog.Error("getting JetStream", "error", err)
		os.Exit(1)
	}
	// Create bucket and handle error if it already exists
	if _, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: bucket,
	}); err != nil {
		if !errors.Is(err, nats.ErrKeyExists) {
			slog.Error("creating bucket", "error", err)
			os.Exit(1)
		}
	}
	bucket, err := js.KeyValue(bucket)
	if err != nil {
		slog.Error("getting bucket", "error", err)
		os.Exit(1)
	}
	watcher, err := bucket.WatchAll(nats.IgnoreDeletes())
	if err != nil {
		slog.Error("watching bucket", "error", err)
		os.Exit(1)
	}
	defer watcher.Stop()

	slog.Info(
		"watching nats bucket",
		"bucket", bucket.Bucket(),
	)
	for kve := range watcher.Updates() {
		if kve == nil {
			continue
		}
		switch op := kve.Operation(); op {
		case nats.KeyValuePut:
			slog.Info("put", "key", kve.Key(), "value", kve.Value())

			workdir := filepath.Join("out", kve.Key(), strconv.FormatUint(kve.Revision(), 10))
			err := ingest.KoBuild(ingest.KoBuildConfig{
				Workdir:        workdir,
				GoServiceTxtar: templates.IngestionTxtar,
				Schema: ingest.Schema{
					Name:  kve.Key(),
					Txtar: kve.Value(),
				},
			})
			if err != nil {
				slog.Error("building", "error", err)
			}
		default:
			slog.Info("unknown operation", "operation", op)
		}
	}
}
