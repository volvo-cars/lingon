// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmk8s

import (
	"github.com/VictoriaMetrics/operator/api/victoriametrics/v1beta1"
	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	KSMVersion  = "2.9.2"
	KSMPort     = 8080
	KSMPortName = "http"
)

var KSM = &Metadata{
	Name:      "kube-state-metrics",
	Namespace: namespace,
	Instance:  "kube-state-metrics-" + namespace,
	Component: "metrics",
	PartOf:    appName,
	Version:   KSMVersion,
	ManagedBy: "lingon",
	Registry:  "registry.k8s.io",
	Image:     "kube-state-metrics/kube-state-metrics",
	Tag:       "v" + KSMVersion,
}

type KubeStateMetrics struct {
	kube.App

	CR     *rbacv1.ClusterRole
	CRB    *rbacv1.ClusterRoleBinding
	Deploy *appsv1.Deployment
	SA     *corev1.ServiceAccount
	SVC    *corev1.Service
	Rules  *v1beta1.VMRule
	Scrape *v1beta1.VMServiceScrape
}

func NewKubeStateMetrics() *KubeStateMetrics {
	return &KubeStateMetrics{
		CR:     KubeStateMetricsCR,
		CRB:    KubeStateMetricsCRB,
		Deploy: KubeStateMetricsDeploy,
		SA:     KubeStateMetricsSA,
		SVC:    KubeStateMetricsSVC,
		Scrape: KubeStateMetricsScrape,
		Rules:  KubeStateMetricsRules,
	}
}

var KubeStateMetricsDeploy = &appsv1.Deployment{
	TypeMeta:   ku.TypeDeploymentV1,
	ObjectMeta: KSM.ObjectMeta(),
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(1)),
		Selector: &metav1.LabelSelector{
			MatchLabels: KSM.MatchLabels(),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: KSM.Labels(),
			},
			Spec: corev1.PodSpec{
				ServiceAccountName: KubeStateMetricsSA.Name,
				Containers: []corev1.Container{
					{
						Args: []string{
							"--port=" + d(KSMPort),
							"--resources=certificatesigningrequests," +
								"configmaps,cronjobs,daemonsets,deployments," +
								"endpoints,horizontalpodautoscalers,ingresses," +
								"jobs,leases,limitranges," +
								"mutatingwebhookconfigurations,namespaces," +
								"networkpolicies,nodes,persistentvolumeclaims," +
								"persistentvolumes,poddisruptionbudgets,pods," +
								"replicasets,replicationcontrollers," +
								"resourcequotas,secrets,services," +
								"statefulsets,storageclasses," +
								"validatingwebhookconfigurations," +
								"volumeattachments",
						},
						Image:           KSM.ContainerURL(),
						ImagePullPolicy: corev1.PullIfNotPresent,
						LivenessProbe: &corev1.Probe{
							InitialDelaySeconds: int32(5),
							ProbeHandler: ku.ProbeHTTP(
								ku.PathHealthz, KSMPort,
							),
							TimeoutSeconds: int32(5),
						},
						Name: KSM.Name,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: int32(KSMPort),
								Name:          KSMPortName,
							},
						},
						ReadinessProbe: &corev1.Probe{
							InitialDelaySeconds: int32(5),
							ProbeHandler:        ku.ProbeHTTP("/", KSMPort),
							TimeoutSeconds:      int32(5),
						},
					},
				},
				SecurityContext: &corev1.PodSecurityContext{
					FSGroup:    P(int64(65534)),
					RunAsGroup: P(int64(65534)),
					RunAsUser:  P(int64(65534)),
				},
			},
		},
	},
}

var KubeStateMetricsCRB = ku.BindClusterRole(
	KSM.Name,
	KubeStateMetricsSA,
	KubeStateMetricsCR,
	KSM.Labels(),
)

var KubeStateMetricsSVC = &corev1.Service{
	TypeMeta: ku.TypeServiceV1,
	ObjectMeta: KSM.ObjectMetaAnnotations(
		map[string]string{"prometheus.io/scrape": "true"},
		// TODO: check! probably a leftover from kube-prometheus-stack
		// ku.AnnotationPrometheus("/",KSMPort),
	),
	Spec: corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			{
				Name:       KSMPortName,
				Port:       int32(KSMPort),
				Protocol:   corev1.ProtocolTCP,
				TargetPort: intstr.FromInt(KSMPort),
			},
		},
		Selector: KSM.MatchLabels(),
		Type:     corev1.ServiceTypeClusterIP,
	},
}

var KubeStateMetricsSA = ku.ServiceAccount(
	KSM.Name,
	KSM.Namespace,
	KSM.Labels(),
	nil,
)

