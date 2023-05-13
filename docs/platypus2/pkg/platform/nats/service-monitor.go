// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package nats

import (
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ServiceMonitor = &v1.ServiceMonitor{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "nats",
		Namespace: "monitoring",
	},
	Spec: v1.ServiceMonitorSpec{
		Endpoints: []v1.Endpoint{
			{
				Path: "/metrics",
				Port: "metrics",
			},
		},
		NamespaceSelector: v1.NamespaceSelector{Any: true},
		Selector: metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app.kubernetes.io/instance": "nats",
				"app.kubernetes.io/name":     "nats",
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "monitoring.coreos.com/v1",
		Kind:       "ServiceMonitor",
	},
}
