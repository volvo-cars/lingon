// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmk8s

import (
	"github.com/VictoriaMetrics/operator/api/victoriametrics/v1beta1"
	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	"github.com/volvo-cars/lingoneks/meta"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	KCMPort     = 10257
	KCMPortName = "http-metrics"
)

var KCM = &meta.Metadata{
	Name:      "kube-controller-manager", // linked to the name of the JobLabel
	Namespace: namespace,
	Instance:  "kube-controller-manager-" + namespace,
	Component: "monitoring",
	PartOf:    appName,
	Version:   version,
	ManagedBy: "lingon",
}

type MonKubeController struct {
	kube.App

	SVC        *corev1.Service
	Scrape     *v1beta1.VMServiceScrape
	AlertRules *v1beta1.VMRule
}

func NewMonKubeController() *MonKubeController {
	return &MonKubeController{
		SVC:        KubeControllerSVC,
		AlertRules: KubeControllerAlertRules,
		Scrape:     KubeControllerScrape,
	}
}

var KubeControllerScrape = &v1beta1.VMServiceScrape{
	TypeMeta:   TypeVMServiceScrapeV1Beta1,
	ObjectMeta: KCM.ObjectMeta(),
	Spec: v1beta1.VMServiceScrapeSpec{
		Endpoints: []v1beta1.Endpoint{
			{
				BearerTokenFile: PathSA + "/token",
				Port:            KCMPortName,
				Scheme:          "https",
				TLSConfig: &v1beta1.TLSConfig{
					CAFile:     PathSA + "/ca.crt",
					ServerName: "kubernetes",
				},
			},
		},
		// JobLabel:          "jobLabel",
		JobLabel: ku.AppLabelName,
		NamespaceSelector: v1beta1.NamespaceSelector{
			MatchNames: []string{ku.NSKubeSystem}, // kube-system
		},
		Selector: metav1.LabelSelector{MatchLabels: KCM.MatchLabels()},
	},
}

var KubeControllerSVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Labels:    KCM.Labels(),
		Name:      KCM.Name,
		Namespace: ku.NSKubeSystem, // "kube-system",
	},
	Spec: corev1.ServiceSpec{
		ClusterIP: corev1.ClusterIPNone,
		Ports: []corev1.ServicePort{
			{
				Name:       KCMPortName,
				Port:       int32(KCMPort),
				Protocol:   corev1.ProtocolTCP,
				TargetPort: intstr.FromInt(KCMPort),
			},
		},
		Selector: map[string]string{"component": "kube-controller-manager"},
		Type:     corev1.ServiceTypeClusterIP,
	},
	TypeMeta: ku.TypeServiceV1,
}

var KubeControllerAlertRules = &v1beta1.VMRule{
	ObjectMeta: KCM.ObjectMeta(),
	Spec: v1beta1.VMRuleSpec{
		Groups: []v1beta1.RuleGroup{
			{
				Name: "kubernetes-system-controller-manager",
				Rules: []v1beta1.Rule{
					{
						Alert: "KubeControllerManagerDown",
						Annotations: map[string]string{
							"description": "KubeControllerManager has disappeared from Prometheus target discovery.",
							"runbook_url": "https://runbooks.prometheus-operator.dev/runbooks/kubernetes/kubecontrollermanagerdown",
							"summary":     "Target disappeared from Prometheus target discovery.",
						},
						Expr:   `absent(up{job="` + KCM.Name + `"} == 1)`,
						For:    "15m",
						Labels: map[string]string{"severity": "critical"},
					},
				},
			},
		},
	},
	TypeMeta: TypeVMRuleV1Beta1,
}
