// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmk8s

import (
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

var O = Core()

func Core() VMOpMeta {
	OpPort := 8080
	OpPortName := "http"
	OpWHPort := 9443
	OpWHPortName := "webhook"

	return VMOpMeta{
		Metadata: meta.Metadata{
			Name:      "victoria-metrics-operator",
			Namespace: namespace,
			Instance:  "victoria-metrics-operator-" + namespace,
			Component: "operator",
			PartOf:    appName,
			Version:   "0.34.1",
			ManagedBy: "lingon",
			Img: meta.ContainerImg{
				Registry: "docker.io",
				Image:    "victoriametrics/operator",
				Tag:      "v0.34.1",
			},
		},
		P: meta.NetPort{
			Container: corev1.ContainerPort{
				ContainerPort: int32(OpPort),
				Name:          OpPortName,
			},
			Service: corev1.ServicePort{
				Name:       OpPortName,
				Port:       int32(OpPort),
				TargetPort: intstr.FromInt(OpPort),
			},
		},
		WH: meta.NetPort{
			Container: corev1.ContainerPort{
				Name:          OpWHPortName,
				ContainerPort: int32(OpWHPort),
			},
			Service: corev1.ServicePort{
				Name:       OpWHPortName,
				Port:       int32(443),
				TargetPort: intstr.FromInt(OpWHPort),
			},
		},
	}
}

type VMOpMeta struct {
	meta.Metadata
	WH meta.NetPort
	P  meta.NetPort
}

type Operator struct {
	kube.App

	CR          *rbacv1.ClusterRole
	CRB         *rbacv1.ClusterRoleBinding
	Deploy      *appsv1.Deployment
	RB          *rbacv1.RoleBinding
	Role        *rbacv1.Role
	SA          *corev1.ServiceAccount
	SVC         *corev1.Service
	DashboardCM *corev1.ConfigMap
	Scrape      *v1beta1.VMServiceScrape
}

func NewOperator() *Operator {
	return &Operator{
		CR: OperatorCR,
		CRB: ku.BindClusterRole(
			O.Name,
			OperatorSA,
			OperatorCR,
			O.Labels(),
		),
		Deploy: OperatorDeploy,
		RB: ku.BindRole(
			O.Name,
			OperatorSA,
			OperatorRole,
			O.Labels(),
		),
		Role:        OperatorRole,
		SA:          OperatorSA,
		SVC:         OperatorSVC,
		DashboardCM: DashboardVMOperator,
		Scrape:      OperatorScrape,
	}
}

var OperatorDeploy = &appsv1.Deployment{
	TypeMeta:   ku.TypeDeploymentV1,
	ObjectMeta: O.ObjectMeta(),
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(1)),
		Selector: &metav1.LabelSelector{
			MatchLabels: O.MatchLabels(),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: O.MatchLabels(),
			},
			Spec: corev1.PodSpec{
				ServiceAccountName: OperatorSA.Name,
				Containers: []corev1.Container{
					{
						Image:           O.Img.URL(),
						ImagePullPolicy: corev1.PullIfNotPresent,
						Name:            O.Name,
						Args: []string{
							"--zap-log-level=info",
							"--enable-leader-election",
						},
						Command: []string{"manager"},
						Env: []corev1.EnvVar{
							{Name: "WATCH_NAMESPACE"},
							ku.EnvVarDownAPI("POD_NAME", "metadata.name"),
							{Name: "OPERATOR_NAME", Value: O.Name},
							// See https://github.com/VictoriaMetrics/operator/blob/master/vars.MD

						},

						Ports: []corev1.ContainerPort{
							O.P.Container,
							O.WH.Container,
						},
					},
				},
			},
		},
	},
}

var OperatorSA = O.ServiceAccount()

