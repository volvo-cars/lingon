// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

//go:build inttest

package terraform_test

import (
	"bytes"
	"fmt"

	"github.com/volvo-cars/lingon/pkg/terra"
	"golang.org/x/exp/slog"
)

type MinimalStack struct {
	terra.Stack
}

func Example_minimalStack() {
	// Initialise the minimal stack
	stack := MinimalStack{}
	// Export the stack to Terraform HCL
	var b bytes.Buffer
	if err := terra.Export(&stack, terra.WithExportWriter(&b)); err != nil {
		slog.Error("exporting stack", "err", err)
		return
	}
	fmt.Println(b.String())

	// Output:
	// terraform {
	// }
}
