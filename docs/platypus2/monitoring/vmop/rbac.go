// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package vmop

import (
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	rbacv1 "k8s.io/api/rbac/v1"
)

var Role = &rbacv1.Role{
	TypeMeta:   ku.TypeRoleV1,
	ObjectMeta: O.ObjectMeta(),
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{""},
			Resources: []string{"configmaps"},
			Verbs: []string{
				"get",
				"list",
				"watch",
				"create",
				"update",
				"patch",
				"delete",
			},
		}, {
			APIGroups: []string{""},
			Resources: []string{"configmaps/status"},
			Verbs:     []string{"get", "update", "patch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"events"},
			Verbs:     []string{"create", "patch"},
		}, {
			APIGroups: []string{"coordination.k8s.io"},
			Resources: []string{"leases"},
			Verbs:     []string{"create", "get", "update"},
		},
	},
}

var CR = &rbacv1.ClusterRole{
	TypeMeta:   ku.TypeClusterRoleV1,
	ObjectMeta: O.ObjectMetaNoNS(),
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{""},
			Resources: []string{"configmaps", "configmaps/finalizers"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"endpoints"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"events"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"namespaces"},
			Verbs:     []string{"get", "list", "watch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{
				"persistentvolumeclaims",
				"persistentvolumeclaims/finalizers",
			},
			Verbs: []string{"*"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"pods"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"secrets", "secrets/finalizers"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"services"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"services/finalizers"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{"apps"},
			Resources: []string{"deployments", "deployments/finalizers"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{"apps"},
			Resources: []string{"replicasets"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{"apps"},
			Resources: []string{
				"statefulsets",
				"statefulsets/finalizers",
				"statefulsets/status",
			},
			Verbs: []string{"*"},
		}, {
			APIGroups: []string{"monitoring.coreos.com"},
			Resources: []string{"*"},
			Verbs:     []string{"*"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmagents", "vmagents/finalizers"},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmagents/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{
				"vmalertmanagers",
				"vmalertmanagers/finalizers",
			},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmalertmanagers/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{
				"vmalertmanagerconfigs",
				"vmalertmanagerconfigs/finalizers",
			},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmalertmanagerconfigss/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmalerts", "vmalerts/finalizers"},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmalerts/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmclusters", "vmclusters/finalizers"},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmclusters/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmpodscrapes", "vmprobscrapes/finalizers"},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmpodscrapes/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmrules", "vmrules/finalizers"},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmrules/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{
				"vmservicescrapes",
				"vmservicescrapes/finalizers",
			},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmservicescrapes/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmprobes"},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmprobes/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmsingles", "vmsingles/finalizers"},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmsingles/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{""},
			Resources: []string{
				"nodes",
				"nodes/proxy",
				"services",
				"endpoints",
				"pods",
				"endpointslices",
				"configmaps",
				"nodes/metrics",
				"namespaces",
			},
			Verbs: []string{"get", "list", "watch"},
		}, {
			APIGroups: []string{"extensions", "networking.k8s.io"},
			Resources: []string{"ingresses"},
			Verbs:     []string{"get", "list", "watch"},
		}, {
			NonResourceURLs: []string{"/metrics", "/metrics/resources"},
			Verbs:           []string{"get", "watch", "list"},
		}, {
			APIGroups: []string{"rbac.authorization.k8s.io"},
			Resources: []string{
				"clusterrolebindings",
				"clusterrolebindings/finalizers",
				"clusterroles",
				"clusterroles/finalizers",
				"roles",
				"rolebindings",
			},
			Verbs: []string{
				"get",
				"list",
				"create",
				"patch",
				"update",
				"watch",
				"delete",
			},
		}, {
			APIGroups: []string{"policy"},
			Resources: []string{
				"podsecuritypolicies",
				"podsecuritypolicies/finalizers",
			},
			Verbs: []string{
				"get",
				"list",
				"create",
				"patch",
				"update",
				"use",
				"watch",
				"delete",
			},
		}, {
			APIGroups: []string{""},
			Resources: []string{
				"serviceaccounts",
				"serviceaccounts/finalizers",
			},
			Verbs: []string{
				"get",
				"list",
				"create",
				"watch",
				"update",
				"delete",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmnodescrapes", "vmnodescrapes/finalizers"},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmnodescrapes/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{
				"vmstaticscrapes",
				"vmstaticscrapes/finalizers",
			},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmstaticscrapes/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{
				"vmauths",
				"vmauths/finalizers",
				"vmusers",
				"vmusers/finalizers",
			},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"list",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"operator.victoriametrics.com"},
			Resources: []string{"vmusers/status", "vmauths/status"},
			Verbs:     []string{"get", "patch", "update"},
		}, {
			APIGroups: []string{"storage.k8s.io"},
			Resources: []string{"storageclasses"},
			Verbs:     []string{"list", "get", "watch"},
		}, {
			APIGroups: []string{"policy"},
			Resources: []string{
				"poddisruptionbudgets",
				"poddisruptionbudgets/finalizers",
			},
			Verbs: []string{"*"},
		}, {
			APIGroups: []string{"route.openshift.io", "image.openshift.io"},
			Resources: []string{"routers/metrics", "registry/metrics"},
			Verbs:     []string{"get"},
		}, {
			APIGroups: []string{"autoscaling"},
			Resources: []string{"horizontalpodautoscalers"},
			Verbs: []string{
				"list",
				"get",
				"delete",
				"create",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"networking.k8s.io", "extensions"},
			Resources: []string{"ingresses", "ingresses/finalizers"},
			Verbs: []string{
				"create",
				"delete",
				"get",
				"patch",
				"update",
				"watch",
			},
		}, {
			APIGroups: []string{"apiextensions.k8s.io"},
			Resources: []string{"customresourcedefinitions"},
			Verbs:     []string{"get", "list"},
		},
	},
}
