// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmk8s

import (
	"testing"

	"github.com/VictoriaMetrics/metricsql"
	"github.com/VictoriaMetrics/operator/api/victoriametrics/v1beta1"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestRules(t *testing.T) {
	tests := []struct {
		name  string
		rules *v1beta1.VMRule
	}{
		// node-exporter
		{
			name:  "NodeExporterRules",
			rules: NodeExporterRules,
		},
		{
			name:  "NodeExporterAlertRules",
			rules: NodeExporterAlertRules,
		},
		// kube-scheduler
		{
			name:  "KubeSchedulerRecordingRules",
			rules: KubeSchedulerRecordingRules,
		},
		{
			name:  "KubeSchedulerAlertRules",
			rules: KubeSchedulerAlertRules,
		},
		// apiserver
		{
			name:  "APIServerBurnRateRules",
			rules: APIServerBurnRateRules,
		},
		{
			name:  "APIServerAvailabilityRules",
			rules: APIServerAvailabilityRules,
		},
		{
			name:  "APIServerHistogramRules",
			rules: APIServerHistogramRules,
		},
		{
			name:  "APIServerSLOsRules",
			rules: APIServerSLOsRules,
		},
		{
			name:  "APIServerRules",
			rules: APIServerRules,
		},
		// kube-controller-manager
		{
			name:  "KubeControllerAlertRules",
			rules: KubeControllerAlertRules,
		},
		// kube-etcd
		{
			name:  "ETCDRules",
			rules: ETCDRules,
		},
		// rules_k8s
		{
			name:  "K8SGeneralAlertRules",
			rules: K8SGeneralAlertRules,
		},
		{
			name:  "K8SRecordingRules",
			rules: K8SRecordingRules,
		},
		{
			name:  "PromGeneralRules",
			rules: PromGeneralRules,
		},
		{
			name:  "NodeRecordingRules",
			rules: NodeRecordingRules,
		},
		{
			name:  "KubeletRecordingRules",
			rules: KubeletRecordingRules,
		},
		{
			name:  "KubernetesAppsAlertRules",
			rules: KubernetesAppsAlertRules,
		},
		{
			name:  "KubernetesResourcesAlertRules",
			rules: KubernetesResourcesAlertRules,
		},
		{
			name:  "KubernetesStorageAlertRules",
			rules: KubernetesStorageAlertRules,
		},
		{
			name:  "KubeletAlertRules",
			rules: KubeletAlertRules,
		},
		{
			name:  "KubernetesSystemAlertRules",
			rules: KubernetesSystemAlertRules,
		},
		{
			name:  "NodeNetworkAlertRules",
			rules: NodeNetworkAlertRules,
		},
		{
			name:  "K8SNodeRules",
			rules: K8SNodeRules,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				for _, g := range tt.rules.Spec.Groups {
					for _, r := range g.Rules {
						_, err := metricsql.Parse(r.Expr)
						tu.AssertNoError(t, err, "parsing rules for "+r.Expr)
					}
				}
			},
		)
	}
}
