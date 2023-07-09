// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package certmanager

import (
	"fmt"

	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Webhook struct {
	kube.App

	WebhookCM                      *corev1.ConfigMap
	WebhookDeploy                  *appsv1.Deployment
	WebhookDynamicServingRB        *rbacv1.RoleBinding
	WebhookDynamicServingRole      *rbacv1.Role
	WebhookMutatingWC              *admissionregistrationv1.MutatingWebhookConfiguration
	WebhookSA                      *corev1.ServiceAccount
	WebhookSVC                     *corev1.Service
	WebhookSubjectAccessReviewsCR  *rbacv1.ClusterRole
	WebhookSubjectAccessReviewsCRB *rbacv1.ClusterRoleBinding
	WebhookValidatingWC            *admissionregistrationv1.ValidatingWebhookConfiguration
}

func d(i int32) string { return fmt.Sprintf("%d", i) }

var WebhookDeploy = &appsv1.Deployment{
	TypeMeta:   ku.TypeDeploymentV1,
	ObjectMeta: CM.Webhook.ObjectMeta(),
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(1)),
		Selector: &metav1.LabelSelector{MatchLabels: CM.Webhook.MatchLabels()},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: CM.Webhook.Labels()},
			Spec: corev1.PodSpec{
				ServiceAccountName: CM.Webhook.ServiceAccount().Name,
				Containers: []corev1.Container{
					{
						Name:            CM.Webhook.Name,
						Image:           CM.Webhook.Img.URL(),
						ImagePullPolicy: corev1.PullIfNotPresent,
						Args: []string{
							"--v=2",
							"--secure-port=" + d(CM.WebhookPort.Container.ContainerPort),
							"--dynamic-serving-ca-secret-namespace=$(POD_NAMESPACE)",
							"--dynamic-serving-ca-secret-name=cert-manager-webhook-ca",
							"--dynamic-serving-dns-names=" + CM.Webhook.Name,
							"--dynamic-serving-dns-names=" + CM.Webhook.Name + ".$(POD_NAMESPACE)",
							"--dynamic-serving-dns-names=" + CM.Webhook.Name + ".$(POD_NAMESPACE).svc",
						},
						Env: []corev1.EnvVar{
							ku.EnvVarDownAPI(
								"POD_NAMESPACE", "metadata.namespace",
							),
						},
						LivenessProbe: &corev1.Probe{
							FailureThreshold:    int32(3),
							InitialDelaySeconds: int32(60),
							PeriodSeconds:       int32(10),
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "/livez",
									Port:   intstr.FromInt(6080),
									Scheme: corev1.URISchemeHTTP,
								},
							},
							SuccessThreshold: int32(1),
							TimeoutSeconds:   int32(1),
						},
						Ports: []corev1.ContainerPort{
							CM.WebhookPort.Container,
							{Name: "healthcheck", ContainerPort: int32(6080)},
						},
						ReadinessProbe: &corev1.Probe{
							FailureThreshold:    int32(3),
							InitialDelaySeconds: int32(5),
							PeriodSeconds:       int32(5),
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "/healthz",
									Port:   intstr.FromInt(6080),
									Scheme: corev1.URISchemeHTTP,
								},
							},
							SuccessThreshold: int32(1),
							TimeoutSeconds:   int32(1),
						},
						SecurityContext: &corev1.SecurityContext{Capabilities: &corev1.Capabilities{Drop: []corev1.Capability{corev1.Capability("ALL")}}},
					},
				},
				NodeSelector: map[string]string{ku.LabelOSStable: "linux"},
				SecurityContext: &corev1.PodSecurityContext{
					RunAsNonRoot:   P(true),
					SeccompProfile: &corev1.SeccompProfile{Type: corev1.SeccompProfileType("RuntimeDefault")},
				},
			},
		},
	},
}

var WebhookSVC = &corev1.Service{
	TypeMeta:   ku.TypeServiceV1,
	ObjectMeta: CM.Webhook.ObjectMeta(),
	Spec: corev1.ServiceSpec{
		Ports:    []corev1.ServicePort{CM.WebhookPort.Service},
		Selector: CM.Webhook.MatchLabels(),
		Type:     corev1.ServiceTypeClusterIP,
	},
}

var WebhookCM = &corev1.ConfigMap{
	Data: nil,
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "webhook",
			"app.kubernetes.io/component":  "webhook",
			"app.kubernetes.io/instance":   "cert-manager",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "webhook",
			"app.kubernetes.io/version":    "v1.12.2",
			"helm.sh/chart":                "cert-manager-v1.12.2",
		},
		Name:      "cert-manager-webhook",
		Namespace: "cert-manager",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ConfigMap",
	},
}

