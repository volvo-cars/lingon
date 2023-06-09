// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package tekton

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var PipelinesControllerLeaderelectionRB = &rbacv1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/component": "controller",
			"app.kubernetes.io/instance":  "default",
			"app.kubernetes.io/part-of":   "tekton-pipelines",
		},
		Name:      "tekton-pipelines-controller-leaderelection",
		Namespace: "tekton-pipelines",
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "Role",
		Name:     "tekton-pipelines-leader-election",
	},
	Subjects: []rbacv1.Subject{{
		Kind:      "ServiceAccount",
		Name:      "tekton-pipelines-controller",
		Namespace: "tekton-pipelines",
	}},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "RoleBinding",
	},
}

var PipelinesControllerRB = &rbacv1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/component": "controller",
			"app.kubernetes.io/instance":  "default",
			"app.kubernetes.io/part-of":   "tekton-pipelines",
		},
		Name:      "tekton-pipelines-controller",
		Namespace: "tekton-pipelines",
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "Role",
		Name:     "tekton-pipelines-controller",
	},
	Subjects: []rbacv1.Subject{{
		Kind:      "ServiceAccount",
		Name:      "tekton-pipelines-controller",
		Namespace: "tekton-pipelines",
	}},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "RoleBinding",
	},
}

var PipelinesInfoRB = &rbacv1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance": "default",
			"app.kubernetes.io/part-of":  "tekton-pipelines",
		},
		Name:      "tekton-pipelines-info",
		Namespace: "tekton-pipelines",
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "Role",
		Name:     "tekton-pipelines-info",
	},
	Subjects: []rbacv1.Subject{{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "Group",
		Name:     "system:authenticated",
	}},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "RoleBinding",
	},
}

var PipelinesResolversNamespaceRbacRB = &rbacv1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/component": "resolvers",
			"app.kubernetes.io/instance":  "default",
			"app.kubernetes.io/part-of":   "tekton-pipelines",
		},
		Name:      "tekton-pipelines-resolvers-namespace-rbac",
		Namespace: "tekton-pipelines-resolvers",
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "Role",
		Name:     "tekton-pipelines-resolvers-namespace-rbac",
	},
	Subjects: []rbacv1.Subject{{
		Kind:      "ServiceAccount",
		Name:      "tekton-pipelines-resolvers",
		Namespace: "tekton-pipelines-resolvers",
	}},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "RoleBinding",
	},
}

var PipelinesWebhookLeaderelectionRB = &rbacv1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/component": "webhook",
			"app.kubernetes.io/instance":  "default",
			"app.kubernetes.io/part-of":   "tekton-pipelines",
		},
		Name:      "tekton-pipelines-webhook-leaderelection",
		Namespace: "tekton-pipelines",
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "Role",
		Name:     "tekton-pipelines-leader-election",
	},
	Subjects: []rbacv1.Subject{{
		Kind:      "ServiceAccount",
		Name:      "tekton-pipelines-webhook",
		Namespace: "tekton-pipelines",
	}},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "RoleBinding",
	},
}

var PipelinesWebhookRB = &rbacv1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/component": "webhook",
			"app.kubernetes.io/instance":  "default",
			"app.kubernetes.io/part-of":   "tekton-pipelines",
		},
		Name:      "tekton-pipelines-webhook",
		Namespace: "tekton-pipelines",
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "Role",
		Name:     "tekton-pipelines-webhook",
	},
	Subjects: []rbacv1.Subject{{
		Kind:      "ServiceAccount",
		Name:      "tekton-pipelines-webhook",
		Namespace: "tekton-pipelines",
	}},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "RoleBinding",
	},
}
