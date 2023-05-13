// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Lingon is a collection of libraries and tools for building platforms using Go.
//
// The following technologies are currently supported:
//   - Terraform
//   - Kubernetes
//
// The only dependencies you need are:
//
//   - Go
//   - Terraform CLI
//   - kubectl
//
// The packages currently available are:
//
// # Kube
//
//   - [kube.App] struct that is embedded to mark kubernetes applications
//   - [kube.Export] converts kubernetes objects defined as Go struct to kubernetes manifests in YAML.
//   - [kube.Explode] kubernetes manifests in YAML to multiple files, organized by namespace.
//   - [kube.Import] converts kubernetes manifests in YAML to Go structs.
//
// # Kubeconfig
//
// utility package to read and merge kubeconfig files without any dependencies on `go-client`.
//
// # Kubeutil
//
// Reusable functions used to create and validate kubernetes objects in Go.
//
// # Terra
//
// Core functionality for working with Terraform.
//   - [terra.Export] converts a [terra.Stack] to HCL files.
//
// # Terragen
//
// Generate Go code for Terraform providers.
//   - [terragen.GenerateProviderSchema] generates a terraform provider JSON schema from terraform provider registry.
//   - [terragen.GenerateGoCode] converts a terraform provider schema to Go structs.
//
// # Testutils
//
// Reusable test functions.
//
// [kube.Export]: https://pkg.go.dev/github.com/volvo-cars/lingon/pkg/kube#Export
// [kube.Explode]: https://pkg.go.dev/github.com/volvo-cars/lingon/pkg/kube#Explode
// [kube.Import]: https://pkg.go.dev/github.com/volvo-cars/lingon/pkg/kube#Import
// [kube.App]: https://pkg.go.dev/github.com/volvo-cars/lingon/pkg/kube#App
// [terra.Export]: https://pkg.go.dev/github.com/volvo-cars/lingon/pkg/terra#Export
// [terragen.GenerateProviderSchema]: https://pkg.go.dev/github.com/volvo-cars/lingon/pkg/terragen#GenerateProviderSchema
// [terragen.GenerateGoCode]: https://pkg.go.dev/github.com/volvo-cars/lingon/pkg/terragen#GenerateGoCode
package main
