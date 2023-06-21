// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmk8s

import (
	"github.com/VictoriaMetrics/operator/api/victoriametrics/v1beta1"
	"github.com/volvo-cars/lingon/pkg/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type MonCoreDNS struct {
	kube.App

	SVC         *corev1.Service
	Scrape      *v1beta1.VMServiceScrape
	DashboardCM *corev1.ConfigMap
}

func NewMonCoreDNS() *MonCoreDNS {
	return &MonCoreDNS{
		SVC:         CoreDNSSVC,
		Scrape:      CoreDNSScrape,
		DashboardCM: CoreDNSDashboardCM,
	}
}

var CoreDNSSVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "vmk8s-victoria-metrics-k8s-stack-coredns",
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
			"jobLabel":                     "coredns",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-coredns",
		Namespace: "kube-system",
	},
	Spec: corev1.ServiceSpec{
		ClusterIP: "None",
		Ports: []corev1.ServicePort{
			{
				Name:       "http-metrics",
				Port:       int32(9153),
				Protocol:   corev1.ProtocolTCP,
				TargetPort: intstr.IntOrString{IntVal: int32(9153)},
			},
		},
		Selector: map[string]string{"k8s-app": "kube-dns"},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	},
}

var CoreDNSScrape = &v1beta1.VMServiceScrape{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-coredns",
		Namespace: "monitoring",
	},
	Spec: v1beta1.VMServiceScrapeSpec{
		Endpoints: []v1beta1.Endpoint{
			{
				BearerTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
				Port:            "http-metrics",
			},
		},
		NamespaceSelector: v1beta1.NamespaceSelector{MatchNames: []string{"kube-system"}},
		Selector: metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app":                        "vmk8s-victoria-metrics-k8s-stack-coredns",
				"app.kubernetes.io/instance": "vmk8s",
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "operator.victoriametrics.com/v1beta1",
		Kind:       "VMServiceScrape",
	},
}

