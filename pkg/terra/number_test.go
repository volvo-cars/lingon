// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"fmt"
	"log/slog"
)

func ExampleNumber() {
	n := Number(1)
	toks, err := n.InternalTokens()
	if err != nil {
		slog.Error("getting tokens", "err", err)
		return
	}
	fmt.Println(string(toks.Bytes()))
	// 	Output: 1
}
