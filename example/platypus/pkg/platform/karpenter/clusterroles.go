// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package karpenter

import (
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var CanUpdateWebhooks = &rbacv1.ClusterRole{
	TypeMeta: metav1.TypeMeta{
		Kind:       "ClusterRole",
		APIVersion: "rbac.authorization.k8s.io/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:   "karpenter",
		Labels: commonLabels,
	},
	Rules: []rbacv1.PolicyRule{
		{
			Verbs: []string{
				"get",
				"list",
				"watch",
			},
			APIGroups: []string{"karpenter.k8s.aws"},
			Resources: []string{"awsnodetemplates"},
		},
		{
			Verbs:         []string{"update"},
			APIGroups:     []string{"admissionregistration.k8s.io"},
			Resources:     []string{"validatingwebhookconfigurations"},
			ResourceNames: []string{"validation.webhook.karpenter.k8s.aws"},
		},
		{
			Verbs:         []string{"update"},
			APIGroups:     []string{"admissionregistration.k8s.io"},
			Resources:     []string{"mutatingwebhookconfigurations"},
			ResourceNames: []string{"defaulting.webhook.karpenter.k8s.aws"},
		},
		{
			Verbs:     []string{"patch", "update"},
			APIGroups: []string{"karpenter.k8s.aws"},
			Resources: []string{"awsnodetemplates/status"},
		},
	},
}

var CoreCr = &rbacv1.ClusterRole{
	TypeMeta: kubeutil.TypeClusterRoleV1,
	ObjectMeta: metav1.ObjectMeta{
		Name:   "karpenter-core",
		Labels: commonLabels,
	},
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{"karpenter.sh"},
			Resources: []string{
				"provisioners",
				"provisioners/status",
				"machines",
				"machines/status",
			},
			Verbs: []string{"get", "list", "watch"},
		},
		{
			APIGroups: []string{""},
			Resources: []string{
				"pods",
				"nodes",
				"persistentvolumes",
				"persistentvolumeclaims",
				"replicationcontrollers",
				"namespaces",
			},
			Verbs: []string{"get", "list", "watch"},
		},
		{
			APIGroups: []string{"storage.k8s.io"},
			Resources: []string{"storageclasses", "csinodes"},
			Verbs:     []string{"get", "watch", "list"},
		},
		{
			APIGroups: []string{"apps"},
			Resources: []string{
				"daemonsets",
				"deployments",
				"replicasets",
				"statefulsets",
			},
			Verbs: []string{"list", "watch"},
		},
		{
			APIGroups: []string{"admissionregistration.k8s.io"},
			Resources: []string{
				"validatingwebhookconfigurations",
				"mutatingwebhookconfigurations",
			},
			Verbs: []string{"get", "watch", "list"},
		},
		{
			APIGroups: []string{"policy"},
			Resources: []string{"poddisruptionbudgets"},
			Verbs:     []string{"get", "list", "watch"},
		},
		{
			APIGroups: []string{"karpenter.sh"},
			Resources: []string{
				"provisioners/status",
				"machines",
				"machines/status",
			},
			Verbs: []string{"create", "delete", "patch"},
		},
		{
			APIGroups: []string{""},
			Resources: []string{"events"},
			Verbs:     []string{"create", "patch"},
		},
		{
			APIGroups: []string{""},
			Resources: []string{"nodes"},
			Verbs:     []string{"create", "patch", "delete"},
		},
		{
			APIGroups: []string{""},
			Resources: []string{"pods/eviction"},
			Verbs:     []string{"create"},
		},
		{
			APIGroups: []string{"admissionregistration.k8s.io"},
			ResourceNames: []string{
				"validation.webhook.karpenter.sh",
				"validation.webhook.config.karpenter.sh",
			},
			Resources: []string{"validatingwebhookconfigurations"},
			Verbs:     []string{"update"},
		},
		{
			APIGroups:     []string{"admissionregistration.k8s.io"},
			ResourceNames: []string{"defaulting.webhook.karpenter.sh"},
			Resources:     []string{"mutatingwebhookconfigurations"},
			Verbs:         []string{"update"},
		},
	},
}

var AdminCr = &rbacv1.ClusterRole{
	TypeMeta: kubeutil.TypeClusterRoleV1,
	ObjectMeta: metav1.ObjectMeta{
		Name: "karpenter-admin",
		Labels: appendCommonLabels(
			map[string]string{
				// Add these permissions to the "admin" default roles
				"rbac.authorization.k8s.io/aggregate-to-admin": "true",
			},
		),
	},
	Rules: []rbacv1.PolicyRule{
		{
			Verbs: []string{
				"get",
				"list",
				"watch",
				"create",
				"delete",
				"patch",
			},
			APIGroups: []string{"karpenter.sh"},
			Resources: []string{"provisioners", "provisioners/status"},
		},
		{
			Verbs: []string{
				"get",
				"list",
				"watch",
				"create",
				"delete",
				"patch",
			},
			APIGroups: []string{"karpenter.k8s.aws"},
			Resources: []string{"awsnodetemplates"},
		},
	},
}
