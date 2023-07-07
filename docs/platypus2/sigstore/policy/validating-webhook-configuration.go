// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package policy

import (
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var PolicySigstoreDevValidatingwebhookconfigurations = &admissionregistrationv1.ValidatingWebhookConfiguration{
	ObjectMeta: metav1.ObjectMeta{Name: "policy.sigstore.dev"},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "admissionregistration.k8s.io/v1",
		Kind:       "ValidatingWebhookConfiguration",
	},
	Webhooks: []admissionregistrationv1.ValidatingWebhook{
		{
			AdmissionReviewVersions: []string{"v1"},
			ClientConfig: admissionregistrationv1.WebhookClientConfig{
				Service: &admissionregistrationv1.ServiceReference{
					Name:      "webhook",
					Namespace: "sigstore",
				},
			},
			FailurePolicy: P(admissionregistrationv1.FailurePolicyType("Fail")),
			Name:          "policy.sigstore.dev",
			NamespaceSelector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "policy.sigstore.dev/include",
						Operator: metav1.LabelSelectorOperator("In"),
						Values:   []string{"true"},
					},
				},
			},
			SideEffects: P(admissionregistrationv1.SideEffectClass("None")),
		},
	},
}

var ValidatingClusterimagepolicySigstoreDevValidatingwebhookconfigurations = &admissionregistrationv1.ValidatingWebhookConfiguration{
	ObjectMeta: metav1.ObjectMeta{Name: "validating.clusterimagepolicy.sigstore.dev"},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "admissionregistration.k8s.io/v1",
		Kind:       "ValidatingWebhookConfiguration",
	},
	Webhooks: []admissionregistrationv1.ValidatingWebhook{
		{
			AdmissionReviewVersions: []string{"v1"},
			ClientConfig: admissionregistrationv1.WebhookClientConfig{
				Service: &admissionregistrationv1.ServiceReference{
					Name:      "webhook",
					Namespace: "sigstore",
				},
			},
			FailurePolicy: P(admissionregistrationv1.FailurePolicyType("Fail")),
			MatchPolicy:   P(admissionregistrationv1.MatchPolicyType("Equivalent")),
			Name:          "validating.clusterimagepolicy.sigstore.dev",
			SideEffects:   P(admissionregistrationv1.SideEffectClass("None")),
		},
	},
}
