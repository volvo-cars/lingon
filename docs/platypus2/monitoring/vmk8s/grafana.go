// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmk8s

import (
	"github.com/VictoriaMetrics/operator/api/victoriametrics/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Grafana struct {
	CM                 *corev1.ConfigMap
	CR                 *rbacv1.ClusterRole
	CRB                *rbacv1.ClusterRoleBinding
	ConfigDashboardsCM *corev1.ConfigMap
	Deploy             *appsv1.Deployment
	RB                 *rbacv1.RoleBinding
	Role               *rbacv1.Role
	SA                 *corev1.ServiceAccount
	SVC                *corev1.Service
	Secrets            *corev1.Secret

	DataSourceCM        *corev1.ConfigMap
	OverviewDashboardCM *corev1.ConfigMap
	GrafanaScrape       *v1beta1.VMServiceScrape
}

type GrafanaTest struct {
	GrafanaTestCM   *corev1.ConfigMap
	GrafanaTestPO   *corev1.Pod
	GrafanaTestRB   *rbacv1.RoleBinding
	GrafanaTestRole *rbacv1.Role
	GrafanaTestSA   *corev1.ServiceAccount
}

func NewGrafanaTest() *GrafanaTest {
	return &GrafanaTest{
		GrafanaTestCM:   GrafanaTestCM,
		GrafanaTestPO:   GrafanaTestPO,
		GrafanaTestRB:   GrafanaTestRB,
		GrafanaTestRole: GrafanaTestRole,
		GrafanaTestSA:   GrafanaTestSA,
	}
}

func NewGrafana() *Grafana {
	return &Grafana{
		CM:                  GrafanaCM,
		CR:                  GrafanaCR,
		CRB:                 GrafanaCRB,
		ConfigDashboardsCM:  GrafanaProviderCM,
		Deploy:              GrafanaDeploy,
		RB:                  GrafanaRB,
		Role:                GrafanaRole,
		SA:                  GrafanaSA,
		SVC:                 GrafanaSVC,
		Secrets:             GrafanaSecrets,
		DataSourceCM:        GrafanaDataSourceCM,
		OverviewDashboardCM: GrafanaOverviewDashCM,
		GrafanaScrape:       GrafanaScrape,
	}
}

