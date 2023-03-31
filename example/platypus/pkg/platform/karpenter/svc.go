// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package karpenter

import (
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var Svc = &corev1.Service{
	TypeMeta: kubeutil.TypeServiceV1,
	ObjectMeta: metav1.ObjectMeta{
		Name:      "karpenter",
		Namespace: "karpenter",
		Labels:    commonLabels,
	},
	Spec: corev1.ServiceSpec{
		Type:     corev1.ServiceType("ClusterIP"),
		Selector: matchLabels,
		Ports: []corev1.ServicePort{
			{
				Name:     "http-metrics",
				Protocol: corev1.ProtocolTCP,
				Port:     8080,
				TargetPort: intstr.IntOrString{
					Type:   intstr.Type(1),
					StrVal: "http-metrics",
				},
			},
			{
				Name:     "https-webhook",
				Protocol: corev1.ProtocolTCP,
				Port:     443,
				TargetPort: intstr.IntOrString{
					Type:   intstr.Type(1),
					StrVal: "https-webhook",
				},
			},
		},
	},
}
