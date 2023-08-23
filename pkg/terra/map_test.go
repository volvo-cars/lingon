// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"fmt"
	"log/slog"
)

func ExampleMap_string() {
	s := Map(
		map[string]StringValue{
			"a": String("a"),
			"b": String("b"),
		},
	)

	toks, err := s.InternalTokens()
	if err != nil {
		slog.Error("getting tokens", "err", err)
		return
	}
	fmt.Println(string(toks.Bytes()))
	// Output:
	// {
	//   "a" = "a"
	//   "b" = "b"
	// }
}

func ExampleMap_number() {
	s := Map(
		map[string]NumberValue{
			"0": Number(0),
			"1": Number(1),
		},
	)

	toks, err := s.InternalTokens()
	if err != nil {
		slog.Error("getting tokens", "err", err)
		return
	}
	fmt.Println(string(toks.Bytes()))
	// Output:
	// {
	//   "0" = 0
	//   "1" = 1
	// }
}

func ExampleMap_mixed() {
	s := Map(
		map[string]StringValue{
			"a": String("a"),
			"b": ReferenceAsString(ReferenceResource(&dummyResource{})),
		},
	)

	toks, err := s.InternalTokens()
	if err != nil {
		slog.Error("getting tokens", "err", err)
		return
	}
	fmt.Println(string(toks.Bytes()))
	// Output:
	// {
	//   "a" = "a"
	//   "b" = dummy.dummy
	// }
}
