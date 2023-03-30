// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package tekton

import (
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var WebhookPipelineDevMutatingwebhookconfigurations = &admissionregistrationv1.MutatingWebhookConfiguration{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/component": "webhook",
			"app.kubernetes.io/instance":  "default",
			"app.kubernetes.io/part-of":   "tekton-pipelines",
			"pipeline.tekton.dev/release": "v0.45.0",
		},
		Name: "webhook.pipeline.tekton.dev",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "admissionregistration.k8s.io/v1",
		Kind:       "MutatingWebhookConfiguration",
	},
	Webhooks: []admissionregistrationv1.MutatingWebhook{admissionregistrationv1.MutatingWebhook{
		AdmissionReviewVersions: []string{"v1"},
		ClientConfig: admissionregistrationv1.WebhookClientConfig{Service: &admissionregistrationv1.ServiceReference{
			Name:      "tekton-pipelines-webhook",
			Namespace: "tekton-pipelines",
		}},
		FailurePolicy: P(admissionregistrationv1.FailurePolicyType("Fail")),
		Name:          "webhook.pipeline.tekton.dev",
		SideEffects:   P(admissionregistrationv1.SideEffectClass("None")),
	}},
}
