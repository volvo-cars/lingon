// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmk8s

import (
	"github.com/VictoriaMetrics/operator/api/victoriametrics/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type MonKubeController struct {
	SVC        *corev1.Service
	Scrape     *v1beta1.VMServiceScrape
	AlertRules *v1beta1.VMRule
}

func NewMonKubeController() *MonKubeController {
	return &MonKubeController{
		SVC:        KubeControllerSVC,
		Scrape:     KubeControllerScrape,
		AlertRules: KubeControllerAlertRules,
	}
}

var KubeControllerScrape = &v1beta1.VMServiceScrape{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-kube-controller-manager",
		Namespace: "monitoring",
	},
	Spec: v1beta1.VMServiceScrapeSpec{
		Endpoints: []v1beta1.Endpoint{
			{
				BearerTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
				Port:            "http-metrics",
				Scheme:          "https",
				TLSConfig: &v1beta1.TLSConfig{
					CAFile:     "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
					ServerName: "kubernetes",
				},
			},
		},
		JobLabel:          "jobLabel",
		NamespaceSelector: v1beta1.NamespaceSelector{MatchNames: []string{"kube-system"}},
		Selector: metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app":                        "vmk8s-victoria-metrics-k8s-stack-kube-controller-manager",
				"app.kubernetes.io/instance": "vmk8s",
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "operator.victoriametrics.com/v1beta1",
		Kind:       "VMServiceScrape",
	},
}

var KubeControllerSVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "vmk8s-victoria-metrics-k8s-stack-kube-controller-manager",
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
			"jobLabel":                     "kube-controller-manager",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-kube-controller-manager",
		Namespace: "kube-system",
	},
	Spec: corev1.ServiceSpec{
		ClusterIP: "None",
		Ports: []corev1.ServicePort{
			{
				Name:       "http-metrics",
				Port:       int32(10257),
				Protocol:   corev1.Protocol("TCP"),
				TargetPort: intstr.IntOrString{IntVal: int32(10257)},
			},
		},
		Selector: map[string]string{"component": "kube-controller-manager"},
		Type:     corev1.ServiceType("ClusterIP"),
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	},
}

var KubeControllerAlertRules = &v1beta1.VMRule{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "victoria-metrics-k8s-stack",
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-kubernetes-system-controller-m",
		Namespace: "monitoring",
	},
	Spec: v1beta1.VMRuleSpec{
		Groups: []v1beta1.RuleGroup{
			{
				Name: "kubernetes-system-controller-manager",
				Rules: []v1beta1.Rule{
					{
						Alert: "KubeControllerManagerDown",
						Annotations: map[string]string{
							"description": "KubeControllerManager has disappeared from Prometheus target discovery.",
							"runbook_url": "https://runbooks.prometheus-operator.dev/runbooks/kubernetes/kubecontrollermanagerdown",
							"summary":     "Target disappeared from Prometheus target discovery.",
						},
						Expr:   "absent(up{job=\"kube-controller-manager\"} == 1)",
						For:    "15m",
						Labels: map[string]string{"severity": "critical"},
					},
				},
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "operator.victoriametrics.com/v1beta1",
		Kind:       "VMRule",
	},
}
