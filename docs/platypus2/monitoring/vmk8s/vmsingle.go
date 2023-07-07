// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package vmk8s

import (
	_ "embed"
	"fmt"

	"github.com/VictoriaMetrics/operator/api/victoriametrics/v1beta1"
	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	"github.com/volvo-cars/lingoneks/meta"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Single = &meta.Metadata{
	Name:      "victoria-metrics",
	Namespace: namespace,
	Instance:  "victoria-metrics-" + namespace,
	Component: "tsmdb",
	PartOf:    appName,
	Version:   version,
	ManagedBy: "lingon",
}

type VicMet struct {
	kube.App

	DB              *v1beta1.VMSingle
	Agent           *v1beta1.VMAgent
	SA              *corev1.ServiceAccount
	SingleAlerts    *v1beta1.VMRule
	HealthAlerts    *v1beta1.VMRule
	AgentAlertRules *v1beta1.VMRule
}

type VicMetOption func(server *VicMet) *VicMet

func NewVicMet(opts ...VicMetOption) *VicMet {
	vm := &VicMet{
		DB:              VMDB,
		Agent:           VMAgent,
		SA:              VictoriaMetricsSA,
		AgentAlertRules: VMAgentAlertRules,
		HealthAlerts:    VMHealthAlertRules,
		SingleAlerts:    VMSingleAlertRules,
	}

	for _, o := range opts {
		vm = o(vm)
	}
	return vm
}

// VMDB is a single instance of Victoria Metrics DB.
// Note that a "vmsingle-" prefix is added to the name.
// See https://github.com/VictoriaMetrics/operator/blob/4c97f70c9a775d2bfff401862acabd5452ef0cf8/api/v1beta1/vmsingle_types.go#L326
var VMDB = &v1beta1.VMSingle{
	TypeMeta:   TypeVMSingleV1Beta1,
	ObjectMeta: Single.ObjectMeta(),
	Spec: v1beta1.VMSingleSpec{
		Image:        v1beta1.Image{Tag: "v" + Single.Version},
		ReplicaCount: P(int32(1)),
		ExtraArgs: map[string]string{
			"vmalert.proxyURL": fmt.Sprintf(
				"%s.%s.svc:8080", // TODO: extract port
				AlertManager.PrefixedName(),
				AlertManager.Namespace,
			),
		},
		RetentionPeriod: "14",
		Resources:       ku.Resources("4", "8Gi", "4", "8Gi"),
		Port:            d(VMSinglePort),
		Storage: &corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceName("storage"): resource.MustParse("20Gi"),
				},
			},
		},
	},
}

