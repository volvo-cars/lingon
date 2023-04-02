// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package externalsecrets

import (
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var CertControllerCR = &rbacv1.ClusterRole{
	TypeMeta: kubeutil.TypeClusterRoleV1,
	ObjectMeta: metav1.ObjectMeta{
		Labels: CertControllerLabels,
		Name:   certControllerName,
	},
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{"apiextensions.k8s.io"},
			Resources: []string{"customresourcedefinitions"},
			Verbs:     []string{"get", "list", "watch", "update", "patch"},
		}, {
			APIGroups: []string{"admissionregistration.k8s.io"},
			Resources: []string{"validatingwebhookconfigurations"},
			Verbs:     []string{"get", "list", "watch", "update", "patch"},
		}, {
			APIGroups: []string{},
			Resources: []string{"endpoints"},
			Verbs:     []string{"list", "get", "watch"},
		}, {
			APIGroups: []string{},
			Resources: []string{"events"},
			Verbs:     []string{"create", "patch"},
		}, {
			APIGroups: []string{},
			Resources: []string{"secrets"},
			Verbs:     []string{"get", "list", "watch", "update", "patch"},
		},
	},
}

var ControllerCr = &rbacv1.ClusterRole{
	TypeMeta: kubeutil.TypeClusterRoleV1,
	ObjectMeta: metav1.ObjectMeta{
		Labels: ESLabels,
		Name:   controllerName,
	},
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{"external-secrets.io"},
			Resources: []string{
				"secretstores",
				"clustersecretstores",
				"externalsecrets",
				"clusterexternalsecrets",
				"pushsecrets",
			},
			Verbs: []string{"get", "list", "watch"},
		}, {
			APIGroups: []string{"external-secrets.io"},
			Resources: []string{
				"externalsecrets",
				"externalsecrets/status",
				"externalsecrets/finalizers",
				"secretstores",
				"secretstores/status",
				"secretstores/finalizers",
				"clustersecretstores",
				"clustersecretstores/status",
				"clustersecretstores/finalizers",
				"clusterexternalsecrets",
				"clusterexternalsecrets/status",
				"clusterexternalsecrets/finalizers",
				"pushsecrets",
				"pushsecrets/status",
				"pushsecrets/finalizers",
			},
			Verbs: []string{"update", "patch"},
		}, {
			APIGroups: []string{"generators.external-secrets.io"},
			Resources: []string{
				"fakes",
				"passwords",
				"acraccesstokens",
				"gcraccesstokens",
				"ecrauthorizationtokens",
			},
			Verbs: []string{"get", "list", "watch"},
		}, {
			APIGroups: []string{},
			Resources: []string{"serviceaccounts", "namespaces"},
			Verbs:     []string{"get", "list", "watch"},
		}, {
			APIGroups: []string{},
			Resources: []string{"configmaps"},
			Verbs:     []string{"get", "list", "watch"},
		}, {
			APIGroups: []string{},
			Resources: []string{"secrets"},
			Verbs: []string{
				"get",
				"list",
				"watch",
				"create",
				"update",
				"delete",
				"patch",
			},
		}, {
			APIGroups: []string{},
			Resources: []string{"serviceaccounts/token"},
			Verbs:     []string{"create"},
		}, {
			APIGroups: []string{},
			Resources: []string{"events"},
			Verbs:     []string{"create", "patch"},
		}, {
			APIGroups: []string{"external-secrets.io"},
			Resources: []string{"externalsecrets"},
			Verbs:     []string{"create", "update", "delete"},
		},
	},
}

var EditCR = &rbacv1.ClusterRole{
	TypeMeta: kubeutil.TypeClusterRoleV1,
	ObjectMeta: metav1.ObjectMeta{
		Labels: kubeutil.MergeLabels(
			ESLabels, map[string]string{
				"rbac.authorization.k8s.io/aggregate-to-admin": "true",
				"rbac.authorization.k8s.io/aggregate-to-edit":  "true",
			},
		),
		Name: AppName + "-edit",
	},
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{"external-secrets.io"},
			Resources: []string{
				"externalsecrets",
				"secretstores",
				"clustersecretstores",
				"pushsecrets",
			},
			Verbs: []string{
				"create",
				"delete",
				"deletecollection",
				"patch",
				"update",
			},
		},
	},
}

var ViewCR = &rbacv1.ClusterRole{
	TypeMeta: kubeutil.TypeClusterRoleV1,
	ObjectMeta: metav1.ObjectMeta{
		Labels: kubeutil.MergeLabels(
			ESLabels, map[string]string{
				"rbac.authorization.k8s.io/aggregate-to-admin": "true",
				"rbac.authorization.k8s.io/aggregate-to-edit":  "true",
				"rbac.authorization.k8s.io/aggregate-to-view":  "true",
			},
		),
		Name: AppName + "-view",
	},
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{"external-secrets.io"},
			Resources: []string{
				"externalsecrets",
				"secretstores",
				"clustersecretstores",
				"pushsecrets",
			},
			Verbs: []string{"get", "watch", "list"},
		},
	},
}
