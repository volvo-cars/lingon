// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmk8s

import (
	_ "embed"

	"github.com/VictoriaMetrics/operator/api/victoriametrics/v1beta1"
	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	"github.com/volvo-cars/lingoneks/meta"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	CDNSPort     = 9153
	CDNSPortName = "http-metrics"
)

var CDNS = &meta.Metadata{
	Name:      "coredns", // linked to the name of the JobLabel
	Namespace: namespace,
	Instance:  "coredns-" + namespace,
	Component: "monitoring",
	PartOf:    appName,
	Version:   version,
	ManagedBy: "lingon",
}

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
		Labels:    CDNS.Labels(),
		Name:      CDNS.Name,
		Namespace: CDNS.Namespace,
	},
	Spec: corev1.ServiceSpec{
		ClusterIP: corev1.ClusterIPNone,
		Ports: []corev1.ServicePort{
			{
				Name:       CDNSPortName,
				Port:       int32(CDNSPort),
				Protocol:   corev1.ProtocolTCP,
				TargetPort: intstr.FromInt(CDNSPort),
			},
		},
		Selector: map[string]string{"k8s-app": "kube-dns"},
	},
	TypeMeta: ku.TypeServiceV1,
}

var CoreDNSScrape = &v1beta1.VMServiceScrape{
	TypeMeta:   TypeVMServiceScrapeV1Beta1,
	ObjectMeta: CDNS.ObjectMeta(),
	Spec: v1beta1.VMServiceScrapeSpec{
		Endpoints: []v1beta1.Endpoint{
			{
				BearerTokenFile: PathSA + "/token",
				Port:            CDNSPortName,
			},
		},
		NamespaceSelector: v1beta1.NamespaceSelector{
			MatchNames: []string{ku.NSKubeSystem}, // kube-system
		},
		JobLabel: ku.AppLabelName,
		Selector: metav1.LabelSelector{MatchLabels: CDNS.MatchLabels()},
	},
}
