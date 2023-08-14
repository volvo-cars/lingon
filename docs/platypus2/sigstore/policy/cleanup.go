// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package policy

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var SigstoreWebhookCleanupSA = &corev1.ServiceAccount{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "post-delete",
			"helm.sh/hook-delete-policy": "hook-succeeded",
			"helm.sh/hook-weight":        "2",
		},
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "sigstore",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "policy-controller",
			"app.kubernetes.io/version":    "0.8.0",
			"control-plane":                "sigstore-policy-controller-cleanup",
			"helm.sh/chart":                "policy-controller-0.6.0",
		},
		Name:      "sigstore-policy-controller-webhook-cleanup",
		Namespace: "sigstore",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ServiceAccount",
	},
}

var LeasesCleanupJOBS = &batchv1.Job{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "post-delete",
			"helm.sh/hook-delete-policy": "hook-succeeded",
			"helm.sh/hook-weight":        "3",
		},
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "sigstore",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "policy-controller",
			"app.kubernetes.io/version":    "0.8.0",
			"control-plane":                "sigstore-policy-controller-webhook",
			"helm.sh/chart":                "policy-controller-0.6.0",
		},
		Name:      "leases-cleanup",
		Namespace: "sigstore",
	},
	Spec: batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Name: "leases-cleanup"},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Command: []string{
							"/bin/sh",
							"-c",
							"kubectl delete leases --all --ignore-not-found -n sigstore",
						},
						Image:           "cgr.dev/chainguard/kubectl:latest",
						ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
						Name:            "kubectl",
					},
				},
				RestartPolicy:      corev1.RestartPolicy("OnFailure"),
				ServiceAccountName: "sigstore-policy-controller-webhook-cleanup",
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "batch/v1",
		Kind:       "Job",
	},
}

var SigstoreCleanupRole = &rbacv1.Role{
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
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{"coordination.k8s.io"},
			Resources: []string{"leases"},
			Verbs:     []string{"list", "delete"},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "Role",
	},
}

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
