// Copyright (c) Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import "fmt"

func ExampleNumber() {
	n := Number(1)
	fmt.Println(string(n.InternalTokens().Bytes()))
	// 	Output: 1
}
