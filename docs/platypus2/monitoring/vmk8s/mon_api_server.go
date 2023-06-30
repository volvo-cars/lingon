// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmk8s

import (
	"github.com/VictoriaMetrics/operator/api/victoriametrics/v1beta1"
	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	"github.com/volvo-cars/lingoneks/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var AS = &meta.Metadata{
	Name:      "apiserver", // name linked to the job name in queries
	Namespace: namespace,
	Instance:  "apiserver" + namespace,
	Component: "monitoring",
	PartOf:    appName,
	Version:   version,
	ManagedBy: "lingon",
}

type MonAPIServer struct {
	kube.App

	APIServerAvailabilityRules *v1beta1.VMRule
	APIServerBurnRateRules     *v1beta1.VMRule
	APIServerHistogramRules    *v1beta1.VMRule
	APIServerSLOsRules         *v1beta1.VMRule
	APIServerRules             *v1beta1.VMRule
	APIServerScrape            *v1beta1.VMServiceScrape
}

func NewMonAPIServer() *MonAPIServer {
	return &MonAPIServer{
		APIServerAvailabilityRules: APIServerAvailabilityRules,
		APIServerBurnRateRules:     APIServerBurnRateRules,
		APIServerHistogramRules:    APIServerHistogramRules,
		APIServerSLOsRules:         APIServerSLOsRules,
		APIServerRules:             APIServerRules,
		APIServerScrape:            APIServerScrape,
	}
}

var APIServerScrape = &v1beta1.VMServiceScrape{
	ObjectMeta: AS.ObjectMeta(),
	Spec: v1beta1.VMServiceScrapeSpec{
		Endpoints: []v1beta1.Endpoint{
			{
				BearerTokenFile: PathSA + "/token",
				Port:            "https",
				Scheme:          "https",
				TLSConfig: &v1beta1.TLSConfig{
					CAFile:     PathSA + "/ca.crt",
					ServerName: "kubernetes",
				},
			},
		},
		JobLabel:          "component",
		NamespaceSelector: v1beta1.NamespaceSelector{MatchNames: []string{ku.NSKubeSystem}},
		Selector: metav1.LabelSelector{
			MatchLabels: map[string]string{
				"component": "kube-apiserver",
			},
		},
	},
	TypeMeta: TypeVMServiceScrapeV1Beta1,
}

