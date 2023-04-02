// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package localfile

import (
	"bytes"
	"fmt"

	"github.com/volvo-cars/lingon/pkg/terra"
)

func ExampleStack() {
	stack := NewLocalFileStack("myfile.txt")
	var b bytes.Buffer
	if err := terra.Export(stack, terra.WithExportWriter(&b)); err != nil {
		fmt.Printf("Error: exporting stack: %s", err)
		return
	}
	fmt.Println(b.String())

	// Output:
	// terraform {
	//   backend "local" {
	//     path = "terraform.tfstate"
	//   }
	//   required_providers {
	//     local = {
	//       source  = "hashicorp/local"
	//       version = "2.4.0"
	//     }
	//   }
	// }
	//
	// // Provider blocks
	// provider "local" {
	// }
	//
	//
	// // Resource blocks
	// resource "local_file" "file" {
	//   content  = "contents"
	//   filename = "myfile.txt"
	// }
}
