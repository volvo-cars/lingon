// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"

	conf "github.com/ardanlabs/conf/v3"
	_ "go.uber.org/automaxprocs"
	"golang.org/x/exp/slog"
)

const serviceName = "narrow"

func main() {
	log := makeLogger(os.Stderr)

	if err := run(log); err != nil {
		log.Error("run", "err", err)
		os.Exit(1) //nolint:gocritic
	}
}

func run(log *slog.Logger) error {
	cfg := struct {
		conf.Version
		NatsServers         []string      `conf:"default:nats://127.0.0.1:4222"`
		GRPCPort            int           `conf:"default:7015"`
		HTTPPort            int           `conf:"default:9000,env:PORT"`
		HTTPHost            string        `conf:"default:0.0.0.0"`
		PathHealth          string        `conf:"default:/healthz"`
		PathVersion         string        `conf:"default:/version"`
		HTTPReadTimeout     time.Duration `conf:"default:5s"`
		HTTPWriteTimeout    time.Duration `conf:"default:10s"`
		HTTPIdleTimeout     time.Duration `conf:"default:120s"`
		HTTPShutdownTimeout time.Duration `conf:"default:5s"`
	}{
		Version: conf.Version{
			Build: commit,
			Desc:  serviceName + "web service",
		},
	}

	const prefix = serviceName
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}

		return fmt.Errorf("parsing config: %w", err)
	}

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	ms := &runtime.MemStats{}
	cfgStr, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	runtime.ReadMemStats(ms)
	log.Info(
		fmt.Sprintf("starting service... %d", time.Now().UTC().Unix()),
		slog.Int("CPU cores", runtime.NumCPU()),
		slog.String("available memory", fmt.Sprintf("%d MB", ms.Sys/1024)),
		slog.String("config", cfgStr),
	)
	defer log.Info("Service stopped")

	// NATS
	opts := nats.GetDefaultOptions()
	opts.Servers = cfg.NatsServers
	opts.AllowReconnect = true
	opts.MaxReconnect = -1
	opts.AsyncErrorCB = func(conn *nats.Conn, c *nats.Subscription, err error) {
		log.Error(
			"error processing NATS message",
			slog.Any("conn", conn),
			slog.String("subject", c.Subject),
			slog.Any("error", err),
		)
	}
	opts.ReconnectedCB = func(conn *nats.Conn) {
		log.Info("reconnection happened", slog.Any("conn", conn))
	}
	opts.DisconnectedErrCB = func(conn *nats.Conn, err error) {
		log.Info(
			"disconnection happened",
			slog.Any("conn", conn),
			slog.Any("disconnect", err),
		)
	}
	nc, err := opts.Connect()
	if err != nil {
		log.Error("connect to NATS", "err", err)
	}
	// TODO: code to handle nats msgs

	<-ctx.Done()
	log.Info("stopping service...")
	// stop nats connections
	if err := nc.Flush(); err != nil {
		return fmt.Errorf("nats flush: %w", err)
	}
	nc.Close()

	return nil
}

// makeLogger returns a logger that writes to w [io.Writer]. If w is nil, os.Stderr is used.
// Timestamp is removed and directory from the source's filename is shown.
func makeLogger(w io.Writer) *slog.Logger {
	if w == nil {
		w = os.Stderr
	}
	return slog.New(
		slog.NewTextHandler(w, &slog.HandlerOptions{AddSource: true, ReplaceAttr: logReplace}).
			WithAttrs([]slog.Attr{slog.String("app", serviceName)}),
	)
}

func logReplace(_ []string, a slog.Attr) slog.Attr {
	// Remove the directory from the source's filename.
	if a.Key == slog.SourceKey {
		source := a.Value.Any().(*slog.Source)
		source.File = filepath.Base(source.File)
	}
	return a
}
