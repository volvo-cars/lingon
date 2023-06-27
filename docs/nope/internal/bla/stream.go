// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package bla

import (
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
)

type StreamRequest struct {
	Name     string
	Subjects []string
}

type Stream struct {
	Name   string
	Config nats.StreamConfig
	// TODO: are subjects needed to be stored?
	// Subjects []string
}

func SyncStream(
	nc *nats.Conn,
	stream *Stream,
	req StreamRequest,
) (*Stream, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("getting jetstream: %w", err)
	}

	if stream == nil {
		// Stream does not exist, so we create it.
		stream, err := js.AddStream(&nats.StreamConfig{
			Name:     req.Name,
			Subjects: req.Subjects,
		})
		if err != nil {
			return nil, fmt.Errorf("creating stream: %w", err)
		}

		return &Stream{
			Name:   stream.Config.Name,
			Config: stream.Config,
			// Subjects: stream.Config.Subjects,
		}, nil
	}

	// Stream should exist, so update it
	upStream, err := js.UpdateStream(&nats.StreamConfig{
		Name:     req.Name,
		Subjects: req.Subjects,
	})
	if err != nil {
		if errors.Is(err, nats.ErrStreamNotFound) {
			// Create the stream if it does not exist.
			return SyncStream(nc, nil, req)
		}
		return nil, fmt.Errorf("updating stream: %w", err)
	}

	return &Stream{
		Name:   upStream.Config.Name,
		Config: upStream.Config,
		// Subjects: upStream.Config.Subjects,
	}, nil
}

func DeleteStream(nc *nats.Conn, name string) error {
	js, err := nc.JetStream()
	if err != nil {
		return fmt.Errorf("getting jetstream: %w", err)
	}

	if err := js.DeleteStream(name); err != nil {
		return fmt.Errorf("deleting stream: %w", err)
	}

	return nil
}
