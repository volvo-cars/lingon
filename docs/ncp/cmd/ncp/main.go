package main

import (
	"ncp/bla"
	"os"

	"log/slog"
)

func main() {
	ts, err := bla.NewTestServer(
		bla.WithDir(os.TempDir()),
	)
	if err != nil {
		slog.Error("creating test server", "error", err)
		os.Exit(1)
	}
	if err := ts.StartUntilReady(); err != nil {
		slog.Error("starting test server", "error", err)
		os.Exit(1)
	}
	defer ts.NS.Shutdown()

	// TODO: add actors
	// ctx := context.Background()

	ts.NS.WaitForShutdown()
}
