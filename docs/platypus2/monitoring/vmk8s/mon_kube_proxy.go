// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmk8s

import (
	"github.com/VictoriaMetrics/operator/api/victoriametrics/v1beta1"
	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	KPPort     = 10249
	KPPortName = "http-metrics"
)

var KP = &Metadata{
	Name:      "kube-proxy", // linked to the name of the JobLabel
	Namespace: namespace,
	Instance:  "kube-proxy-" + namespace,
	Component: "monitoring",
	PartOf:    appName,
	Version:   version,
	ManagedBy: "lingon",
}

type MonKubeProxy struct {
	kube.App

	KubeProxySVC    *corev1.Service
	KubeProxyScrape *v1beta1.VMServiceScrape
}

func NewMonKubeProxy() *MonKubeProxy {
	return &MonKubeProxy{
		KubeProxySVC:    KubeProxySVC,
		KubeProxyScrape: KubeProxyScrape,
	}
}

var KubeProxySVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Labels:    KP.Labels(),
		Name:      KP.Name,
		Namespace: ku.NSKubeSystem,
	},
	Spec: corev1.ServiceSpec{
		ClusterIP: corev1.ClusterIPNone,
		Ports: []corev1.ServicePort{
			{
				Name:       KPPortName,
				Port:       int32(KPPort),
				Protocol:   corev1.ProtocolTCP,
				TargetPort: intstr.FromInt(KPPort),
			},
		},
		Selector: map[string]string{"k8s-app": "kube-proxy"},
		Type:     corev1.ServiceTypeClusterIP,
	},
	TypeMeta: ku.TypeServiceV1,
}

var KubeProxyScrape = &v1beta1.VMServiceScrape{
	ObjectMeta: KP.ObjectMeta(),
	Spec: v1beta1.VMServiceScrapeSpec{
		Endpoints: []v1beta1.Endpoint{
			{
				BearerTokenFile: PathSA + "/token",
				Port:            KPPortName,
				Scheme:          "https",
				TLSConfig:       &v1beta1.TLSConfig{CAFile: PathSA + "/ca.crt"},
			},
		},
		// JobLabel:          "jobLabel",
		JobLabel: ku.AppLabelName,
		NamespaceSelector: v1beta1.NamespaceSelector{
			MatchNames: []string{ku.NSKubeSystem}, // kube-system
		},
		Selector: metav1.LabelSelector{MatchLabels: KP.MatchLabels()},
	},
	TypeMeta: TypeVMServiceScrapeV1Beta1,
}
