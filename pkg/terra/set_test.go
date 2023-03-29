// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"fmt"
)

func ExampleSet_string() {
	s := Set(
		String("a"),
		String("b"),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: ["a", "b"]
}

func ExampleSet_number() {
	s := Set(
		Number(0),
		Number(1),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: [0, 1]
}

func ExampleSet_bool() {
	s := Set(
		Bool(false),
		Bool(true),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: [false, true]
}

func ExampleSet_ref() {
	s := Set(
		Reference("a").AsString(),
		Reference("b").AsString(),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: [a, b]
}

func ExampleSet_mixed() {
	s := Set(
		String("a"),
		Reference("a").AsString(),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: ["a", a]
}