var KubeStateMetricsCR = &rbacv1.ClusterRole{
	TypeMeta:   ku.TypeClusterRoleV1,
	ObjectMeta: KSM.ObjectMetaNoNS(),
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{"certificates.k8s.io"},
			Resources: []string{"certificatesigningrequests"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"configmaps"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{"batch"},
			Resources: []string{"cronjobs"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{"extensions", "apps"},
			Resources: []string{"daemonsets"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{"extensions", "apps"},
			Resources: []string{"deployments"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"endpoints"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{"autoscaling"},
			Resources: []string{"horizontalpodautoscalers"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{"extensions", "networking.k8s.io"},
			Resources: []string{"ingresses"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{"batch"},
			Resources: []string{"jobs"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{"coordination.k8s.io"},
			Resources: []string{"leases"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"limitranges"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{"admissionregistration.k8s.io"},
			Resources: []string{"mutatingwebhookconfigurations"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"namespaces"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{"networking.k8s.io"},
			Resources: []string{"networkpolicies"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"nodes"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"persistentvolumeclaims"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"persistentvolumes"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{"policy"},
			Resources: []string{"poddisruptionbudgets"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"pods"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{"extensions", "apps"},
			Resources: []string{"replicasets"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"replicationcontrollers"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"resourcequotas"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"secrets"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"services"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{"apps"},
			Resources: []string{"statefulsets"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{"storage.k8s.io"},
			Resources: []string{"storageclasses"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{"admissionregistration.k8s.io"},
			Resources: []string{"validatingwebhookconfigurations"},
			Verbs:     []string{"list", "watch"},
		}, {
			APIGroups: []string{"storage.k8s.io"},
			Resources: []string{"volumeattachments"},
			Verbs:     []string{"list", "watch"},
		},
	},
}

var KubeStateMetricsRules = &v1beta1.VMRule{
	TypeMeta:   TypeVMRuleV1Beta1,
	ObjectMeta: KSM.ObjectMeta(),
	Spec: v1beta1.VMRuleSpec{
		Groups: []v1beta1.RuleGroup{
			{
				Name: KSM.Name,
				Rules: []v1beta1.Rule{
					{
						Alert: "KubeStateMetricsListErrors",
						Annotations: map[string]string{
							"description": "kube-state-metrics is experiencing errors at an elevated rate in list operations. This is likely causing it to not be able to expose metrics about Kubernetes objects correctly or at all.",
							"runbook_url": "https://runbooks.prometheus-operator.dev/runbooks/kube-state-metrics/kubestatemetricslisterrors",
							"summary":     "kube-state-metrics is experiencing errors in list operations.",
						},
						Expr: `
(sum(rate(kube_state_metrics_list_total{job="kube-state-metrics",result="error"}[5m]))
  /
sum(rate(kube_state_metrics_list_total{job="kube-state-metrics"}[5m])))
> 0.01
`,
						For:    "15m",
						Labels: map[string]string{"severity": "critical"},
					}, {
						Alert: "KubeStateMetricsWatchErrors",
						Annotations: map[string]string{
							"description": "kube-state-metrics is experiencing errors at an elevated rate in watch operations. This is likely causing it to not be able to expose metrics about Kubernetes objects correctly or at all.",
							"runbook_url": "https://runbooks.prometheus-operator.dev/runbooks/kube-state-metrics/kubestatemetricswatcherrors",
							"summary":     "kube-state-metrics is experiencing errors in watch operations.",
						},
						Expr: `
(sum(rate(kube_state_metrics_watch_total{job="kube-state-metrics",result="error"}[5m]))
  /
sum(rate(kube_state_metrics_watch_total{job="kube-state-metrics"}[5m])))
> 0.01
`,
						For:    "15m",
						Labels: map[string]string{"severity": "critical"},
					}, {
						Alert: "KubeStateMetricsShardingMismatch",
						Annotations: map[string]string{
							"description": "kube-state-metrics pods are running with different --total-shards configuration, some Kubernetes objects may be exposed multiple times or not exposed at all.",
							"runbook_url": "https://runbooks.prometheus-operator.dev/runbooks/kube-state-metrics/kubestatemetricsshardingmismatch",
							"summary":     "kube-state-metrics sharding is misconfigured.",
						},
						Expr:   `stdvar (kube_state_metrics_total_shards{job="kube-state-metrics"}) != 0`,
						For:    "15m",
						Labels: map[string]string{"severity": "critical"},
					}, {
						Alert: "KubeStateMetricsShardsMissing",
						Annotations: map[string]string{
							"description": "kube-state-metrics shards are missing, some Kubernetes objects are not being exposed.",
							"runbook_url": "https://runbooks.prometheus-operator.dev/runbooks/kube-state-metrics/kubestatemetricsshardsmissing",
							"summary":     "kube-state-metrics shards are missing.",
						},
						Expr: `
2^max(kube_state_metrics_total_shards{job="kube-state-metrics"}) - 1
  -
sum( 2 ^ max by (shard_ordinal) (kube_state_metrics_shard_ordinal{job="kube-state-metrics"}) )
!= 0
`,
						For:    "15m",
						Labels: map[string]string{"severity": "critical"},
					},
				},
			},
		},
	},
}

var KubeStateMetricsScrape = &v1beta1.VMServiceScrape{
	TypeMeta:   TypeVMServiceScrapeV1Beta1,
	ObjectMeta: KSM.ObjectMeta(),
	Spec: v1beta1.VMServiceScrapeSpec{
		Endpoints: []v1beta1.Endpoint{
			{
				HonorLabels: true,
				MetricRelabelConfigs: []*v1beta1.RelabelConfig{
					{
						Action: "labeldrop",
						Regex:  "(uid|container_id|image_id)",
					},
				},
				Port: KSMPortName,
			},
		},
		JobLabel: ku.AppLabelName, // "app.kubernetes.io/name",
		Selector: metav1.LabelSelector{MatchLabels: KSM.MatchLabels()},
	},
}
