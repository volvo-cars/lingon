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
		Name:      AppName,
		Namespace: Namespace,
		Labels:    commonLabels,
	},
	Spec: corev1.ServiceSpec{
		Type:     corev1.ServiceTypeClusterIP,
		Selector: matchLabels,
		Ports: []corev1.ServicePort{
			{
				Name:       PortNameMetrics,
				Port:       PortMetrics,
				TargetPort: intstr.FromString(PortNameMetrics),
			},
			{
				Name:       PortNameWebhook,
				Port:       PortWebhookSvc,
				TargetPort: intstr.FromString(PortNameWebhook),
			},
		},
	},
}
