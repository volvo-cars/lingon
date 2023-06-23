// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmk8s

import (
	"fmt"

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
	GrafanaVersion             = "9.3.0"
	GrafanaSideCarImg          = "quay.io/kiwigrid/k8s-sidecar:1.19.2"
	GrafanaPort                = 3000
	GrafanaPortName            = "service"
	DashboardLabel             = "grafana_dashboard"
	DataSourceLabel            = "grafana_datasource"
	defaultDashboardConfigName = "grafana-default-dashboards"
)

var Graf = &Metadata{
	Name:      "grafana",
	Namespace: namespace,
	Instance:  "grafana" + namespace,
	Component: "dashboards",
	PartOf:    appName,
	Version:   GrafanaVersion,
	ManagedBy: "lingon",
	Registry:  "",
	Image:     "grafana/grafana",
	Tag:       GrafanaVersion,
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
	ObjectMeta: Graf.ObjectMetaNameSuffixNoNS("-cr"),
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
					"checksum/config":     ku.HashConfig(GrafanaCM),
					"checksum/provider":   ku.HashConfig(GrafanaProviderCM),
					"checksum/datasource": ku.HashConfig(GrafanaDataSourceCM),
					"checksum/secret":     ku.HashSecret(GrafanaSecrets),
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
							{
								Name:  "METHOD",
								Value: "WATCH",
							}, {
								Name:  "LABEL",
								Value: DashboardLabel,
							}, {
								Name:  "FOLDER",
								Value: "/tmp/dashboards",
							}, {
								Name:  "RESOURCE",
								Value: "both",
							},
						},

						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/tmp/dashboards",
								Name:      "sc-dashboard-volume",
							},
						},
					}, {
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
					}, {
						Name:            Graf.Name,
						Image:           Graf.ContainerURL(),
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
						Image:           "curlimages/curl:7.85.0",
						ImagePullPolicy: corev1.PullIfNotPresent,
						Name:            "download-dashboards",
						Command:         []string{"/bin/sh"},
						Args: []string{
							"-c",
							"mkdir -p /var/lib/grafana/dashboards/default && " +
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
				},
				SecurityContext: &corev1.PodSecurityContext{
					FSGroup:    P(int64(472)),
					RunAsGroup: P(int64(472)),
					RunAsUser:  P(int64(472)),
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

var GrafanaSVC = Graf.Service(80, GrafanaPort, GrafanaPortName)

// GrafanaSVC = &corev1.Service{
// 	TypeMeta:   ku.TypeServiceV1,
// 	ObjectMeta: Graf.ObjectMeta(),
// 	Spec: corev1.ServiceSpec{
// 		Ports: []corev1.ServicePort{
// 			{
// 				Name:       GrafanaPortName,
// 				Port:       int32(80),
// 				Protocol:   corev1.ProtocolTCP,
// 				TargetPort: intstr.FromInt(GrafanaPort),
// 			},
// 		},
// 		Selector: Graf.MatchLabels(),
// 		Type:     corev1.ServiceTypeClusterIP,
// 	},
// }

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
	ObjectMeta: Graf.ObjectMeta(),
	TypeMeta:   ku.TypeConfigMapV1,
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
		"datasource.yaml": `
apiVersion: 1
datasources:
- name: VictoriaMetrics
  type: prometheus
  url: ` + fmt.Sprintf(
			"http://%s.%s.svc:8429/",
			VMDB.PrefixedName(), namespace,
		) + `
  access: proxy
  isDefault: true
  jsonData: 
    {}
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
		Labels: ku.MergeLabels(
			Graf.Labels(),
			map[string]string{DashboardLabel: "1"},
		),
		Name:      Graf.Name + "-dash-overview",
		Namespace: Graf.Namespace,
	},
	TypeMeta: ku.TypeConfigMapV1,
}