var WebhookDynamicServingRole = &rbacv1.Role{
	TypeMeta:   ku.TypeRoleV1,
	ObjectMeta: CM.Webhook.ObjectMetaNameSuffix(":dynamic-serving"),
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups:     []string{""},
			ResourceNames: []string{"cert-manager-webhook-ca"},
			Resources:     []string{"secrets"},
			Verbs:         []string{"get", "list", "watch", "update"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"secrets"},
			Verbs:     []string{"create"},
		},
	},
}

var WebhookSubjectAccessReviewsCR = &rbacv1.ClusterRole{
	TypeMeta:   ku.TypeClusterRoleV1,
	ObjectMeta: CM.Webhook.ObjectMetaNameSuffixNoNS(":subjectaccessreviews"),
	// ObjectMeta: metav1.ObjectMeta{
	// 	Labels: map[string]string{
	// 		"app":                          "webhook",
	// 		"app.kubernetes.io/component":  "webhook",
	// 		"app.kubernetes.io/instance":   "cert-manager",
	// 		"app.kubernetes.io/managed-by": "Helm",
	// 		"app.kubernetes.io/name":       "webhook",
	// 		"app.kubernetes.io/version":    "v1.12.2",
	// 		"helm.sh/chart":                "cert-manager-v1.12.2",
	// 	},
	// 	Name: "cert-manager-webhook:subjectaccessreviews",
	// },
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{"authorization.k8s.io"},
			Resources: []string{"subjectaccessreviews"},
			Verbs:     []string{"create"},
		},
	},
}

var WebhookSubjectAccessReviewsCRB = ku.BindClusterRole(
	WebhookSubjectAccessReviewsCR.Name,
	CM.WebhookSA,
	WebhookSubjectAccessReviewsCR,
	CM.Webhook.Labels(),
)

// var WebhookSubjectAccessReviewsCRB = &rbacv1.ClusterRoleBinding{
// 	ObjectMeta: metav1.ObjectMeta{
// 		Labels: map[string]string{
// 			"app":                          "webhook",
// 			"app.kubernetes.io/component":  "webhook",
// 			"app.kubernetes.io/instance":   "cert-manager",
// 			"app.kubernetes.io/managed-by": "Helm",
// 			"app.kubernetes.io/name":       "webhook",
// 			"app.kubernetes.io/version":    "v1.12.2",
// 			"helm.sh/chart":                "cert-manager-v1.12.2",
// 		},
// 		Name: "cert-manager-webhook:subjectaccessreviews",
// 	},
// 	RoleRef: rbacv1.RoleRef{
// 		APIGroup: "rbac.authorization.k8s.io",
// 		Kind:     "ClusterRole",
// 		Name:     "cert-manager-webhook:subjectaccessreviews",
// 	},
// 	Subjects: []rbacv1.Subject{
// 		{
// 			Kind:      "ServiceAccount",
// 			Name:      "cert-manager-webhook",
// 			Namespace: "cert-manager",
// 		},
// 	},
// 	TypeMeta: metav1.TypeMeta{
// 		APIVersion: "rbac.authorization.k8s.io/v1",
// 		Kind:       "ClusterRoleBinding",
// 	},
// }

var WebhookDynamicServingRB = ku.BindRole(
	WebhookDynamicServingRole.Name,
	CM.WebhookSA,
	WebhookDynamicServingRole,
	CM.Webhook.Labels(),
)

// var WebhookDynamicServingRB = &rbacv1.RoleBinding{
// 	ObjectMeta: metav1.ObjectMeta{
// 		Labels: map[string]string{
// 			"app":                          "webhook",
// 			"app.kubernetes.io/component":  "webhook",
// 			"app.kubernetes.io/instance":   "cert-manager",
// 			"app.kubernetes.io/managed-by": "Helm",
// 			"app.kubernetes.io/name":       "webhook",
// 			"app.kubernetes.io/version":    "v1.12.2",
// 			"helm.sh/chart":                "cert-manager-v1.12.2",
// 		},
// 		Name:      "cert-manager-webhook:dynamic-serving",
// 		Namespace: "cert-manager",
// 	},
// 	RoleRef: rbacv1.RoleRef{
// 		APIGroup: "rbac.authorization.k8s.io",
// 		Kind:     "Role",
// 		Name:     "cert-manager-webhook:dynamic-serving",
// 	},
// 	Subjects: []rbacv1.Subject{
// 		{
// 			Kind:      "ServiceAccount",
// 			Name:      "cert-manager-webhook",
// 			Namespace: "cert-manager",
// 		},
// 	},
// 	TypeMeta: metav1.TypeMeta{
// 		APIVersion: "rbac.authorization.k8s.io/v1",
// 		Kind:       "RoleBinding",
// 	},
// }
