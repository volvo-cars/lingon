// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmk8s

import (
	"os"
	"testing"

	"github.com/volvo-cars/lingon/pkg/kube"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestApp(t *testing.T) {
	os.RemoveAll("out")

	app := New()
	tu.AssertNoError(
		t,
		kube.Export(app, kube.WithExportOutputDirectory("out")),
		"export",
	)
}
