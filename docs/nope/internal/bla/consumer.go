// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package bla

import (
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
)

type ConsumerRequest struct {
	Stream string
	Name   string
}

type Consumer struct {
	Name string
}

func SyncConsumer(nc *nats.Conn, consumer *Consumer, req ConsumerRequest) (*Consumer, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("getting jetstream: %w", err)
	}

	if consumer == nil {
		consumerInfo, err := js.AddConsumer(req.Stream, &nats.ConsumerConfig{
			Durable:   req.Name,
			AckPolicy: nats.AckExplicitPolicy,
		})
		if err != nil {
			return nil, fmt.Errorf("creating consumer: %w", err)
		}
		return &Consumer{
			Name: consumerInfo.Name,
		}, nil
	}
	consumerInfo, err := js.UpdateConsumer(req.Stream, &nats.ConsumerConfig{
		Durable:   req.Name,
		AckPolicy: nats.AckExplicitPolicy,
	})
	if err != nil {
		if errors.Is(err, nats.ErrConsumerNotFound) {
			return SyncConsumer(nc, nil, req)
		}
		return nil, fmt.Errorf("updating consumer: %w", err)
	}
	return &Consumer{
		Name: consumerInfo.Name,
	}, nil
}

func DeleteConsumer(nc *nats.Conn, req ConsumerRequest) error {
	js, err := nc.JetStream()
	if err != nil {
		return fmt.Errorf("getting jetstream: %w", err)
	}
	if err := js.DeleteConsumer(req.Stream, req.Name); err != nil {
		return fmt.Errorf("deleting consumer: %w", err)
	}
	return nil
}
