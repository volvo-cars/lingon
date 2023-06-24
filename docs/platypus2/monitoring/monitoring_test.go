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
)

func TestMonitoringExport(t *testing.T) {
	tests := map[string]kube.Exporter{
		"out/1_promcrd":      promcrd.New(),
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

// // TODO: THIS IS INTEGRATION and needs KWOK
// func TestMonitoringDeploy(t *testing.T) {
// 	ctx := context.Background()
//
// 	// pcrd := promcrd.New()
// 	// tu.AssertNoError(t, pcrd.Apply(ctx), "prometheus crd")
//
// 	// ms := metricsserver.New()
// 	// tu.AssertNoError(t, ms.Apply(ctx), "metrics-server")
// 	//
// 	// ps := promstack.New()
// 	// tu.AssertNoError(t, ps.Apply(ctx), "prometheus stack")
//
// 	vmcrds := vmcrd.New()
// 	tu.AssertNoError(t, vmcrds.Apply(ctx), "victoria metrics crds")
//
// 	vm := vmk8s.New()
// 	tu.AssertNoError(t, vm.Apply(ctx), "victoria metrics stack")
// }
