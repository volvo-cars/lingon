// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmk8s

import (
	_ "embed"
	"fmt"

	"github.com/VictoriaMetrics/operator/api/victoriametrics/v1beta1"
	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	"github.com/volvo-cars/lingoneks/meta"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	kyaml "sigs.k8s.io/yaml"
)

const (
	GrafanaVersion                   = "9.5.3"
	GrafanaSideCarImg                = "quay.io/kiwigrid/k8s-sidecar:1.19.2"
	GrafanaPort                      = 3000
	GrafanaPortName                  = "service"
	DashboardLabel                   = "grafana_dashboard"
	DashboardFolderLabel             = "k8s-sidecar-target-directory"
	DataSourceLabel                  = "grafana_datasource"
	defaultDashboardConfigName       = "grafana-default-dashboards"
	curlImg                          = "curlimages/curl:7.85.0"
	userID                     int64 = 472
	pluginsPath                      = "/var/lib/grafana/plugins"
)

func PatchDataSourceLabels(o metav1.ObjectMeta) metav1.ObjectMeta {
	o.Labels = ku.MergeLabels(o.Labels, map[string]string{DataSourceLabel: "1"})
	return o
}

func PatchDashLabels(o metav1.ObjectMeta) metav1.ObjectMeta {
	o.Labels = ku.MergeLabels(o.Labels, map[string]string{DashboardLabel: "1"})
	return o
}

func PatchDashFolder(o metav1.ObjectMeta, folder string) metav1.ObjectMeta {
	o.Annotations = ku.MergeLabels(
		o.Annotations,
		map[string]string{DashboardFolderLabel: folder},
	)
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

	GrafanaScrape *v1beta1.VMServiceScrape

	CM                  *corev1.ConfigMap
	ProviderCM          *corev1.ConfigMap
	DataSourceCM        *corev1.ConfigMap
	OverviewDashboardCM *corev1.ConfigMap
	DefaultDashboardCM  *corev1.ConfigMap
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

		CM:                  GrafCAM.ConfigMap(),
		ProviderCM:          SideCarProvider.ConfigMap(),
		DataSourceCM:        GrafanaDataSource.ConfigMap(),
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

var SideCarDashboard = corev1.Container{
	Name:            "grafana-sc-dashboard",
	Image:           GrafanaSideCarImg,
	ImagePullPolicy: corev1.PullIfNotPresent,
	Env: []corev1.EnvVar{
		{Name: "METHOD", Value: "WATCH"},
		{Name: "LABEL", Value: DashboardLabel},
		{Name: "FOLDER", Value: "/tmp/dashboards"},
		{Name: "RESOURCE", Value: "both"},
	},
	VolumeMounts: []corev1.VolumeMount{VolumeDashboards.VolumeMount},
}

var SideCarDataSource = corev1.Container{
	Name:            "grafana-sc-datasources",
	Image:           GrafanaSideCarImg,
	ImagePullPolicy: corev1.PullIfNotPresent,
	Env: []corev1.EnvVar{
		{Name: "METHOD", Value: "WATCH"},
		{Name: "LABEL", Value: DataSourceLabel},
		{Name: "RESOURCE", Value: "both"},
		{Name: "REQ_METHOD", Value: "POST"},
		{
			Name: "FOLDER",
			// Value: "/etc/grafana/provisioning/datasources",
			Value: VolumeDataSource.VolumeMount.MountPath,
		},
		{
			Name: "REQ_URL",
			Value: fmt.Sprintf(
				"http://localhost:%d/api/admin/provisioning/datasources/reload",
				GrafanaPort,
			),
		},
		ku.SecretEnvVar(
			"REQ_USERNAME", "admin-user", GrafanaSecrets.Name,
		),
		ku.SecretEnvVar(
			"REQ_PASSWORD", "admin-password", GrafanaSecrets.Name,
		),
	},
	VolumeMounts: []corev1.VolumeMount{VolumeDataSource.VolumeMount},
}

var InitContainerPlugins = corev1.Container{
	Name:            "load-vm-ds-plugin",
	Image:           curlImg,
	ImagePullPolicy: corev1.PullIfNotPresent,
	Command:         []string{"/bin/sh"},
	Args: []string{
		"-c",
		"mkdir -p " + pluginsPath + " && " +
			"/bin/sh -x " + VolumeMountScripts.MountPath,
	},
	VolumeMounts: []corev1.VolumeMount{
		VolumeStorage.VolumeMount,
		VolumeMountScripts,
	},
	// BUG: this causes the folder to be owned by root
	// WorkingDir: "/var/lib/grafana/plugins",
}

var GrafanaContainer = corev1.Container{
	Name:            Graf.Name,
	Image:           Graf.Img.URL(),
	ImagePullPolicy: corev1.PullIfNotPresent,
	Env: []corev1.EnvVar{
		ku.SecretEnvVar(
			"GF_SECURITY_ADMIN_USER", "admin-user", GrafanaSecrets.Name,
		),
		ku.SecretEnvVar(
			"GF_SECURITY_ADMIN_PASSWORD", "admin-password", GrafanaSecrets.Name,
		),
		GrafCAM.EnvConfigMapRef("GF_INSTALL_PLUGINS", "plugins"),
		{Name: "GF_PATHS_DATA", Value: "/var/lib/grafana/"},
		{Name: "GF_PATHS_LOGS", Value: "/var/log/grafana"},
		{Name: "GF_PATHS_PLUGINS", Value: pluginsPath},
		{Name: "GF_PATHS_PROVISIONING", Value: "/etc/grafana/provisioning"},
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
		"500m", "256Mi", "500m", "256Mi",
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
		VolumeMountGrafanaIni,
		VolumeStorage.VolumeMount,
		SideCarProvider.VolumeMount,
		VolumeDashboards.VolumeMount,
		VolumeDashProvider,
		VolumeDataSource.VolumeMount,
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
					"checksum/config":                         GrafCAM.Hash(),
					"checksum/provider":                       SideCarProvider.Hash(),
					"checksum/datasource":                     GrafanaDataSource.Hash(),
					"checksum/secret":                         ku.HashSecret(GrafanaSecrets),
					"kubectl.kubernetes.io/default-container": GrafanaContainer.Name,
				},
				Labels: Graf.MatchLabels(),
			},
			Spec: corev1.PodSpec{
				AutomountServiceAccountToken: P(true),
				Containers: []corev1.Container{
					SideCarDashboard,
					SideCarDataSource,
					GrafanaContainer,
				},
				EnableServiceLinks: P(true),
				InitContainers:     []corev1.Container{InitContainerPlugins},
				SecurityContext: &corev1.PodSecurityContext{
					FSGroup:      P(userID),
					RunAsGroup:   P(userID),
					RunAsNonRoot: P(true),
					RunAsUser:    P(userID),
				},
				ServiceAccountName: GrafanaSA.Name,
				Volumes: []corev1.Volume{
					{
						Name: "dashboards-default",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: defaultDashboardConfigName,
								},
							},
						},
					},
					GrafCAM.VolumeAndMount().Volume(),
					SideCarProvider.VolumeAndMount().Volume(),
					VolumeStorage.Volume(),
					VolumeDashboards.Volume(),
					VolumeDataSource.Volume(),
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
		Selector:  metav1.LabelSelector{MatchLabels: Graf.MatchLabels()},
	},
	TypeMeta: TypeVMServiceScrapeV1Beta1,
}