var APIServerAvailabilityRules = &v1beta1.VMRule{
	ObjectMeta: AS.ObjectMetaNameSuffix("rules"),
	Spec: v1beta1.VMRuleSpec{
		Groups: []v1beta1.RuleGroup{
			{
				Interval: "3m",
				Name:     "kube-apiserver-availability.rules",
				Rules: []v1beta1.Rule{
					{
						Expr:   "avg_over_time(code_verb:apiserver_request_total:increase1h[30d]) * 24 * 30",
						Record: "code_verb:apiserver_request_total:increase30d",
					}, {
						Expr:   `sum by (cluster, code) (code_verb:apiserver_request_total:increase30d{verb=~"LIST|GET"})`,
						Labels: map[string]string{"verb": "read"},
						Record: "code:apiserver_request_total:increase30d",
					}, {
						Expr:   `sum by (cluster, code) (code_verb:apiserver_request_total:increase30d{verb=~"POST|PUT|PATCH|DELETE"})`,
						Labels: map[string]string{"verb": "write"},
						Record: "code:apiserver_request_total:increase30d",
					}, {
						Expr:   "sum by (cluster, verb, scope) (increase(apiserver_request_slo_duration_seconds_count[1h]))",
						Record: "cluster_verb_scope:apiserver_request_slo_duration_seconds_count:increase1h",
					}, {
						Expr:   "sum by (cluster, verb, scope) (avg_over_time(cluster_verb_scope:apiserver_request_slo_duration_seconds_count:increase1h[30d]) * 24 * 30)",
						Record: "cluster_verb_scope:apiserver_request_slo_duration_seconds_count:increase30d",
					}, {
						Expr:   "sum by (cluster, verb, scope, le) (increase(apiserver_request_slo_duration_seconds_bucket[1h]))",
						Record: "cluster_verb_scope_le:apiserver_request_slo_duration_seconds_bucket:increase1h",
					}, {
						Expr:   "sum by (cluster, verb, scope, le) (avg_over_time(cluster_verb_scope_le:apiserver_request_slo_duration_seconds_bucket:increase1h[30d]) * 24 * 30)",
						Record: "cluster_verb_scope_le:apiserver_request_slo_duration_seconds_bucket:increase30d",
					}, {
						Expr: `
1 - (
  (
    # write too slow
    sum by (cluster) (cluster_verb_scope:apiserver_request_slo_duration_seconds_count:increase30d{verb=~"POST|PUT|PATCH|DELETE"})
    -
    sum by (cluster) (cluster_verb_scope_le:apiserver_request_slo_duration_seconds_bucket:increase30d{verb=~"POST|PUT|PATCH|DELETE",le="1"})
  ) +
  (
    # read too slow
    sum by (cluster) (cluster_verb_scope:apiserver_request_slo_duration_seconds_count:increase30d{verb=~"LIST|GET"})
    -
    (
      (
        sum by (cluster) (cluster_verb_scope_le:apiserver_request_slo_duration_seconds_bucket:increase30d{verb=~"LIST|GET",scope=~"resource|",le="1"})
        or
        vector(0)
      )
      +
      sum by (cluster) (cluster_verb_scope_le:apiserver_request_slo_duration_seconds_bucket:increase30d{verb=~"LIST|GET",scope="namespace",le="5"})
      +
      sum by (cluster) (cluster_verb_scope_le:apiserver_request_slo_duration_seconds_bucket:increase30d{verb=~"LIST|GET",scope="cluster",le="30"})
    )
  ) +
  # errors
  sum by (cluster) (code:apiserver_request_total:increase30d{code=~"5.."} or vector(0))
)
/
sum by (cluster) (code:apiserver_request_total:increase30d)
`,
						Labels: map[string]string{"verb": "all"},
						Record: "apiserver_request:availability30d",
					}, {
						Expr: `
1 - (
  sum by (cluster) (cluster_verb_scope:apiserver_request_slo_duration_seconds_count:increase30d{verb=~"LIST|GET"})
  -
  (
    # too slow
    (
      sum by (cluster) (cluster_verb_scope_le:apiserver_request_slo_duration_seconds_bucket:increase30d{verb=~"LIST|GET",scope=~"resource|",le="1"})
      or
      vector(0)
    )
    +
    sum by (cluster) (cluster_verb_scope_le:apiserver_request_slo_duration_seconds_bucket:increase30d{verb=~"LIST|GET",scope="namespace",le="5"})
    +
    sum by (cluster) (cluster_verb_scope_le:apiserver_request_slo_duration_seconds_bucket:increase30d{verb=~"LIST|GET",scope="cluster",le="30"})
  )
  +
  # errors
  sum by (cluster) (code:apiserver_request_total:increase30d{verb="read",code=~"5.."} or vector(0))
)
/
sum by (cluster) (code:apiserver_request_total:increase30d{verb="read"})
`,
						Labels: map[string]string{"verb": "read"},
						Record: "apiserver_request:availability30d",
					}, {
						Expr: `
1 - (
  (
    # too slow
    sum by (cluster) (cluster_verb_scope:apiserver_request_slo_duration_seconds_count:increase30d{verb=~"POST|PUT|PATCH|DELETE"})
    -
    sum by (cluster) (cluster_verb_scope_le:apiserver_request_slo_duration_seconds_bucket:increase30d{verb=~"POST|PUT|PATCH|DELETE",le="1"})
  )
  +
  # errors
  sum by (cluster) (code:apiserver_request_total:increase30d{verb="write",code=~"5.."} or vector(0))
)
/
sum by (cluster) (code:apiserver_request_total:increase30d{verb="write"})
`,
						Labels: map[string]string{"verb": "write"},
						Record: "apiserver_request:availability30d",
					}, {
						Expr:   `sum by (cluster,code,resource) (rate(apiserver_request_total{job="` + AS.Name + `",verb=~"LIST|GET"}[5m]))`,
						Labels: map[string]string{"verb": "read"},
						Record: "code_resource:apiserver_request_total:rate5m",
					}, {
						Expr:   `sum by (cluster,code,resource) (rate(apiserver_request_total{job="` + AS.Name + `",verb=~"POST|PUT|PATCH|DELETE"}[5m]))`,
						Labels: map[string]string{"verb": "write"},
						Record: "code_resource:apiserver_request_total:rate5m",
					}, {
						Expr:   `sum by (cluster, code, verb) (increase(apiserver_request_total{job="` + AS.Name + `",verb=~"LIST|GET|POST|PUT|PATCH|DELETE",code=~"2.."}[1h]))`,
						Record: "code_verb:apiserver_request_total:increase1h",
					}, {
						Expr:   `sum by (cluster, code, verb) (increase(apiserver_request_total{job="` + AS.Name + `",verb=~"LIST|GET|POST|PUT|PATCH|DELETE",code=~"3.."}[1h]))`,
						Record: "code_verb:apiserver_request_total:increase1h",
					}, {
						Expr:   `sum by (cluster, code, verb) (increase(apiserver_request_total{job="` + AS.Name + `",verb=~"LIST|GET|POST|PUT|PATCH|DELETE",code=~"4.."}[1h]))`,
						Record: "code_verb:apiserver_request_total:increase1h",
					}, {
						Expr:   `sum by (cluster, code, verb) (increase(apiserver_request_total{job="` + AS.Name + `",verb=~"LIST|GET|POST|PUT|PATCH|DELETE",code=~"5.."}[1h]))`,
						Record: "code_verb:apiserver_request_total:increase1h",
					},
				},
			},
		},
	},
	TypeMeta: TypeVMRuleV1Beta1,
}