var GrafanaDeploy = &appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name:      "vmk8s-grafana",
		Namespace: "monitoring",
	},
	Spec: appsv1.DeploymentSpec{
		Replicas:             P(int32(1)),
		RevisionHistoryLimit: P(int32(10)),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app.kubernetes.io/instance": "vmk8s",
				"app.kubernetes.io/name":     "grafana",
			},
		},
		Strategy: appsv1.DeploymentStrategy{Type: appsv1.DeploymentStrategyType("RollingUpdate")},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					"checksum/config":                       "2ececdbdce31ca013f19ea0000c079aba2a5b9fb01ec29c10c1270f073852744",
					"checksum/dashboards-json-config":       "967e7b73376b048fc74810c99a9fda0454f0064a6b81d916c84f4d827ad8e8c7",
					"checksum/sc-dashboard-provider-config": "fa4ef62d42c42e06a0fe021f7c9faceecdb28b48d2bc59fc6f61c46c79936f86",
					"checksum/secret":                       "156a06b1b8ff51c3069373c417218597a672708542ce79e8513fc9aa6bade9f2",
				},
				Labels: map[string]string{
					"app.kubernetes.io/instance": "vmk8s",
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
								Name:  "FOLDER",
								Value: "/tmp/dashboards",
							}, {
								Name:  "RESOURCE",
								Value: "both",
							},
						},
						Image:           "quay.io/kiwigrid/k8s-sidecar:1.19.2",
						ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
						Name:            "grafana-sc-dashboard",
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
										LocalObjectReference: corev1.LocalObjectReference{Name: "vmk8s-grafana"},
									},
								},
							}, {
								Name: "REQ_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										Key:                  "admin-password",
										LocalObjectReference: corev1.LocalObjectReference{Name: "vmk8s-grafana"},
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
						Image:           "quay.io/kiwigrid/k8s-sidecar:1.19.2",
						ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
						Name:            "grafana-sc-datasources",
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/etc/grafana/provisioning/datasources",
								Name:      "sc-datasources-volume",
							},
						},
					}, {
						Env: []corev1.EnvVar{
							{
								Name: "GF_SECURITY_ADMIN_USER",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										Key:                  "admin-user",
										LocalObjectReference: corev1.LocalObjectReference{Name: "vmk8s-grafana"},
									},
								},
							}, {
								Name: "GF_SECURITY_ADMIN_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										Key:                  "admin-password",
										LocalObjectReference: corev1.LocalObjectReference{Name: "vmk8s-grafana"},
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
						Image:           "grafana/grafana:9.3.0",
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
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/etc/grafana/grafana.ini",
								Name:      "config",
								SubPath:   "grafana.ini",
							}, {
								MountPath: "/var/lib/grafana",
								Name:      "storage",
							}, {
								MountPath: "/etc/grafana/provisioning/dashboards/dashboardproviders.yaml",
								Name:      "config",
								SubPath:   "dashboardproviders.yaml",
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
				InitContainers: []corev1.Container{
					{
						Args: []string{
							"-c",
							"mkdir -p /var/lib/grafana/dashboards/default && /bin/sh -x /etc/grafana/download_dashboards.sh",
						},
						Command:         []string{"/bin/sh"},
						Image:           "curlimages/curl:7.85.0",
						ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
						Name:            "download-dashboards",
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/etc/grafana/download_dashboards.sh",
								Name:      "config",
								SubPath:   "download_dashboards.sh",
							}, {
								MountPath: "/var/lib/grafana",
								Name:      "storage",
							},
						},
					},
				},
				SecurityContext: &corev1.PodSecurityContext{
					FSGroup:    P(int64(472)),
					RunAsGroup: P(int64(472)),
					RunAsUser:  P(int64(472)),
				},
				ServiceAccountName: "vmk8s-grafana",
				Volumes: []corev1.Volume{
					{
						Name:         "config",
						VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: "vmk8s-grafana"}}},
					}, {
						Name:         "dashboards-default",
						VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: "vmk8s-grafana-dashboards-default"}}},
					}, {
						Name:         "storage",
						VolumeSource: corev1.VolumeSource{},
					}, {
						Name:         "sc-dashboard-volume",
						VolumeSource: corev1.VolumeSource{},
					}, {
						Name:         "sc-dashboard-provider",
						VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: "vmk8s-grafana-config-dashboards"}}},
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

var GrafanaCM = &corev1.ConfigMap{
	Data: map[string]string{
		"dashboardproviders.yaml": `
apiVersion: 1
providers:
- disableDeletion: false
  editable: true
  folder: ""
  name: default
  options:
    path: /var/lib/grafana/dashboards/default
  orgId: 1
  type: file

`,
		"download_dashboards.sh": `
#!/usr/bin/env sh
set -euf
mkdir -p /var/lib/grafana/dashboards/default
curl -skf \
--connect-timeout 60 \
--max-time 60 \
-H "Accept: application/json" \
-H "Content-Type: application/json;charset=UTF-8" \
  "https://grafana.com/api/dashboards/1860/revisions/22/download" \
  | sed '/-- .* --/! s/"datasource":.*,/"datasource": "VictoriaMetrics",/g' \
> "/var/lib/grafana/dashboards/default/nodeexporter.json"

`,
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
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name:      "vmk8s-grafana",
		Namespace: "monitoring",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ConfigMap",
	},
}

var GrafanaProviderCM = &corev1.ConfigMap{
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
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name:      "vmk8s-grafana-config-dashboards",
		Namespace: "monitoring",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ConfigMap",
	},
}

var GrafanaDataSourceCM = &corev1.ConfigMap{
	Data: map[string]string{
		"datasource.yaml": `
apiVersion: 1
datasources:
- name: VictoriaMetrics
  type: prometheus
  url: http://vmsingle-vmk8s-victoria-metrics-k8s-stack.monitoring.svc:8429/
  access: proxy
  isDefault: true
  jsonData: 
    {}
`,
	},
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "victoria-metrics-k8s-stack-grafana",
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"grafana_datasource":           "1",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-grafana-ds",
		Namespace: "monitoring",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ConfigMap",
	},
}

