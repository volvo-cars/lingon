// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmk8s

import (
	_ "embed"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/VictoriaMetrics/operator/api/victoriametrics/v1beta1"
	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	"github.com/volvo-cars/lingoneks/meta"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	GrafanaVersion                   = "9.5.3"
	GrafanaSideCarImg                = "quay.io/kiwigrid/k8s-sidecar:1.19.2"
	GrafanaPort                      = 3000
	GrafanaPortName                  = "service"
	DashboardLabel                   = "grafana_dashboard"
	DataSourceLabel                  = "grafana_datasource"
	defaultDashboardConfigName       = "grafana-default-dashboards"
	curlImg                          = "curlimages/curl:7.85.0"
	userID                     int64 = 472
)

func PatchDashLabels(o metav1.ObjectMeta) metav1.ObjectMeta {
	o.Labels = ku.MergeLabels(o.Labels, map[string]string{DashboardLabel: "1"})
	return o
}

var Graf = &meta.Metadata{
	Name:      "grafana",
	Namespace: namespace,
	Instance:  "grafana" + namespace,
	Component: "dashboards",
	PartOf:    appName,
	Version:   GrafanaVersion,
	ManagedBy: "lingon",
	Img: meta.ContainerImg{
		Registry: "docker.io",
		Image:    "grafana/grafana",
		Tag:      GrafanaVersion,
	},
}

type Grafana struct {
	kube.App

	Deploy  *appsv1.Deployment
	SVC     *corev1.Service
	Secrets *corev1.Secret

	SA   *corev1.ServiceAccount
	CR   *rbacv1.ClusterRole
	CRB  *rbacv1.ClusterRoleBinding
	Role *rbacv1.Role
	RB   *rbacv1.RoleBinding

	CM                  *corev1.ConfigMap
	ProviderCM          *corev1.ConfigMap
	DataSourceCM        *corev1.ConfigMap
	OverviewDashboardCM *corev1.ConfigMap
	DefaultDashboardCM  *corev1.ConfigMap
	GrafanaScrape       *v1beta1.VMServiceScrape
}

func NewGrafana() *Grafana {
	return &Grafana{
		Deploy:        GrafanaDeploy,
		SVC:           GrafanaSVC,
		Secrets:       GrafanaSecrets,
		GrafanaScrape: GrafanaScrape,

		SA:   GrafanaSA,
		Role: GrafanaRole,
		RB:   ku.BindRole(Graf.Name, GrafanaSA, GrafanaRole, Graf.Labels()),
		CR:   GrafanaCR,
		CRB: ku.BindClusterRole(
			Graf.Name, GrafanaSA, GrafanaCR, Graf.Labels(),
		),

		CM:                  GrafanaCM,
		ProviderCM:          GrafanaProviderCM,
		DataSourceCM:        GrafanaDataSourceCM,
		OverviewDashboardCM: GrafanaOverviewDashCM,
		DefaultDashboardCM: ku.DataConfigMap(
			defaultDashboardConfigName,
			Graf.Namespace, Graf.Labels(), nil, map[string]string{},
		),
	}
}

var GrafanaSA = ku.ServiceAccount(Graf.Name, Graf.Namespace, Graf.Labels(), nil)

var GrafanaCR = &rbacv1.ClusterRole{
	TypeMeta:   ku.TypeClusterRoleV1,
	ObjectMeta: Graf.ObjectMetaNameSuffixNoNS("cr"),
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{""},
			Resources: []string{"configmaps", "secrets"},
			Verbs:     []string{"get", "watch", "list"},
		},
	},
}

var GrafanaRole = &rbacv1.Role{
	TypeMeta:   ku.TypeRoleV1,
	ObjectMeta: Graf.ObjectMeta(),
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups:     []string{"extensions"},
			ResourceNames: []string{Graf.Name},
			Resources:     []string{"podsecuritypolicies"},
			Verbs:         []string{"use"},
		},
	},
}