func burnrateRead(name, dur string) string {
	return `
(
  (
    # too slow
    sum by (cluster) (rate(apiserver_request_slo_duration_seconds_count{job="` + name + `",verb=~"LIST|GET",subresource!~"proxy|attach|log|exec|portforward"}[` + dur + `]))
    -
    (
      (
        sum by (cluster) (rate(apiserver_request_slo_duration_seconds_bucket{job="` + name + `",verb=~"LIST|GET",subresource!~"proxy|attach|log|exec|portforward",scope=~"resource|",le="1"}[` + dur + `]))
        or
        vector(0)
      )
      +
      sum by (cluster) (rate(apiserver_request_slo_duration_seconds_bucket{job="` + name + `",verb=~"LIST|GET",subresource!~"proxy|attach|log|exec|portforward",scope="namespace",le="5"}[` + dur + `]))
      +
      sum by (cluster) (rate(apiserver_request_slo_duration_seconds_bucket{job="` + name + `",verb=~"LIST|GET",subresource!~"proxy|attach|log|exec|portforward",scope="cluster",le="30"}[` + dur + `]))
    )
  )
  +
  # errors
  sum by (cluster) (rate(apiserver_request_total{job="` + name + `",verb=~"LIST|GET",code=~"5.."}[` + dur + `]))
)
/
sum by (cluster) (rate(apiserver_request_total{job="` + name + `",verb=~"LIST|GET"}[` + dur + `]))
`
}

func burnrateWrite(name, dur string) string {
	return `
(
  (
    # too slow
    sum by (cluster) (rate(apiserver_request_slo_duration_seconds_count{job="` + name + `",verb=~"POST|PUT|PATCH|DELETE",subresource!~"proxy|attach|log|exec|portforward"}[` + dur + `]))
    -
    sum by (cluster) (rate(apiserver_request_slo_duration_seconds_bucket{job="` + name + `",verb=~"POST|PUT|PATCH|DELETE",subresource!~"proxy|attach|log|exec|portforward",le="1"}[` + dur + `]))
  )
  +
  sum by (cluster) (rate(apiserver_request_total{job="` + name + `",verb=~"POST|PUT|PATCH|DELETE",code=~"5.."}[` + dur + `]))
)
/
sum by (cluster) (rate(apiserver_request_total{job="` + name + `",verb=~"POST|PUT|PATCH|DELETE"}[` + dur + `]))
`
}

func burnRateRules(name string, windows []string) []v1beta1.Rule {
	size := len(windows)
	res := make([]v1beta1.Rule, 2*size)
	for i, w := range windows {
		res[i] = v1beta1.Rule{
			Expr:   burnrateRead(name, w),
			Labels: map[string]string{"verb": "read"},
			Record: "apiserver_request:burnrate" + w,
		}
		res[i+size] = v1beta1.Rule{
			Expr:   burnrateWrite(name, w),
			Labels: map[string]string{"verb": "write"},
			Record: "apiserver_request:burnrate" + w,
		}
	}
	return res
}

var APIServerBurnRateRules = &v1beta1.VMRule{
	ObjectMeta: AS.ObjectMetaNameSuffix("burnrate"),
	Spec: v1beta1.VMRuleSpec{
		Groups: []v1beta1.RuleGroup{
			{
				Name: "kube-apiserver-burnrate.rules",
				Rules: burnRateRules(
					AS.Name,
					[]string{"5m", "30m", "1h", "2h", "6h", "1d", "3d"},
				),
			},
		},
	},
	TypeMeta: TypeVMRuleV1Beta1,
}

