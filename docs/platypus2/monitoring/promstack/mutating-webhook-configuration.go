// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package promstack

import (
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var KubePromtheusStackKubeAdmissionMutatingwebhookconfigurations = &admissionregistrationv1.MutatingWebhookConfiguration{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-admission",
			"app.kubernetes.io/instance":   "kube-promtheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "46.8.0",
			"chart":                        "kube-prometheus-stack-46.8.0",
			"heritage":                     "Helm",
			"release":                      "kube-promtheus-stack",
		},
		Name: "kube-promtheus-stack-kube-admission",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "admissionregistration.k8s.io/v1",
		Kind:       "MutatingWebhookConfiguration",
	},
	Webhooks: []admissionregistrationv1.MutatingWebhook{{
		AdmissionReviewVersions: []string{"v1", "v1beta1"},
		ClientConfig: admissionregistrationv1.WebhookClientConfig{Service: &admissionregistrationv1.ServiceReference{
			Name:      "kube-promtheus-stack-kube-operator",
			Namespace: "monitoring",
			Path:      P("/admission-prometheusrules/mutate"),
		}},
		FailurePolicy: P(admissionregistrationv1.FailurePolicyType("Ignore")),
		Name:          "prometheusrulemutate.monitoring.coreos.com",
		Rules: []admissionregistrationv1.RuleWithOperations{{
			Operations: []admissionregistrationv1.OperationType{admissionregistrationv1.OperationType("CREATE"), admissionregistrationv1.OperationType("UPDATE")},
			Rule: admissionregistrationv1.Rule{
				APIGroups:   []string{"monitoring.coreos.com"},
				APIVersions: []string{"*"},
				Resources:   []string{"prometheusrules"},
			},
		}},
		SideEffects:    P(admissionregistrationv1.SideEffectClass("None")),
		TimeoutSeconds: P(int32(10)),
	}},
}
