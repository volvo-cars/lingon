// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package karpenter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/volvo-cars/lingon/pkg/kube"
)

func TestExport(t *testing.T) {
	app := New(
		Opts{
			ClusterName:            "CLUSTER_NAME",
			ClusterEndpoint:        "CLUSTER_ENDPOINT",
			IAMRoleArn:             "ROLE_ARN",
			DefaultInstanceProfile: "DEFAULT_INSTANCE_PROFILE",
			InterruptQueue:         "INTERRUPT_QUEUE",
		},
	)
	err := kube.Export(app, kube.WithExportOutputDirectory("out"))
	require.NoError(t, err)
}