var APIServerHistogramRules = &v1beta1.VMRule{
	ObjectMeta: AS.ObjectMetaNameSuffix("histogram"),
	Spec: v1beta1.VMRuleSpec{
		Groups: []v1beta1.RuleGroup{
			{
				Name: "kube-apiserver-histogram.rules",
				Rules: []v1beta1.Rule{
					{
						Expr: `histogram_quantile(0.99, sum by (cluster, le, resource) (rate(apiserver_request_slo_duration_seconds_bucket{job="` + AS.Name + `",verb=~"LIST|GET",subresource!~"proxy|attach|log|exec|portforward"}[5m]))) > 0`,
						Labels: map[string]string{
							"quantile": "0.99",
							"verb":     "read",
						},
						Record: "cluster_quantile:apiserver_request_slo_duration_seconds:histogram_quantile",
					}, {
						Expr: `histogram_quantile(0.99, sum by (cluster, le, resource) (rate(apiserver_request_slo_duration_seconds_bucket{job="` + AS.Name + `",verb=~"POST|PUT|PATCH|DELETE",subresource!~"proxy|attach|log|exec|portforward"}[5m]))) > 0`,
						Labels: map[string]string{
							"quantile": "0.99",
							"verb":     "write",
						},
						Record: "cluster_quantile:apiserver_request_slo_duration_seconds:histogram_quantile",
					},
				},
			},
		},
	},
	TypeMeta: TypeVMRuleV1Beta1,
}

var APIServerSLOsRules = &v1beta1.VMRule{
	ObjectMeta: AS.ObjectMetaNameSuffix("slos"),
	Spec: v1beta1.VMRuleSpec{
		Groups: []v1beta1.RuleGroup{
			{
				Name: "kube-apiserver-slos",
				Rules: []v1beta1.Rule{
					{
						Alert: "KubeAPIErrorBudgetBurn",
						Annotations: map[string]string{
							"description": "The API server is burning too much error budget.",
							"runbook_url": "https://runbooks.prometheus-operator.dev/runbooks/kubernetes/kubeapierrorbudgetburn",
							"summary":     "The API server is burning too much error budget.",
						},
						Expr: `
sum(apiserver_request:burnrate1h) > (14.40 * 0.01000)
and
sum(apiserver_request:burnrate5m) > (14.40 * 0.01000)
`,
						For: "2m",
						Labels: map[string]string{
							"long":     "1h",
							"severity": "critical",
							"short":    "5m",
						},
					}, {
						Alert: "KubeAPIErrorBudgetBurn",
						Annotations: map[string]string{
							"description": "The API server is burning too much error budget.",
							"runbook_url": "https://runbooks.prometheus-operator.dev/runbooks/kubernetes/kubeapierrorbudgetburn",
							"summary":     "The API server is burning too much error budget.",
						},
						Expr: `
sum(apiserver_request:burnrate6h) > (6.00 * 0.01000)
and
sum(apiserver_request:burnrate30m) > (6.00 * 0.01000)
`,
						For: "15m",
						Labels: map[string]string{
							"long":     "6h",
							"severity": "critical",
							"short":    "30m",
						},
					}, {
						Alert: "KubeAPIErrorBudgetBurn",
						Annotations: map[string]string{
							"description": "The API server is burning too much error budget.",
							"runbook_url": "https://runbooks.prometheus-operator.dev/runbooks/kubernetes/kubeapierrorbudgetburn",
							"summary":     "The API server is burning too much error budget.",
						},
						Expr: `
sum(apiserver_request:burnrate1d) > (3.00 * 0.01000)
and
sum(apiserver_request:burnrate2h) > (3.00 * 0.01000)
`,
						For: "1h",
						Labels: map[string]string{
							"long":     "1d",
							"severity": "warning",
							"short":    "2h",
						},
					}, {
						Alert: "KubeAPIErrorBudgetBurn",
						Annotations: map[string]string{
							"description": "The API server is burning too much error budget.",
							"runbook_url": "https://runbooks.prometheus-operator.dev/runbooks/kubernetes/kubeapierrorbudgetburn",
							"summary":     "The API server is burning too much error budget.",
						},
						Expr: `
sum(apiserver_request:burnrate3d) > (1.00 * 0.01000)
and
sum(apiserver_request:burnrate6h) > (1.00 * 0.01000)
`,
						For: "3h",
						Labels: map[string]string{
							"long":     "3d",
							"severity": "warning",
							"short":    "6h",
						},
					},
				},
			},
		},
	},
	TypeMeta: TypeVMRuleV1Beta1,
}

