// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubernetes

//go:generate echo ">>>> generating kubernetes readme\n"
//go:generate go run github.com/dave/rebecca/cmd/becca@v0.9.2  -package=github.com/volvo-cars/lingon/docs/kubernetes/crd -input readme.md.tpl
//go:generate go run github.com/dave/rebecca/cmd/becca@v0.9.2  -package=github.com/volvo-cars/lingon/docs/kubernetes/crd -input=optionsimport.md.tpl -output=options-import.md
//go:generate go run github.com/dave/rebecca/cmd/becca@v0.9.2  -package=github.com/volvo-cars/lingon/docs/kubernetes/crd -input=optionsexport.md.tpl -output=options-export.md
