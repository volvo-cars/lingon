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
)

func TestMonitoringExport(t *testing.T) {
	folders := []string{
		"out/1_promcrd",
		"out/2_metrics-server",
		"out/3_promstack",
		"out/4_surveyor",
	}
	for _, f := range folders {
		_ = os.RemoveAll(f)
	}

	pcrd := promcrd.New()
	if err := pcrd.Export(folders[0]); err != nil {
		tu.AssertNoError(t, err, "prometheus crd")
	}
	ms := metricsserver.New()
	if err := ms.Export(folders[1]); err != nil {
		tu.AssertNoError(t, err, "metrics-server")
	}

	ps := promstack.New()
	if err := ps.Export(folders[2]); err != nil {
		tu.AssertNoError(t, err, "prometheus stack")
	}
}

// TODO: THIS IS INTEGRATION and needs KWOK
// func TestMonitoringDeploy(t *testing.T) {
// 	ctx := context.Background()
//
// 	pcrd := promcrd.New()
// 	if err := pcrd.Apply(ctx); err != nil {
// 		tu.AssertNoError(t, err, "prometheus crd")
// 	}
// 	ms := metricsserver.New()
// 	if err := ms.Apply(ctx); err != nil {
// 		tu.AssertNoError(t, err, "metrics-server")
// 	}
//
// 	ps := promstack.New()
// 	if err := ps.Apply(ctx); err != nil {
// 		tu.AssertNoError(t, err, "prometheus stack")
// 	}
// }
