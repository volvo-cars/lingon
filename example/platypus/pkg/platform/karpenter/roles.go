package karpenter

import (
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var DnsRole = &rbacv1.Role{
	TypeMeta: kubeutil.TypeRoleV1,
	ObjectMeta: metav1.ObjectMeta{
		Name:      "karpenter-dns",
		Namespace: "kube-system",
		Labels:    commonLabels,
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
		Name:      "karpenter-dns",
		Namespace: "kube-system",
		Labels:    commonLabels,
	},
	Subjects: kubeutil.RoleSubject("karpenter", "karpenter"),
	RoleRef:  kubeutil.RoleRef("karpenter-dns"),
}

var Role = &rbacv1.Role{
	TypeMeta: kubeutil.TypeRoleV1,
	ObjectMeta: metav1.ObjectMeta{
		Name:      "karpenter",
		Namespace: "karpenter",
		Labels:    commonLabels,
	},
	Rules: []rbacv1.PolicyRule{
		{
			Verbs:     []string{"get", "watch"},
			APIGroups: []string{"coordination.k8s.io"},
			Resources: []string{"leases"},
		},
		{
			Verbs:     []string{"get", "list", "watch"},
			APIGroups: []string{""},
			Resources: []string{"configmaps", "namespaces", "secrets"},
		},
		{
			Verbs:         []string{"update"},
			APIGroups:     []string{""},
			Resources:     []string{"secrets"},
			ResourceNames: []string{"karpenter-cert"},
		},
		{
			Verbs:     []string{"update", "patch", "delete"},
			APIGroups: []string{""},
			Resources: []string{"configmaps"},
			ResourceNames: []string{
				"karpenter-global-settings",
				"config-logging",
			},
		},
		{
			Verbs:     []string{"patch", "update"},
			APIGroups: []string{"coordination.k8s.io"},
			Resources: []string{"leases"},
			ResourceNames: []string{
				"karpenter-leader-election",
				"webhook.configmapwebhook.00-of-01",
				"webhook.defaultingwebhook.00-of-01",
				"webhook.validationwebhook.00-of-01",
				"webhook.webhookcertificates.00-of-01",
			},
		},
		{
			Verbs:     []string{"create"},
			APIGroups: []string{"coordination.k8s.io"},
			Resources: []string{"leases"},
		},
		{
			Verbs:     []string{"create"},
			APIGroups: []string{""},
			Resources: []string{"configmaps"},
		},
	},
}