var OperatorCR = &rbacv1.ClusterRole{
	TypeMeta:   ku.TypeClusterRoleV1,
	ObjectMeta: O.ObjectMetaNoNS(),
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{""},
			Resources: []string{"configmaps", "configmaps/finalizers"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"endpoints"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"events"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"namespaces"},
			Verbs:     []string{"get", "list", "watch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{
				"persistentvolumeclaims",
				"persistentvolumeclaims/finalizers",
			},
			Verbs: []string{"*"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"pods"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"secrets", "secrets/finalizers"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"services"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"services/finalizers"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{"apps"},
			Resources: []string{"deployments", "deployments/finalizers"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{"apps"},
			Resources: []string{"replicasets"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{"apps"},
			Resources: []string{
				"statefulsets",
				"statefulsets/finalizers",
				"statefulsets/status",
			},
			Verbs: []string{"*"},
		}, {
			APIGroups: []string{"monitoring.coreos.com"},
			Resources: []string{"*"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmagents", "vmagents/finalizers"},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmagents/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{
				"vmalertmanagers",
				"vmalertmanagers/finalizers",
			},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmalertmanagers/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{
				"vmalertmanagerconfigs",
				"vmalertmanagerconfigs/finalizers",
			},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmalertmanagerconfigss/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmalerts", "vmalerts/finalizers"},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmalerts/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmclusters", "vmclusters/finalizers"},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmclusters/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmpodscrapes", "vmprobscrapes/finalizers"},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmpodscrapes/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmrules", "vmrules/finalizers"},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmrules/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{
				"vmservicescrapes",
				"vmservicescrapes/finalizers",
			},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmservicescrapes/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmprobes"},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmprobes/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmsingles", "vmsingles/finalizers"},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmsingles/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{""},
			Resources: []string{
				"nodes",
				"nodes/proxy",
				"services",
				"endpoints",
				"pods",
				"endpointslices",
				"configmaps",
				"nodes/metrics",
				"namespaces",
			},
			Verbs: []string{"get", "list", "watch"},
		}, {
			APIGroups: []string{"extensions", "networking.k8s.io"},
			Resources: []string{"ingresses"},
			Verbs:     []string{"get", "list", "watch"},
		}, {
			NonResourceURLs: []string{"/metrics", "/metrics/resources"},
			Verbs:           []string{"get", "watch", "list"},
		}, {
			APIGroups: []string{"rbac.authorization.k8s.io"},
			Resources: []string{
				"clusterrolebindings",
				"clusterrolebindings/finalizers",
				"clusterroles",
				"clusterroles/finalizers",
				"roles",
				"rolebindings",
			},
			Verbs: []string{
				"get",
				"list",
				"create",
				"patch",
				"update",
				"watch",
				"delete",
			},
		}, {
			APIGroups: []string{"policy"},
			Resources: []string{
				"podsecuritypolicies",
				"podsecuritypolicies/finalizers",
			},
			Verbs: []string{
				"get",
				"list",
				"create",
				"patch",
				"update",
				"use",
				"watch",
				"delete",
			},
		}, {
			APIGroups: []string{""},
			Resources: []string{
				"serviceaccounts",
				"serviceaccounts/finalizers",
			},
			Verbs: []string{
				"get",
				"list",
				"create",
				"watch",
				"update",
				"delete",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmnodescrapes", "vmnodescrapes/finalizers"},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmnodescrapes/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{
				"vmstaticscrapes",
				"vmstaticscrapes/finalizers",
			},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmstaticscrapes/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{
				"vmauths",
				"vmauths/finalizers",
				"vmusers",
				"vmusers/finalizers",
			},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmusers/status", "vmauths/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"storage.k8s.io"},
			Resources: []string{"storageclasses"},
			Verbs:     []string{"list", "get", "watch"},
		}, {
			APIGroups: []string{"policy"},
			Resources: []string{
				"poddisruptionbudgets",
				"poddisruptionbudgets/finalizers",
			},
			Verbs: []string{"*"},
		}, {
			APIGroups: []string{"route.openshift.io", "image.openshift.io"},
			Resources: []string{"routers/metrics", "registry/metrics"},
			Verbs:     []string{"get"},
		}, {
			APIGroups: []string{"autoscaling"},
			Resources: []string{"horizontalpodautoscalers"},
			Verbs: []string{
				"list",
				"get",
				"delete",
				"create",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"networking.k8s.io", "extensions"},
			Resources: []string{"ingresses", "ingresses/finalizers"},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"apiextensions.k8s.io"},
			Resources: []string{"customresourcedefinitions"},
			Verbs:     []string{"get", "list"},
		},
	},
}