var GrafCAM = ku.ConfigAndMount{
	ObjectMeta:  Graf.ObjectMeta(),
	VolumeMount: corev1.VolumeMount{Name: "config"},
	Data: map[string]string{
		"grafana.ini":     grafanaINI,
		"plugins":         "https://grafana.com/api/plugins/marcusolsson-json-datasource/versions/1.3.6/download;marcusolsson-json-datasource",
		filenameProviders: dashboardProvidersYaml,
		scriptsName:       downloadVMDSsh,
	},
}

const (
	PrometheusDataSourceName      = "Prometheus"
	PrometheusDataSourceID        = "prometheus"
	VictoriaMetricsDataSourceName = "VictoriaMetrics"
	VictoriaMetricsDataSourceID   = "victoriametrics-datasource"
)

const grafanaINI = `
[plugins]
allow_loading_unsigned_plugins = ` + VictoriaMetricsDataSourceID + `
[analytics]
check_for_updates = true
[grafana_net]
url = https://grafana.net
[log]
mode = console
[paths]
data = /var/lib/grafana/
logs = /var/log/grafana
plugins = ` + pluginsPath + `
provisioning = /etc/grafana/provisioning
[server]
domain = ''

`

func d64(i int64) string { return fmt.Sprintf("%d", i) }

