// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import "fmt"

func ExampleString() {
	s := String("hello world")
	fmt.Println(exampleTokensOrError(s))
	// 	Output: "hello world"
}
