package promstack

import (
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var GrafanaSecrets = &corev1.Secret{
	Data: map[string][]byte{
		"admin-password": []byte("prom-operator"),
		"admin-user":     []byte("admin"),
		"ldap-toml":      []byte(""),
	},
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.5.1",
			"helm.sh/chart":                "grafana-6.56.2",
		},
		Name:      "kube-prometheus-stack-grafana",
		Namespace: namespace,
	},
	Type: corev1.SecretType("Opaque"),
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Secret",
	},
} // TODO: SECRETS SHOULD BE STORED ELSEWHERE THAN IN THE CODE!!!!

var GrafanaTestPO = &corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "test-success",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
		},
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.5.1",
			"helm.sh/chart":                "grafana-6.56.2",
		},
		Name:      "kube-prometheus-stack-grafana-test",
		Namespace: namespace,
	},
	Spec: corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Command: []string{
					"/opt/bats/bin/bats",
					"-t",
					"/tests/run.sh",
				},
				Image:           "docker.io/bats/bats:v1.4.1",
				ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
				Name:            "kube-prometheus-stack-test",
				VolumeMounts: []corev1.VolumeMount{
					{
						MountPath: "/tests",
						Name:      "tests",
						ReadOnly:  true,
					},
				},
			},
		},
		RestartPolicy:      corev1.RestartPolicy("Never"),
		ServiceAccountName: "kube-prometheus-stack-grafana-test",
		Volumes: []corev1.Volume{
			{
				Name:         "tests",
				VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: "kube-prometheus-stack-grafana-test"}}},
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Pod",
	},
}

