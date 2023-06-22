// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package bla

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

func UserJWTOption(userJWT string, userKeyPair nkeys.KeyPair) nats.Option {
	return nats.UserJWT(
		func() (string, error) {
			return userJWT, nil
		},
		func(bytes []byte) ([]byte, error) {
			return userKeyPair.Sign(bytes)
		},
	)
}