var GrafanaDeploy = &appsv1.Deployment{
	TypeMeta:   ku.TypeDeploymentV1,
	ObjectMeta: Graf.ObjectMeta(),
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(1)),
		Selector: &metav1.LabelSelector{MatchLabels: Graf.MatchLabels()},
		Strategy: appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					"checksum/config":                         ku.HashConfig(GrafanaCM),
					"checksum/provider":                       ku.HashConfig(GrafanaProviderCM),
					"checksum/datasource":                     ku.HashConfig(GrafanaDataSourceCM),
					"checksum/secret":                         ku.HashSecret(GrafanaSecrets),
					"kubectl.kubernetes.io/default-container": Graf.Name,
				},
				Labels: Graf.MatchLabels(),
			},
			Spec: corev1.PodSpec{
				AutomountServiceAccountToken: P(true),
				Containers: []corev1.Container{
					{
						Name:            "grafana-sc-dashboard",
						Image:           GrafanaSideCarImg,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Env: []corev1.EnvVar{
							{Name: "METHOD", Value: "WATCH"},
							{Name: "LABEL", Value: DashboardLabel},
							{Name: "FOLDER", Value: "/tmp/dashboards"},
							{Name: "RESOURCE", Value: "both"},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/tmp/dashboards",
								Name:      "sc-dashboard-volume",
							},
						},
					},
					{
						Name:            "grafana-sc-datasources",
						Image:           GrafanaSideCarImg,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Env: []corev1.EnvVar{
							{
								Name:  "METHOD",
								Value: "WATCH",
							},
							{
								Name:  "LABEL",
								Value: DataSourceLabel,
							},
							{
								Name:  "FOLDER",
								Value: "/etc/grafana/provisioning/datasources",
							},
							{
								Name:  "RESOURCE",
								Value: "both",
							},
							{
								Name: "REQ_URL",
								Value: fmt.Sprintf(
									"http://localhost:%d/api/admin/provisioning/datasources/reload",
									GrafanaPort,
								),
							},
							{
								Name:  "REQ_METHOD",
								Value: "POST",
							},
							ku.SecretEnvVar(
								"REQ_USERNAME",
								"admin-user",
								GrafanaSecrets.Name,
							),
							ku.SecretEnvVar(
								"REQ_PASSWORD",
								"admin-password",
								GrafanaSecrets.Name,
							),
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/etc/grafana/provisioning/datasources",
								Name:      "sc-datasources-volume",
							},
						},
					},
					{
						Name:            Graf.Name,
						Image:           Graf.Img.URL(),
						ImagePullPolicy: corev1.PullIfNotPresent,
						Env: []corev1.EnvVar{
							ku.SecretEnvVar(
								"GF_SECURITY_ADMIN_USER",
								"admin-user",
								GrafanaSecrets.Name,
							),
							ku.SecretEnvVar(
								"GF_SECURITY_ADMIN_PASSWORD",
								"admin-password",
								GrafanaSecrets.Name,
							),
							{
								Name: "GF_INSTALL_PLUGINS",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										Key:                  "plugins",
										LocalObjectReference: corev1.LocalObjectReference{Name: Graf.Name},
									},
								},
							},
							{
								Name:  "GF_PATHS_DATA",
								Value: "/var/lib/grafana/",
							},
							{
								Name:  "GF_PATHS_LOGS",
								Value: "/var/log/grafana",
							},
							{
								Name:  "GF_PATHS_PLUGINS",
								Value: "/var/lib/grafana/plugins",
							},
							{
								Name:  "GF_PATHS_PROVISIONING",
								Value: "/etc/grafana/provisioning",
							},
						},

						LivenessProbe: &corev1.Probe{
							FailureThreshold:    int32(10),
							InitialDelaySeconds: int32(60),
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/api/health",
									Port: intstr.FromInt(GrafanaPort),
								},
							},
							TimeoutSeconds: int32(30),
						},
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: int32(GrafanaPort),
								Name:          Graf.Name,
								Protocol:      corev1.ProtocolTCP,
							},
						},
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/api/health",
									Port: intstr.FromInt(GrafanaPort),
								},
							},
						},
						Resources: ku.Resources(
							"500m", "128Mi", "500m", "128Mi",
						),
						SecurityContext: &corev1.SecurityContext{
							Capabilities: &corev1.Capabilities{
								Drop: []corev1.Capability{corev1.Capability("ALL")},
							},
							SeccompProfile: &corev1.SeccompProfile{
								Type: corev1.SeccompProfileType("RuntimeDefault"),
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
						Name:            "download-dashboards",
						Image:           curlImg,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Command:         []string{"/bin/sh"},
						Args: []string{
							"-c",
							"mkdir -p /var/lib/grafana/dashboards/default && " +
								// If it is not created here,
								// it will be assigned root ownership ¯\_(ツ)_/¯.
								"mkdir -p /var/lib/grafana/plugins && " +
								"/bin/sh -x /etc/grafana/download_dashboards.sh",
						},
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
					{
						Name:            "load-vm-ds-plugin",
						Image:           curlImg,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Command:         []string{"/bin/sh"},
						Args: []string{
							"-c", "mkdir -p /var/lib/grafana/plugins && " +
								"/bin/sh -x /etc/grafana/download_vm_ds.sh",
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/var/lib/grafana",
								Name:      "storage",
							}, {
								MountPath: "/etc/grafana/download_vm_ds.sh",
								Name:      "config",
								SubPath:   "download_vm_ds.sh",
							},
						},
						WorkingDir: "/var/lib/grafana/plugins",
					},
				},
				SecurityContext: &corev1.PodSecurityContext{
					FSGroup:      P(userID),
					RunAsGroup:   P(userID),
					RunAsNonRoot: P(true),
					RunAsUser:    P(userID),
				},
				ServiceAccountName: GrafanaSA.Name,
				Volumes: []corev1.Volume{
					{
						Name:         "config",
						VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: GrafanaCM.Name}}},
					}, {
						Name:         "dashboards-default",
						VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: defaultDashboardConfigName}}},
					}, {
						Name:         "storage",
						VolumeSource: corev1.VolumeSource{},
					}, {
						Name:         "sc-dashboard-volume",
						VolumeSource: corev1.VolumeSource{},
					}, {
						Name:         "sc-dashboard-provider",
						VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: GrafanaProviderCM.Name}}},
					}, {
						Name:         "sc-datasources-volume",
						VolumeSource: corev1.VolumeSource{},
					},
				},
			},
		},
	},
}

