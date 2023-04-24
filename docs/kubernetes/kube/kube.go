// Copyright (c) Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

//go:generate rm -rf out/
//go:generate go run github.com/dave/rebecca/cmd/becca@v0.9.2  -package=github.com/volvo-cars/lingon/docs/kubernetes/kube -input=readme.md.tpl -output=readme.md
//go:generate echo ">>>> exploding manifests"
//go:generate go run github.com/volvo-cars/lingon/cmd/explode -in ../../../pkg/kube/testdata/tekton.yaml -out out/explode
//go:generate echo ">>>> importing YAML to Go"
//go:generate go run github.com/volvo-cars/lingon/cmd/kygo -in ../../../pkg/kube/testdata/tekton.yaml -out out/tekton -app=tekton -group -clean-name
//go:generate echo ">>>> exporting Go to YAML"
//go:generate go test -v .
