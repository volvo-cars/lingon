// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package certmanager

import (
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var WebhookValidatingwebhookconfigurations = &admissionregistrationv1.ValidatingWebhookConfiguration{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{"cert-manager.io/inject-ca-from-secret": "cert-manager/cert-manager-webhook-ca"},
		Labels: map[string]string{
			"app":                          "webhook",
			"app.kubernetes.io/component":  "webhook",
			"app.kubernetes.io/instance":   "cert-manager",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "webhook",
			"app.kubernetes.io/version":    "v1.12.2",
			"helm.sh/chart":                "cert-manager-v1.12.2",
		},
		Name: "cert-manager-webhook",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "admissionregistration.k8s.io/v1",
		Kind:       "ValidatingWebhookConfiguration",
	},
	Webhooks: []admissionregistrationv1.ValidatingWebhook{
		{
			AdmissionReviewVersions: []string{"v1"},
			ClientConfig: admissionregistrationv1.WebhookClientConfig{
				Service: &admissionregistrationv1.ServiceReference{
					Name:      "cert-manager-webhook",
					Namespace: "cert-manager",
					Path:      P("/validate"),
				},
			},
			FailurePolicy: P(admissionregistrationv1.FailurePolicyType("Fail")),
			MatchPolicy:   P(admissionregistrationv1.MatchPolicyType("Equivalent")),
			Name:          "webhook.cert-manager.io",
			NamespaceSelector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "cert-manager.io/disable-validation",
						Operator: metav1.LabelSelectorOperator("NotIn"),
						Values:   []string{"true"},
					}, {
						Key:      "name",
						Operator: metav1.LabelSelectorOperator("NotIn"),
						Values:   []string{"cert-manager"},
					},
				},
			},
			Rules: []admissionregistrationv1.RuleWithOperations{
				{
					Operations: []admissionregistrationv1.OperationType{
						admissionregistrationv1.OperationType("CREATE"),
						admissionregistrationv1.OperationType("UPDATE"),
					},
					Rule: admissionregistrationv1.Rule{
						APIGroups: []string{
							"cert-manager.io",
							"acme.cert-manager.io",
						},
						APIVersions: []string{"v1"},
						Resources:   []string{"*/*"},
					},
				},
			},
			SideEffects:    P(admissionregistrationv1.SideEffectClass("None")),
			TimeoutSeconds: P(int32(10)),
		},
	},
}
