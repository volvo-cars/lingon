// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package crd

//go:generate rm -rf out/
//go:generate go run github.com/dave/rebecca/cmd/becca@v0.9.2  -package=github.com/volvo-cars/lingon/docs/kubernetes/crd -input readme.md.tpl

//go:generate echo "\n>>>> CRD: importing YAML to Go\n"
//go:generate go test -v example_import_test.go
//go:generate go test -v import_options_test.go

//go:generate echo "\n>>>> CRD: exporting Go to YAML\n"
//go:generate go test -v example_export_test.go