var GrafanaDeploy = &appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.5.1",
			"helm.sh/chart":                "grafana-6.56.2",
		},
		Name:      "kube-prometheus-stack-grafana",
		Namespace: namespace,
	},
	Spec: appsv1.DeploymentSpec{
		Replicas:             P(int32(1)),
		RevisionHistoryLimit: P(int32(10)),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app.kubernetes.io/instance": "kube-prometheus-stack",
				"app.kubernetes.io/name":     "grafana",
			},
		},
		Strategy: appsv1.DeploymentStrategy{Type: appsv1.DeploymentStrategyType("RollingUpdate")},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					"checksum/config":                         "f3980bde957e01c9ff3b8301e865a08da4592cb2093f0fda2476cf95af9cae55",
					"checksum/dashboards-json-config":         "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
					"checksum/sc-dashboard-provider-config":   "992c887d7b187a2bbc6ffd68310aa1ccbbb9d03f814379692ef23cca0cab35a1",
					"checksum/secret":                         "dcf75285610f32ea17bf2f05b0a41f29330c7dd46109e8f866b5c48c5fea03c3",
					"kubectl.kubernetes.io/default-container": "grafana",
				},
				Labels: map[string]string{
					"app.kubernetes.io/instance": "kube-prometheus-stack",
					"app.kubernetes.io/name":     "grafana",
				},
			},
			Spec: corev1.PodSpec{
				AutomountServiceAccountToken: P(true),
				Containers: []corev1.Container{
					{
						Env: []corev1.EnvVar{
							{
								Name:  "METHOD",
								Value: "WATCH",
							}, {
								Name:  "LABEL",
								Value: "grafana_dashboard",
							}, {
								Name:  "LABEL_VALUE",
								Value: "1",
							}, {
								Name:  "FOLDER",
								Value: "/tmp/dashboards",
							}, {
								Name:  "RESOURCE",
								Value: "both",
							}, {
								Name:  "NAMESPACE",
								Value: "ALL",
							}, {
								Name: "REQ_USERNAME",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										Key:                  "admin-user",
										LocalObjectReference: corev1.LocalObjectReference{Name: "kube-prometheus-stack-grafana"},
									},
								},
							}, {
								Name: "REQ_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										Key:                  "admin-password",
										LocalObjectReference: corev1.LocalObjectReference{Name: "kube-prometheus-stack-grafana"},
									},
								},
							}, {
								Name:  "REQ_URL",
								Value: "http://localhost:3000/api/admin/provisioning/dashboards/reload",
							}, {
								Name:  "REQ_METHOD",
								Value: "POST",
							},
						},
						Image:           "quay.io/kiwigrid/k8s-sidecar:1.22.0",
						ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
						Name:            "grafana-sc-dashboard",
						SecurityContext: &corev1.SecurityContext{
							Capabilities:   &corev1.Capabilities{Drop: []corev1.Capability{corev1.Capability("ALL")}},
							SeccompProfile: &corev1.SeccompProfile{Type: corev1.SeccompProfileType("RuntimeDefault")},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/tmp/dashboards",
								Name:      "sc-dashboard-volume",
							},
						},
					}, {
						Env: []corev1.EnvVar{
							{
								Name:  "METHOD",
								Value: "WATCH",
							}, {
								Name:  "LABEL",
								Value: "grafana_datasource",
							}, {
								Name:  "LABEL_VALUE",
								Value: "1",
							}, {
								Name:  "FOLDER",
								Value: "/etc/grafana/provisioning/datasources",
							}, {
								Name:  "RESOURCE",
								Value: "both",
							}, {
								Name: "REQ_USERNAME",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										Key:                  "admin-user",
										LocalObjectReference: corev1.LocalObjectReference{Name: "kube-prometheus-stack-grafana"},
									},
								},
							}, {
								Name: "REQ_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										Key:                  "admin-password",
										LocalObjectReference: corev1.LocalObjectReference{Name: "kube-prometheus-stack-grafana"},
									},
								},
							}, {
								Name:  "REQ_URL",
								Value: "http://localhost:3000/api/admin/provisioning/datasources/reload",
							}, {
								Name:  "REQ_METHOD",
								Value: "POST",
							},
						},
						Image:           "quay.io/kiwigrid/k8s-sidecar:1.22.0",
						ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
						Name:            "grafana-sc-datasources",
						SecurityContext: &corev1.SecurityContext{
							Capabilities:   &corev1.Capabilities{Drop: []corev1.Capability{corev1.Capability("ALL")}},
							SeccompProfile: &corev1.SeccompProfile{Type: corev1.SeccompProfileType("RuntimeDefault")},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/etc/grafana/provisioning/datasources",
								Name:      "sc-datasources-volume",
							},
						},
					}, {
						Env: []corev1.EnvVar{
							{
								Name:      "POD_IP",
								ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.podIP"}},
							}, {
								Name: "GF_SECURITY_ADMIN_USER",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										Key:                  "admin-user",
										LocalObjectReference: corev1.LocalObjectReference{Name: "kube-prometheus-stack-grafana"},
									},
								},
							}, {
								Name: "GF_SECURITY_ADMIN_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										Key:                  "admin-password",
										LocalObjectReference: corev1.LocalObjectReference{Name: "kube-prometheus-stack-grafana"},
									},
								},
							}, {
								Name:  "GF_PATHS_DATA",
								Value: "/var/lib/grafana/",
							}, {
								Name:  "GF_PATHS_LOGS",
								Value: "/var/log/grafana",
							}, {
								Name:  "GF_PATHS_PLUGINS",
								Value: "/var/lib/grafana/plugins",
							}, {
								Name:  "GF_PATHS_PROVISIONING",
								Value: "/etc/grafana/provisioning",
							},
						},
						Image:           "docker.io/grafana/grafana:9.5.1",
						ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
						LivenessProbe: &corev1.Probe{
							FailureThreshold:    int32(10),
							InitialDelaySeconds: int32(60),
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/api/health",
									Port: intstr.IntOrString{IntVal: int32(3000)},
								},
							},
							TimeoutSeconds: int32(30),
						},
						Name: "grafana",
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: int32(3000),
								Name:          "grafana",
								Protocol:      corev1.Protocol("TCP"),
							}, {
								ContainerPort: int32(9094),
								Name:          "gossip-tcp",
								Protocol:      corev1.Protocol("TCP"),
							}, {
								ContainerPort: int32(9094),
								Name:          "gossip-udp",
								Protocol:      corev1.Protocol("UDP"),
							},
						},
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/api/health",
									Port: intstr.IntOrString{IntVal: int32(3000)},
								},
							},
						},
						SecurityContext: &corev1.SecurityContext{
							Capabilities:   &corev1.Capabilities{Drop: []corev1.Capability{corev1.Capability("ALL")}},
							SeccompProfile: &corev1.SeccompProfile{Type: corev1.SeccompProfileType("RuntimeDefault")},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/etc/grafana/grafana.ini",
								Name:      "config",
								SubPath:   "grafana.ini",
							}, {
								MountPath: "/var/lib/grafana",
								Name:      "storage",
							}, {
								MountPath: "/tmp/dashboards",
								Name:      "sc-dashboard-volume",
							}, {
								MountPath: "/etc/grafana/provisioning/dashboards/sc-dashboardproviders.yaml",
								Name:      "sc-dashboard-provider",
								SubPath:   "provider.yaml",
							}, {
								MountPath: "/etc/grafana/provisioning/datasources",
								Name:      "sc-datasources-volume",
							},
						},
					},
				},
				EnableServiceLinks: P(true),
				SecurityContext: &corev1.PodSecurityContext{
					FSGroup:      P(int64(472)),
					RunAsGroup:   P(int64(472)),
					RunAsNonRoot: P(true),
					RunAsUser:    P(int64(472)),
				},
				ServiceAccountName: "kube-prometheus-stack-grafana",
				Volumes: []corev1.Volume{
					{
						Name:         "config",
						VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: "kube-prometheus-stack-grafana"}}},
					}, {
						Name:         "storage",
						VolumeSource: corev1.VolumeSource{},
					}, {
						Name:         "sc-dashboard-volume",
						VolumeSource: corev1.VolumeSource{},
					}, {
						Name:         "sc-dashboard-provider",
						VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: "kube-prometheus-stack-grafana-config-dashboards"}}},
					}, {
						Name:         "sc-datasources-volume",
						VolumeSource: corev1.VolumeSource{},
					},
				},
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "apps/v1",
		Kind:       "Deployment",
	},
}