var downloadVMDSsh = `
#!/usr/bin/env sh

set -euxf

ls -R -l /var/lib/grafana/
id
mkdir -p ` + pluginsPath + `/
chown -R ` + d64(userID) + `:` + d64(userID) + ` ` + pluginsPath + `/
mkdir -p /var/lib/grafana/dashboards/default
# getting rate-limited by github
# ver=$(curl -s https://api.github.com/repos/VictoriaMetrics/grafana-datasource/releases/latest | grep -oE 'v\d+\.\d+\.\d+' | head -1)
ver="v0.2.0"
curl -L https://github.com/VictoriaMetrics/grafana-datasource/releases/download/$ver/victoriametrics-datasource-$ver.tar.gz -o /var/lib/grafana/plugins/plugin.tar.gz
tar -xzf ` + pluginsPath + `/plugin.tar.gz -C ` + pluginsPath + `/
rm -f ` + pluginsPath + `/plugin.tar.gz
chown -R ` + d64(userID) + `:` + d64(userID) + ` ` + pluginsPath + `/
`

var VolumeMountScripts = corev1.VolumeMount{
	MountPath: "/etc/grafana/download_vm_ds.sh",
	Name:      GrafCAM.VolumeMount.Name,
	SubPath:   scriptsName,
}

const scriptsName = "download_vm_ds.sh"

var VolumeMountGrafanaIni = corev1.VolumeMount{
	MountPath: "/etc/grafana/grafana.ini",
	Name:      GrafCAM.VolumeMount.Name,
	SubPath:   "grafana.ini",
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

var VolumeDashProvider = corev1.VolumeMount{
	MountPath: "/etc/grafana/provisioning/dashboards/sc-" + filenameProviders,
	Name:      GrafCAM.VolumeMount.Name,
	SubPath:   filenameProviders,
}

const filenameProviders = "dashboardproviders.yaml"

var VolumeDataSource = ku.VolumeAndMount{
	VolumeMount: corev1.VolumeMount{
		MountPath: "/etc/grafana/provisioning/datasources",
		Name:      "sc-datasources-volume",
	},
	VolumeSource: corev1.VolumeSource{},
}

var VolumeStorage = ku.VolumeAndMount{
	VolumeMount: corev1.VolumeMount{
		MountPath: "/var/lib/grafana",
		Name:      "storage",
	},
	VolumeSource: corev1.VolumeSource{},
}

var VolumeDashboards = ku.VolumeAndMount{
	VolumeMount: corev1.VolumeMount{
		MountPath: "/tmp/dashboards",
		Name:      "sc-dashboard-volume",
	},
	VolumeSource: corev1.VolumeSource{},
}

var SideCarProvider = ku.ConfigAndMount{
	ObjectMeta: Graf.ObjectMetaNameSuffix("sidecar-provider"),
	VolumeMount: corev1.VolumeMount{
		Name:      "sc-dashboard-provider",
		MountPath: "/etc/grafana/provisioning/dashboards/" + filenameProviders,
		SubPath:   filenameProviders,
	},
	Data: map[string]string{
		filenameProviders: `
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
      foldersFromFilesStructure: true
      path: /tmp/dashboards
`,
	},
}

type GrafDSConfigFile struct {
	APIVersion  int              `json:"apiVersion"`
	DataSources []GrafDataSource `json:"datasources"`
}

type GrafDataSource struct {
	Access    string `json:"access"`
	IsDefault bool   `json:"isDefault"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	URL       string `json:"url"`
}

var GrafDSConfig = GrafDSConfigFile{
	APIVersion: 1,
	DataSources: []GrafDataSource{
		{
			Name:      VictoriaMetricsDataSourceName,
			Type:      VictoriaMetricsDataSourceID,
			Access:    "proxy",
			IsDefault: true,
			URL: fmt.Sprintf(
				"http://%s.%s.svc:%d/",
				VMDB.PrefixedName(),
				VMDB.Namespace,
				VMSinglePort,
			),
		},
		{
			Name:      PrometheusDataSourceName,
			Type:      PrometheusDataSourceID,
			Access:    "proxy",
			IsDefault: false,
			URL: fmt.Sprintf(
				"http://%s.%s.svc:%d/",
				VMDB.PrefixedName(),
				VMDB.Namespace,
				VMSinglePort,
			),
		},
	},
}

func YamlMust(a any) string {
	res, err := kyaml.Marshal(a)
	if err != nil {
		panic("encoding to yaml: " + err.Error())
	}
	return string(res)
}

var GrafanaDataSource = ku.ConfigAndMount{
	ObjectMeta:  PatchDataSourceLabels(Graf.ObjectMetaNameSuffix("datasources")),
	VolumeMount: corev1.VolumeMount{},
	Data:        map[string]string{"datasource.yaml": YamlMust(GrafDSConfig)},
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
