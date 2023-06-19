package vmk8s

import (
	"github.com/VictoriaMetrics/operator/api/victoriametrics/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type KubeStateMetrics struct {
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
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/component":  "metrics",
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "kube-state-metrics",
			"app.kubernetes.io/part-of":    "kube-state-metrics",
			"app.kubernetes.io/version":    "2.7.0",
			"helm.sh/chart":                "kube-state-metrics-4.24.0",
		},
		Name:      "vmk8s-kube-state-metrics",
		Namespace: "monitoring",
	},
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(1)),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app.kubernetes.io/instance": "vmk8s",
				"app.kubernetes.io/name":     "kube-state-metrics",
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app.kubernetes.io/component":  "metrics",
					"app.kubernetes.io/instance":   "vmk8s",
					"app.kubernetes.io/managed-by": "Helm",
					"app.kubernetes.io/name":       "kube-state-metrics",
					"app.kubernetes.io/part-of":    "kube-state-metrics",
					"app.kubernetes.io/version":    "2.7.0",
					"helm.sh/chart":                "kube-state-metrics-4.24.0",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Args: []string{
							"--port=8080",
							"--resources=certificatesigningrequests,configmaps,cronjobs,daemonsets,deployments,endpoints,horizontalpodautoscalers,ingresses,jobs,leases,limitranges,mutatingwebhookconfigurations,namespaces,networkpolicies,nodes,persistentvolumeclaims,persistentvolumes,poddisruptionbudgets,pods,replicasets,replicationcontrollers,resourcequotas,secrets,services,statefulsets,storageclasses,validatingwebhookconfigurations,volumeattachments",
						},
						Image:           "registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.7.0",
						ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
						LivenessProbe: &corev1.Probe{
							InitialDelaySeconds: int32(5),
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/healthz",
									Port: intstr.IntOrString{IntVal: int32(8080)},
								},
							},
							TimeoutSeconds: int32(5),
						},
						Name: "kube-state-metrics",
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: int32(8080),
								Name:          "http",
							},
						},
						ReadinessProbe: &corev1.Probe{
							InitialDelaySeconds: int32(5),
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/",
									Port: intstr.IntOrString{IntVal: int32(8080)},
								},
							},
							TimeoutSeconds: int32(5),
						},
					},
				},
				SecurityContext: &corev1.PodSecurityContext{
					FSGroup:    P(int64(65534)),
					RunAsGroup: P(int64(65534)),
					RunAsUser:  P(int64(65534)),
				},
				ServiceAccountName: "vmk8s-kube-state-metrics",
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "apps/v1",
		Kind:       "Deployment",
	},
}

var KubeStateMetricsCRB = &rbacv1.ClusterRoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/component":  "metrics",
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "kube-state-metrics",
			"app.kubernetes.io/part-of":    "kube-state-metrics",
			"app.kubernetes.io/version":    "2.7.0",
			"helm.sh/chart":                "kube-state-metrics-4.24.0",
		},
		Name: "vmk8s-kube-state-metrics",
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "ClusterRole",
		Name:     "vmk8s-kube-state-metrics",
	},
	Subjects: []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      "vmk8s-kube-state-metrics",
			Namespace: "monitoring",
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "ClusterRoleBinding",
	},
}

var KubeStateMetricsSVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{"prometheus.io/scrape": "true"},
		Labels: map[string]string{
			"app.kubernetes.io/component":  "metrics",
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "kube-state-metrics",
			"app.kubernetes.io/part-of":    "kube-state-metrics",
			"app.kubernetes.io/version":    "2.7.0",
			"helm.sh/chart":                "kube-state-metrics-4.24.0",
		},
		Name:      "vmk8s-kube-state-metrics",
		Namespace: "monitoring",
	},
	Spec: corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			{
				Name:       "http",
				Port:       int32(8080),
				Protocol:   corev1.Protocol("TCP"),
				TargetPort: intstr.IntOrString{IntVal: int32(8080)},
			},
		},
		Selector: map[string]string{
			"app.kubernetes.io/instance": "vmk8s",
			"app.kubernetes.io/name":     "kube-state-metrics",
		},
		Type: corev1.ServiceType("ClusterIP"),
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	},
}

var KubeStateMetricsSA = &corev1.ServiceAccount{
	ImagePullSecrets: []corev1.LocalObjectReference{},
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/component":  "metrics",
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "kube-state-metrics",
			"app.kubernetes.io/part-of":    "kube-state-metrics",
			"app.kubernetes.io/version":    "2.7.0",
			"helm.sh/chart":                "kube-state-metrics-4.24.0",
		},
		Name:      "vmk8s-kube-state-metrics",
		Namespace: "monitoring",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ServiceAccount",
	},
}

var KubeStateMetricsCR = &rbacv1.ClusterRole{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/component":  "metrics",
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "kube-state-metrics",
			"app.kubernetes.io/part-of":    "kube-state-metrics",
			"app.kubernetes.io/version":    "2.7.0",
			"helm.sh/chart":                "kube-state-metrics-4.24.0",
		},
		Name: "vmk8s-kube-state-metrics",
	},
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
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "ClusterRole",
	},
}

var KubeStateMetricsRules = &v1beta1.VMRule{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "victoria-metrics-k8s-stack",
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-kube-state-metrics",
		Namespace: "monitoring",
	},
	Spec: v1beta1.VMRuleSpec{
		Groups: []v1beta1.RuleGroup{
			{
				Name: "kube-state-metrics",
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
						Expr:   "stdvar (kube_state_metrics_total_shards{job=\"kube-state-metrics\"}) != 0",
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
	TypeMeta: metav1.TypeMeta{
		APIVersion: "operator.victoriametrics.com/v1beta1",
		Kind:       "VMRule",
	},
}

var KubeStateMetricsScrape = &v1beta1.VMServiceScrape{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-kube-state-metrics",
		Namespace: "monitoring",
	},
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
				Port: "http",
			},
		},
		JobLabel: "app.kubernetes.io/name",
		Selector: metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app.kubernetes.io/instance": "vmk8s",
				"app.kubernetes.io/name":     "kube-state-metrics",
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "operator.victoriametrics.com/v1beta1",
		Kind:       "VMServiceScrape",
	},
}
