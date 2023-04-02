// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package crd

//go:generate echo ">>>> cleaning previous output"
//go:generate rm -rf out/
//go:generate echo ">>>> importing YAML to Go"
//go:generate go test -v example_import_test.go
//go:generate echo ">>>> exporting Go to YAML"
//go:generate go test -v example_export_test.go

//go:generate echo ">>>> generating terraform readme\n"
//go:generate go run github.com/dave/rebecca/cmd/becca@v0.9.2  -package=github.com/volvo-cars/lingon/docs/kubernetes/crd -input readme.md.tpl
