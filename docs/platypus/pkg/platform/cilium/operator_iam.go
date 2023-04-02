// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package cilium

import (
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var OperatorCR = &rbacv1.ClusterRole{
	TypeMeta:   kubeutil.TypeClusterRoleV1,
	ObjectMeta: metav1.ObjectMeta{Name: "cilium-operator"},
	Rules: []rbacv1.PolicyRule{
		{
			Verbs: []string{
				"get",
				"list",
				"watch",
				"delete",
			},
			APIGroups: []string{""},
			Resources: []string{"pods"},
		},
		{
			Verbs: []string{
				"list",
				"watch",
			},
			APIGroups: []string{""},
			Resources: []string{"nodes"},
		},
		{
			Verbs:     []string{"patch"},
			APIGroups: []string{""},
			Resources: []string{
				"nodes",
				"nodes/status",
			},
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
			Verbs:     []string{"update"},
			APIGroups: []string{""},
			Resources: []string{"services/status"},
		},
		{
			Verbs: []string{
				"get",
				"list",
				"watch",
			},
			APIGroups: []string{""},
			Resources: []string{"namespaces"},
		},
		{
			Verbs: []string{
				"get",
				"list",
				"watch",
			},
			APIGroups: []string{""},
			Resources: []string{
				"services",
				"endpoints",
			},
		},
		{
			Verbs: []string{
				"create",
				"update",
				"deletecollection",
				"patch",
				"get",
				"list",
				"watch",
			},
			APIGroups: []string{"cilium.io"},
			Resources: []string{
				"ciliumnetworkpolicies",
				"ciliumclusterwidenetworkpolicies",
			},
		},
		{
			Verbs: []string{
				"patch",
				"update",
			},
			APIGroups: []string{"cilium.io"},
			Resources: []string{
				"ciliumnetworkpolicies/status",
				"ciliumclusterwidenetworkpolicies/status",
			},
		},
		{
			Verbs: []string{
				"delete",
				"list",
				"watch",
			},
			APIGroups: []string{"cilium.io"},
			Resources: []string{
				"ciliumendpoints",
				"ciliumidentities",
			},
		},
		{
			Verbs:     []string{"update"},
			APIGroups: []string{"cilium.io"},
			Resources: []string{"ciliumidentities"},
		},
		{
			Verbs: []string{
				"create",
				"update",
				"get",
				"list",
				"watch",
				"delete",
			},
			APIGroups: []string{"cilium.io"},
			Resources: []string{"ciliumnodes"},
		},
		{
			Verbs:     []string{"update"},
			APIGroups: []string{"cilium.io"},
			Resources: []string{"ciliumnodes/status"},
		},
		{
			Verbs: []string{
				"create",
				"update",
				"get",
				"list",
				"watch",
				"delete",
			},
			APIGroups: []string{"cilium.io"},
			Resources: []string{
				"ciliumendpointslices",
				"ciliumenvoyconfigs",
			},
		},
		{
			Verbs: []string{
				"create",
				"get",
				"list",
				"watch",
			},
			APIGroups: []string{"apiextensions.k8s.io"},
			Resources: []string{"customresourcedefinitions"},
		},
		{
			Verbs:     []string{"update"},
			APIGroups: []string{"apiextensions.k8s.io"},
			Resources: []string{"customresourcedefinitions"},
			ResourceNames: []string{
				"ciliumbgploadbalancerippools.cilium.io",
				"ciliumbgppeeringpolicies.cilium.io",
				"ciliumclusterwideenvoyconfigs.cilium.io",
				"ciliumclusterwidenetworkpolicies.cilium.io",
				"ciliumegressgatewaypolicies.cilium.io",
				"ciliumegressnatpolicies.cilium.io",
				"ciliumendpoints.cilium.io",
				"ciliumendpointslices.cilium.io",
				"ciliumenvoyconfigs.cilium.io",
				"ciliumexternalworkloads.cilium.io",
				"ciliumidentities.cilium.io",
				"ciliumlocalredirectpolicies.cilium.io",
				"ciliumnetworkpolicies.cilium.io",
				"ciliumnodes.cilium.io",
			},
		},
		{
			Verbs: []string{
				"create",
				"get",
				"update",
			},
			APIGroups: []string{"coordination.k8s.io"},
			Resources: []string{"leases"},
		},
	},
}
