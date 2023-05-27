// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package externalsecrets

import (
	"os"
	"strings"
	"testing"

	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestExport(t *testing.T) {
	_ = os.RemoveAll("out")

	app := New()
	err := kube.Export(app, kube.WithExportOutputDirectory("out"))
	tu.AssertNoError(t, err)

	ly, err := ku.ListYAMLFiles("out")

	err = kube.Import(
		kube.WithImportAppName(AppName),
		kube.WithImportManifestFiles(ly),
		kube.WithImportPackageName(strings.ReplaceAll(AppName, "-", "")),
		kube.WithImportRemoveAppName(true),
		kube.WithImportOutputDirectory("out/go"),
	)
}
