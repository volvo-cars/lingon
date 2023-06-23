// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package vmop

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var VictoriaMetricsOperatorRB = &rbacv1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmop",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-operator",
			"helm.sh/chart":                "victoria-metrics-operator-0.23.1",
		},
		Name:      "vmop-victoria-metrics-operator",
		Namespace: "monitoring",
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "Role",
		Name:     "vmop-victoria-metrics-operator",
	},
	Subjects: []rbacv1.Subject{{
		Kind:      "ServiceAccount",
		Name:      "vmop-victoria-metrics-operator",
		Namespace: "monitoring",
	}},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "RoleBinding",
	},
}
