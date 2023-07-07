// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package certmanager

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var SVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "cert-manager",
			"app.kubernetes.io/component":  "controller",
			"app.kubernetes.io/instance":   "cert-manager",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "cert-manager",
			"app.kubernetes.io/version":    "v1.12.2",
			"helm.sh/chart":                "cert-manager-v1.12.2",
		},
		Name:      "cert-manager",
		Namespace: "cert-manager",
	},
	Spec: corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			{
				Name:       "tcp-prometheus-servicemonitor",
				Port:       int32(9402),
				Protocol:   corev1.Protocol("TCP"),
				TargetPort: intstr.IntOrString{IntVal: int32(9402)},
			},
		},
		Selector: map[string]string{
			"app.kubernetes.io/component": "controller",
			"app.kubernetes.io/instance":  "cert-manager",
			"app.kubernetes.io/name":      "cert-manager",
		},
		Type: corev1.ServiceType("ClusterIP"),
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	},
}

var WebhookSVC = &corev1.Service{
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
	Spec: corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			{
				Name:     "https",
				Port:     int32(443),
				Protocol: corev1.Protocol("TCP"),
				TargetPort: intstr.IntOrString{
					StrVal: "https",
					Type:   intstr.Type(int64(1)),
				},
			},
		},
		Selector: map[string]string{
			"app.kubernetes.io/component": "webhook",
			"app.kubernetes.io/instance":  "cert-manager",
			"app.kubernetes.io/name":      "webhook",
		},
		Type: corev1.ServiceType("ClusterIP"),
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	},
}