var APIServerRules = &v1beta1.VMRule{
	ObjectMeta: AS.ObjectMetaNameSuffix("system"),
	Spec: v1beta1.VMRuleSpec{
		Groups: []v1beta1.RuleGroup{
			{
				Name: "kubernetes-system-apiserver",
				Rules: []v1beta1.Rule{
					{
						Alert: "KubeClientCertificateExpiration",
						Annotations: map[string]string{
							"description": "A client certificate used to authenticate to kubernetes apiserver is expiring in less than 7.0 days.",
							"runbook_url": "https://runbooks.prometheus-operator.dev/runbooks/kubernetes/kubeclientcertificateexpiration",
							"summary":     "Client certificate is about to expire.",
						},
						Expr:   `apiserver_client_certificate_expiration_seconds_count{job="` + AS.Name + `"} > 0 and on(job) histogram_quantile(0.01, sum by (job, le) (rate(apiserver_client_certificate_expiration_seconds_bucket{job="` + AS.Name + `"}[5m]))) < 604800`,
						For:    "5m",
						Labels: map[string]string{"severity": "warning"},
					}, {
						Alert: "KubeClientCertificateExpiration",
						Annotations: map[string]string{
							"description": "A client certificate used to authenticate to kubernetes apiserver is expiring in less than 24.0 hours.",
							"runbook_url": "https://runbooks.prometheus-operator.dev/runbooks/kubernetes/kubeclientcertificateexpiration",
							"summary":     "Client certificate is about to expire.",
						},
						Expr:   `apiserver_client_certificate_expiration_seconds_count{job="` + AS.Name + `"} > 0 and on(job) histogram_quantile(0.01, sum by (job, le) (rate(apiserver_client_certificate_expiration_seconds_bucket{job="` + AS.Name + `"}[5m]))) < 86400`,
						For:    "5m",
						Labels: map[string]string{"severity": "critical"},
					}, {
						Alert: "KubeAggregatedAPIErrors",
						Annotations: map[string]string{
							"description": "Kubernetes aggregated API {{ $labels.name }}/{{ $labels.namespace }} has reported errors. It has appeared unavailable {{ $value | humanize }} times averaged over the past 10m.",
							"runbook_url": "https://runbooks.prometheus-operator.dev/runbooks/kubernetes/kubeaggregatedapierrors",
							"summary":     "Kubernetes aggregated API has reported errors.",
						},
						Expr:   "sum by(name, namespace, cluster)(increase(aggregator_unavailable_apiservice_total[10m])) > 4",
						Labels: map[string]string{"severity": "warning"},
					}, {
						Alert: "KubeAggregatedAPIDown",
						Annotations: map[string]string{
							"description": "Kubernetes aggregated API {{ $labels.name }}/{{ $labels.namespace }} has been only {{ $value | humanize }}% available over the last 10m.",
							"runbook_url": "https://runbooks.prometheus-operator.dev/runbooks/kubernetes/kubeaggregatedapidown",
							"summary":     "Kubernetes aggregated API is down.",
						},
						Expr:   "(1 - max by(name, namespace, cluster)(avg_over_time(aggregator_unavailable_apiservice[10m]))) * 100 < 85",
						For:    "5m",
						Labels: map[string]string{"severity": "warning"},
					}, {
						Alert: "KubeAPIDown",
						Annotations: map[string]string{
							"description": "KubeAPI has disappeared from Prometheus target discovery.",
							"runbook_url": "https://runbooks.prometheus-operator.dev/runbooks/kubernetes/kubeapidown",
							"summary":     "Target disappeared from Prometheus target discovery.",
						},
						Expr:   `absent(up{job="` + AS.Name + `"} == 1)`,
						For:    "15m",
						Labels: map[string]string{"severity": "critical"},
					}, {
						Alert: "KubeAPITerminatedRequests",
						Annotations: map[string]string{
							"description": "The kubernetes apiserver has terminated {{ $value | humanizePercentage }} of its incoming requests.",
							"runbook_url": "https://runbooks.prometheus-operator.dev/runbooks/kubernetes/kubeapiterminatedrequests",
							"summary":     "The kubernetes apiserver has terminated {{ $value | humanizePercentage }} of its incoming requests.",
						},
						Expr:   `sum(rate(apiserver_request_terminations_total{job="` + AS.Name + `"}[10m]))  / (  sum(rate(apiserver_request_total{job="` + AS.Name + `"}[10m])) + sum(rate(apiserver_request_terminations_total{job="` + AS.Name + `"}[10m])) ) > 0.20`,
						For:    "5m",
						Labels: map[string]string{"severity": "warning"},
					},
				},
			},
		},
	},
	TypeMeta: TypeVMRuleV1Beta1,
}