var GrafanaSecrets = &corev1.Secret{
	Data: map[string][]byte{
		"admin-password": []byte("HT56XNIyTRJcajA5dPY8K2atkoyFHOsbq4l60oTH"),
		"admin-user":     []byte("admin"),
		"ldap-toml":      []byte(""),
	},
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name:      "vmk8s-grafana",
		Namespace: "monitoring",
	},
	Type: corev1.SecretType("Opaque"),
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Secret",
	},
} // TODO: SECRETS SHOULD BE STORED ELSEWHERE THAN IN THE CODE!!!!

var GrafanaSVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name:      "vmk8s-grafana",
		Namespace: "monitoring",
	},
	Spec: corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			{
				Name:       "service",
				Port:       int32(80),
				Protocol:   corev1.Protocol("TCP"),
				TargetPort: intstr.IntOrString{IntVal: int32(3000)},
			},
		},
		Selector: map[string]string{
			"app.kubernetes.io/instance": "vmk8s",
			"app.kubernetes.io/name":     "grafana",
		},
		Type: corev1.ServiceType("ClusterIP"),
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	},
}

var GrafanaCR = &rbacv1.ClusterRole{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name: "vmk8s-grafana-clusterrole",
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

var GrafanaCRB = &rbacv1.ClusterRoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name: "vmk8s-grafana-clusterrolebinding",
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "ClusterRole",
		Name:     "vmk8s-grafana-clusterrole",
	},
	Subjects: []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      "vmk8s-grafana",
			Namespace: "monitoring",
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "ClusterRoleBinding",
	},
}

var GrafanaSA = &corev1.ServiceAccount{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name:      "vmk8s-grafana",
		Namespace: "monitoring",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ServiceAccount",
	},
}

var GrafanaRole = &rbacv1.Role{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name:      "vmk8s-grafana",
		Namespace: "monitoring",
	},
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups:     []string{"extensions"},
			ResourceNames: []string{"vmk8s-grafana"},
			Resources:     []string{"podsecuritypolicies"},
			Verbs:         []string{"use"},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "Role",
	},
}

var GrafanaRB = &rbacv1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name:      "vmk8s-grafana",
		Namespace: "monitoring",
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "Role",
		Name:     "vmk8s-grafana",
	},
	Subjects: []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      "vmk8s-grafana",
			Namespace: "monitoring",
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "RoleBinding",
	},
}

