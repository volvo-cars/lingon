package promstack

import (
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var KubeStateMetricsCR = &rbacv1.ClusterRole{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/component":  "metrics",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "kube-state-metrics",
			"app.kubernetes.io/part-of":    "kube-state-metrics",
			"app.kubernetes.io/version":    "2.8.2",
			"helm.sh/chart":                "kube-state-metrics-5.5.0",
			"release":                      "kube-prometheus-stack",
		},
		Name: "kube-prometheus-stack-kube-state-metrics",
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

var KubeStateMetricsSVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{"prometheus.io/scrape": "true"},
		Labels: map[string]string{
			"app.kubernetes.io/component":  "metrics",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "kube-state-metrics",
			"app.kubernetes.io/part-of":    "kube-state-metrics",
			"app.kubernetes.io/version":    "2.8.2",
			"helm.sh/chart":                "kube-state-metrics-5.5.0",
			"release":                      "kube-prometheus-stack",
		},
		Name:      "kube-prometheus-stack-kube-state-metrics",
		Namespace: namespace,
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
			"app.kubernetes.io/instance": "kube-prometheus-stack",
			"app.kubernetes.io/name":     "kube-state-metrics",
		},
		Type: corev1.ServiceType("ClusterIP"),
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	},
}

var KubeStateMetricsDeploy = &appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/component":  "metrics",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "kube-state-metrics",
			"app.kubernetes.io/part-of":    "kube-state-metrics",
			"app.kubernetes.io/version":    "2.8.2",
			"helm.sh/chart":                "kube-state-metrics-5.5.0",
			"release":                      "kube-prometheus-stack",
		},
		Name:      "kube-prometheus-stack-kube-state-metrics",
		Namespace: namespace,
	},
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(1)),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app.kubernetes.io/instance": "kube-prometheus-stack",
				"app.kubernetes.io/name":     "kube-state-metrics",
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app.kubernetes.io/component":  "metrics",
					"app.kubernetes.io/instance":   "kube-prometheus-stack",
					"app.kubernetes.io/managed-by": "Helm",
					"app.kubernetes.io/name":       "kube-state-metrics",
					"app.kubernetes.io/part-of":    "kube-state-metrics",
					"app.kubernetes.io/version":    "2.8.2",
					"helm.sh/chart":                "kube-state-metrics-5.5.0",
					"release":                      "kube-prometheus-stack",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Args: []string{
							"--port=8080",
							"--resources=certificatesigningrequests,configmaps,cronjobs,daemonsets,deployments,endpoints,horizontalpodautoscalers,ingresses,jobs,leases,limitranges,mutatingwebhookconfigurations,namespaces,networkpolicies,nodes,persistentvolumeclaims,persistentvolumes,poddisruptionbudgets,pods,replicasets,replicationcontrollers,resourcequotas,secrets,services,statefulsets,storageclasses,validatingwebhookconfigurations,volumeattachments",
						},
						Image:           "registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.8.2",
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
						SecurityContext: &corev1.SecurityContext{Capabilities: &corev1.Capabilities{Drop: []corev1.Capability{corev1.Capability("ALL")}}},
					},
				},
				SecurityContext: &corev1.PodSecurityContext{
					FSGroup:        P(int64(65534)),
					RunAsGroup:     P(int64(65534)),
					RunAsNonRoot:   P(true),
					RunAsUser:      P(int64(65534)),
					SeccompProfile: &corev1.SeccompProfile{Type: corev1.SeccompProfileType("RuntimeDefault")},
				},
				ServiceAccountName: "kube-prometheus-stack-kube-state-metrics",
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
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "kube-state-metrics",
			"app.kubernetes.io/part-of":    "kube-state-metrics",
			"app.kubernetes.io/version":    "2.8.2",
			"helm.sh/chart":                "kube-state-metrics-5.5.0",
			"release":                      "kube-prometheus-stack",
		},
		Name: "kube-prometheus-stack-kube-state-metrics",
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "ClusterRole",
		Name:     "kube-prometheus-stack-kube-state-metrics",
	},
	Subjects: []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      "kube-prometheus-stack-kube-state-metrics",
			Namespace: namespace,
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "ClusterRoleBinding",
	},
}

var KubeStateMetricsSA = &corev1.ServiceAccount{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/component":  "metrics",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "kube-state-metrics",
			"app.kubernetes.io/part-of":    "kube-state-metrics",
			"app.kubernetes.io/version":    "2.8.2",
			"helm.sh/chart":                "kube-state-metrics-5.5.0",
			"release":                      "kube-prometheus-stack",
		},
		Name:      "kube-prometheus-stack-kube-state-metrics",
		Namespace: namespace,
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ServiceAccount",
	},
}

var KubeStateMetricsServiceMonitor = &v1.ServiceMonitor{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/component":  "metrics",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "kube-state-metrics",
			"app.kubernetes.io/part-of":    "kube-state-metrics",
			"app.kubernetes.io/version":    "2.8.2",
			"helm.sh/chart":                "kube-state-metrics-5.5.0",
			"release":                      "kube-prometheus-stack",
		},
		Name:      "kube-prometheus-stack-kube-state-metrics",
		Namespace: namespace,
	},
	Spec: v1.ServiceMonitorSpec{
		Endpoints: []v1.Endpoint{
			{
				HonorLabels: true,
				Port:        "http",
			},
		},
		JobLabel: "app.kubernetes.io/name",
		Selector: metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app.kubernetes.io/instance": "kube-prometheus-stack",
				"app.kubernetes.io/name":     "kube-state-metrics",
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "monitoring.coreos.com/v1",
		Kind:       "ServiceMonitor",
	},
}
