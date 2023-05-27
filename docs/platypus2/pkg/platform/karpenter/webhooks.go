// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package karpenter

import (
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	ar "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var createUpdateOps = []ar.OperationType{
	ar.Create,
	ar.Update,
}

// MUTATION WEBHOOK

var WebhookMutatingKarpenterAws = &ar.MutatingWebhookConfiguration{
	TypeMeta: kubeutil.TypeMutatingWebhookConfigurationV1,
	ObjectMeta: metav1.ObjectMeta{
		Name:   "defaulting.webhook.karpenter.k8s.aws",
		Labels: commonLabels,
	},
	Webhooks: []ar.MutatingWebhook{
		{
			Name:                    "defaulting.webhook.karpenter.k8s.aws",
			FailurePolicy:           P(ar.Fail),
			SideEffects:             P(ar.SideEffectClassNone),
			AdmissionReviewVersions: []string{"v1"},
			ClientConfig: ar.WebhookClientConfig{
				Service: &ar.ServiceReference{
					Namespace: Svc.Namespace,
					Name:      Svc.Name,
				},
			},
			Rules: []ar.RuleWithOperations{
				{
					Operations: createUpdateOps,
					Rule: ar.Rule{
						APIGroups:   []string{"karpenter.k8s.aws"},
						APIVersions: []string{"v1alpha1"},
						Resources: []string{
							"awsnodetemplates",
							"awsnodetemplates/status",
						},
						Scope: P(ar.AllScopes),
					},
				},
				{
					Operations: createUpdateOps,
					Rule: ar.Rule{
						APIGroups:   []string{"karpenter.sh"},
						APIVersions: []string{"v1alpha5"},
						Resources: []string{
							"provisioners",
							"provisioners/status",
						},
					},
				},
			},
		},
	},
}

// VALIDATION WEBHOOKS

var WebhookValidationKarpenterConfig = &ar.ValidatingWebhookConfiguration{
	TypeMeta: kubeutil.TypeValidatingWebhookConfigurationV1,
	ObjectMeta: metav1.ObjectMeta{
		Name:   "validation.webhook.config.karpenter.sh",
		Labels: commonLabels,
	},
	Webhooks: []ar.ValidatingWebhook{
		{
			Name:                    "validation.webhook.config.karpenter.sh",
			FailurePolicy:           P(ar.Fail),
			SideEffects:             P(ar.SideEffectClassNone),
			AdmissionReviewVersions: []string{"v1"},
			ClientConfig: ar.WebhookClientConfig{
				Service: &ar.ServiceReference{
					Namespace: Svc.Namespace,
					Name:      Svc.Name,
				},
			},
			ObjectSelector: &metav1.LabelSelector{
				// FIXME: not set on anything ??
				MatchLabels: map[string]string{"app.kubernetes.io/part-of": "karpenter"},
			},
		},
	},
}

var WebhookValidationKarpenter = &ar.ValidatingWebhookConfiguration{
	TypeMeta: kubeutil.TypeValidatingWebhookConfigurationV1,
	ObjectMeta: metav1.ObjectMeta{
		Name:   "validation.webhook.karpenter.sh",
		Labels: commonLabels,
	},
	Webhooks: []ar.ValidatingWebhook{
		{
			Name:                    "validation.webhook.karpenter.sh",
			FailurePolicy:           P(ar.Fail),
			SideEffects:             P(ar.SideEffectClassNone),
			AdmissionReviewVersions: []string{"v1"},
			ClientConfig: ar.WebhookClientConfig{
				Service: &ar.ServiceReference{
					Namespace: Svc.Namespace,
					Name:      Svc.Name,
				},
			},
			Rules: []ar.RuleWithOperations{
				{
					Operations: createUpdateOps,
					Rule: ar.Rule{
						APIGroups:   []string{"karpenter.sh"},
						APIVersions: []string{"v1alpha5"},
						Resources: []string{
							"provisioners",
							"provisioners/status",
						},
					},
				},
			},
		},
	},
}

var WebhookValidationKarpenterAWS = &ar.ValidatingWebhookConfiguration{
	TypeMeta: kubeutil.TypeValidatingWebhookConfigurationV1,
	ObjectMeta: metav1.ObjectMeta{
		Name:   "validation.webhook.karpenter.k8s.aws",
		Labels: commonLabels,
	},
	Webhooks: []ar.ValidatingWebhook{
		{
			Name:                    "validation.webhook.karpenter.k8s.aws",
			FailurePolicy:           P(ar.Fail),
			SideEffects:             P(ar.SideEffectClassNone),
			AdmissionReviewVersions: []string{"v1"},
			ClientConfig: ar.WebhookClientConfig{
				Service: &ar.ServiceReference{
					Namespace: Svc.Namespace,
					Name:      Svc.Name,
				},
			},

			Rules: []ar.RuleWithOperations{
				{
					Operations: createUpdateOps,
					Rule: ar.Rule{
						APIGroups:   []string{"karpenter.k8s.aws"},
						APIVersions: []string{"v1alpha1"},
						Resources: []string{
							"awsnodetemplates",
							"awsnodetemplates/status",
						},
						Scope: P(ar.AllScopes),
					},
				}, {
					Operations: createUpdateOps,
					Rule: ar.Rule{
						APIGroups:   []string{"karpenter.sh"},
						APIVersions: []string{"v1alpha5"},
						Resources: []string{
							"provisioners",
							"provisioners/status",
						},
					},
				},
			},
		},
	},
}
