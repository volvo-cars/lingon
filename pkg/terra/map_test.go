// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"fmt"
)

func ExampleMap_string() {
	s := Map(
		map[string]StringValue{
			"a": String("a"),
			"b": String("b"),
		},
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
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

	fmt.Println(string(s.InternalTokens().Bytes()))
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
			"b": Reference("b").AsString(),
		},
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output:
	// {
	//   "a" = "a"
	//   "b" = b
	// }
}