var VMSingleAlertRules = &v1beta1.VMRule{
	ObjectMeta: metav1.ObjectMeta{
		Labels:    Single.Labels(),
		Name:      Single.Name + "-alerting-rules",
		Namespace: namespace,
	},
	Spec: v1beta1.VMRuleSpec{
		Groups: []v1beta1.RuleGroup{
			{
				Concurrency: 2,
				Interval:    "30s",
				// Name:        "vmsingle",
				Name: Single.Name,
				Rules: []v1beta1.Rule{
					{
						Alert: "DiskRunsOutOfSpaceIn3Days",
						Annotations: map[string]string{
							"dashboard": "grafana.domain.com/d/wNf0q_kZk?viewPanel=73&var-instance={{ $labels.instance }}",
							"description": `
Taking into account current ingestion rate, free disk space will be enough only for {{ $value | humanizeDuration }} on instance {{ $labels.instance }}.
 Consider to limit the ingestion rate, decrease retention or scale the disk space if possible.
`,
							"summary": "Instance {{ $labels.instance }} will run out of disk space soon",
						},
						Expr: `
vm_free_disk_space_bytes / ignoring(path)
(
   (
    rate(vm_rows_added_to_storage_total[1d]) -
    ignoring(type) rate(vm_deduplicated_samples_total{type="merge"}[1d])
   )
  * scalar(
    sum(vm_data_size_bytes{type!~"indexdb.*"}) /
    sum(vm_rows{type!~"indexdb.*"})
   )
) < 3 * 24 * 3600 > 0
`,
						For:    "30m",
						Labels: map[string]string{"severity": "critical"},
					}, {
						Alert: "DiskRunsOutOfSpace",
						Annotations: map[string]string{
							"dashboard": "grafana.domain.com/d/wNf0q_kZk?viewPanel=53&var-instance={{ $labels.instance }}",
							"description": `
Disk utilisation on instance {{ $labels.instance }} is more than 80%.
 Having less than 20% of free disk space could cripple merges processes and overall performance. Consider to limit the ingestion rate, decrease retention or scale the disk space if possible.
`,
							"summary": "Instance {{ $labels.instance }} will run out of disk space soon",
						},
						Expr: `
sum(vm_data_size_bytes) by(instance) /
(
 sum(vm_free_disk_space_bytes) by(instance) +
 sum(vm_data_size_bytes) by(instance)
) > 0.8
`,
						For:    "30m",
						Labels: map[string]string{"severity": "critical"},
					}, {
						Alert: "RequestErrorsToAPI",
						Annotations: map[string]string{
							"dashboard":   "grafana.domain.com/d/wNf0q_kZk?viewPanel=35&var-instance={{ $labels.instance }}",
							"description": "Requests to path {{ $labels.path }} are receiving errors. Please verify if clients are sending correct requests.",
							"summary":     "Too many errors served for path {{ $labels.path }} (instance {{ $labels.instance }})",
						},
						Expr:   "increase(vm_http_request_errors_total[5m]) > 0",
						For:    "15m",
						Labels: map[string]string{"severity": "warning"},
					}, {
						Alert: "ConcurrentFlushesHitTheLimit",
						Annotations: map[string]string{
							"dashboard": "grafana.domain.com/d/wNf0q_kZk?viewPanel=59&var-instance={{ $labels.instance }}",
							"description": `
The limit of concurrent flushes on instance {{ $labels.instance }} is equal to number of CPUs.
 When VictoriaMetrics constantly hits the limit it means that storage is overloaded and requires more CPU.
`,
							"summary": "VictoriaMetrics on instance {{ $labels.instance }} is constantly hitting concurrent flushes limit",
						},
						Expr: "avg_over_time(vm_concurrent_insert_current[1m]) >= vm_concurrent_insert_capacity",
						For:  "15m",
						Labels: map[string]string{
							"severity": "warning",
							"show_at":  "dashboard",
						},
					}, {
						Alert: "RowsRejectedOnIngestion",
						Annotations: map[string]string{
							"dashboard":   "grafana.domain.com/d/wNf0q_kZk?viewPanel=58&var-instance={{ $labels.instance }}",
							"description": `VM is rejecting to ingest rows on "{{ $labels.instance }}" due to the following reason: "{{ $labels.reason }}"`,
							"summary":     `Some rows are rejected on "{{ $labels.instance }}" on ingestion attempt`,
						},
						Expr:   "sum(rate(vm_rows_ignored_total[5m])) by (instance, reason) > 0",
						For:    "15m",
						Labels: map[string]string{"severity": "warning"},
					}, {
						Alert: "TooHighChurnRate",
						Annotations: map[string]string{
							"dashboard": "grafana.domain.com/d/wNf0q_kZk?viewPanel=66&var-instance={{ $labels.instance }}",
							"description": `
VM constantly creates new time series on "{{ $labels.instance }}".
 This effect is known as Churn Rate.
 High Churn Rate tightly connected with database performance and may result in unexpected OOM's or slow queries.
`,
							"summary": `Churn rate is more than 10% on "{{ $labels.instance }}" for the last 15m`,
						},
						Expr: `
(
   sum(rate(vm_new_timeseries_created_total[5m])) by(instance)
   /
   sum(rate(vm_rows_inserted_total[5m])) by (instance)
 ) > 0.1
`,
						For:    "15m",
						Labels: map[string]string{"severity": "warning"},
					}, {
						Alert: "TooHighChurnRate24h",
						Annotations: map[string]string{
							"dashboard": "grafana.domain.com/d/wNf0q_kZk?viewPanel=66&var-instance={{ $labels.instance }}",
							"description": `
The number of created new time series over last 24h is 3x times higher than current number of active series on "{{ $labels.instance }}".
 This effect is known as Churn Rate.
 High Churn Rate tightly connected with database performance and may result in unexpected OOM's or slow queries.
`,
							"summary": `Too high number of new series on "{{ $labels.instance }}" created over last 24h`,
						},
						Expr: `
sum(increase(vm_new_timeseries_created_total[24h])) by(instance)
>
(sum(vm_cache_entries{type="storage/hour_metric_ids"}) by(instance) * 3)
`,
						For:    "15m",
						Labels: map[string]string{"severity": "warning"},
					}, {
						Alert: "TooHighSlowInsertsRate",
						Annotations: map[string]string{
							"dashboard":   "grafana.domain.com/d/wNf0q_kZk?viewPanel=68&var-instance={{ $labels.instance }}",
							"description": `High rate of slow inserts on "{{ $labels.instance }}" may be a sign of resource exhaustion for the current load. It is likely more RAM is needed for optimal handling of the current number of active time series.`,
							"summary":     `Percentage of slow inserts is more than 5% on "{{ $labels.instance }}" for the last 15m`,
						},
						Expr: `
(
   sum(rate(vm_slow_row_inserts_total[5m])) by(instance)
   /
   sum(rate(vm_rows_inserted_total[5m])) by (instance)
 ) > 0.05
`,
						For:    "15m",
						Labels: map[string]string{"severity": "warning"},
					}, {
						Alert: "LabelsLimitExceededOnIngestion",
						Annotations: map[string]string{
							"dashboard":   "grafana.domain.com/d/wNf0q_kZk?viewPanel=74&var-instance={{ $labels.instance }}",
							"description": "VictoriaMetrics limits the number of labels per each metric with `-maxLabelsPerTimeseries` command-line flag.\n This prevents from ingesting metrics with too many labels. Please verify that `-maxLabelsPerTimeseries` is configured correctly or that clients which send these metrics aren't misbehaving.",
							"summary":     "Metrics ingested in ({{ $labels.instance }}) are exceeding labels limit",
						},
						Expr:   "sum(increase(vm_metrics_with_dropped_labels_total[5m])) by (instance) > 0",
						For:    "15m",
						Labels: map[string]string{"severity": "warning"},
					},
				},
			},
		},
	},
	TypeMeta: TypeVMRuleV1Beta1,
}

