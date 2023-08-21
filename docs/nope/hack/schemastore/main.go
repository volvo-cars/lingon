package main

import (
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/volvo-cars/nope/internal/bla"
	"golang.org/x/exp/slog"
	"golang.org/x/tools/txtar"
)

const (
	bucket = "schema"
)

func main() {
	var (
		schema string
		dir    string
		in     fileList
	)

	flag.StringVar(&schema, "schema", "", "name of the schema")
	flag.StringVar(&dir, "dir", "", "relative directory for the input files")
	flag.Var(&in, "in", "file to include in the schema")
	flag.Parse()

	var cErr error
	if schema == "" {
		cErr = errors.Join(cErr, errors.New("schema name is required"))
	}
	if len(in) == 0 {
		cErr = errors.Join(cErr, errors.New("input file name is required"))
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
	ar, err := bla.PackTxtar(dir, in)
	if err != nil {
		slog.Error("packing txtar", "error", err)
		os.Exit(1)
	}

	rev, err := bucket.Put(schema, txtar.Format(ar))
	if err != nil {
		slog.Error("putting schema", "error", err)
		os.Exit(1)
	}
	slog.Info(
		"schema stored",
		"bucket", bucket,
		"schema", schema,
		"revision", rev,
	)
}

type fileList []string

func (i *fileList) String() string {
	return strings.Join(*i, ", ")
}

func (i *fileList) Set(value string) error {
	*i = append(*i, value)
	return nil
}
