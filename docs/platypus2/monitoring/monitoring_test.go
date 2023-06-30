// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package monitoring

import (
	"os"
	"testing"

	"github.com/volvo-cars/lingon/pkg/kube"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	"github.com/volvo-cars/lingoneks/monitoring/metricsserver"
	"github.com/volvo-cars/lingoneks/monitoring/promcrd"
	"github.com/volvo-cars/lingoneks/monitoring/vmk8s"
	"github.com/volvo-cars/lingoneks/monitoring/vmop"
	"github.com/volvo-cars/lingoneks/monitoring/vmop/vmcrd"
)

func TestMonitoringExport(t *testing.T) {
	tests := map[string]kube.Exporter{
		"out/1_promcrd":      promcrd.New(),
		"out/1_vmcrd":        vmcrd.New(),
		"out/2_vmop":         vmop.New(),
		"out/metrics-server": metricsserver.New(),
		"out/vmk8s":          vmk8s.New(),
	}
	for f, km := range tests {
		_ = os.RemoveAll(f)

		tu.AssertNoError(
			t,
			kube.Export(km, kube.WithExportOutputDirectory(f)),
			f,
		)
	}
}
