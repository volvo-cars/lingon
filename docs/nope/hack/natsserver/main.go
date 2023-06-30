// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nats-io/jwt/v2"
	"github.com/volvo-cars/nope/internal/natsutil"
	"golang.org/x/exp/slog"
)

func main() {
	var (
		out string
	)
	flag.StringVar(&out, "out", "out", "output directory")
	flag.Parse()

	absOut, err := filepath.Abs(out)
	if err != nil {
		slog.Error("getting absolute path for out: ", "error", err)
		os.Exit(1)
	}

	jwtDir := filepath.Join(absOut, "jwt")
	operatorNKeyFile := filepath.Join(absOut, "operator.nk")
	sysUserCredsFile := filepath.Join(absOut, "sys_user.creds")
	dotEnvFile := filepath.Join(absOut, ".env")

	ts, err := natsutil.NewTestServer(natsutil.WithDir(jwtDir))
	if err != nil {
		slog.Error("creating NATS test server: ", "error", err)
		os.Exit(1)
	}

	slog.Info("writing operator NKey", "file", operatorNKeyFile)
	if err := os.WriteFile(operatorNKeyFile, ts.Auth.OperatorNKey, 0o600); err != nil {
		slog.Error("writing operator nk: ", "error", err)
		os.Exit(1)
	}

	slog.Info("writing sys user creds", "file", sysUserCredsFile)
	userCreds, err := jwt.FormatUserConfig(ts.Auth.SysUserJWT, ts.Auth.SysUserNKey)
	if err != nil {
		slog.Error("formatting sys user creds: ", "error", err)
		os.Exit(1)
	}
	if err := os.WriteFile(sysUserCredsFile, userCreds, 0o600); err != nil {
		slog.Error("writing sys user creds: ", "error", err)
		os.Exit(1)
	}

	slog.Info("writing env file", "file", dotEnvFile)
	contents := fmt.Sprintf(
		"export NATS_URL=%s\nexport NATS_CREDS=%s\nexport NATS_OPERATOR_SEED=%s\n",
		ts.NS.ClientURL(),
		sysUserCredsFile,
		operatorNKeyFile,
	)
	if err := os.WriteFile(dotEnvFile, []byte(contents), 0o600); err != nil {
		slog.Error("writing env file: ", "error", err)
		os.Exit(1)
	}

	if err := ts.StartUntilReady(); err != nil {
		slog.Error("starting NATS test server: ", "error", err)
		os.Exit(1)
	}
	slog.Info("NATS test server started", "url", ts.NS.ClientURL())

	ts.NS.WaitForShutdown()
}
