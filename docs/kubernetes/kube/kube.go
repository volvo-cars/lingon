// Copyright (c) Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

//go:generate rm -rf out/
//go:generate go run github.com/dave/rebecca/cmd/becca@v0.9.2  -package=github.com/volvo-cars/lingon/docs/kubernetes/kube -input=readme.md.tpl -output=readme.md

//go:generate echo "\n>>>> KUBE: exploding manifests\n"
//go:generate go run github.com/volvo-cars/lingon/cmd/explode -in ../../../pkg/kube/testdata/tekton.yaml -out out/explode

//go:generate echo "\n>>>> KUBE: importing YAML to Go\n"
//go:generate go run github.com/volvo-cars/lingon/cmd/kygo -in ../../../pkg/kube/testdata/tekton.yaml -out out/tekton -app=tekton -group -clean-name

//go:generate echo "\n>>>> KUBE: exporting Go to YAML\n"
//go:generate go test -v .
