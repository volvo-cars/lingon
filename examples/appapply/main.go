// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package main

//go:generate go run github.com/volvo-cars/lingon/cmd/kygo -in ../../pkg/kube/testdata/grafana.yaml -out out/grafana -app grafana -group=false -clean-name=false
//go:generate go run github.com/volvo-cars/lingon/cmd/explode -in ../../pkg/kube/testdata/grafana.yaml -out out/explode

import (
	"fmt"

	"appapply/out/grafana"

	"github.com/volvo-cars/lingon/pkg/kube"
)

func main() {
	g := grafana.New()
	if err := kube.Export(g, "export"); err != nil {
		fmt.Println(err)
		return
	}
}
