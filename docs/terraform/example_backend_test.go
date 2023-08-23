// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

//go:build inttest

package terraform_test

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/volvo-cars/lingon/pkg/terra"
)

var _ terra.Backend = (*BackendLocal)(nil)

// BackendLocal implements the Terraform local backend type.
// https://developer.hashicorp.com/terraform/language/settings/backends/local
type BackendLocal struct {
	Path string `hcl:"path,attr" validate:"required"`
}

// BackendType defines the type of the backend
func (b *BackendLocal) BackendType() string {
	return "local"
}

type BackendStack struct {
	terra.Stack

	Backend *BackendLocal
}

func Example_backendLocal() {
	stack := BackendStack{
		Backend: &BackendLocal{Path: "terraform.tfstate"},
	}
	// Export the stack to Terraform HCL
	var b bytes.Buffer
	if err := terra.Export(&stack, terra.WithExportWriter(&b)); err != nil {
		slog.Error("exporting stack", "err", err)
		return
	}
	fmt.Println(b.String())

	// Output:
	// terraform {
	//   backend "local" {
	//     path = "terraform.tfstate"
	//   }
	// }
}
