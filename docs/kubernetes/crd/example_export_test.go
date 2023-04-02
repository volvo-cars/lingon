// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package crd_test

import (
	"fmt"
	"os"
	"path/filepath"

	team "github.com/volvo-cars/lingon/docs/kubernetes/crd/out"
)

var defaultOut = "out"

// Example_export to shows how to export to kubernetes manifests files in YAML.
func Example_export() {
	out := filepath.Join(defaultOut, "manifests")

	// Remove previously generated output directory
	_ = os.RemoveAll(out)

	tm := team.New()
	if err := tm.Export(out); err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println("successfully exported CRDs to manifests")

	// Output:
	// successfully exported CRDs to manifests
}
