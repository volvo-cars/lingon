// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package policy

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var SigstoreCleanupRB = &rbacv1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "post-delete",
			"helm.sh/hook-delete-policy": "hook-succeeded",
			"helm.sh/hook-weight":        "1",
		},
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "sigstore",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "policy-controller",
			"app.kubernetes.io/version":    "0.8.0",
			"control-plane":                "sigstore-policy-controller-cleanup",
			"helm.sh/chart":                "policy-controller-0.6.0",
		},
		Name:      "sigstore-policy-controller-cleanup",
		Namespace: "sigstore",
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "Role",
		Name:     "sigstore-policy-controller-cleanup",
	},
	Subjects: []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      "sigstore-policy-controller-webhook-cleanup",
			Namespace: "sigstore",
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "RoleBinding",
	},
}

var SigstoreWebhookRB = &rbacv1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "sigstore",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "policy-controller",
			"app.kubernetes.io/version":    "0.8.0",
			"control-plane":                "sigstore-policy-controller-webhook",
			"helm.sh/chart":                "policy-controller-0.6.0",
		},
		Name:      "sigstore-policy-controller-webhook",
		Namespace: "sigstore",
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "Role",
		Name:     "sigstore-policy-controller-webhook",
	},
	Subjects: []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      "sigstore-policy-controller-webhook",
			Namespace: "sigstore",
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "RoleBinding",
	},
}
