// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

//go:generate rm -rf out/
//go:generate echo ">>>> exploding manifests"
//go:generate go run github.com/volvo-cars/lingon/cmd/explode -in ../../pkg/kube/testdata/tekton.yaml -out out/explode
//go:generate echo ">>>> importing YAML to Go"
//go:generate go run github.com/volvo-cars/lingon/cmd/kygo -in ../../pkg/kube/testdata/tekton.yaml -out out/tekton -app=tekton -group -clean-name
//go:generate echo ">>>> exporting Go to YAML"
//go:generate go run main.go

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"kube/out/tekton"

	"github.com/volvo-cars/lingon/pkg/kube"
)

func main() {
	fmt.Println("running main.go")
	tk := tekton.New()
	out := filepath.Join("out", "export")
	fmt.Printf("exporting to %s\n", out)
	if err := kube.Export(tk, out); err != nil {
		panic(err)
	}

	// or use ExportWriter
	var buf bytes.Buffer
	_ = kube.ExportWriter(tk, &buf)

	manifests := strings.Split(buf.String(), "---")

	fmt.Printf("\nexported %d manifests\n\n", len(manifests))
	fmt.Println(">>> first manifest <<<")
	if len(manifests) > 0 {
		fmt.Printf("%s\n", manifests[0])
	}

	// _ = tk.Apply(context.TODO())
}
