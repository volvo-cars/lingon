// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"fmt"
	"log/slog"
)

func ExampleBool() {
	b := Bool(true)
	toks, err := b.InternalTokens()
	if err != nil {
		slog.Error("getting tokens", "err", err)
		return
	}
	fmt.Println(string(toks.Bytes()))
	// 	Output: true
}
