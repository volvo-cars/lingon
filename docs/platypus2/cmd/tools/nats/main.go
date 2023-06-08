// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/nats-io/nats.go"
	"github.com/volvo-cars/lingoneks/cmd/tools/nats/natspb"
	_ "go.uber.org/automaxprocs"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
		NatsServers []string `conf:"default:nats://0.0.0.0:4222"`
		GRPCPort    int      `conf:"default:7015"`
		// HTTPPort            int           `conf:"default:9000,env:PORT"`
		// HTTPHost            string        `conf:"default:0.0.0.0"`
		// PathHealth          string        `conf:"default:/healthz"`
		// PathVersion         string        `conf:"default:/version"`
		// HTTPReadTimeout     time.Duration `conf:"default:5s"`
		// HTTPWriteTimeout    time.Duration `conf:"default:10s"`
		// HTTPIdleTimeout     time.Duration `conf:"default:120s"`
		// HTTPShutdownTimeout time.Duration `conf:"default:5s"`
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

	// Closing signal

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(
		ctx,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
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
	no := nats.GetDefaultOptions()
	no.Servers = cfg.NatsServers
	no.AllowReconnect = true
	no.MaxReconnect = -1
	no.AsyncErrorCB = AsyncErrorCB(log)
	no.ReconnectedCB = ReconnectedCB(log)
	no.DisconnectedErrCB = DisconnectedErrCB(log)
	no.ClosedCB = func(c *nats.Conn) { log.Info("nats connection closed") }

	log.Info("connecting to NATS")
	nc, err := no.Connect()
	if err != nil {
		log.Error("connect to NATS", "err", err)
	}

	errCh := make(chan error, 1)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.GRPCPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	s := NewServer(nc, log)
	natspb.RegisterEnvelopeServiceServer(grpcServer, s)
	reflection.Register(grpcServer)
	go func(lis net.Listener) {
		errCh <- grpcServer.Serve(lis)
	}(lis)

outter:
	for {
		select {
		case e := <-errCh:
			log.Error("errCh", "err", e)
		case <-ctx.Done():
			log.Info("stopping service...")

			grpcServer.GracefulStop()
			break outter
		}
	}
	return nil
}