var CoreDNSDashboardCM = &corev1.ConfigMap{
	Data: map[string]string{
		"k8s-system-coredns.json": `
{
    "__inputs": [
        {
            "name": "DS_PROMETHEUS",
            "label": "Prometheus",
            "description": "",
            "type": "datasource",
            "pluginId": "prometheus",
            "pluginName": "Prometheus"
        }
    ],
    "__elements": [],
    "__requires": [
        {
            "type": "grafana",
            "id": "grafana",
            "name": "Grafana",
            "version": "8.4.4"
        },
        {
            "type": "datasource",
            "id": "prometheus",
            "name": "Prometheus",
            "version": "5.0.0"
        },
        {
            "type": "panel",
            "id": "timeseries",
            "name": "Time series",
            "version": ""
        },
        {
            "type": "panel",
            "id": "stat",
            "name": "Stat",
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
    "description": "This is a modern CoreDNS dashboard for your Kubernetes cluster(s). Made for kube-prometheus-stack and take advantage of the latest Grafana features. GitHub repository: https://github.com/dotdc/grafana-dashboards-kubernetes",
    "editable": true,
    "fiscalYearStartMonth": 0,
    "graphTooltip": 1,
    "links": [],
    "liveNow": false,
    "panels": [
        {
            "datasource": {
                "type": "prometheus",
                "uid": "${datasource}"
            },
            "fieldConfig": {
                "defaults": {
                    "mappings": [
                        {
                            "options": {
                                "0": {
                                    "text": "DOWN"
                                },
                                "1": {
                                    "text": "UP"
                                }
                            },
                            "type": "value"
                        }
                    ],
                    "thresholds": {
                        "mode": "absolute",
                        "steps": [
                            {
                                "color": "red",
                                "value": null
                            },
                            {
                                "color": "green",
                                "value": 1
                            }
                        ]
                    }
                },
                "overrides": []
            },
            "gridPos": {
                "h": 3,
                "w": 24,
                "x": 0,
                "y": 0
            },
            "id": 25,
            "options": {
                "colorMode": "background",
                "graphMode": "none",
                "justifyMode": "auto",
                "orientation": "vertical",
                "reduceOptions": {
                    "calcs": [
                        "lastNotNull"
                    ],
                    "fields": "",
                    "values": false
                },
                "textMode": "value_and_name"
            },
            "pluginVersion": "9.2.4",
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "${datasource}"
                    },
                    "exemplar": true,
                    "expr": "up{job=\"coredns\", instance=~\"$instance\"}",
                    "interval": "",
                    "legendFormat": "{{ instance }}",
                    "refId": "A"
                }
            ],
            "title": "CoreDNS - Health Status",
            "type": "stat"
        },
        {
            "datasource": {
                "type": "prometheus",
                "uid": "${datasource}"
            },
            "fieldConfig": {
                "defaults": {
                    "color": {
                        "mode": "palette-classic"
                    },
                    "custom": {
                        "axisCenteredZero": false,
                        "axisColorMode": "text",
                        "axisLabel": "",
                        "axisPlacement": "auto",
                        "barAlignment": 0,
                        "drawStyle": "line",
                        "fillOpacity": 25,
                        "gradientMode": "opacity",
                        "hideFrom": {
                            "legend": false,
                            "tooltip": false,
                            "viz": false
                        },
                        "lineInterpolation": "smooth",
                        "lineWidth": 2,
                        "pointSize": 5,
                        "scaleDistribution": {
                            "type": "linear"
                        },
                        "showPoints": "never",
                        "spanNulls": true,
                        "stacking": {
                            "group": "A",
                            "mode": "none"
                        },
                        "thresholdsStyle": {
                            "mode": "off"
                        }
                    },
                    "decimals": 2,
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
                    "unit": "percentunit"
                },
                "overrides": []
            },
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 0,
                "y": 3
            },
            "id": 19,
            "options": {
                "legend": {
                    "calcs": [],
                    "displayMode": "list",
                    "placement": "bottom",
                    "showLegend": true
                },
                "tooltip": {
                    "mode": "single",
                    "sort": "none"
                }
            },
            "pluginVersion": "8.3.3",
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "${datasource}"
                    },
                    "exemplar": true,
                    "expr": "rate(process_cpu_seconds_total{job=\"coredns\", instance=~\"$instance\"}[$__rate_interval])",
                    "interval": "$resolution",
                    "legendFormat": "{{ instance }}",
                    "refId": "A"
                }
            ],
            "title": "CoreDNS - CPU Usage by instance",
            "type": "timeseries"
        },
        {
            "datasource": {
                "type": "prometheus",
                "uid": "${datasource}"
            },
            "fieldConfig": {
                "defaults": {
                    "color": {
                        "mode": "palette-classic"
                    },
                    "custom": {
                        "axisCenteredZero": false,
                        "axisColorMode": "text",
                        "axisLabel": "",
                        "axisPlacement": "auto",
                        "barAlignment": 0,
                        "drawStyle": "line",
                        "fillOpacity": 25,
                        "gradientMode": "opacity",
                        "hideFrom": {
                            "legend": false,
                            "tooltip": false,
                            "viz": false
                        },
                        "lineInterpolation": "smooth",
                        "lineWidth": 2,
                        "pointSize": 5,
                        "scaleDistribution": {
                            "type": "linear"
                        },
                        "showPoints": "never",
                        "spanNulls": true,
                        "stacking": {
                            "group": "A",
                            "mode": "none"
                        },
                        "thresholdsStyle": {
                            "mode": "off"
                        }
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
                    "unit": "bytes"
                },
                "overrides": []
            },
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 12,
                "y": 3
            },
            "id": 21,
            "options": {
                "legend": {
                    "calcs": [],
                    "displayMode": "list",
                    "placement": "bottom",
                    "showLegend": true
                },
                "tooltip": {
                    "mode": "single",
                    "sort": "none"
                }
            },
            "pluginVersion": "8.3.3",
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "${datasource}"
                    },
                    "exemplar": true,
                    "expr": "process_resident_memory_bytes{job=\"coredns\", instance=~\"$instance\"}",
                    "interval": "",
                    "legendFormat": "{{ instance }}",
                    "refId": "A"
                }
            ],
            "title": "CoreDNS - Memory Usage by instance",
            "type": "timeseries"
        },
        {
            "datasource": {
                "type": "prometheus",
                "uid": "${datasource}"
            },
            "fieldConfig": {
                "defaults": {
                    "color": {
                        "mode": "palette-classic"
                    },
                    "custom": {
                        "axisCenteredZero": false,
                        "axisColorMode": "text",
                        "axisLabel": "",
                        "axisPlacement": "auto",
                        "barAlignment": 0,
                        "drawStyle": "line",
                        "fillOpacity": 25,
                        "gradientMode": "opacity",
                        "hideFrom": {
                            "legend": false,
                            "tooltip": false,
                            "viz": false
                        },
                        "lineInterpolation": "smooth",
                        "lineWidth": 2,
                        "pointSize": 5,
                        "scaleDistribution": {
                            "type": "linear"
                        },
                        "showPoints": "never",
                        "spanNulls": true,
                        "stacking": {
                            "group": "A",
                            "mode": "none"
                        },
                        "thresholdsStyle": {
                            "mode": "off"
                        }
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
                    "unit": "short"
                },
                "overrides": []
            },
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 0,
                "y": 11
            },
            "id": 9,
            "options": {
                "legend": {
                    "calcs": [],
                    "displayMode": "list",
                    "placement": "bottom",
                    "showLegend": false
                },
                "tooltip": {
                    "mode": "multi",
                    "sort": "none"
                }
            },
            "pluginVersion": "8.3.3",
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "${datasource}"
                    },
                    "exemplar": true,
                    "expr": "sum(rate(coredns_dns_requests_total{instance=~\"$instance\",proto=\"$protocol\"}[$__rate_interval]))",
                    "interval": "$resolution",
                    "legendFormat": "total $protocol requests",
                    "refId": "A"
                }
            ],
            "title": "CoreDNS - Total DNS Requests ($protocol)",
            "type": "timeseries"
        },
        {
            "datasource": {
                "type": "prometheus",
                "uid": "${datasource}"
            },
            "fieldConfig": {
                "defaults": {
                    "color": {
                        "mode": "palette-classic"
                    },
                    "custom": {
                        "axisCenteredZero": false,
                        "axisColorMode": "text",
                        "axisLabel": "",
                        "axisPlacement": "auto",
                        "barAlignment": 0,
                        "drawStyle": "line",
                        "fillOpacity": 25,
                        "gradientMode": "opacity",
                        "hideFrom": {
                            "legend": false,
                            "tooltip": false,
                            "viz": false
                        },
                        "lineInterpolation": "smooth",
                        "lineWidth": 2,
                        "pointSize": 5,
                        "scaleDistribution": {
                            "type": "linear"
                        },
                        "showPoints": "never",
                        "spanNulls": true,
                        "stacking": {
                            "group": "A",
                            "mode": "none"
                        },
                        "thresholdsStyle": {
                            "mode": "off"
                        }
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
                    "unit": "bytes"
                },
                "overrides": []
            },
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 12,
                "y": 11
            },
            "id": 7,
            "options": {
                "legend": {
                    "calcs": [],
                    "displayMode": "list",
                    "placement": "bottom",
                    "showLegend": false
                },
                "tooltip": {
                    "mode": "multi",
                    "sort": "none"
                }
            },
            "pluginVersion": "8.3.3",
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "${datasource}"
                    },
                    "exemplar": true,
                    "expr": "avg(rate(coredns_dns_request_size_bytes_sum{instance=~\"$instance\",proto=\"$protocol\"}[$__rate_interval])) by (proto)",
                    "interval": "$resolution",
                    "legendFormat": "average $protocol packet size",
                    "refId": "A"
                }
            ],
            "title": "CoreDNS - Average Packet Size  ($protocol)",
            "type": "timeseries"
        },
        {
            "datasource": {
                "type": "prometheus",
                "uid": "${datasource}"
            },
            "fieldConfig": {
                "defaults": {
                    "color": {
                        "mode": "palette-classic"
                    },
                    "custom": {
                        "axisCenteredZero": false,
                        "axisColorMode": "text",
                        "axisLabel": "",
                        "axisPlacement": "auto",
                        "barAlignment": 0,
                        "drawStyle": "line",
                        "fillOpacity": 25,
                        "gradientMode": "opacity",
                        "hideFrom": {
                            "legend": false,
                            "tooltip": false,
                            "viz": false
                        },
                        "lineInterpolation": "smooth",
                        "lineWidth": 2,
                        "pointSize": 5,
                        "scaleDistribution": {
                            "type": "linear"
                        },
                        "showPoints": "never",
                        "spanNulls": true,
                        "stacking": {
                            "group": "A",
                            "mode": "none"
                        },
                        "thresholdsStyle": {
                            "mode": "off"
                        }
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
                    "unit": "short"
                },
                "overrides": []
            },
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 0,
                "y": 19
            },
            "id": 2,
            "options": {
                "legend": {
                    "calcs": [],
                    "displayMode": "list",
                    "placement": "bottom",
                    "showLegend": true
                },
                "tooltip": {
                    "mode": "multi",
                    "sort": "desc"
                }
            },
            "pluginVersion": "8.3.3",
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "${datasource}"
                    },
                    "exemplar": true,
                    "expr": "sum(rate(coredns_dns_requests_total{instance=~\"$instance\"}[$__rate_interval])) by (type)",
                    "interval": "$resolution",
                    "legendFormat": "{{ type }}",
                    "refId": "A"
                }
            ],
            "title": "CoreDNS - Requests by type",
            "type": "timeseries"
        },
        {
            "datasource": {
                "type": "prometheus",
                "uid": "${datasource}"
            },
            "fieldConfig": {
                "defaults": {
                    "color": {
                        "mode": "palette-classic"
                    },
                    "custom": {
                        "axisCenteredZero": false,
                        "axisColorMode": "text",
                        "axisLabel": "",
                        "axisPlacement": "auto",
                        "barAlignment": 0,
                        "drawStyle": "line",
                        "fillOpacity": 25,
                        "gradientMode": "opacity",
                        "hideFrom": {
                            "legend": false,
                            "tooltip": false,
                            "viz": false
                        },
                        "lineInterpolation": "smooth",
                        "lineWidth": 2,
                        "pointSize": 5,
                        "scaleDistribution": {
                            "type": "linear"
                        },
                        "showPoints": "never",
                        "spanNulls": true,
                        "stacking": {
                            "group": "A",
                            "mode": "none"
                        },
                        "thresholdsStyle": {
                            "mode": "off"
                        }
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
                    "unit": "short"
                },
                "overrides": []
            },
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 12,
                "y": 19
            },
            "id": 4,
            "options": {
                "legend": {
                    "calcs": [],
                    "displayMode": "list",
                    "placement": "bottom",
                    "showLegend": true
                },
                "tooltip": {
                    "mode": "multi",
                    "sort": "desc"
                }
            },
            "pluginVersion": "8.3.3",
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "${datasource}"
                    },
                    "exemplar": true,
                    "expr": "sum(rate(coredns_dns_responses_total{instance=~\"$instance\"}[$__rate_interval])) by (rcode)",
                    "interval": "$resolution",
                    "legendFormat": "{{ rcode }}",
                    "refId": "A"
                }
            ],
            "title": "CoreDNS - Requests by return code",
            "type": "timeseries"
        },
        {
            "datasource": {
                "type": "prometheus",
                "uid": "${datasource}"
            },
            "fieldConfig": {
                "defaults": {
                    "color": {
                        "mode": "palette-classic"
                    },
                    "custom": {
                        "axisCenteredZero": false,
                        "axisColorMode": "text",
                        "axisLabel": "",
                        "axisPlacement": "auto",
                        "barAlignment": 0,
                        "drawStyle": "line",
                        "fillOpacity": 25,
                        "gradientMode": "opacity",
                        "hideFrom": {
                            "legend": false,
                            "tooltip": false,
                            "viz": false
                        },
                        "lineInterpolation": "smooth",
                        "lineWidth": 2,
                        "pointSize": 5,
                        "scaleDistribution": {
                            "type": "linear"
                        },
                        "showPoints": "never",
                        "spanNulls": true,
                        "stacking": {
                            "group": "A",
                            "mode": "none"
                        },
                        "thresholdsStyle": {
                            "mode": "off"
                        }
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
                    "unit": "short"
                },
                "overrides": []
            },
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 0,
                "y": 27
            },
            "id": 23,
            "options": {
                "legend": {
                    "calcs": [],
                    "displayMode": "list",
                    "placement": "bottom",
                    "showLegend": false
                },
                "tooltip": {
                    "mode": "multi",
                    "sort": "none"
                }
            },
            "pluginVersion": "8.3.3",
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "${datasource}"
                    },
                    "exemplar": true,
                    "expr": "sum(rate(coredns_forward_requests_total[$__rate_interval]))",
                    "interval": "$resolution",
                    "legendFormat": "total forward requests",
                    "refId": "A"
                }
            ],
            "title": "CoreDNS - Total Forward Requests",
            "type": "timeseries"
        },
        {
            "datasource": {
                "type": "prometheus",
                "uid": "${datasource}"
            },
            "fieldConfig": {
                "defaults": {
                    "color": {
                        "mode": "palette-classic"
                    },
                    "custom": {
                        "axisCenteredZero": false,
                        "axisColorMode": "text",
                        "axisLabel": "",
                        "axisPlacement": "auto",
                        "barAlignment": 0,
                        "drawStyle": "line",
                        "fillOpacity": 25,
                        "gradientMode": "opacity",
                        "hideFrom": {
                            "legend": false,
                            "tooltip": false,
                            "viz": false
                        },
                        "lineInterpolation": "smooth",
                        "lineWidth": 2,
                        "pointSize": 5,
                        "scaleDistribution": {
                            "type": "linear"
                        },
                        "showPoints": "never",
                        "spanNulls": true,
                        "stacking": {
                            "group": "A",
                            "mode": "none"
                        },
                        "thresholdsStyle": {
                            "mode": "off"
                        }
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
                    "unit": "short"
                },
                "overrides": []
            },
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 12,
                "y": 27
            },
            "id": 13,
            "options": {
                "legend": {
                    "calcs": [],
                    "displayMode": "list",
                    "placement": "bottom",
                    "showLegend": true
                },
                "tooltip": {
                    "mode": "multi",
                    "sort": "none"
                }
            },
            "pluginVersion": "8.3.3",
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "${datasource}"
                    },
                    "exemplar": true,
                    "expr": "sum(rate(coredns_forward_responses_total{rcode=~\"SERVFAIL|REFUSED\"}[$__rate_interval])) by (rcode)",
                    "interval": "$resolution",
                    "legendFormat": "{{ rcode }}",
                    "refId": "A"
                }
            ],
            "title": "CoreDNS - DNS Errors",
            "type": "timeseries"
        },
        {
            "datasource": {
                "type": "prometheus",
                "uid": "${datasource}"
            },
            "fieldConfig": {
                "defaults": {
                    "color": {
                        "mode": "palette-classic"
                    },
                    "custom": {
                        "axisCenteredZero": false,
                        "axisColorMode": "text",
                        "axisLabel": "",
                        "axisPlacement": "auto",
                        "barAlignment": 0,
                        "drawStyle": "line",
                        "fillOpacity": 25,
                        "gradientMode": "opacity",
                        "hideFrom": {
                            "legend": false,
                            "tooltip": false,
                            "viz": false
                        },
                        "lineInterpolation": "smooth",
                        "lineWidth": 2,
                        "pointSize": 5,
                        "scaleDistribution": {
                            "type": "linear"
                        },
                        "showPoints": "never",
                        "spanNulls": true,
                        "stacking": {
                            "group": "A",
                            "mode": "none"
                        },
                        "thresholdsStyle": {
                            "mode": "off"
                        }
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
                    "unit": "short"
                },
                "overrides": []
            },
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 0,
                "y": 35
            },
            "id": 17,
            "options": {
                "legend": {
                    "calcs": [],
                    "displayMode": "list",
                    "placement": "bottom",
                    "showLegend": true
                },
                "tooltip": {
                    "mode": "multi",
                    "sort": "desc"
                }
            },
            "pluginVersion": "8.3.3",
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "${datasource}"
                    },
                    "exemplar": true,
                    "expr": "sum(rate(coredns_cache_hits_total{instance=~\"$instance\"}[$__rate_interval])) by (type)",
                    "interval": "$resolution",
                    "legendFormat": "{{ type }}",
                    "refId": "A"
                },
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "${datasource}"
                    },
                    "exemplar": true,
                    "expr": "sum(rate(coredns_cache_misses_total{instance=~\"$instance\"}[$__rate_interval])) by (type)",
                    "interval": "$resolution",
                    "legendFormat": "misses",
                    "refId": "B"
                }
            ],
            "title": "CoreDNS - Cache Hits / Misses",
            "type": "timeseries"
        },
        {
            "datasource": {
                "type": "prometheus",
                "uid": "${datasource}"
            },
            "fieldConfig": {
                "defaults": {
                    "color": {
                        "mode": "palette-classic"
                    },
                    "custom": {
                        "axisCenteredZero": false,
                        "axisColorMode": "text",
                        "axisLabel": "",
                        "axisPlacement": "auto",
                        "barAlignment": 0,
                        "drawStyle": "line",
                        "fillOpacity": 25,
                        "gradientMode": "opacity",
                        "hideFrom": {
                            "legend": false,
                            "tooltip": false,
                            "viz": false
                        },
                        "lineInterpolation": "smooth",
                        "lineWidth": 2,
                        "pointSize": 5,
                        "scaleDistribution": {
                            "type": "linear"
                        },
                        "showPoints": "never",
                        "spanNulls": true,
                        "stacking": {
                            "group": "A",
                            "mode": "none"
                        },
                        "thresholdsStyle": {
                            "mode": "off"
                        }
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
                    "unit": "bytes"
                },
                "overrides": []
            },
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 12,
                "y": 35
            },
            "id": 15,
            "options": {
                "legend": {
                    "calcs": [],
                    "displayMode": "list",
                    "placement": "bottom",
                    "showLegend": true
                },
                "tooltip": {
                    "mode": "multi",
                    "sort": "desc"
                }
            },
            "pluginVersion": "8.3.3",
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "${datasource}"
                    },
                    "exemplar": true,
                    "expr": "sum(coredns_cache_entries) by (type)",
                    "interval": "",
                    "legendFormat": "{{ type }}",
                    "refId": "A"
                }
            ],
            "title": "CoreDNS - Cache Size",
            "type": "timeseries"
        },
        {
            "datasource": {
                "type": "prometheus",
                "uid": "${datasource}"
            },
            "fieldConfig": {
                "defaults": {
                    "custom": {
                        "hideFrom": {
                            "legend": false,
                            "tooltip": false,
                            "viz": false
                        },
                        "scaleDistribution": {
                            "type": "linear"
                        }
                    }
                },
                "overrides": []
            },
            "gridPos": {
                "h": 10,
                "w": 12,
                "x": 0,
                "y": 43
            },
            "id": 27,
            "options": {
                "calculate": false,
                "cellGap": 1,
                "color": {
                    "exponent": 0.5,
                    "fill": "dark-orange",
                    "mode": "scheme",
                    "reverse": false,
                    "scale": "exponential",
                    "scheme": "RdYlBu",
                    "steps": 64
                },
                "exemplars": {
                    "color": "rgba(255,0,255,0.7)"
                },
                "filterValues": {
                    "le": 1e-09
                },
                "legend": {
                    "show": true
                },
                "rowsFrame": {
                    "layout": "auto"
                },
                "tooltip": {
                    "show": true,
                    "yHistogram": false
                },
                "yAxis": {
                    "axisPlacement": "left",
                    "reverse": false,
                    "unit": "s"
                }
            },
            "pluginVersion": "9.3.1",
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "${datasource}"
                    },
                    "editorMode": "code",
                    "expr": "sum(increase(coredns_dns_request_duration_seconds_bucket{instance=~\"$instance\"}[$__rate_interval])) by (le)",
                    "format": "heatmap",
                    "legendFormat": "{{le}}",
                    "range": true,
                    "refId": "A"
                }
            ],
            "title": "CoreDNS - DNS request duration",
            "type": "heatmap"
        },
        {
            "datasource": {
                "type": "prometheus",
                "uid": "${datasource}"
            },
            "fieldConfig": {
                "defaults": {
                    "custom": {
                        "hideFrom": {
                            "legend": false,
                            "tooltip": false,
                            "viz": false
                        },
                        "scaleDistribution": {
                            "type": "linear"
                        }
                    }
                },
                "overrides": []
            },
            "gridPos": {
                "h": 10,
                "w": 12,
                "x": 12,
                "y": 43
            },
            "id": 30,
            "options": {
                "calculate": false,
                "cellGap": 1,
                "color": {
                    "exponent": 0.5,
                    "fill": "dark-orange",
                    "mode": "scheme",
                    "reverse": false,
                    "scale": "exponential",
                    "scheme": "RdYlBu",
                    "steps": 64
                },
                "exemplars": {
                    "color": "rgba(255,0,255,0.7)"
                },
                "filterValues": {
                    "le": 1e-09
                },
                "legend": {
                    "show": true
                },
                "rowsFrame": {
                    "layout": "auto"
                },
                "tooltip": {
                    "show": true,
                    "yHistogram": false
                },
                "yAxis": {
                    "axisPlacement": "left",
                    "reverse": false,
                    "unit": "s"
                }
            },
            "pluginVersion": "9.3.1",
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "${datasource}"
                    },
                    "editorMode": "code",
                    "expr": "sum(increase(coredns_forward_request_duration_seconds_bucket{instance=~\"$instance\"}[$__rate_interval])) by (le)",
                    "format": "heatmap",
                    "legendFormat": "{{le}}",
                    "range": true,
                    "refId": "A"
                }
            ],
            "title": "CoreDNS - Forward request duration",
            "type": "heatmap"
        },
        {
            "datasource": {
                "type": "prometheus",
                "uid": "${datasource}"
            },
            "fieldConfig": {
                "defaults": {
                    "custom": {
                        "hideFrom": {
                            "legend": false,
                            "tooltip": false,
                            "viz": false
                        },
                        "scaleDistribution": {
                            "type": "linear"
                        }
                    }
                },
                "overrides": []
            },
            "gridPos": {
                "h": 10,
                "w": 12,
                "x": 0,
                "y": 53
            },
            "id": 28,
            "options": {
                "calculate": false,
                "cellGap": 1,
                "color": {
                    "exponent": 0.5,
                    "fill": "dark-orange",
                    "mode": "scheme",
                    "reverse": false,
                    "scale": "exponential",
                    "scheme": "RdYlBu",
                    "steps": 64
                },
                "exemplars": {
                    "color": "rgba(255,0,255,0.7)"
                },
                "filterValues": {
                    "le": 1e-09
                },
                "legend": {
                    "show": true
                },
                "rowsFrame": {
                    "layout": "auto"
                },
                "tooltip": {
                    "show": true,
                    "yHistogram": false
                },
                "yAxis": {
                    "axisPlacement": "left",
                    "reverse": false,
                    "unit": "decbytes"
                }
            },
            "pluginVersion": "9.3.1",
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "${datasource}"
                    },
                    "editorMode": "code",
                    "expr": "sum(increase(coredns_dns_request_size_bytes_bucket{instance=~\"$instance\", le!=\"0\"}[$__rate_interval])) by (le)",
                    "format": "heatmap",
                    "legendFormat": "{{le}}",
                    "range": true,
                    "refId": "A"
                }
            ],
            "title": "CoreDNS - DNS request size",
            "type": "heatmap"
        },
        {
            "datasource": {
                "type": "prometheus",
                "uid": "${datasource}"
            },
            "fieldConfig": {
                "defaults": {
                    "custom": {
                        "hideFrom": {
                            "legend": false,
                            "tooltip": false,
                            "viz": false
                        },
                        "scaleDistribution": {
                            "type": "linear"
                        }
                    }
                },
                "overrides": []
            },
            "gridPos": {
                "h": 10,
                "w": 12,
                "x": 12,
                "y": 53
            },
            "id": 29,
            "options": {
                "calculate": false,
                "cellGap": 1,
                "color": {
                    "exponent": 0.5,
                    "fill": "dark-orange",
                    "mode": "scheme",
                    "reverse": false,
                    "scale": "exponential",
                    "scheme": "RdYlBu",
                    "steps": 64
                },
                "exemplars": {
                    "color": "rgba(255,0,255,0.7)"
                },
                "filterValues": {
                    "le": 1e-09
                },
                "legend": {
                    "show": true
                },
                "rowsFrame": {
                    "layout": "auto"
                },
                "tooltip": {
                    "show": true,
                    "yHistogram": false
                },
                "yAxis": {
                    "axisPlacement": "left",
                    "reverse": false,
                    "unit": "decbytes"
                }
            },
            "pluginVersion": "9.3.1",
            "targets": [
                {
                    "datasource": {
                        "type": "prometheus",
                        "uid": "${datasource}"
                    },
                    "editorMode": "code",
                    "expr": "sum(increase(coredns_dns_response_size_bytes_bucket{instance=~\"$instance\", le!=\"0\"}[$__rate_interval])) by (le)",
                    "format": "heatmap",
                    "legendFormat": "{{le}}",
                    "range": true,
                    "refId": "A"
                }
            ],
            "title": "CoreDNS - DNS response size",
            "type": "heatmap"
        }
    ],
    "refresh": "30s",
    "schemaVersion": 37,
    "style": "dark",
    "tags": [
        "Kubernetes",
        "Prometheus"
    ],
    "templating": {
        "list": [
            {
                "current": {
                    "selected": false,
                    "text": "Prometheus",
                    "value": "Prometheus"
                },
                "hide": 0,
                "includeAll": false,
                "multi": false,
                "name": "datasource",
                "options": [],
                "query": "prometheus",
                "queryValue": "",
                "refresh": 1,
                "regex": "",
                "skipUrlSync": false,
                "type": "datasource"
            },
            {
                "current": {
                    "selected": false,
                    "text": "All",
                    "value": "$__all"
                },
                "datasource": {
                    "type": "prometheus",
                    "uid": "${datasource}"
                },
                "definition": "label_values(up{job=\"coredns\"}, instance)",
                "hide": 0,
                "includeAll": true,
                "label": "",
                "multi": false,
                "name": "instance",
                "options": [],
                "query": {
                    "query": "label_values(up{job=\"coredns\"}, instance)",
                    "refId": "StandardVariableQuery"
                },
                "refresh": 1,
                "regex": "",
                "skipUrlSync": false,
                "sort": 1,
                "tagValuesQuery": "",
                "tagsQuery": "",
                "type": "query",
                "useTags": false
            },
            {
                "allValue": "udp,tcp",
                "current": {
                    "selected": false,
                    "text": "udp",
                    "value": "udp"
                },
                "datasource": {
                    "type": "prometheus",
                    "uid": "${datasource}"
                },
                "definition": "label_values(coredns_dns_requests_total, proto)",
                "hide": 0,
                "includeAll": false,
                "label": "",
                "multi": false,
                "name": "protocol",
                "options": [],
                "query": {
                    "query": "label_values(coredns_dns_requests_total, proto)",
                    "refId": "StandardVariableQuery"
                },
                "refresh": 1,
                "regex": "",
                "skipUrlSync": false,
                "sort": 1,
                "tagValuesQuery": "",
                "tagsQuery": "",
                "type": "query",
                "useTags": false
            },
            {
                "current": {
                    "selected": true,
                    "text": "30s",
                    "value": "30s"
                },
                "hide": 0,
                "includeAll": false,
                "multi": false,
                "name": "resolution",
                "options": [
                    {
                        "selected": false,
                        "text": "1s",
                        "value": "1s"
                    },
                    {
                        "selected": false,
                        "text": "15s",
                        "value": "15s"
                    },
                    {
                        "selected": true,
                        "text": "30s",
                        "value": "30s"
                    },
                    {
                        "selected": false,
                        "text": "1m",
                        "value": "1m"
                    },
                    {
                        "selected": false,
                        "text": "3m",
                        "value": "3m"
                    },
                    {
                        "selected": false,
                        "text": "5m",
                        "value": "5m"
                    }
                ],
                "query": "1s, 15s, 30s, 1m, 3m, 5m",
                "queryValue": "",
                "skipUrlSync": false,
                "type": "custom"
            }
        ]
    },
    "time": {
        "from": "now-1h",
        "to": "now"
    },
    "timepicker": {},
    "timezone": "utc",
    "title": "Kubernetes / System / CoreDNS",
    "uid": "k8s_system_coredns",
    "version": 11,
    "weekStart": ""
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
		Name:      "vmk8s-victoria-metrics-k8s-stack-k8s-system-coredns",
		Namespace: "monitoring",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ConfigMap",
	},
}
