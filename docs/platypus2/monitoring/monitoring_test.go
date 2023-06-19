// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package monitoring

import (
	"os"
	"testing"

	tu "github.com/volvo-cars/lingon/pkg/testutil"
	"github.com/volvo-cars/lingoneks/monitoring/metricsserver"
	"github.com/volvo-cars/lingoneks/monitoring/promcrd"
	"github.com/volvo-cars/lingoneks/monitoring/promstack"
	"github.com/volvo-cars/lingoneks/monitoring/vmcrd"
	"github.com/volvo-cars/lingoneks/monitoring/vmk8s"
)

func TestMonitoringExport(t *testing.T) {
	folders := []string{
		"out/1_promcrd",
		"out/2_metrics-server",
		"out/3_promstack",
		"out/4_vmcrd",
		"out/5_vmk8s",
	}
	for _, f := range folders {
		_ = os.RemoveAll(f)
	}

	pcrd := promcrd.New()
	tu.AssertNoError(t, pcrd.Export(folders[0]), "prometheus crd")

	ms := metricsserver.New()
	tu.AssertNoError(t, ms.Export(folders[1]), "metrics-server")

	ps := promstack.New()
	tu.AssertNoError(t, ps.Export(folders[2]), "prometheus stack")

	vmcrds := vmcrd.New()
	tu.AssertNoError(t, vmcrds.Export(folders[3]), "victoria metrics crds")

	vm := vmk8s.New()
	tu.AssertNoError(t, vm.Export(folders[4]), "victoria metrics stack")
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
