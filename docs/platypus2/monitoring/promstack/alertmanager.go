// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package promstack

import (
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var AlertmanagerAlertmanagerSecrets = &corev1.Secret{
	Data: map[string][]byte{
		"alertmanager.yaml": []byte(`global:
  resolve_timeout: 5m
inhibit_rules:
- equal:
  - namespace
  - alertname
  source_matchers:
  - severity = critical
  target_matchers:
  - severity =~ warning|info
- equal:
  - namespace
  - alertname
  source_matchers:
  - severity = warning
  target_matchers:
  - severity = info
- equal:
  - namespace
  source_matchers:
  - alertname = InfoInhibitor
  target_matchers:
  - severity = info
receivers:
- name: "null"
route:
  group_by:
  - namespace
  group_interval: 5m
  group_wait: 30s
  receiver: "null"
  repeat_interval: 12h
  routes:
  - matchers:
    - alertname =~ "InfoInhibitor|Watchdog"
    receiver: "null"
templates:
- /etc/alertmanager/config/*.tmpl
`),
	},
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-alertmanager",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "45.27.2",
			"chart":                        "kube-prometheus-stack-45.27.2",
			"heritage":                     "Helm",
			"release":                      "kube-prometheus-stack",
		},
		Name:      "alertmanager-kube-prometheus-stack-alertmanager",
		Namespace: namespace,
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Secret",
	},
} // TODO: SECRETS SHOULD BE STORED ELSEWHERE THAN IN THE CODE!!!!

var AlertmanagerAlertmanager = &v1.Alertmanager{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-alertmanager",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "45.27.2",
			"chart":                        "kube-prometheus-stack-45.27.2",
			"heritage":                     "Helm",
			"release":                      "kube-prometheus-stack",
		},
		Name:      "kube-prometheus-stack-alertmanager",
		Namespace: namespace,
	},
	Spec: v1.AlertmanagerSpec{
		ExternalURL: "http://kube-prometheus-stack-alertmanager.monitoring:9093",
		Image:       P("quay.io/prometheus/alertmanager:v0.25.0"),
		LogFormat:   "logfmt",
		LogLevel:    "info",
		PortName:    "http-web",
		Replicas:    P(int32(1)),
		Retention:   v1.GoDuration("120h"),
		RoutePrefix: "/",
		SecurityContext: &corev1.PodSecurityContext{
			FSGroup:      P(int64(2000)),
			RunAsGroup:   P(int64(2000)),
			RunAsNonRoot: P(true),
			RunAsUser:    P(int64(1000)),
		},
		ServiceAccountName: "kube-prometheus-stack-alertmanager",
		Version:            "v0.25.0",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "monitoring.coreos.com/v1",
		Kind:       "Alertmanager",
	},
}

var AlertmanagerSA = &corev1.ServiceAccount{
	AutomountServiceAccountToken: P(true),
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-alertmanager",
			"app.kubernetes.io/component":  "alertmanager",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "kube-prometheus-stack-alertmanager",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "45.27.2",
			"chart":                        "kube-prometheus-stack-45.27.2",
			"heritage":                     "Helm",
			"release":                      "kube-prometheus-stack",
		},
		Name:      "kube-prometheus-stack-alertmanager",
		Namespace: namespace,
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ServiceAccount",
	},
}

var AlertmanagerSVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-alertmanager",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "45.27.2",
			"chart":                        "kube-prometheus-stack-45.27.2",
			"heritage":                     "Helm",
			"release":                      "kube-prometheus-stack",
			"self-monitor":                 "true",
		},
		Name:      "kube-prometheus-stack-alertmanager",
		Namespace: namespace,
	},
	Spec: corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			{
				Name:       "http-web",
				Port:       int32(9093),
				Protocol:   corev1.Protocol("TCP"),
				TargetPort: intstr.IntOrString{IntVal: int32(9093)},
			},
		},
		Selector: map[string]string{
			"alertmanager":           "kube-prometheus-stack-alertmanager",
			"app.kubernetes.io/name": "alertmanager",
		},
		Type: corev1.ServiceType("ClusterIP"),
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	},
}
