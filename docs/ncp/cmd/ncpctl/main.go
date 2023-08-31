package main

import (
	"errors"
	"flag"
	"log/slog"
	"ncp/bla"
	"os"
	"strings"

	"github.com/nats-io/nats.go"
	"golang.org/x/tools/txtar"
)

func main() {
	var (
		dir     string
		files   protoFiles
		name    string
		version string
	)
	flag.StringVar(&dir, "dir", ".", "output dir")
	flag.Var(&files, "file", "proto files")
	flag.StringVar(&name, "schema-name", "", "schema name")
	flag.StringVar(&version, "schema-version", "", "schema version")
	flag.Parse()

	var vErr error
	if len(files) == 0 {
		vErr = errors.Join(vErr, errors.New("no proto files given"))
	}
	if name == "" {
		vErr = errors.Join(vErr, errors.New("no schema name given"))
	}
	if version == "" {
		vErr = errors.Join(vErr, errors.New("no schema version given"))
	}
	natsCreds, ok := os.LookupEnv("NATS_CREDS")
	if !ok {
		vErr = errors.Join(vErr, errors.New("NATS_CREDS not set"))
	}
	if vErr != nil {
		slog.Error("validate", "error", vErr)
		os.Exit(1)
	}

	arch, err := bla.PackTxtar(dir, files)
	if err != nil {
		slog.Error("pack txtar", "error", err)
		os.Exit(1)
	}

	nc, err := nats.Connect(nats.DefaultURL, nats.UserCredentials(natsCreds))
	if err != nil {
		slog.Error("nats connect", "error", err)
		os.Exit(1)
	}

	reply, err := bla.SendSchemaPublishMsg(nc, bla.SchemaPublishMsg{
		Name:    name,
		Version: version,
		Schema:  txtar.Format(arch),
	})
	if err != nil {
		slog.Error("send schema publish", "error", err)
		os.Exit(1)
	}
	slog.Info("schema published", "reply", reply)
	// bla.UserJWTOption()
	// nats.Connect(nats.DefaultURL, nats.UserJWTAndSeed(arch.JWTAuth, arch.OperatorSeed)))
}

type protoFiles []string

func (i *protoFiles) String() string {
	return strings.Join(*i, ",")
}

func (i *protoFiles) Set(value string) error {
	*i = append(*i, value)
	return nil
}
