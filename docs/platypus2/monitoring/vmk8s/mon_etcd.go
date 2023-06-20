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

type MonETCD struct {
	kube.App

	ETCDRules  *v1beta1.VMRule
	ETCDSvc    *corev1.Service
	ETCDScrape *v1beta1.VMServiceScrape
}

func NewMonETCD() *MonETCD {
	return &MonETCD{
		ETCDRules:  ETCDRules,
		ETCDSvc:    ETCDSvc,
		ETCDScrape: ETCDScrape,
	}
}

var ETCDRules = &v1beta1.VMRule{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "victoria-metrics-k8s-stack",
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-etcd",
		Namespace: "monitoring",
	},
	Spec: v1beta1.VMRuleSpec{
		Groups: []v1beta1.RuleGroup{
			{
				Name: "etcd",
				Rules: []v1beta1.Rule{
					{
						Alert:       "etcdInsufficientMembers",
						Annotations: map[string]string{"message": "etcd cluster \"{{ $labels.job }}\": insufficient members ({{ $value }})."},
						Expr:        "sum(up{job=~\".*etcd.*\"} == bool 1) by (job) < ((count(up{job=~\".*etcd.*\"}) by (job) + 1) / 2)",
						For:         "3m",
						Labels:      map[string]string{"severity": "critical"},
					}, {
						Alert:       "etcdNoLeader",
						Annotations: map[string]string{"message": "etcd cluster \"{{ $labels.job }}\": member {{ $labels.instance }} has no leader."},
						Expr:        "etcd_server_has_leader{job=~\".*etcd.*\"} == 0",
						For:         "1m",
						Labels:      map[string]string{"severity": "critical"},
					}, {
						Alert:       "etcdHighNumberOfLeaderChanges",
						Annotations: map[string]string{"message": "etcd cluster \"{{ $labels.job }}\": instance {{ $labels.instance }} has seen {{ $value }} leader changes within the last hour."},
						Expr:        "rate(etcd_server_leader_changes_seen_total{job=~\".*etcd.*\"}[15m]) > 3",
						For:         "15m",
						Labels:      map[string]string{"severity": "warning"},
					}, {
						Alert:       "etcdHighNumberOfFailedGRPCRequests",
						Annotations: map[string]string{"message": "etcd cluster \"{{ $labels.job }}\": {{ $value }}% of requests for {{ $labels.grpc_method }} failed on etcd instance {{ $labels.instance }}."},
						Expr: `
								100 * sum(rate(grpc_server_handled_total{job=~".*etcd.*", grpc_code!="OK"}[5m])) BY (job, instance, grpc_service, grpc_method)
								  /
								sum(rate(grpc_server_handled_total{job=~".*etcd.*"}[5m])) BY (job, instance, grpc_service, grpc_method)
								  > 1
								`,
						For:    "10m",
						Labels: map[string]string{"severity": "warning"},
					}, {
						Alert:       "etcdHighNumberOfFailedGRPCRequests",
						Annotations: map[string]string{"message": "etcd cluster \"{{ $labels.job }}\": {{ $value }}% of requests for {{ $labels.grpc_method }} failed on etcd instance {{ $labels.instance }}."},
						Expr: `
								100 * sum(rate(grpc_server_handled_total{job=~".*etcd.*", grpc_code!="OK"}[5m])) BY (job, instance, grpc_service, grpc_method)
								  /
								sum(rate(grpc_server_handled_total{job=~".*etcd.*"}[5m])) BY (job, instance, grpc_service, grpc_method)
								  > 5
								`,
						For:    "5m",
						Labels: map[string]string{"severity": "critical"},
					}, {
						Alert:       "etcdGRPCRequestsSlow",
						Annotations: map[string]string{"message": "etcd cluster \"{{ $labels.job }}\": gRPC requests to {{ $labels.grpc_method }} are taking {{ $value }}s on etcd instance {{ $labels.instance }}."},
						Expr: `
								histogram_quantile(0.99, sum(rate(grpc_server_handling_seconds_bucket{job=~".*etcd.*", grpc_type="unary"}[5m])) by (job, instance, grpc_service, grpc_method, le))
								> 0.15
								`,
						For:    "10m",
						Labels: map[string]string{"severity": "critical"},
					}, {
						Alert:       "etcdMemberCommunicationSlow",
						Annotations: map[string]string{"message": "etcd cluster \"{{ $labels.job }}\": member communication with {{ $labels.To }} is taking {{ $value }}s on etcd instance {{ $labels.instance }}."},
						Expr: `
								histogram_quantile(0.99, rate(etcd_network_peer_round_trip_time_seconds_bucket{job=~".*etcd.*"}[5m]))
								> 0.15
								`,
						For:    "10m",
						Labels: map[string]string{"severity": "warning"},
					}, {
						Alert:       "etcdHighNumberOfFailedProposals",
						Annotations: map[string]string{"message": "etcd cluster \"{{ $labels.job }}\": {{ $value }} proposal failures within the last hour on etcd instance {{ $labels.instance }}."},
						Expr:        "rate(etcd_server_proposals_failed_total{job=~\".*etcd.*\"}[15m]) > 5",
						For:         "15m",
						Labels:      map[string]string{"severity": "warning"},
					}, {
						Alert:       "etcdHighFsyncDurations",
						Annotations: map[string]string{"message": "etcd cluster \"{{ $labels.job }}\": 99th percentile fync durations are {{ $value }}s on etcd instance {{ $labels.instance }}."},
						Expr: `
								histogram_quantile(0.99, rate(etcd_disk_wal_fsync_duration_seconds_bucket{job=~".*etcd.*"}[5m]))
								> 0.5
								`,
						For:    "10m",
						Labels: map[string]string{"severity": "warning"},
					}, {
						Alert:       "etcdHighCommitDurations",
						Annotations: map[string]string{"message": "etcd cluster \"{{ $labels.job }}\": 99th percentile commit durations {{ $value }}s on etcd instance {{ $labels.instance }}."},
						Expr: `
								histogram_quantile(0.99, rate(etcd_disk_backend_commit_duration_seconds_bucket{job=~".*etcd.*"}[5m]))
								> 0.25
								`,
						For:    "10m",
						Labels: map[string]string{"severity": "warning"},
					}, {
						Alert:       "etcdHighNumberOfFailedHTTPRequests",
						Annotations: map[string]string{"message": "{{ $value }}% of requests for {{ $labels.method }} failed on etcd instance {{ $labels.instance }}"},
						Expr: `
								sum(rate(etcd_http_failed_total{job=~".*etcd.*", code!="404"}[5m])) BY (method) / sum(rate(etcd_http_received_total{job=~".*etcd.*"}[5m]))
								BY (method) > 0.01
								`,
						For:    "10m",
						Labels: map[string]string{"severity": "warning"},
					}, {
						Alert:       "etcdHighNumberOfFailedHTTPRequests",
						Annotations: map[string]string{"message": "{{ $value }}% of requests for {{ $labels.method }} failed on etcd instance {{ $labels.instance }}."},
						Expr: `
								sum(rate(etcd_http_failed_total{job=~".*etcd.*", code!="404"}[5m])) BY (method) / sum(rate(etcd_http_received_total{job=~".*etcd.*"}[5m]))
								BY (method) > 0.05
								`,
						For:    "10m",
						Labels: map[string]string{"severity": "critical"},
					}, {
						Alert:       "etcdHTTPRequestsSlow",
						Annotations: map[string]string{"message": "etcd instance {{ $labels.instance }} HTTP requests to {{ $labels.method }} are slow."},
						Expr: `
								histogram_quantile(0.99, rate(etcd_http_successful_duration_seconds_bucket[5m]))
								> 0.15
								`,
						For:    "10m",
						Labels: map[string]string{"severity": "warning"},
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

var ETCDSvc = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "vmk8s-victoria-metrics-k8s-stack-kube-etcd",
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
			"jobLabel":                     "kube-etcd",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-kube-etcd",
		Namespace: "kube-system",
	},
	Spec: corev1.ServiceSpec{
		ClusterIP: "None",
		Ports: []corev1.ServicePort{
			{
				Name:       "http-metrics",
				Port:       int32(2379),
				Protocol:   corev1.Protocol("TCP"),
				TargetPort: intstr.IntOrString{IntVal: int32(2379)},
			},
		},
		Selector: map[string]string{"component": "etcd"},
		Type:     corev1.ServiceType("ClusterIP"),
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	},
}

var ETCDScrape = &v1beta1.VMServiceScrape{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-kube-etcd",
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
				"app":                        "vmk8s-victoria-metrics-k8s-stack-kube-etcd",
				"app.kubernetes.io/instance": "vmk8s",
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "operator.victoriametrics.com/v1beta1",
		Kind:       "VMServiceScrape",
	},
}
