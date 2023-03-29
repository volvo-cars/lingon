// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

/*
Package kube provides APIs to convert Kubernetes manifests from and to Go code.
The main value of this package is to provide a way to manage kubernetes objects
by converting them from YAML to Go code, and back again.
This is useful when modifying and managing Go code as a source of truth
for your kubernetes platform while still being able to export to YAML which
is the defacto standard for kubernetes manifests.

Have a look at the tests for more examples.

# From YAML to Go code

The *Import* function converts Kubernetes manifests from YAML to Go code.
It is flexible as it accepts a number of options to customize the output.
All the functions starting with `With` are options.

# From Go to YAML

The Export function converts Kubernetes manifests from Go code to YAML.

# Explode manifests into separate files

The Explode function organizes Kubernetes manifests as separate files
in a directory structure. It represents more closely the way they appear in
a kubernetes cluster.
*/
package kube
