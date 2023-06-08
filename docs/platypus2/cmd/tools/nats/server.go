// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"

	"github.com/nats-io/nats.go"
	"github.com/volvo-cars/lingoneks/cmd/tools/nats/natspb"
	"golang.org/x/exp/slog"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	natspb.UnimplementedEnvelopeServiceServer
	Nats *nats.Conn
	log  *slog.Logger
}

func NewServer(nc *nats.Conn, log *slog.Logger) *Server {
	server := &Server{
		Nats: nc,
		log:  log,
		// TODO prometheus, health check
	}
	return server
}

func (s *Server) Ingest(ctx context.Context, m *natspb.Envelope) (
	*natspb.Ack, error,
) {
	data, err := proto.Marshal(m.Msg)
	if err != nil {
		return nil, err
	}
	if m.Topic == "" {
		return nil, errors.New("topic missing")
	}
	if err = s.Nats.Publish(m.Topic, data); err != nil {
		return nil, err
	}

	resp := &natspb.Ack{
		Topic: m.Topic,
		Id:    m.Msg.Id,
	}
	return resp, nil
}

func AsyncErrorCB(log *slog.Logger) func(
	*nats.Conn,
	*nats.Subscription,
	error,
) {
	return func(conn *nats.Conn, c *nats.Subscription, err error) {
		log.Error(
			"error processing NATS message",
			slog.Any("conn", conn),
			slog.String("subject", c.Subject),
			slog.Any("error", err),
		)
	}
}

func ReconnectedCB(log *slog.Logger) func(conn *nats.Conn) {
	return func(conn *nats.Conn) {
		log.Info("reconnection happened", slog.Any("conn", conn))
	}
}

func DisconnectedErrCB(log *slog.Logger) func(conn *nats.Conn, err error) {
	return func(conn *nats.Conn, err error) {
		log.Info(
			"disconnection happened",
			slog.Any("conn", conn),
			slog.Any("disconnect", err),
		)
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	// stop nats connections

	// drain publishers
	if err := s.Nats.Drain(); err != nil {
		s.log.Error("nats drain", "err", err)
	}
	// flush subscriptions
	if err := s.Nats.Flush(); err != nil {
		s.log.Error("nats flush", "err", err)
	}
	s.Nats.Close()
}
