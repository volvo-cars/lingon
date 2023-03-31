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
// # Kube
//
//   - [kube.App] struct that is embedded to mark kubernetes applications
//   - [kube.Export] kubernetes objects defined as Go struct to kubernetes manifests in YAML.
//   - [kube.Explode] kubernetes manifests in YAML to multiple files, organized by namespace.
//   - [kube.Import] kubernetes manifests in YAML to Go structs.
//
// # Kubeconfig
//
// Manipulate kubeconfig files **without** any dependencies on `go-client`.
//
// # KubeUtil
//
// Reusable functions used to create kubernetes objects in Go.
//
// # Testutils
//
// Reusable test functions.
//
// [kube.Export]: https://pkg.go.dev/github.com/lingon-lang/lingon/pkg/kube#Export
// [kube.Explode]: https://pkg.go.dev/github.com/lingon-lang/lingon/pkg/kube#Explode
// [kube.Import]: https://pkg.go.dev/github.com/lingon-lang/lingon/pkg/kube#Import
// [kube.App]: https://pkg.go.dev/github.com/lingon-lang/lingon/pkg/kube#App
package main
