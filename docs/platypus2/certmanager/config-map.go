// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package certmanager

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var WebhookCM = &corev1.ConfigMap{
	Data: nil,
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "webhook",
			"app.kubernetes.io/component":  "webhook",
			"app.kubernetes.io/instance":   "cert-manager",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "webhook",
			"app.kubernetes.io/version":    "v1.12.2",
			"helm.sh/chart":                "cert-manager-v1.12.2",
		},
		Name:      "cert-manager-webhook",
		Namespace: "cert-manager",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ConfigMap",
	},
}