// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package karpenter

import (
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var DnsRole = &rbacv1.Role{
	TypeMeta: kubeutil.TypeRoleV1,
	ObjectMeta: metav1.ObjectMeta{
		Name:      KA.Name + "-dns",
		Namespace: kubeutil.NSKubeSystem,
		Labels:    KA.Labels(),
	},
	Rules: []rbacv1.PolicyRule{
		{
			Verbs:         []string{"get"},
			APIGroups:     []string{""},
			Resources:     []string{"services"},
			ResourceNames: []string{"kube-dns"},
		},
	},
}

var DnsRoleBinding = &rbacv1.RoleBinding{
	TypeMeta: kubeutil.TypeRoleBindingV1,
	ObjectMeta: metav1.ObjectMeta{
		Name:      KA.Name + "-dns",
		Namespace: kubeutil.NSKubeSystem,
		Labels:    KA.Labels(),
	},
	Subjects: kubeutil.RoleSubject(KA.Name, KA.Namespace),
	RoleRef:  kubeutil.RoleRef(DnsRole.Name),
}

var Role = &rbacv1.Role{
	TypeMeta:   kubeutil.TypeRoleV1,
	ObjectMeta: KA.ObjectMeta(),
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{"coordination.k8s.io"},
			Resources: []string{"leases"},
			Verbs:     []string{"get", "watch"},
		},
		{
			APIGroups: []string{""},
			Resources: []string{"configmaps", "namespaces", "secrets"},
			Verbs:     []string{"get", "list", "watch"},
		},
		{
			APIGroups:     []string{""},
			Resources:     []string{"secrets"},
			Verbs:         []string{"update"},
			ResourceNames: []string{CertSecret.Name},
		},
		{
			APIGroups: []string{""},
			Resources: []string{"configmaps"},
			Verbs:     []string{"update", "patch", "delete"},
			ResourceNames: []string{
				KA.ConfigName,
				LoggingConfig.Name,
			},
		},
		{
			APIGroups: []string{"coordination.k8s.io"},
			Resources: []string{"leases"},
			Verbs:     []string{"patch", "update"},
			ResourceNames: []string{
				"karpenter-leader-election",
				"webhook.configmapwebhook.00-of-01",
				"webhook.defaultingwebhook.00-of-01",
				"webhook.validationwebhook.00-of-01",
				"webhook.webhookcertificates.00-of-01",
			},
		},
		{
			APIGroups: []string{"coordination.k8s.io"},
			Resources: []string{"leases"},
			Verbs:     []string{"create"},
		},
		{
			APIGroups: []string{""},
			Resources: []string{"configmaps"},
			Verbs:     []string{"create"},
		},
	},
}