var GrafanaTestCM = &corev1.ConfigMap{
	Data: map[string]string{
		"run.sh": `
@test "Test Health" {
  url="http://vmk8s-grafana/api/health"
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
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name:      "vmk8s-grafana-test",
		Namespace: "monitoring",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ConfigMap",
	},
}

var GrafanaTestRole = &rbacv1.Role{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "test-success",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
		},
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name:      "vmk8s-grafana-test",
		Namespace: "monitoring",
	},
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups:     []string{"policy"},
			ResourceNames: []string{"vmk8s-grafana-test"},
			Resources:     []string{"podsecuritypolicies"},
			Verbs:         []string{"use"},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "Role",
	},
}

var GrafanaTestRB = &rbacv1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "test-success",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
		},
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name:      "vmk8s-grafana-test",
		Namespace: "monitoring",
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "Role",
		Name:     "vmk8s-grafana-test",
	},
	Subjects: []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      "vmk8s-grafana-test",
			Namespace: "monitoring",
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "RoleBinding",
	},
}

var GrafanaTestSA = &corev1.ServiceAccount{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "test-success",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
		},
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name:      "vmk8s-grafana-test",
		Namespace: "monitoring",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ServiceAccount",
	},
}

var GrafanaTestPO = &corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "test-success",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
		},
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name:      "vmk8s-grafana-test",
		Namespace: "monitoring",
	},
	Spec: corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Command: []string{
					"/opt/bats/bin/bats",
					"-t",
					"/tests/run.sh",
				},
				Image:           "bats/bats:v1.4.1",
				ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
				Name:            "vmk8s-test",
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
		ServiceAccountName: "vmk8s-grafana-test",
		Volumes: []corev1.Volume{
			{
				Name:         "tests",
				VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: "vmk8s-grafana-test"}}},
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Pod",
	},
}

var GrafanaScrape = &v1beta1.VMServiceScrape{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-grafana",
		Namespace: "monitoring",
	},
	Spec: v1beta1.VMServiceScrapeSpec{
		Endpoints: []v1beta1.Endpoint{{Port: "service"}},
		Selector: metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app.kubernetes.io/instance": "vmk8s",
				"app.kubernetes.io/name":     "grafana",
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "operator.victoriametrics.com/v1beta1",
		Kind:       "VMServiceScrape",
	},
}

var GrafanaOverviewDashCM = &corev1.ConfigMap{
	Data: map[string]string{
		"grafana-overview.json": `
{
    "annotations": {
        "list": [
            {
                "builtIn": 1,
                "datasource": "-- Grafana --",
                "enable": true,
                "hide": true,
                "iconColor": "rgba(0, 211, 255, 1)",
                "name": "Annotations & Alerts",
                "target": {
                    "limit": 100,
                    "matchAny": false,
                    "tags": [
                    ],
                    "type": "dashboard"
                },
                "type": "dashboard"
            }
        ]
    },
    "editable": true,
    "gnetId": null,
    "graphTooltip": 0,
    "id": 3085,
    "iteration": 1631554945276,
    "links": [
    ],
    "panels": [
        {
            "datasource": "$datasource",
            "fieldConfig": {
                "defaults": {
                    "mappings": [
                    ],
                    "noValue": "0",
                    "thresholds": {
                        "mode": "absolute",
                        "steps": [
                            {
                                "color": "green",
                                "value": null
                            },
                            {
                                "color": "red",
                                "value": 80
                            }
                        ]
                    }
                },
                "overrides": [
                ]
            },
            "gridPos": {
                "h": 5,
                "w": 6,
                "x": 0,
                "y": 0
            },
            "id": 6,
            "options": {
                "colorMode": "value",
                "graphMode": "area",
                "justifyMode": "auto",
                "orientation": "auto",
                "reduceOptions": {
                    "calcs": [
                        "mean"
                    ],
                    "fields": "",
                    "values": false
                },
                "text": {
                },
                "textMode": "auto"
            },
            "pluginVersion": "8.1.3",
            "targets": [
                {
                    "expr": "grafana_alerting_result_total{job=~\"$job\", instance=~\"$instance\", state=\"alerting\"}",
                    "instant": true,
                    "interval": "",
                    "legendFormat": "",
                    "refId": "A"
                }
            ],
            "timeFrom": null,
            "timeShift": null,
            "title": "Firing Alerts",
            "type": "stat"
        },
        {
            "datasource": "$datasource",
            "fieldConfig": {
                "defaults": {
                    "mappings": [
                    ],
                    "thresholds": {
                        "mode": "absolute",
                        "steps": [
                            {
                                "color": "green",
                                "value": null
                            },
                            {
                                "color": "red",
                                "value": 80
                            }
                        ]
                    }
                },
                "overrides": [
                ]
            },
            "gridPos": {
                "h": 5,
                "w": 6,
                "x": 6,
                "y": 0
            },
            "id": 8,
            "options": {
                "colorMode": "value",
                "graphMode": "area",
                "justifyMode": "auto",
                "orientation": "auto",
                "reduceOptions": {
                    "calcs": [
                        "mean"
                    ],
                    "fields": "",
                    "values": false
                },
                "text": {
                },
                "textMode": "auto"
            },
            "pluginVersion": "8.1.3",
            "targets": [
                {
                    "expr": "sum(grafana_stat_totals_dashboard{job=~\"$job\", instance=~\"$instance\"})",
                    "interval": "",
                    "legendFormat": "",
                    "refId": "A"
                }
            ],
            "timeFrom": null,
            "timeShift": null,
            "title": "Dashboards",
            "type": "stat"
        },
        {
            "datasource": "$datasource",
            "fieldConfig": {
                "defaults": {
                    "custom": {
                        "align": null,
                        "displayMode": "auto"
                    },
                    "mappings": [
                    ],
                    "thresholds": {
                        "mode": "absolute",
                        "steps": [
                            {
                                "color": "green",
                                "value": null
                            },
                            {
                                "color": "red",
                                "value": 80
                            }
                        ]
                    }
                },
                "overrides": [
                ]
            },
            "gridPos": {
                "h": 5,
                "w": 12,
                "x": 12,
                "y": 0
            },
            "id": 10,
            "options": {
                "showHeader": true
            },
            "pluginVersion": "8.1.3",
            "targets": [
                {
                    "expr": "grafana_build_info{job=~\"$job\", instance=~\"$instance\"}",
                    "instant": true,
                    "interval": "",
                    "legendFormat": "",
                    "refId": "A"
                }
            ],
            "timeFrom": null,
            "timeShift": null,
            "title": "Build Info",
            "transformations": [
                {
                    "id": "labelsToFields",
                    "options": {
                    }
                },
                {
                    "id": "organize",
                    "options": {
                        "excludeByName": {
                            "Time": true,
                            "Value": true,
                            "branch": true,
                            "container": true,
                            "goversion": true,
                            "namespace": true,
                            "pod": true,
                            "revision": true
                        },
                        "indexByName": {
                            "Time": 7,
                            "Value": 11,
                            "branch": 4,
                            "container": 8,
                            "edition": 2,
                            "goversion": 6,
                            "instance": 1,
                            "job": 0,
                            "namespace": 9,
                            "pod": 10,
                            "revision": 5,
                            "version": 3
                        },
                        "renameByName": {
                        }
                    }
                }
            ],
            "type": "table"
        },
        {
            "aliasColors": {
            },
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": "$datasource",
            "fieldConfig": {
                "defaults": {
                    "links": [
                    ]
                },
                "overrides": [
                ]
            },
            "fill": 1,
            "fillGradient": 0,
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 0,
                "y": 5
            },
            "hiddenSeries": false,
            "id": 2,
            "legend": {
                "avg": false,
                "current": false,
                "max": false,
                "min": false,
                "show": true,
                "total": false,
                "values": false
            },
            "lines": true,
            "linewidth": 1,
            "nullPointMode": "null",
            "options": {
                "alertThreshold": true
            },
            "percentage": false,
            "pluginVersion": "8.1.3",
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [
            ],
            "spaceLength": 10,
            "stack": true,
            "steppedLine": false,
            "targets": [
                {
                    "expr": "sum by (status_code) (irate(grafana_http_request_duration_seconds_count{job=~\"$job\", instance=~\"$instance\"}[1m])) ",
                    "interval": "",
                    "legendFormat": "{{status_code}}",
                    "refId": "A"
                }
            ],
            "thresholds": [
            ],
            "timeFrom": null,
            "timeRegions": [
            ],
            "timeShift": null,
            "title": "RPS",
            "tooltip": {
                "shared": true,
                "sort": 0,
                "value_type": "individual"
            },
            "type": "graph",
            "xaxis": {
                "buckets": null,
                "mode": "time",
                "name": null,
                "show": true,
                "values": [
                ]
            },
            "yaxes": [
                {
                    "$$hashKey": "object:157",
                    "format": "reqps",
                    "label": null,
                    "logBase": 1,
                    "max": null,
                    "min": null,
                    "show": true
                },
                {
                    "$$hashKey": "object:158",
                    "format": "short",
                    "label": null,
                    "logBase": 1,
                    "max": null,
                    "min": null,
                    "show": false
                }
            ],
            "yaxis": {
                "align": false,
                "alignLevel": null
            }
        },
        {
            "aliasColors": {
            },
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": "$datasource",
            "fieldConfig": {
                "defaults": {
                    "links": [
                    ]
                },
                "overrides": [
                ]
            },
            "fill": 1,
            "fillGradient": 0,
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 12,
                "y": 5
            },
            "hiddenSeries": false,
            "id": 4,
            "legend": {
                "avg": false,
                "current": false,
                "max": false,
                "min": false,
                "show": true,
                "total": false,
                "values": false
            },
            "lines": true,
            "linewidth": 1,
            "nullPointMode": "null",
            "options": {
                "alertThreshold": true
            },
            "percentage": false,
            "pluginVersion": "8.1.3",
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [
            ],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "exemplar": true,
                    "expr": "histogram_quantile(0.99, sum(irate(grafana_http_request_duration_seconds_bucket{instance=~\"$instance\", job=~\"$job\"}[$__rate_interval])) by (le)) * 1",
                    "interval": "",
                    "legendFormat": "99th Percentile",
                    "refId": "A"
                },
                {
                    "exemplar": true,
                    "expr": "histogram_quantile(0.50, sum(irate(grafana_http_request_duration_seconds_bucket{instance=~\"$instance\", job=~\"$job\"}[$__rate_interval])) by (le)) * 1",
                    "interval": "",
                    "legendFormat": "50th Percentile",
                    "refId": "B"
                },
                {
                    "exemplar": true,
                    "expr": "sum(irate(grafana_http_request_duration_seconds_sum{instance=~\"$instance\", job=~\"$job\"}[$__rate_interval])) * 1 / sum(irate(grafana_http_request_duration_seconds_count{instance=~\"$instance\", job=~\"$job\"}[$__rate_interval]))",
                    "interval": "",
                    "legendFormat": "Average",
                    "refId": "C"
                }
            ],
            "thresholds": [
            ],
            "timeFrom": null,
            "timeRegions": [
            ],
            "timeShift": null,
            "title": "Request Latency",
            "tooltip": {
                "shared": true,
                "sort": 0,
                "value_type": "individual"
            },
            "type": "graph",
            "xaxis": {
                "buckets": null,
                "mode": "time",
                "name": null,
                "show": true,
                "values": [
                ]
            },
            "yaxes": [
                {
                    "$$hashKey": "object:210",
                    "format": "ms",
                    "label": null,
                    "logBase": 1,
                    "max": null,
                    "min": null,
                    "show": true
                },
                {
                    "$$hashKey": "object:211",
                    "format": "short",
                    "label": null,
                    "logBase": 1,
                    "max": null,
                    "min": null,
                    "show": true
                }
            ],
            "yaxis": {
                "align": false,
                "alignLevel": null
            }
        }
    ],
    "schemaVersion": 30,
    "style": "dark",
    "tags": [
    ],
    "templating": {
        "list": [
            {
                "current": {
                    "selected": true,
                    "text": "dev-cortex",
                    "value": "dev-cortex"
                },
                "description": null,
                "error": null,
                "hide": 0,
                "includeAll": false,
                "label": null,
                "multi": false,
                "name": "datasource",
                "options": [
                ],
                "query": "prometheus",
                "queryValue": "",
                "refresh": 1,
                "regex": "",
                "skipUrlSync": false,
                "type": "datasource"
            },
            {
                "allValue": ".*",
                "current": {
                    "selected": false,
                    "text": [
                        "default/grafana"
                    ],
                    "value": [
                        "default/grafana"
                    ]
                },
                "datasource": "$datasource",
                "definition": "label_values(grafana_build_info, job)",
                "description": null,
                "error": null,
                "hide": 0,
                "includeAll": true,
                "label": null,
                "multi": true,
                "name": "job",
                "options": [
                ],
                "query": {
                    "query": "label_values(grafana_build_info, job)",
                    "refId": "Billing Admin-job-Variable-Query"
                },
                "refresh": 1,
                "regex": "",
                "skipUrlSync": false,
                "sort": 0,
                "tagValuesQuery": "",
                "tagsQuery": "",
                "type": "query",
                "useTags": false
            },
            {
                "allValue": ".*",
                "current": {
                    "selected": false,
                    "text": "All",
                    "value": "$__all"
                },
                "datasource": "$datasource",
                "definition": "label_values(grafana_build_info, instance)",
                "description": null,
                "error": null,
                "hide": 0,
                "includeAll": true,
                "label": null,
                "multi": true,
                "name": "instance",
                "options": [
                ],
                "query": {
                    "query": "label_values(grafana_build_info, instance)",
                    "refId": "Billing Admin-instance-Variable-Query"
                },
                "refresh": 1,
                "regex": "",
                "skipUrlSync": false,
                "sort": 0,
                "tagValuesQuery": "",
                "tagsQuery": "",
                "type": "query",
                "useTags": false
            }
        ]
    },
    "time": {
        "from": "now-6h",
        "to": "now"
    },
    "timepicker": {
        "refresh_intervals": [
            "10s",
            "30s",
            "1m",
            "5m",
            "15m",
            "30m",
            "1h",
            "2h",
            "1d"
        ]
    },
    "timezone": "utc",
    "title": "Grafana Overview",
    "uid": "6be0s85Mk",
    "version": 2
}
`,
	},
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "victoria-metrics-k8s-stack-grafana",
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"grafana_dashboard":            "1",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-grafana-overview",
		Namespace: "monitoring",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ConfigMap",
	},
}
