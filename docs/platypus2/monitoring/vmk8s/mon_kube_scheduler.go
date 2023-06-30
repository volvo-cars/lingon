// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmk8s

import (
	"github.com/VictoriaMetrics/operator/api/victoriametrics/v1beta1"
	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	"github.com/volvo-cars/lingoneks/meta"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	KSCPort     = 10251
	KSCPortName = "http-metrics"
	PathSA      = "/var/run/secrets/kubernetes.io/serviceaccount"
)

var KSC = &meta.Metadata{
	Name:      "kube-scheduler",
	Namespace: namespace,
	Instance:  "kube-scheduler-" + namespace,
	Component: "monitoring",
	PartOf:    appName,
	Version:   version,
	ManagedBy: "lingon",
}

type MonKubeScheduler struct {
	kube.App

	KubeSchedulerScrape     *v1beta1.VMServiceScrape
	KubeSchedulerRules      *v1beta1.VMRule
	KubeSchedulerAlertRules *v1beta1.VMRule
	KubeSchedulerSVC        *corev1.Service
}

func NewMonKubeScheduler() *MonKubeScheduler {
	return &MonKubeScheduler{
		KubeSchedulerScrape:     KubeSchedulerScrape,
		KubeSchedulerRules:      KubeSchedulerRecordingRules,
		KubeSchedulerAlertRules: KubeSchedulerAlertRules,
		KubeSchedulerSVC:        KubeSchedulerSVC,
	}
}

var KubeSchedulerScrape = &v1beta1.VMServiceScrape{
	TypeMeta:   TypeVMServiceScrapeV1Beta1,
	ObjectMeta: KSC.ObjectMeta(),
	Spec: v1beta1.VMServiceScrapeSpec{
		Endpoints: []v1beta1.Endpoint{
			{
				BearerTokenFile: PathSA + "/token",
				Port:            KSCPortName,
				Scheme:          "https",
				TLSConfig:       &v1beta1.TLSConfig{CAFile: PathSA + "/ca.crt"},
			},
		},
		JobLabel: "component",
		NamespaceSelector: v1beta1.NamespaceSelector{
			MatchNames: []string{ku.NSKubeSystem}, // kube-system
		},
		Selector: metav1.LabelSelector{
			MatchLabels: map[string]string{"component": "kube-scheduler"},
		},
	},
}

var KubeSchedulerAlertRules = &v1beta1.VMRule{
	TypeMeta:   TypeVMRuleV1Beta1,
	ObjectMeta: KSC.ObjectMeta(),
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
						Expr:   `absent(up{job="` + KSC.Name + `"} == 1)`,
						For:    "15m",
						Labels: map[string]string{"severity": "critical"},
					},
				},
			},
		},
	},
}

var KubeSchedulerRecordingRules = &v1beta1.VMRule{
	TypeMeta:   TypeVMRuleV1Beta1,
	ObjectMeta: KSC.ObjectMetaNameSuffix("rules"),
	Spec: v1beta1.VMRuleSpec{
		Groups: []v1beta1.RuleGroup{
			{
				Name: "kube-scheduler.rules",
				Rules: []v1beta1.Rule{
					{
						Expr:   `histogram_quantile(0.99, sum(rate(scheduler_e2e_scheduling_duration_seconds_bucket{job="` + KSC.Name + `"}[5m])) without(instance, pod))`,
						Labels: map[string]string{"quantile": "0.99"},
						Record: "cluster_quantile:scheduler_e2e_scheduling_duration_seconds:histogram_quantile",
					}, {
						Expr:   `histogram_quantile(0.99, sum(rate(scheduler_scheduling_algorithm_duration_seconds_bucket{job="` + KSC.Name + `"}[5m])) without(instance, pod))`,
						Labels: map[string]string{"quantile": "0.99"},
						Record: "cluster_quantile:scheduler_scheduling_algorithm_duration_seconds:histogram_quantile",
					}, {
						Expr:   `histogram_quantile(0.99, sum(rate(scheduler_binding_duration_seconds_bucket{job="` + KSC.Name + `"}[5m])) without(instance, pod))`,
						Labels: map[string]string{"quantile": "0.99"},
						Record: "cluster_quantile:scheduler_binding_duration_seconds:histogram_quantile",
					}, {
						Expr:   `histogram_quantile(0.9, sum(rate(scheduler_e2e_scheduling_duration_seconds_bucket{job="` + KSC.Name + `"}[5m])) without(instance, pod))`,
						Labels: map[string]string{"quantile": "0.9"},
						Record: "cluster_quantile:scheduler_e2e_scheduling_duration_seconds:histogram_quantile",
					}, {
						Expr:   `histogram_quantile(0.9, sum(rate(scheduler_scheduling_algorithm_duration_seconds_bucket{job="` + KSC.Name + `"}[5m])) without(instance, pod))`,
						Labels: map[string]string{"quantile": "0.9"},
						Record: "cluster_quantile:scheduler_scheduling_algorithm_duration_seconds:histogram_quantile",
					}, {
						Expr:   `histogram_quantile(0.9, sum(rate(scheduler_binding_duration_seconds_bucket{job="` + KSC.Name + `"}[5m])) without(instance, pod))`,
						Labels: map[string]string{"quantile": "0.9"},
						Record: "cluster_quantile:scheduler_binding_duration_seconds:histogram_quantile",
					}, {
						Expr:   `histogram_quantile(0.5, sum(rate(scheduler_e2e_scheduling_duration_seconds_bucket{job="` + KSC.Name + `"}[5m])) without(instance, pod))`,
						Labels: map[string]string{"quantile": "0.5"},
						Record: "cluster_quantile:scheduler_e2e_scheduling_duration_seconds:histogram_quantile",
					}, {
						Expr:   `histogram_quantile(0.5, sum(rate(scheduler_scheduling_algorithm_duration_seconds_bucket{job="` + KSC.Name + `"}[5m])) without(instance, pod))`,
						Labels: map[string]string{"quantile": "0.5"},
						Record: "cluster_quantile:scheduler_scheduling_algorithm_duration_seconds:histogram_quantile",
					}, {
						Expr:   `histogram_quantile(0.5, sum(rate(scheduler_binding_duration_seconds_bucket{job="` + KSC.Name + `"}[5m])) without(instance, pod))`,
						Labels: map[string]string{"quantile": "0.5"},
						Record: "cluster_quantile:scheduler_binding_duration_seconds:histogram_quantile",
					},
				},
			},
		},
	},
}

var KubeSchedulerSVC = &corev1.Service{
	TypeMeta: ku.TypeServiceV1,
	ObjectMeta: metav1.ObjectMeta{
		Labels:    KSC.Labels(),
		Name:      KSC.Name,
		Namespace: ku.NSKubeSystem, // "kube-system",
	},
	Spec: corev1.ServiceSpec{
		ClusterIP: corev1.ClusterIPNone,
		Ports: []corev1.ServicePort{
			{
				Name:       KSCPortName,
				Port:       int32(KSCPort),
				Protocol:   corev1.ProtocolTCP,
				TargetPort: intstr.FromInt(KSCPort),
			},
		},
		Selector: map[string]string{"component": "kube-scheduler"},
		Type:     corev1.ServiceTypeClusterIP,
	},
}
