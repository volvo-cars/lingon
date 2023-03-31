// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package cilium

import (
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// REPLACE BY KUBEUTILS
//
// var SvcAcc = &corev1.ServiceAccount{
// 	TypeMeta: kubeutils.TypeMeta("ServiceAccount"),
// 	ObjectMeta: metav1.ObjectMeta{
// 		Name:      "cilium",
// 		Namespace: "kube-system",
// 	},
// }

var CR = &rbacv1.ClusterRole{
	TypeMeta:   kubeutil.TypeClusterRoleV1,
	ObjectMeta: metav1.ObjectMeta{Name: "cilium"},
	Rules: []rbacv1.PolicyRule{
		{
			Verbs: []string{
				"get",
				"list",
				"watch",
			},
			APIGroups: []string{"networking.k8s.io"},
			Resources: []string{"networkpolicies"},
		},
		{
			Verbs: []string{
				"get",
				"list",
				"watch",
			},
			APIGroups: []string{"discovery.k8s.io"},
			Resources: []string{"endpointslices"},
		},
		{
			Verbs: []string{
				"get",
				"list",
				"watch",
			},
			APIGroups: []string{""},
			Resources: []string{
				"namespaces",
				"services",
				"pods",
				"endpoints",
				"nodes",
			},
		},
		{
			Verbs: []string{
				"list",
				"watch",
				"get",
			},
			APIGroups: []string{"apiextensions.k8s.io"},
			Resources: []string{"customresourcedefinitions"},
		},
		{
			Verbs: []string{
				"list",
				"watch",
			},
			APIGroups: []string{"cilium.io"},
			Resources: []string{
				"ciliumbgploadbalancerippools",
				"ciliumbgppeeringpolicies",
				"ciliumclusterwideenvoyconfigs",
				"ciliumclusterwidenetworkpolicies",
				"ciliumegressgatewaypolicies",
				"ciliumegressnatpolicies",
				"ciliumendpoints",
				"ciliumendpointslices",
				"ciliumenvoyconfigs",
				"ciliumidentities",
				"ciliumlocalredirectpolicies",
				"ciliumnetworkpolicies",
				"ciliumnodes",
			},
		},
		{
			Verbs:     []string{"create"},
			APIGroups: []string{"cilium.io"},
			Resources: []string{
				"ciliumidentities",
				"ciliumendpoints",
				"ciliumnodes",
			},
		},
		{
			Verbs:     []string{"update"},
			APIGroups: []string{"cilium.io"},
			Resources: []string{"ciliumidentities"},
		},
		{
			Verbs: []string{
				"delete",
				"get",
			},
			APIGroups: []string{"cilium.io"},
			Resources: []string{"ciliumendpoints"},
		},
		{
			Verbs: []string{
				"get",
				"update",
			},
			APIGroups: []string{"cilium.io"},
			Resources: []string{
				"ciliumnodes",
				"ciliumnodes/status",
			},
		},
		{
			Verbs:     []string{"patch"},
			APIGroups: []string{"cilium.io"},
			Resources: []string{
				"ciliumnetworkpolicies/status",
				"ciliumclusterwidenetworkpolicies/status",
				"ciliumendpoints/status",
				"ciliumendpoints",
			},
		},
	},
}