var OperatorRole = &rbacv1.Role{
	TypeMeta:   ku.TypeRoleV1,
	ObjectMeta: O.ObjectMeta(),
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{""},
			Resources: []string{"configmaps"},
			Verbs: []string{
				"get",
				"list",
				"watch",
				"create",
				"update",
				"patch",
				"delete",
			},
		}, {
			APIGroups: []string{""},
			Resources: []string{"configmaps/status"},
			Verbs:     []string{"get", "update", "patch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"events"},
			Verbs:     []string{"create", "patch"},
		}, {
			APIGroups: []string{"coordination.k8s.io"},
			Resources: []string{"leases"},
			Verbs:     []string{"create", "get", "update"},
		},
	},
}

var OperatorSVC = &corev1.Service{
	TypeMeta:   ku.TypeServiceV1,
	ObjectMeta: O.ObjectMeta(),
	Spec: corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			O.P.Service,
			O.WH.Service,
		},
		Selector: O.MatchLabels(),
		Type:     corev1.ServiceTypeClusterIP,
	},
}

var OperatorScrape = &v1beta1.VMServiceScrape{
	TypeMeta:   TypeVMServiceScrapeV1Beta1,
	ObjectMeta: O.ObjectMeta(),
	Spec: v1beta1.VMServiceScrapeSpec{
		Endpoints: []v1beta1.Endpoint{{Port: O.P.Service.Name}},
		NamespaceSelector: v1beta1.NamespaceSelector{
			MatchNames: []string{O.Namespace},
		},
		Selector: metav1.LabelSelector{MatchLabels: O.MatchLabels()},
	},
}