var GrafanaServiceMonitor = &v1.ServiceMonitor{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.5.1",
			"helm.sh/chart":                "grafana-6.56.2",
		},
		Name:      "kube-prometheus-stack-grafana",
		Namespace: namespace,
	},
	Spec: v1.ServiceMonitorSpec{
		Endpoints: []v1.Endpoint{
			{
				HonorLabels:   true,
				Path:          "/metrics",
				Port:          "http-web",
				Scheme:        "http",
				ScrapeTimeout: v1.Duration("30s"),
			},
		},
		JobLabel:          "kube-prometheus-stack",
		NamespaceSelector: v1.NamespaceSelector{MatchNames: []string{"monitoring"}},
		Selector: metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app.kubernetes.io/instance": "kube-prometheus-stack",
				"app.kubernetes.io/name":     "grafana",
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "monitoring.coreos.com/v1",
		Kind:       "ServiceMonitor",
	},
}

var GrafanaSA = &corev1.ServiceAccount{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.5.1",
			"helm.sh/chart":                "grafana-6.56.2",
		},
		Name:      "kube-prometheus-stack-grafana",
		Namespace: namespace,
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ServiceAccount",
	},
}

var GrafanaTestSA = &corev1.ServiceAccount{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "test-success",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
		},
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.5.1",
			"helm.sh/chart":                "grafana-6.56.2",
		},
		Name:      "kube-prometheus-stack-grafana-test",
		Namespace: namespace,
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ServiceAccount",
	},
}

var GrafanaCRB = &rbacv1.ClusterRoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.5.1",
			"helm.sh/chart":                "grafana-6.56.2",
		},
		Name: "kube-prometheus-stack-grafana-clusterrolebinding",
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "ClusterRole",
		Name:     "kube-prometheus-stack-grafana-clusterrole",
	},
	Subjects: []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      "kube-prometheus-stack-grafana",
			Namespace: namespace,
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "ClusterRoleBinding",
	},
}

var GrafanaCR = &rbacv1.ClusterRole{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.5.1",
			"helm.sh/chart":                "grafana-6.56.2",
		},
		Name: "kube-prometheus-stack-grafana-clusterrole",
	},
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{""},
			Resources: []string{"configmaps", "secrets"},
			Verbs:     []string{"get", "watch", "list"},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "ClusterRole",
	},
}

