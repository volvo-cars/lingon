// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package terra_test

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/volvo-cars/lingon/pkg/terra"
	"golang.org/x/exp/slog"
)

// StackConfig defines a reusable stack configuration containing things like
// the backend and any providers that need to be initialised
type StackConfig struct {
	terra.Stack
	Backend  *BackendS3
	Provider *MockProvider
}

var _ terra.Backend = (*BackendS3)(nil)

// BackendS3 defines a backend configuration for our stacks to use.
// For backends, the `hcl` struct tags are required.
type BackendS3 struct {
	Bucket string `hcl:"bucket" validate:"required"`
	Key    string `hcl:"key" validate:"required"`
	// Add any other options needed
}

func (b *BackendS3) BackendType() string {
	return "s3"
}

type MyStack struct {
	// Embed our reusable custom StackConfig
	StackConfig
	// Add a mock resource
	Resource *MockResource `validate:"required"`
}

func Example_stackConfig() {
	stack := MyStack{
		StackConfig: StackConfig{
			Backend: &BackendS3{
				Bucket: "my-s3-bucket",
				Key:    "some/path/to/state",
			},
			Provider: &MockProvider{},
		},
		Resource: &MockResource{},
	}

	// Typically you would use terra.Export() and write to a file. We will
	// write to a buffer for test purposes
	var b bytes.Buffer
	if err := terra.ExportWriter(&stack, &b); err != nil {
		slog.Error("exporting stack", "err", err.Error())
		os.Exit(1)
	}
	fmt.Println(b.String())
	// Output:
	// terraform {
	//   backend "s3" {
	//     bucket = "my-s3-bucket"
	//     key    = "some/path/to/state"
	//   }
	//
	//   required_providers {
	//     mock = {
	//       source  = "mock/mock"
	//       version = "0"
	//     }
	//   }
	// }
	//
	// // Provider blocks
	// provider "mock" {
	// }
	//
	// // Resource blocks
	// resource "mock_resource" "mock" {
	// }
}

var _ terra.Provider = (*MockProvider)(nil)

type MockProvider struct{}

func (m MockProvider) LocalName() string {
	return "mock"
}

func (m MockProvider) Source() string {
	return "mock/mock"
}

func (m MockProvider) Version() string {
	return "0"
}

func (m MockProvider) Configuration() interface{} {
	return struct{}{}
}

var _ terra.Resource = (*MockResource)(nil)

// MockResource implements a dummy resource.
// Users do not need to write resources themselves,
// they should be generated using terragen.
type MockResource struct{}

func (m MockResource) Type() string {
	return "mock_resource"
}

func (m MockResource) LocalName() string {
	return "mock"
}

func (m MockResource) Configuration() interface{} {
	return struct{}{}
}

func (m MockResource) ImportState(attributes io.Reader) error {
	return nil
}
