// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmk8s

import (
	"github.com/VictoriaMetrics/operator/api/victoriametrics/v1beta1"
	"github.com/volvo-cars/lingon/pkg/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type MonKubeProxy struct {
	kube.App

	KubeProxySVC    *corev1.Service
	KubeProxyScrape *v1beta1.VMServiceScrape
}

func NewMonKubeProxy() *MonKubeProxy {
	return &MonKubeProxy{
		KubeProxySVC:    KubeProxySVC,
		KubeProxyScrape: KubeProxyScrape,
	}
}

var KubeProxySVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "vmk8s-victoria-metrics-k8s-stack-kube-proxy",
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
			"jobLabel":                     "kube-proxy",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-kube-proxy",
		Namespace: "kube-system",
	},
	Spec: corev1.ServiceSpec{
		ClusterIP: "None",
		Ports: []corev1.ServicePort{
			{
				Name:       "http-metrics",
				Port:       int32(10249),
				Protocol:   corev1.Protocol("TCP"),
				TargetPort: intstr.IntOrString{IntVal: int32(10249)},
			},
		},
		Selector: map[string]string{"k8s-app": "kube-proxy"},
		Type:     corev1.ServiceType("ClusterIP"),
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	},
}

var KubeProxyScrape = &v1beta1.VMServiceScrape{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-kube-proxy",
		Namespace: "monitoring",
	},
	Spec: v1beta1.VMServiceScrapeSpec{
		Endpoints: []v1beta1.Endpoint{
			{
				BearerTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
				Port:            "http-metrics",
				Scheme:          "https",
				TLSConfig:       &v1beta1.TLSConfig{CAFile: "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"},
			},
		},
		JobLabel:          "jobLabel",
		NamespaceSelector: v1beta1.NamespaceSelector{MatchNames: []string{"kube-system"}},
		Selector: metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app":                        "vmk8s-victoria-metrics-k8s-stack-kube-proxy",
				"app.kubernetes.io/instance": "vmk8s",
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "operator.victoriametrics.com/v1beta1",
		Kind:       "VMServiceScrape",
	},
}