var DashboardVMOperator = &corev1.ConfigMap{
	Data: map[string]string{
		"operator.json": `
{
    "__inputs": [],
    "__elements": {},
    "__requires": [
        {
            "type": "grafana",
            "id": "grafana",
            "name": "Grafana",
            "version": "9.2.2"
        },
        {
            "type": "panel",
            "id": "graph",
            "name": "Graph (old)",
            "version": ""
        },
        {
            "type": "datasource",
            "id": "prometheus",
            "name": "Prometheus",
            "version": "1.0.0"
        },
        {
            "type": "panel",
            "id": "stat",
            "name": "Stat",
            "version": ""
        },
        {
            "type": "panel",
            "id": "text",
            "name": "Text",
            "version": ""
        }
    ],
    "annotations": {
        "list": [
            {
                "builtIn": 1,
                "datasource": {
                    "type": "datasource",
                    "uid": "grafana"
                },
                "enable": true,
                "hide": true,
                "iconColor": "rgba(0, 211, 255, 1)",
                "name": "Annotations & Alerts",
                "target": {
                    "limit": 100,
                    "matchAny": false,
                    "tags": [],
                    "type": "dashboard"
                },
                "type": "dashboard"
            }
        ]
    },
    "description": "Overview for operator VictoriaMetrics v0.25.0 or higher",
    "editable": true,
    "fiscalYearStartMonth": 0,
    "graphTooltip": 0,
    "id": null,
    "links": [],
    "liveNow": false,
    "panels": [
        {
            "collapsed": false,
            "datasource": {
                "type": "prometheus",
                "uid": "$ds"
            },
            "gridPos": {
                "h": 1,
                "w": 24,
                "x": 0,
                "y": 0
            },
            "id": 8,
            "panels": [],
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "refId": "A"
                }
            ],
            "title": "Overview",
            "type": "row"
        },
        {
            "datasource": {
                "type": "prometheus",
                "uid": "$ds"
            },
            "gridPos": {
                "h": 3,
                "w": 4,
                "x": 0,
                "y": 1
            },
            "id": 24,
            "options": {
                "code": {
                    "language": "plaintext",
                    "showLineNumbers": false,
                    "showMiniMap": false
                },
                "content": "<div style=\"text-align: center;\">$version</div>",
                "mode": "markdown"
            },
            "pluginVersion": "9.2.2",
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "refId": "A"
                }
            ],
            "title": "Version",
            "type": "text"
        },
        {
            "datasource": {
                "type": "prometheus",
                "uid": "$ds"
            },
            "description": "Number of objects at kubernetes cluster per each controller",
            "fieldConfig": {
                "defaults": {
                    "color": {
                        "mode": "thresholds"
                    },
                    "mappings": [],
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
                "overrides": []
            },
            "gridPos": {
                "h": 7,
                "w": 20,
                "x": 4,
                "y": 1
            },
            "id": 14,
            "options": {
                "colorMode": "none",
                "graphMode": "area",
                "justifyMode": "auto",
                "orientation": "auto",
                "reduceOptions": {
                    "calcs": [
                        "lastNotNull"
                    ],
                    "fields": "",
                    "values": false
                },
                "text": {},
                "textMode": "auto"
            },
            "pluginVersion": "9.2.2",
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "editorMode": "code",
                    "expr": "max(operator_controller_objects_count{job=~\"$job\",instance=~\"$instance\"}) by (controller)",
                    "legendFormat": "{{controller}}",
                    "range": true,
                    "refId": "A"
                }
            ],
            "title": "CRD Objects count by controller",
            "type": "stat"
        },
        {
            "datasource": {
                "type": "prometheus",
                "uid": "$ds"
            },
            "fieldConfig": {
                "defaults": {
                    "color": {
                        "mode": "thresholds"
                    },
                    "mappings": [],
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
                    },
                    "unit": "s"
                },
                "overrides": []
            },
            "gridPos": {
                "h": 4,
                "w": 4,
                "x": 0,
                "y": 4
            },
            "id": 22,
            "options": {
                "colorMode": "value",
                "graphMode": "area",
                "justifyMode": "auto",
                "orientation": "auto",
                "reduceOptions": {
                    "calcs": [
                        "lastNotNull"
                    ],
                    "fields": "",
                    "values": false
                },
                "textMode": "auto"
            },
            "pluginVersion": "9.2.2",
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "editorMode": "code",
                    "exemplar": false,
                    "expr": "vm_app_uptime_seconds{job=~\"$job\",instance=~\"$instance\"}",
                    "format": "table",
                    "instant": true,
                    "interval": "",
                    "legendFormat": "{{instance}}",
                    "range": false,
                    "refId": "A"
                }
            ],
            "title": "Uptime",
            "type": "stat"
        },
        {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": {
                "type": "prometheus",
                "uid": "$ds"
            },
            "fill": 1,
            "fillGradient": 0,
            "gridPos": {
                "h": 13,
                "w": 12,
                "x": 0,
                "y": 8
            },
            "hiddenSeries": false,
            "id": 12,
            "legend": {
                "alignAsTable": true,
                "avg": true,
                "current": false,
                "max": true,
                "min": false,
                "show": true,
                "total": false,
                "values": true
            },
            "lines": true,
            "linewidth": 1,
            "nullPointMode": "null",
            "options": {
                "alertThreshold": true
            },
            "percentage": false,
            "pluginVersion": "9.2.2",
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "editorMode": "code",
                    "expr": "sum(rate(controller_runtime_reconcile_total{job=~\"$job\",instance=~\"$instance\",result=~\"requeue_after|requeue|success\"}[$__rate_interval])) by(controller)",
                    "legendFormat": "{{controller}}",
                    "range": true,
                    "refId": "A"
                }
            ],
            "thresholds": [],
            "timeRegions": [],
            "title": "Reconciliation rate by controller",
            "tooltip": {
                "shared": true,
                "sort": 0,
                "value_type": "individual"
            },
            "type": "graph",
            "xaxis": {
                "mode": "time",
                "show": true,
                "values": []
            },
            "yaxes": [
                {
                    "format": "short",
                    "logBase": 1,
                    "show": true
                },
                {
                    "format": "short",
                    "logBase": 1,
                    "show": true
                }
            ],
            "yaxis": {
                "align": false
            }
        },
        {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": {
                "type": "prometheus",
                "uid": "$ds"
            },
            "description": "",
            "fill": 1,
            "fillGradient": 0,
            "gridPos": {
                "h": 13,
                "w": 12,
                "x": 12,
                "y": 8
            },
            "hiddenSeries": false,
            "id": 16,
            "legend": {
                "alignAsTable": true,
                "avg": true,
                "current": false,
                "max": true,
                "min": false,
                "show": true,
                "total": false,
                "values": true
            },
            "lines": true,
            "linewidth": 1,
            "nullPointMode": "null",
            "options": {
                "alertThreshold": true
            },
            "percentage": false,
            "pluginVersion": "9.2.2",
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "editorMode": "code",
                    "expr": "sum(rate(operator_log_messages_total{job=~\"$job\",instance=~\"$instance\"}[$__rate_interval])) by (level)",
                    "legendFormat": "{{label_name}}",
                    "range": true,
                    "refId": "A"
                }
            ],
            "thresholds": [],
            "timeRegions": [],
            "title": "Log message rate",
            "tooltip": {
                "shared": true,
                "sort": 0,
                "value_type": "individual"
            },
            "type": "graph",
            "xaxis": {
                "mode": "time",
                "show": true,
                "values": []
            },
            "yaxes": [
                {
                    "format": "short",
                    "logBase": 1,
                    "show": true
                },
                {
                    "format": "short",
                    "logBase": 1,
                    "show": true
                }
            ],
            "yaxis": {
                "align": false
            }
        },
        {
            "collapsed": false,
            "datasource": {
                "type": "prometheus",
                "uid": "$ds"
            },
            "gridPos": {
                "h": 1,
                "w": 24,
                "x": 0,
                "y": 21
            },
            "id": 6,
            "panels": [],
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "refId": "A"
                }
            ],
            "title": "Troubleshooting",
            "type": "row"
        },
        {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": {
                "type": "prometheus",
                "uid": "$ds"
            },
            "description": "Non zero metrics indicates about error with CR object definition (typos or incorrect values) or errors with kubernetes API connection.",
            "fill": 1,
            "fillGradient": 0,
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 0,
                "y": 22
            },
            "hiddenSeries": false,
            "id": 10,
            "legend": {
                "alignAsTable": true,
                "avg": false,
                "current": true,
                "max": true,
                "min": false,
                "show": true,
                "total": false,
                "values": true
            },
            "lines": true,
            "linewidth": 1,
            "nullPointMode": "null",
            "options": {
                "alertThreshold": true
            },
            "percentage": false,
            "pluginVersion": "8.3.2",
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "editorMode": "code",
                    "exemplar": false,
                    "expr": "sum(rate(controller_runtime_reconcile_errors_total{job=~\"$job\",instance=~\"$instance\"}[$__rate_interval])) by(controller) > 0 ",
                    "instant": false,
                    "legendFormat": "{{controller}}",
                    "range": true,
                    "refId": "A"
                },
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "editorMode": "code",
                    "expr": "sum(rate(controller_runtime_reconcile_total{job=~\"$job\",instance=~\"$instance\",result=\"error\"}[$__rate_interval])) by(controller) > 0",
                    "hide": false,
                    "legendFormat": "{{label_name}}",
                    "range": true,
                    "refId": "B"
                }
            ],
            "thresholds": [],
            "timeRegions": [],
            "title": "reconcile errors by controller",
            "tooltip": {
                "shared": true,
                "sort": 0,
                "value_type": "individual"
            },
            "type": "graph",
            "xaxis": {
                "mode": "time",
                "show": true,
                "values": []
            },
            "yaxes": [
                {
                    "format": "short",
                    "logBase": 1,
                    "show": true
                },
                {
                    "format": "short",
                    "logBase": 1,
                    "show": true
                }
            ],
            "yaxis": {
                "align": false
            }
        },
        {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": {
                "type": "prometheus",
                "uid": "$ds"
            },
            "description": "Operator limits number of reconcilation events to 5 events per 2 seconds.\n For now, this limit is applied only for vmalert and vmagent controllers.\n It should reduce load at kubernetes cluster and increase operator performance.",
            "fill": 1,
            "fillGradient": 0,
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 12,
                "y": 22
            },
            "hiddenSeries": false,
            "id": 18,
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
            "pluginVersion": "8.3.2",
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "editorMode": "code",
                    "expr": "sum(rate(operator_reconcile_throttled_events_total[$__rate_interval])) by(controller)",
                    "legendFormat": "{{controller}}",
                    "range": true,
                    "refId": "A"
                }
            ],
            "thresholds": [],
            "timeRegions": [],
            "title": "throttled reconcilation events",
            "tooltip": {
                "shared": true,
                "sort": 0,
                "value_type": "individual"
            },
            "type": "graph",
            "xaxis": {
                "mode": "time",
                "show": true,
                "values": []
            },
            "yaxes": [
                {
                    "format": "short",
                    "logBase": 1,
                    "show": true
                },
                {
                    "format": "short",
                    "logBase": 1,
                    "show": true
                }
            ],
            "yaxis": {
                "align": false
            }
        },
        {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": {
                "type": "prometheus",
                "uid": "$ds"
            },
            "description": "Number of objects waiting in the queue for reconciliation. Non-zero values indicate that operator cannot process CR objects changes with the given resources.",
            "fill": 1,
            "fillGradient": 0,
            "gridPos": {
                "h": 11,
                "w": 12,
                "x": 0,
                "y": 30
            },
            "hiddenSeries": false,
            "id": 20,
            "legend": {
                "alignAsTable": true,
                "avg": false,
                "current": true,
                "max": true,
                "min": false,
                "show": true,
                "total": false,
                "values": true
            },
            "lines": true,
            "linewidth": 1,
            "nullPointMode": "null",
            "options": {
                "alertThreshold": true
            },
            "percentage": false,
            "pluginVersion": "8.3.2",
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "editorMode": "code",
                    "expr": "max(workqueue_depth{job=~\"$job\",instance=~\"$instance\"}) by (name)",
                    "legendFormat": "{{label_name}}",
                    "range": true,
                    "refId": "A"
                }
            ],
            "thresholds": [],
            "timeRegions": [],
            "title": "Wokring queue depth",
            "tooltip": {
                "shared": true,
                "sort": 0,
                "value_type": "individual"
            },
            "type": "graph",
            "xaxis": {
                "mode": "time",
                "show": true,
                "values": []
            },
            "yaxes": [
                {
                    "format": "short",
                    "logBase": 1,
                    "show": true
                },
                {
                    "format": "short",
                    "logBase": 1,
                    "show": true
                }
            ],
            "yaxis": {
                "align": false
            }
        },
        {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": {
                "type": "prometheus",
                "uid": "$ds"
            },
            "description": " For controllers with StatefulSet it's ok to see latency greater then 3 seconds. It could be vmalertmanager,vmcluster or vmagent in statefulMode.\n\n For other controllers, latency greater then 1 second may indicate issues with kubernetes cluster or operator's performance.\n ",
            "fieldConfig": {
                "defaults": {
                    "unit": "s"
                },
                "overrides": []
            },
            "fill": 1,
            "fillGradient": 0,
            "gridPos": {
                "h": 11,
                "w": 12,
                "x": 12,
                "y": 30
            },
            "hiddenSeries": false,
            "id": 26,
            "legend": {
                "alignAsTable": true,
                "avg": true,
                "current": false,
                "max": true,
                "min": false,
                "rightSide": false,
                "show": true,
                "total": false,
                "values": true
            },
            "lines": true,
            "linewidth": 1,
            "nullPointMode": "null",
            "options": {
                "alertThreshold": true
            },
            "percentage": false,
            "pluginVersion": "8.3.2",
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "editorMode": "code",
                    "expr": "histogram_quantile(0.99,sum(rate(controller_runtime_reconcile_time_seconds_bucket[$__rate_interval])) by(le,controller) )",
                    "legendFormat": "q.99 {{controller}}",
                    "range": true,
                    "refId": "A"
                }
            ],
            "thresholds": [],
            "timeRegions": [],
            "title": "Reconcilation latency by controller",
            "tooltip": {
                "shared": true,
                "sort": 0,
                "value_type": "individual"
            },
            "type": "graph",
            "xaxis": {
                "mode": "time",
                "show": true,
                "values": []
            },
            "yaxes": [
                {
                    "format": "s",
                    "logBase": 1,
                    "show": true
                },
                {
                    "format": "short",
                    "logBase": 1,
                    "show": true
                }
            ],
            "yaxis": {
                "align": false
            }
        },
        {
            "collapsed": false,
            "datasource": {
                "type": "prometheus",
                "uid": "$ds"
            },
            "gridPos": {
                "h": 1,
                "w": 24,
                "x": 0,
                "y": 41
            },
            "id": 4,
            "panels": [],
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "refId": "A"
                }
            ],
            "title": "resources",
            "type": "row"
        },
        {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": {
                "type": "prometheus",
                "uid": "$ds"
            },
            "fieldConfig": {
                "defaults": {
                    "unit": "bytes"
                },
                "overrides": []
            },
            "fill": 1,
            "fillGradient": 0,
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 0,
                "y": 42
            },
            "hiddenSeries": false,
            "id": 28,
            "legend": {
                "alignAsTable": true,
                "avg": true,
                "current": true,
                "max": true,
                "min": false,
                "show": true,
                "total": false,
                "values": true
            },
            "lines": true,
            "linewidth": 1,
            "nullPointMode": "null",
            "options": {
                "alertThreshold": true
            },
            "percentage": false,
            "pluginVersion": "8.3.2",
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "editorMode": "code",
                    "expr": "sum(go_memstats_sys_bytes{job=~\"$job\", instance=~\"$instance\"}) ",
                    "legendFormat": "requested from system",
                    "range": true,
                    "refId": "A"
                },
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "editorMode": "code",
                    "expr": "sum(go_memstats_heap_inuse_bytes{job=~\"$job\", instance=~\"$instance\"}) ",
                    "hide": false,
                    "legendFormat": "heap inuse",
                    "range": true,
                    "refId": "B"
                },
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "editorMode": "code",
                    "expr": "sum(go_memstats_stack_inuse_bytes{job=~\"$job\", instance=~\"$instance\"})",
                    "hide": false,
                    "legendFormat": "stack inuse",
                    "range": true,
                    "refId": "C"
                },
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "editorMode": "code",
                    "expr": "sum(process_resident_memory_bytes{job=~\"$job\", instance=~\"$instance\"})",
                    "hide": false,
                    "legendFormat": "resident",
                    "range": true,
                    "refId": "D"
                }
            ],
            "thresholds": [],
            "timeRegions": [],
            "title": "Memory usage ($instance)",
            "tooltip": {
                "shared": true,
                "sort": 0,
                "value_type": "individual"
            },
            "type": "graph",
            "xaxis": {
                "mode": "time",
                "show": true,
                "values": []
            },
            "yaxes": [
                {
                    "format": "bytes",
                    "logBase": 1,
                    "show": true
                },
                {
                    "format": "short",
                    "logBase": 1,
                    "show": true
                }
            ],
            "yaxis": {
                "align": false
            }
        },
        {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": {
                "type": "prometheus",
                "uid": "$ds"
            },
            "fill": 1,
            "fillGradient": 0,
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 12,
                "y": 42
            },
            "hiddenSeries": false,
            "id": 30,
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
            "pluginVersion": "8.3.2",
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "editorMode": "code",
                    "expr": "rate(process_cpu_seconds_total{job=~\"$job\", instance=~\"$instance\"}[$__rate_interval])",
                    "legendFormat": "CPU cores used",
                    "range": true,
                    "refId": "A"
                }
            ],
            "thresholds": [],
            "timeRegions": [],
            "title": "CPU ($instance)",
            "tooltip": {
                "shared": true,
                "sort": 0,
                "value_type": "individual"
            },
            "type": "graph",
            "xaxis": {
                "mode": "time",
                "show": true,
                "values": []
            },
            "yaxes": [
                {
                    "format": "short",
                    "logBase": 1,
                    "show": true
                },
                {
                    "format": "short",
                    "logBase": 1,
                    "show": true
                }
            ],
            "yaxis": {
                "align": false
            }
        },
        {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": {
                "type": "prometheus",
                "uid": "$ds"
            },
            "fill": 1,
            "fillGradient": 0,
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 0,
                "y": 50
            },
            "hiddenSeries": false,
            "id": 32,
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
            "pluginVersion": "8.3.2",
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "editorMode": "code",
                    "expr": "sum(go_goroutines{job=~\"$job\", instance=~\"$instance\"})",
                    "legendFormat": "goroutines",
                    "range": true,
                    "refId": "A"
                }
            ],
            "thresholds": [],
            "timeRegions": [],
            "title": "Goroutines ($instance)",
            "tooltip": {
                "shared": true,
                "sort": 0,
                "value_type": "individual"
            },
            "type": "graph",
            "xaxis": {
                "mode": "time",
                "show": true,
                "values": []
            },
            "yaxes": [
                {
                    "format": "short",
                    "logBase": 1,
                    "show": true
                },
                {
                    "format": "short",
                    "logBase": 1,
                    "show": true
                }
            ],
            "yaxis": {
                "align": false
            }
        },
        {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": {
                "type": "prometheus",
                "uid": "$ds"
            },
            "fieldConfig": {
                "defaults": {
                    "unit": "s"
                },
                "overrides": []
            },
            "fill": 1,
            "fillGradient": 0,
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 12,
                "y": 50
            },
            "hiddenSeries": false,
            "id": 34,
            "legend": {
                "alignAsTable": true,
                "avg": true,
                "current": false,
                "max": true,
                "min": false,
                "show": true,
                "total": false,
                "values": true
            },
            "lines": true,
            "linewidth": 1,
            "nullPointMode": "null",
            "options": {
                "alertThreshold": true
            },
            "percentage": false,
            "pluginVersion": "8.3.2",
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "$ds"
                    },
                    "editorMode": "code",
                    "expr": "sum(rate(go_gc_duration_seconds_sum{job=~\"$job\", instance=~\"$instance\"}[$__rate_interval]))\n/\nsum(rate(go_gc_duration_seconds_count{job=~\"$job\", instance=~\"$instance\"}[$__rate_interval]))",
                    "legendFormat": "avg gc duration",
                    "range": true,
                    "refId": "A"
                }
            ],
            "thresholds": [],
            "timeRegions": [],
            "title": "GC duration ($instance)",
            "tooltip": {
                "shared": true,
                "sort": 0,
                "value_type": "individual"
            },
            "type": "graph",
            "xaxis": {
                "mode": "time",
                "show": true,
                "values": []
            },
            "yaxes": [
                {
                    "format": "s",
                    "logBase": 1,
                    "show": true
                },
                {
                    "format": "short",
                    "logBase": 1,
                    "show": true
                }
            ],
            "yaxis": {
                "align": false
            }
        }
    ],
    "refresh": "",
    "schemaVersion": 37,
    "style": "dark",
    "tags": [
        "operator",
        "VictoriaMetrics"
    ],
    "templating": {
        "list": [
            {
                "current": {
                    "selected": false,
                    "text": "cloud-c15",
                    "value": "cloud-c15"
                },
                "hide": 0,
                "includeAll": false,
                "multi": false,
                "name": "ds",
                "options": [],
                "query": "prometheus",
                "queryValue": "",
                "refresh": 1,
                "regex": "",
                "skipUrlSync": false,
                "type": "datasource"
            },
            {
                "current": {},
                "datasource": {
                    "type": "prometheus",
                    "uid": "$ds"
                },
                "definition": "label_values(operator_log_messages_total,job)",
                "hide": 0,
                "includeAll": false,
                "multi": false,
                "name": "job",
                "options": [],
                "query": {
                    "query": "label_values(operator_log_messages_total,job)",
                    "refId": "StandardVariableQuery"
                },
                "refresh": 2,
                "regex": "",
                "skipUrlSync": false,
                "sort": 0,
                "type": "query"
            },
            {
                "current": {},
                "datasource": {
                    "type": "prometheus",
                    "uid": "$ds"
                },
                "definition": "label_values(operator_log_messages_total{job=~\"$job\"},instance)",
                "hide": 0,
                "includeAll": true,
                "multi": false,
                "name": "instance",
                "options": [],
                "query": {
                    "query": "label_values(operator_log_messages_total{job=~\"$job\"},instance)",
                    "refId": "StandardVariableQuery"
                },
                "refresh": 2,
                "regex": "",
                "skipUrlSync": false,
                "sort": 0,
                "type": "query"
            },
            {
                "current": {},
                "datasource": {
                    "type": "prometheus",
                    "uid": "$ds"
                },
                "definition": "label_values(vm_app_version{job=\"$job\", instance=\"$instance\"},  version)",
                "hide": 2,
                "includeAll": false,
                "multi": false,
                "name": "version",
                "options": [],
                "query": {
                    "query": "label_values(vm_app_version{job=\"$job\", instance=\"$instance\"},  version)",
                    "refId": "StandardVariableQuery"
                },
                "refresh": 1,
                "regex": "",
                "skipUrlSync": false,
                "sort": 2,
                "type": "query"
            }
        ]
    },
    "time": {
        "from": "now-15m",
        "to": "now"
    },
    "timepicker": {},
    "timezone": "utc",
    "title": "VictoriaMetrics - operator",
    "uid": "1H179hunk",
    "version": 1,
    "weekStart": ""
}
`,
	},
	ObjectMeta: metav1.ObjectMeta{
		Labels: ku.MergeLabels(
			O.Labels(),
			map[string]string{DashboardLabel: "1"},
		),
		Name:      O.Name,
		Namespace: O.Namespace,
	},
	TypeMeta: ku.TypeConfigMapV1,
}
