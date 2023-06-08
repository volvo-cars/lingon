// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package karpenter

import (
	"os"
	"testing"

	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestExport(t *testing.T) {
	_ = os.RemoveAll("out")

	app := New(
		Opts{
			ClusterName:            "REPLACE_ME_CLUSTER_NAME",
			ClusterEndpoint:        "REPLACE_ME_CLUSTER_ENDPOINT",
			IAMRoleArn:             "REPLACE_ME_ROLE_ARN",
			DefaultInstanceProfile: "REPLACE_ME_DEFAULT_INSTANCE_PROFILE",
			InterruptQueue:         "REPLACE_ME_INTERRUPT_QUEUE",
		},
	)
	err := kube.Export(app, kube.WithExportOutputDirectory("out"))
	tu.AssertNoError(t, err)

	ly, err := ku.ListYAMLFiles("out")

	err = kube.Import(
		kube.WithImportAppName("karpenter"),
		kube.WithImportManifestFiles(ly),
		kube.WithImportPackageName("karpenter"),
		kube.WithImportRemoveAppName(true),
		kube.WithImportOutputDirectory("out/go"),
	)
}
