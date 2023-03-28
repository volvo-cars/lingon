// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"fmt"
)

func ExampleList_string() {
	s := List(
		String("a"),
		String("b"),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: ["a", "b"]
}

func ExampleList_number() {
	s := List(
		Number(0),
		Number(1),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: [0, 1]
}

func ExampleList_bool() {
	s := List(
		Bool(false),
		Bool(true),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: [false, true]
}

func ExampleList_ref() {
	s := List(
		Reference("a").AsString(),
		Reference("b").AsString(),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: [a, b]
}

func ExampleList_mixed() {
	s := List(
		String("a"),
		Number(1).AsString(),
		Reference("a").AsString(),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: ["a", "1", a]
}
