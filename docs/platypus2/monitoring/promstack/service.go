// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package promstack

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var CorednsSVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-coredns",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "45.27.2",
			"chart":                        "kube-prometheus-stack-45.27.2",
			"heritage":                     "Helm",
			"jobLabel":                     "coredns",
			"release":                      "kube-prometheus-stack",
		},
		Name:      "kube-prometheus-stack-coredns",
		Namespace: "kube-system",
	},
	Spec: corev1.ServiceSpec{
		ClusterIP: "None",
		Ports: []corev1.ServicePort{
			{
				Name:       "http-metrics",
				Port:       int32(9153),
				Protocol:   corev1.Protocol("TCP"),
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

var KubeControllerManagerSVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-kube-controller-manager",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "45.27.2",
			"chart":                        "kube-prometheus-stack-45.27.2",
			"heritage":                     "Helm",
			"jobLabel":                     "kube-controller-manager",
			"release":                      "kube-prometheus-stack",
		},
		Name:      "kube-prometheus-stack-kube-controller-manager",
		Namespace: "kube-system",
	},
	Spec: corev1.ServiceSpec{
		ClusterIP: "None",
		Ports: []corev1.ServicePort{
			{
				Name:       "http-metrics",
				Port:       int32(10257),
				Protocol:   corev1.Protocol("TCP"),
				TargetPort: intstr.IntOrString{IntVal: int32(10257)},
			},
		},
		Selector: map[string]string{"component": "kube-controller-manager"},
		Type:     corev1.ServiceType("ClusterIP"),
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	},
}

var KubeEtcdSVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-kube-etcd",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "45.27.2",
			"chart":                        "kube-prometheus-stack-45.27.2",
			"heritage":                     "Helm",
			"jobLabel":                     "kube-etcd",
			"release":                      "kube-prometheus-stack",
		},
		Name:      "kube-prometheus-stack-kube-etcd",
		Namespace: "kube-system",
	},
	Spec: corev1.ServiceSpec{
		ClusterIP: "None",
		Ports: []corev1.ServicePort{
			{
				Name:       "http-metrics",
				Port:       int32(2381),
				Protocol:   corev1.Protocol("TCP"),
				TargetPort: intstr.IntOrString{IntVal: int32(2381)},
			},
		},
		Selector: map[string]string{"component": "etcd"},
		Type:     corev1.ServiceType("ClusterIP"),
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	},
}

var KubeProxySVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-kube-proxy",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "45.27.2",
			"chart":                        "kube-prometheus-stack-45.27.2",
			"heritage":                     "Helm",
			"jobLabel":                     "kube-proxy",
			"release":                      "kube-prometheus-stack",
		},
		Name:      "kube-prometheus-stack-kube-proxy",
		Namespace: "kube-system",
	},
	Spec: corev1.ServiceSpec{
		ClusterIP: "None",
		Ports: []corev1.ServicePort{
			{
				Name:       "http-metrics",
				Port:       int32(10249),
				Protocol:   corev1.Protocol("TCP"),
				TargetPort: intstr.IntOrString{IntVal: int32(10249)},
			},
		},
		Selector: map[string]string{"k8s-app": "kube-proxy"},
		Type:     corev1.ServiceType("ClusterIP"),
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	},
}

var KubeSchedulerSVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-kube-scheduler",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "45.27.2",
			"chart":                        "kube-prometheus-stack-45.27.2",
			"heritage":                     "Helm",
			"jobLabel":                     "kube-scheduler",
			"release":                      "kube-prometheus-stack",
		},
		Name:      "kube-prometheus-stack-kube-scheduler",
		Namespace: "kube-system",
	},
	Spec: corev1.ServiceSpec{
		ClusterIP: "None",
		Ports: []corev1.ServicePort{
			{
				Name:       "http-metrics",
				Port:       int32(10259),
				Protocol:   corev1.Protocol("TCP"),
				TargetPort: intstr.IntOrString{IntVal: int32(10259)},
			},
		},
		Selector: map[string]string{"component": "kube-scheduler"},
		Type:     corev1.ServiceType("ClusterIP"),
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	},
}

var OperatorSVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-operator",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "45.27.2",
			"chart":                        "kube-prometheus-stack-45.27.2",
			"heritage":                     "Helm",
			"release":                      "kube-prometheus-stack",
		},
		Name:      "kube-prometheus-stack-operator",
		Namespace: namespace,
	},
	Spec: corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			{
				Name: "https",
				Port: int32(443),
				TargetPort: intstr.IntOrString{
					StrVal: "https",
					Type:   intstr.Type(int64(1)),
				},
			},
		},
		Selector: map[string]string{
			"app":     "kube-prometheus-stack-operator",
			"release": "kube-prometheus-stack",
		},
		Type: corev1.ServiceType("ClusterIP"),
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	},
}

var PrometheusSVC = &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-prometheus",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "45.27.2",
			"chart":                        "kube-prometheus-stack-45.27.2",
			"heritage":                     "Helm",
			"release":                      "kube-prometheus-stack",
			"self-monitor":                 "true",
		},
		Name:      "kube-prometheus-stack-prometheus",
		Namespace: namespace,
	},
	Spec: corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			{
				Name:       "http-web",
				Port:       int32(9090),
				TargetPort: intstr.IntOrString{IntVal: int32(9090)},
			},
		},
		Selector: map[string]string{
			"app.kubernetes.io/name": "prometheus",
			"prometheus":             "kube-prometheus-stack-prometheus",
		},
		Type: corev1.ServiceType("ClusterIP"),
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	},
}
