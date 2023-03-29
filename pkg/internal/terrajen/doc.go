// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Package terrajen implements a Go code generator for terraform.
//
// A Terraform Provider Schemas object is needed to generate code. This can be obtained using the
// [Terraform CLI]. We have a helper for creating a provider schemas object:
// [github.com/volvo-cars/lingon/pkg/terragen.GenerateProvidersSchema]
//
// This package leverages
// [github.com/dave/jennifer] for generating Go code.
//
// [Terraform CLI]: https://developer.hashicorp.com/terraform/cli/commands/providers/schema
package terrajen
