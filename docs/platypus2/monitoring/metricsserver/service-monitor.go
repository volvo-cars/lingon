// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package metricsserver

import (
	promoperator "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ServiceMonitor = &promoperator.ServiceMonitor{
	ObjectMeta: metav1.ObjectMeta{
		Labels:    BaseLabels(),
		Name:      "metrics-server",
		Namespace: namespace,
	},
	Spec: promoperator.ServiceMonitorSpec{
		Endpoints: []promoperator.Endpoint{
			{
				Interval:      promoperator.Duration("1m"),
				Path:          "/metrics",
				Port:          "https",
				Scheme:        "https",
				ScrapeTimeout: promoperator.Duration("10s"),
				TLSConfig:     &promoperator.TLSConfig{SafeTLSConfig: promoperator.SafeTLSConfig{InsecureSkipVerify: true}},
			},
		},
		JobLabel:          "metrics-server",
		NamespaceSelector: promoperator.NamespaceSelector{MatchNames: []string{namespace}},
		Selector: metav1.LabelSelector{
			MatchLabels: matchLabels,
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "monitoring.coreos.com/v1",
		Kind:       "ServiceMonitor",
	},
}
