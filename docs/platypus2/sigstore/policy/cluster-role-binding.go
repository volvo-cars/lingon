// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package policy

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var SigstoreWebhookCRB = &rbacv1.ClusterRoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "sigstore",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "policy-controller",
			"app.kubernetes.io/version":    "0.8.0",
			"control-plane":                "sigstore-policy-controller-webhook",
			"helm.sh/chart":                "policy-controller-0.6.0",
		},
		Name: "sigstore-policy-controller-webhook",
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "ClusterRole",
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
		Kind:       "ClusterRoleBinding",
	},
}