var GrafanaSVC = &corev1.Service{
	TypeMeta:   ku.TypeServiceV1,
	ObjectMeta: Graf.ObjectMeta(),
	Spec: corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			{
				Name:       GrafanaPortName,
				Port:       80,
				Protocol:   corev1.ProtocolTCP,
				TargetPort: intstr.FromInt(GrafanaPort),
			},
		},
		Selector: Graf.MatchLabels(),
		Type:     corev1.ServiceTypeClusterIP,
	},
}

var GrafanaScrape = &v1beta1.VMServiceScrape{
	ObjectMeta: Graf.ObjectMeta(),
	Spec: v1beta1.VMServiceScrapeSpec{
		Endpoints: []v1beta1.Endpoint{{Port: GrafanaPortName}},
		Selector: metav1.LabelSelector{
			MatchLabels: Graf.MatchLabels(),
		},
	},
	TypeMeta: TypeVMServiceScrapeV1Beta1,
}

type DashSource struct {
	Name   string
	URL    string
	Source string
}

const (
	PrometheusDataSourceName      = "Prometheus"
	VictoriaMetricsDataSourceName = "VictoriaMetrics"
)

func (d *DashSource) Validate() error {
	if _, err := url.Parse(d.URL); err != nil {
		return fmt.Errorf("url %s - %s: %w", d.Name, d.URL, err)
	}

	if d.Name == "" {
		return fmt.Errorf("dashboard %s: name undefined", d.URL)
	}
	n := d.Name
	n = strings.ReplaceAll(n, " ", "-")
	n = strings.ReplaceAll(n, "/", "_")

	switch d.Source {
	case PrometheusDataSourceName:
	case VictoriaMetricsDataSourceName:
	default:
		return fmt.Errorf("datasource %v: %s", d.Name, d.Source)
	}
	return nil
}

func downloadDashboards(dss []DashSource) string {
	var buf strings.Builder
	var errs error

	buf.WriteString(
		`
#!/usr/bin/env sh
set -euf
mkdir -p /var/lib/grafana/dashboards/default
ls -Rl /var/lib/grafana/
`,
	)

	for _, ds := range dss {
		if err := ds.Validate(); err != nil {
			errs = errors.Join(errs, err)
		}
		buf.WriteString(
			`

curl -skf \
--connect-timeout 60 \
--max-time 60 \
-H "Accept: application/json" \
-H "Content-Type: application/json;charset=UTF-8" \
  "` + ds.URL + `" \
  | sed '/-- .* --/! s/"datasource":.*,/"datasource": "` + ds.Source + `",/g' \
> "/var/lib/grafana/dashboards/default/` + ds.Name + `.json"
`,
		)
	}
	if errs != nil {
		panic(errs)
	}

	return buf.String()
}

const dashboardProvidersYaml = `
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
`

const grafanaINI = `
[plugins]
allow_loading_unsigned_plugins = victoriametrics-datasource
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

`

