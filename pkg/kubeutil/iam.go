// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SimpleSA(name, namespace string) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta:   TypeServiceAccountV1,
		ObjectMeta: ObjectMeta(name, namespace, nil, nil),
	}
}

func ServiceAccount(
	name, namespace string,
	labels, annotations map[string]string,
) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta:   TypeServiceAccountV1,
		ObjectMeta: ObjectMeta(name, namespace, labels, nil),
	}
}

func Role(
	name, namespace string,
	labels map[string]string,
	rules []rbacv1.PolicyRule,
) *rbacv1.Role {
	return &rbacv1.Role{
		TypeMeta: TypeRoleV1,
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Rules: rules,
	}
}

func ClusterRole(
	name string,
	labels map[string]string,
	rules []rbacv1.PolicyRule,
) *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta:   TypeClusterRoleV1,
		ObjectMeta: ObjectMeta(name, "", labels, nil),
		Rules:      rules,
	}
}