var GrafanaRole = &rbacv1.Role{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.5.1",
			"helm.sh/chart":                "grafana-6.56.2",
		},
		Name:      "kube-prometheus-stack-grafana",
		Namespace: namespace,
	},
	Rules: []rbacv1.PolicyRule{},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "Role",
	},
}

var GrafanaRB = &rbacv1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.5.1",
			"helm.sh/chart":                "grafana-6.56.2",
		},
		Name:      "kube-prometheus-stack-grafana",
		Namespace: namespace,
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "Role",
		Name:     "kube-prometheus-stack-grafana",
	},
	Subjects: []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      "kube-prometheus-stack-grafana",
			Namespace: namespace,
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "RoleBinding",
	},
}

var GrafanaSVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.5.1",
			"helm.sh/chart":                "grafana-6.56.2",
		},
		Name:      "kube-prometheus-stack-grafana",
		Namespace: namespace,
	},
	Spec: corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			{
				Name:       "http-web",
				Port:       int32(80),
				Protocol:   corev1.Protocol("TCP"),
				TargetPort: intstr.IntOrString{IntVal: int32(3000)},
			},
		},
		Selector: map[string]string{
			"app.kubernetes.io/instance": "kube-prometheus-stack",
			"app.kubernetes.io/name":     "grafana",
		},
		Type: corev1.ServiceType("ClusterIP"),
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	},
}

var GrafanaCM = &corev1.ConfigMap{
	Data: map[string]string{
		"grafana.ini": `
[analytics]
check_for_updates = true
[grafana_net]
url = https://grafana.net
[log]
mode = console
[paths]
data = /var/lib/grafana/
logs = /var/log/grafana
plugins = /var/lib/grafana/plugins
provisioning = /etc/grafana/provisioning
[server]
domain = ''

`,
	},
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.5.1",
			"helm.sh/chart":                "grafana-6.56.2",
		},
		Name:      "kube-prometheus-stack-grafana",
		Namespace: namespace,
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ConfigMap",
	},
}

var GrafanaConfigDashboardsCM = &corev1.ConfigMap{
	Data: map[string]string{
		"provider.yaml": `
apiVersion: 1
providers:
  - name: 'sidecarProvider'
    orgId: 1
    folder: ''
    type: file
    disableDeletion: false
    allowUiUpdates: false
    updateIntervalSeconds: 30
    options:
      foldersFromFilesStructure: false
      path: /tmp/dashboards
`,
	},
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.5.1",
			"helm.sh/chart":                "grafana-6.56.2",
		},
		Name:      "kube-prometheus-stack-grafana-config-dashboards",
		Namespace: namespace,
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ConfigMap",
	},
}

var GrafanaDatasourceCM = &corev1.ConfigMap{
	Data: map[string]string{
		"datasource.yaml": `
apiVersion: 1
datasources:
- name: Prometheus
  type: prometheus
  uid: prometheus
  url: http://kube-prometheus-stack-prometheus.monitoring:9090/
  access: proxy
  isDefault: true
  jsonData:
    httpMethod: POST
    timeInterval: 30s
`,
	},
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-grafana",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "45.27.2",
			"chart":                        "kube-prometheus-stack-45.27.2",
			"grafana_datasource":           "1",
			"heritage":                     "Helm",
			"release":                      "kube-prometheus-stack",
		},
		Name:      "kube-prometheus-stack-grafana-datasource",
		Namespace: namespace,
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ConfigMap",
	},
}

var GrafanaTestCM = &corev1.ConfigMap{
	Data: map[string]string{
		"run.sh": `
@test "Test Health" {
  url="http://kube-prometheus-stack-grafana/api/health"
  code=$(wget --server-response --spider --timeout 90 --tries 10 ${url} 2>&1 | awk '/^  HTTP/{print $2}')
  [ "$code" == "200" ]
}
`,
	},
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "test-success",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
		},
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.5.1",
			"helm.sh/chart":                "grafana-6.56.2",
		},
		Name:      "kube-prometheus-stack-grafana-test",
		Namespace: namespace,
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ConfigMap",
	},
}
