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

type MonKubeScheduler struct {
	kube.App

	KubeSchedulerRules      *v1beta1.VMRule
	KubeSchedulerAlertRules *v1beta1.VMRule
	KubeSchedulerSVC        *corev1.Service
	KubeSchedulerScrape     *v1beta1.VMServiceScrape
}

func NewMonKubeScheduler() *MonKubeScheduler {
	return &MonKubeScheduler{
		KubeSchedulerRules:      KubeSchedulerRecordingRules,
		KubeSchedulerSVC:        KubeSchedulerSVC,
		KubeSchedulerScrape:     KubeSchedulerScrape,
		KubeSchedulerAlertRules: KubeSchedulerAlertRules,
	}
}

var KubeSchedulerScrape = &v1beta1.VMServiceScrape{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-kube-scheduler",
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
				"app":                        "vmk8s-victoria-metrics-k8s-stack-kube-scheduler",
				"app.kubernetes.io/instance": "vmk8s",
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "operator.victoriametrics.com/v1beta1",
		Kind:       "VMServiceScrape",
	},
}

var KubeSchedulerAlertRules = &v1beta1.VMRule{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "victoria-metrics-k8s-stack",
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-kubernetes-system-scheduler",
		Namespace: "monitoring",
	},
	Spec: v1beta1.VMRuleSpec{
		Groups: []v1beta1.RuleGroup{
			{
				Name: "kubernetes-system-scheduler",
				Rules: []v1beta1.Rule{
					{
						Alert: "KubeSchedulerDown",
						Annotations: map[string]string{
							"description": "KubeScheduler has disappeared from Prometheus target discovery.",
							"runbook_url": "https://runbooks.prometheus-operator.dev/runbooks/kubernetes/kubeschedulerdown",
							"summary":     "Target disappeared from Prometheus target discovery.",
						},
						Expr:   "absent(up{job=\"kube-scheduler\"} == 1)",
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

var KubeSchedulerRecordingRules = &v1beta1.VMRule{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "victoria-metrics-k8s-stack",
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-kube-scheduler.rules",
		Namespace: "monitoring",
	},
	Spec: v1beta1.VMRuleSpec{
		Groups: []v1beta1.RuleGroup{
			{
				Name: "kube-scheduler.rules",
				Rules: []v1beta1.Rule{
					{
						Expr:   "histogram_quantile(0.99, sum(rate(scheduler_e2e_scheduling_duration_seconds_bucket{job=\"kube-scheduler\"}[5m])) without(instance, pod))",
						Labels: map[string]string{"quantile": "0.99"},
						Record: "cluster_quantile:scheduler_e2e_scheduling_duration_seconds:histogram_quantile",
					}, {
						Expr:   "histogram_quantile(0.99, sum(rate(scheduler_scheduling_algorithm_duration_seconds_bucket{job=\"kube-scheduler\"}[5m])) without(instance, pod))",
						Labels: map[string]string{"quantile": "0.99"},
						Record: "cluster_quantile:scheduler_scheduling_algorithm_duration_seconds:histogram_quantile",
					}, {
						Expr:   "histogram_quantile(0.99, sum(rate(scheduler_binding_duration_seconds_bucket{job=\"kube-scheduler\"}[5m])) without(instance, pod))",
						Labels: map[string]string{"quantile": "0.99"},
						Record: "cluster_quantile:scheduler_binding_duration_seconds:histogram_quantile",
					}, {
						Expr:   "histogram_quantile(0.9, sum(rate(scheduler_e2e_scheduling_duration_seconds_bucket{job=\"kube-scheduler\"}[5m])) without(instance, pod))",
						Labels: map[string]string{"quantile": "0.9"},
						Record: "cluster_quantile:scheduler_e2e_scheduling_duration_seconds:histogram_quantile",
					}, {
						Expr:   "histogram_quantile(0.9, sum(rate(scheduler_scheduling_algorithm_duration_seconds_bucket{job=\"kube-scheduler\"}[5m])) without(instance, pod))",
						Labels: map[string]string{"quantile": "0.9"},
						Record: "cluster_quantile:scheduler_scheduling_algorithm_duration_seconds:histogram_quantile",
					}, {
						Expr:   "histogram_quantile(0.9, sum(rate(scheduler_binding_duration_seconds_bucket{job=\"kube-scheduler\"}[5m])) without(instance, pod))",
						Labels: map[string]string{"quantile": "0.9"},
						Record: "cluster_quantile:scheduler_binding_duration_seconds:histogram_quantile",
					}, {
						Expr:   "histogram_quantile(0.5, sum(rate(scheduler_e2e_scheduling_duration_seconds_bucket{job=\"kube-scheduler\"}[5m])) without(instance, pod))",
						Labels: map[string]string{"quantile": "0.5"},
						Record: "cluster_quantile:scheduler_e2e_scheduling_duration_seconds:histogram_quantile",
					}, {
						Expr:   "histogram_quantile(0.5, sum(rate(scheduler_scheduling_algorithm_duration_seconds_bucket{job=\"kube-scheduler\"}[5m])) without(instance, pod))",
						Labels: map[string]string{"quantile": "0.5"},
						Record: "cluster_quantile:scheduler_scheduling_algorithm_duration_seconds:histogram_quantile",
					}, {
						Expr:   "histogram_quantile(0.5, sum(rate(scheduler_binding_duration_seconds_bucket{job=\"kube-scheduler\"}[5m])) without(instance, pod))",
						Labels: map[string]string{"quantile": "0.5"},
						Record: "cluster_quantile:scheduler_binding_duration_seconds:histogram_quantile",
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

var KubeSchedulerSVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "vmk8s-victoria-metrics-k8s-stack-kube-scheduler",
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
			"jobLabel":                     "kube-scheduler",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-kube-scheduler",
		Namespace: "kube-system",
	},
	Spec: corev1.ServiceSpec{
		ClusterIP: "None",
		Ports: []corev1.ServicePort{
			{
				Name:       "http-metrics",
				Port:       int32(10251),
				Protocol:   corev1.ProtocolTCP,
				TargetPort: intstr.IntOrString{IntVal: int32(10251)},
			},
		},
		Selector: map[string]string{"component": "kube-scheduler"},
		Type:     corev1.ServiceType("ClusterIP"),
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	},
}
