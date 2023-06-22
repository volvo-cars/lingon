package main

import (
	"os"

	"github.com/nats-io/jwt/v2"
	"github.com/volvo-cars/nope/internal/natsutil"
	"golang.org/x/exp/slog"
)

func main() {
	ts, err := natsutil.NewTestServer("./jwt")
	if err != nil {
		slog.Error("creating NATS test server: ", "error", err)
		os.Exit(1)
	}

	slog.Info("writing operator NKey", "file", "./operator.nk")
	if err := os.WriteFile("./operator.nk", ts.Auth.OperatorNKey, 0600); err != nil {
		slog.Error("writing operator nk: ", "error", err)
		os.Exit(1)
	}

	slog.Info("writing sys user creds", "file", "./sys_user.creds")
	userCreds, err := jwt.FormatUserConfig(ts.Auth.SysUserJWT, ts.Auth.SysUserNKey)
	if err != nil {
		slog.Error("formatting sys user creds: ", "error", err)
		os.Exit(1)
	}
	if err := os.WriteFile("./sys_user.creds", userCreds, 0600); err != nil {
		slog.Error("writing sys user creds: ", "error", err)
		os.Exit(1)
	}

	if err := ts.StartUntilReady(); err != nil {
		slog.Error("starting NATS test server: ", "error", err)
		os.Exit(1)
	}
	slog.Info("NATS test server started", "url", ts.NS.ClientURL())
	ts.NS.WaitForShutdown()
}
