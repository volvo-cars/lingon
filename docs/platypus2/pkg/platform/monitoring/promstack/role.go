// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package promstack

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var AdmissionRole = &rbacv1.Role{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "pre-install,pre-upgrade,post-install,post-upgrade",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
		},
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-admission",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "45.27.2",
			"chart":                        "kube-prometheus-stack-45.27.2",
			"heritage":                     "Helm",
			"release":                      "kube-prometheus-stack",
		},
		Name:      "kube-prometheus-stack-admission",
		Namespace: namespace,
	},
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{""},
			Resources: []string{"secrets"},
			Verbs:     []string{"get", "create"},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "Role",
	},
}

var GrafanaRole = &rbacv1.Role{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.5.1",
			"helm.sh/chart":                "grafana-6.56.2",
		},
		Name:      "kube-prometheus-stack-grafana",
		Namespace: namespace,
	},
	Rules: []rbacv1.PolicyRule{},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "Role",
	},
}
