// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra_test

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/volvo-cars/lingon/pkg/terra"
)

// StackConfig defines a reusable stack configuration containing things like
// the backend and any providers that need to be initialised
type StackConfig struct {
	terra.Stack
	Backend  *BackendS3
	Provider *DummyProvider
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
	// Add a dummy resource
	Resource *DummyResource `validate:"required"`
}

func Example_stackConfig() {
	stack := MyStack{
		StackConfig: StackConfig{
			Backend: &BackendS3{
				Bucket: "my-s3-bucket",
				Key:    "some/path/to/state",
			},
			Provider: &DummyProvider{},
		},
		Resource: &DummyResource{},
	}

	// Typically you would use terra.Export() and write to a file. We will
	// write to a buffer for test purposes
	var b bytes.Buffer
	if err := terra.Export(&stack, terra.WithExportWriter(&b)); err != nil {
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
	//   required_providers {
	//     dummy = {
	//       source  = "dummy/dummy"
	//       version = "0"
	//     }
	//   }
	// }
	//
	// // Provider blocks
	// provider "dummy" {
	// }
	//
	// // Resource blocks
	// resource "dummy_resource" "dummy" {
	// }
}

var _ terra.Provider = (*DummyProvider)(nil)

type DummyProvider struct{}

func (m DummyProvider) LocalName() string {
	return "dummy"
}

func (m DummyProvider) Source() string {
	return "dummy/dummy"
}

func (m DummyProvider) Version() string {
	return "0"
}

func (m DummyProvider) Configuration() interface{} {
	return struct{}{}
}

var _ terra.Resource = (*DummyResource)(nil)

// DummyResource implements a dummy resource.
// Users do not need to write resources themselves,
// they should be generated using terragen.
type DummyResource struct{}

func (m DummyResource) Type() string {
	return "dummy_resource"
}

func (m DummyResource) LocalName() string {
	return "dummy"
}

func (m DummyResource) Configuration() interface{} {
	return struct{}{}
}

func (m DummyResource) Dependencies() terra.Dependencies {
	return nil
}

func (m DummyResource) LifecycleManagement() *terra.Lifecycle {
	return nil
}

func (m DummyResource) ImportState(attributes io.Reader) error {
	return nil
}