const downloadVMDSsh = `
set -euxf

ls -R -l /var/lib/grafana/
id
mkdir -p /var/lib/grafana/plugins/
# ver=$(curl -s https://api.github.com/repos/VictoriaMetrics/grafana-datasource/releases/latest | grep -oE 'v\d+\.\d+\.\d+' | head -1)
ver="v0.2.0"
curl -L https://github.com/VictoriaMetrics/grafana-datasource/releases/download/$ver/victoriametrics-datasource-$ver.tar.gz -o /var/lib/grafana/plugins/plugin.tar.gz
tar -xzf /var/lib/grafana/plugins/plugin.tar.gz -C /var/lib/grafana/plugins/
rm -f /var/lib/grafana/plugins/plugin.tar.gz
chown -R 472:472 /var/lib/grafana/plugins/
`

var GrafanaCM = &corev1.ConfigMap{
	ObjectMeta: Graf.ObjectMeta(),
	TypeMeta:   ku.TypeConfigMapV1,
	Data: map[string]string{
		"grafana.ini": grafanaINI,

		"plugins": "https://grafana.com/api/plugins/marcusolsson-json-datasource/versions/1.3.2/download;marcusolsson-json-datasource",

		"dashboardproviders.yaml": dashboardProvidersYaml,

		"download_vm_ds.sh": downloadVMDSsh,

		"download_dashboards.sh": downloadDashboards(
			[]DashSource{
				{
					Name:   "nodeexporter",
					URL:    "https://grafana.com/api/dashboards/1860/revisions/22/download",
					Source: VictoriaMetricsDataSourceName,
				},
				// Karpenter dashboards
				{
					Name:   "karpenter-performance-dashboard",
					URL:    "https://raw.githubusercontent.com/aws/karpenter/main/website/content/en/v0.28/getting-started/getting-started-with-karpenter/karpenter-performance-dashboard.json",
					Source: VictoriaMetricsDataSourceName,
				},
				{
					Name:   "karpenter-controllers",
					URL:    "https://raw.githubusercontent.com/aws/karpenter/main/website/content/en/v0.28/getting-started/getting-started-with-karpenter/karpenter-controllers.json",
					Source: VictoriaMetricsDataSourceName,
				},
				{
					Name:   "karpenter-controllers-allocation",
					URL:    "https://raw.githubusercontent.com/aws/karpenter/main/website/content/en/v0.28/getting-started/getting-started-with-karpenter/karpenter-controllers-allocation.json",
					Source: VictoriaMetricsDataSourceName,
				},
				{
					Name:   "karpenter-capacity-dashboard",
					URL:    "https://raw.githubusercontent.com/aws/karpenter/main/website/content/en/v0.28/getting-started/getting-started-with-karpenter/karpenter-capacity-dashboard.json",
					Source: VictoriaMetricsDataSourceName,
				},
			},
		),
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
	ObjectMeta: Graf.ObjectMetaNameSuffix("-config-dashboards"),
	TypeMeta:   ku.TypeConfigMapV1,
}

var GrafanaDataSourceCM = &corev1.ConfigMap{
	Data: map[string]string{
		// 		"datasources.yaml": `
		// apiVersion: 1
		// datasources:
		// - access: proxy
		//   isDefault: true
		//   name: Prometheus
		//   type: prometheus
		//   url: http://prometheus-prometheus-server
		// - access: proxy
		//   editable: false
		//   jsonData:
		//     authType: default
		//     defaultRegion: us-east-1
		//   name: CloudWatch
		//   type: cloudwatch
		//   uid: cloudwatch
		//
		// `,
		"datasource.yaml": `
apiVersion: 1
datasources:
- name: ` + VictoriaMetricsDataSourceName + `
  type: victoriametrics-datasource
  url: ` + fmt.Sprintf(
			"http://%s.%s.svc:8429/",
			VMDB.PrefixedName(), namespace,
		) + `
  access: proxy
  isDefault: true
`,
	},
	ObjectMeta: metav1.ObjectMeta{
		Labels: ku.MergeLabels(
			Graf.Labels(),
			map[string]string{DataSourceLabel: "1"},
		),
		Name:      Graf.Name + "-ds",
		Namespace: Graf.Namespace,
	},
	TypeMeta: ku.TypeConfigMapV1,
}

var GrafanaSecrets = &corev1.Secret{
	Data: map[string][]byte{
		"admin-password": []byte("HT56XNIyTRJcajA5dPY8K2atkoyFHOsbq4l60oTH"),
		"admin-user":     []byte("admin"),
		"ldap-toml":      []byte(""),
	},
	ObjectMeta: Graf.ObjectMeta(),
	Type:       corev1.SecretTypeOpaque,
	TypeMeta:   ku.TypeSecretV1,
} // TODO: SECRETS SHOULD BE STORED ELSEWHERE THAN IN THE CODE!!!!