var VMHealthAlertRules = &v1beta1.VMRule{
	ObjectMeta: metav1.ObjectMeta{
		Labels: Single.Labels(),
		// Name:      "vmk8s-victoria-metrics-k8s-stack-vm-health",
		Name:      Single.Name + "-health",
		Namespace: namespace,
	},
	Spec: v1beta1.VMRuleSpec{
		Groups: []v1beta1.RuleGroup{
			{
				// Name: "vm-health",
				Name: Single.Name + "-health",
				Rules: []v1beta1.Rule{
					{
						Alert: "TooManyRestarts",
						Annotations: map[string]string{
							"description": "Job {{ $labels.job }} (instance {{ $labels.instance }}) has restarted more than twice in the last 15 minutes. It might be crashlooping.",
							"summary":     "{{ $labels.job }} too many restarts (instance {{ $labels.instance }})",
						},
						Expr:   `changes(process_start_time_seconds{job=~"victoriametrics|vmselect|vminsert|vmstorage|vmagent|vmalert"}[15m]) > 2`,
						Labels: map[string]string{"severity": "critical"},
					}, {
						Alert: "ServiceDown",
						Annotations: map[string]string{
							"description": "{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 2 minutes.",
							"summary":     "Service {{ $labels.job }} is down on {{ $labels.instance }}",
						},
						Expr:   `up{job=~"victoriametrics|vmselect|vminsert|vmstorage|vmagent|vmalert"} == 0`,
						For:    "2m",
						Labels: map[string]string{"severity": "critical"},
					}, {
						Alert: "ProcessNearFDLimits",
						Annotations: map[string]string{
							"description": "Exhausting OS file descriptors limit can cause severe degradation of the process. Consider to increase the limit as fast as possible.",
							"summary":     `Number of free file descriptors is less than 100 for "{{ $labels.job }}"("{{ $labels.instance }}") for the last 5m`,
						},
						Expr:   "(process_max_fds - process_open_fds) < 100",
						For:    "5m",
						Labels: map[string]string{"severity": "critical"},
					}, {
						Alert: "TooHighMemoryUsage",
						Annotations: map[string]string{
							"description": "Too high memory usage may result into multiple issues such as OOMs or degraded performance. Consider to either increase available memory or decrease the load on the process.",
							"summary":     `It is more than 90% of memory used by "{{ $labels.job }}"("{{ $labels.instance }}") during the last 5m`,
						},
						Expr:   "(process_resident_memory_anon_bytes / vm_available_memory_bytes) > 0.9",
						For:    "5m",
						Labels: map[string]string{"severity": "critical"},
					}, {
						Alert: "TooHighCPUUsage",
						Annotations: map[string]string{
							"description": "Too high CPU usage may be a sign of insufficient resources and make process unstable. Consider to either increase available CPU resources or decrease the load on the process.",
							"summary":     `More than 90% of CPU is used by "{{ $labels.job }}"("{{ $labels.instance }}") during the last 5m`,
						},
						Expr:   "rate(process_cpu_seconds_total[5m]) / process_cpu_cores_available > 0.9",
						For:    "5m",
						Labels: map[string]string{"severity": "critical"},
					}, {
						Alert: "TooManyLogs",
						Annotations: map[string]string{
							"description": `
Logging rate for job "{{ $labels.job }}" ({{ $labels.instance }}) is {{ $value }} for last 15m.
 Worth to check logs for specific error messages.
`,
							"summary": `Too many logs printed for job "{{ $labels.job }}" ({{ $labels.instance }})`,
						},
						Expr:   `sum(increase(vm_log_messages_total{level="error"}[5m])) by (job, instance) > 0`,
						For:    "15m",
						Labels: map[string]string{"severity": "warning"},
					},
				},
			},
		},
	},
	TypeMeta: TypeVMRuleV1Beta1,
}
