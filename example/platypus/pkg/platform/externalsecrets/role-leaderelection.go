// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package externalsecrets

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var LeaderElectionRole = &rbacv1.Role{
	ObjectMeta: metav1.ObjectMeta{
		Labels:    ESLabels,
		Name:      "external-secrets-leaderelection",
		Namespace: "external-secrets",
	},
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups:     []string{},
			ResourceNames: []string{"external-secrets-controller"},
			Resources:     []string{"configmaps"},
			Verbs:         []string{"get", "update", "patch"},
		}, {
			APIGroups: []string{},
			Resources: []string{"configmaps"},
			Verbs:     []string{"create"},
		}, {
			APIGroups: []string{"coordination.k8s.io"},
			Resources: []string{"leases"},
			Verbs:     []string{"get", "create", "update", "patch"},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "Role",
	},
}
