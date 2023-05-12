// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package karpenter

import (
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	ar "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var WebhookMutatingKarpenterAws = &ar.MutatingWebhookConfiguration{
	TypeMeta: kubeutil.TypeMutatingWebhookConfigurationV1,
	ObjectMeta: metav1.ObjectMeta{
		Name:   "defaulting.webhook.karpenter.k8s.aws",
		Labels: commonLabels,
	},
	Webhooks: []ar.MutatingWebhook{
		{
			Name:         "defaulting.webhook.karpenter.k8s.aws",
			ClientConfig: webHookClientConfig,
			Rules: []ar.RuleWithOperations{
				awsRuleWithNoDeleteOp,
				karpenterRuleWithNoDeleteOp,
			},
			FailurePolicy:           P(ar.Fail),
			SideEffects:             P(ar.SideEffectClassNone),
			AdmissionReviewVersions: []string{"v1"},
		},
	},
}

var awsRuleWithNoDeleteOp = ar.RuleWithOperations{
	Operations: createUpdateOps,
	Rule:       awsNodeTemplateRule,
}

var WebhookMutatingKarpenter = &ar.MutatingWebhookConfiguration{
	TypeMeta: kubeutil.TypeMutatingWebhookConfigurationV1,
	ObjectMeta: metav1.ObjectMeta{
		Name:   "defaulting.webhook.karpenter.sh",
		Labels: commonLabels,
	},
	Webhooks: []ar.MutatingWebhook{
		{
			Name:         "defaulting.webhook.karpenter.sh",
			ClientConfig: webHookClientConfig,
			Rules: []ar.RuleWithOperations{
				karpenterRuleWithNoDeleteOp,
			},
			FailurePolicy:           P(ar.Fail),
			SideEffects:             P(ar.SideEffectClassNone),
			AdmissionReviewVersions: []string{"v1"},
		},
	},
}

var karpenterRuleWithAllOperations = ar.RuleWithOperations{
	Operations: allOperations,
	Rule:       provisionersRule,
}

var karpenterRuleWithNoDeleteOp = ar.RuleWithOperations{
	Operations: createUpdateOps,
	Rule:       provisionersRule,
}

var webHookClientConfig = ar.WebhookClientConfig{
	Service: &ar.ServiceReference{
		Namespace: AppName,
		Name:      Namespace,
	},
}

var createUpdateOps = []ar.OperationType{
	ar.Create,
	ar.Update,
}

var allOperations = []ar.OperationType{
	ar.Create,
	ar.Update,
	ar.Delete,
}

var provisionersRule = ar.Rule{
	APIGroups:   []string{"karpenter.sh"},
	APIVersions: []string{"v1alpha5"},
	Resources: []string{
		"provisioners",
		"provisioners/status",
	},
}

var awsNodeTemplateRule = ar.Rule{
	APIGroups:   []string{"karpenter.k8s.aws"},
	APIVersions: []string{"v1alpha1"},
	Resources: []string{
		"awsnodetemplates",
		"awsnodetemplates/status",
	},
	Scope: P(ar.AllScopes),
}

var WebhookValidationKarpenter = &ar.ValidatingWebhookConfiguration{
	TypeMeta: kubeutil.TypeValidatingWebhookConfigurationV1,
	ObjectMeta: metav1.ObjectMeta{
		Name:   "validation.webhook.karpenter.sh",
		Labels: commonLabels,
	},
	Webhooks: []ar.ValidatingWebhook{
		{
			Name:         "validation.webhook.karpenter.sh",
			ClientConfig: webHookClientConfig,
			Rules: []ar.RuleWithOperations{
				{
					Operations: createUpdateOps,
					Rule:       provisionersRule,
				},
			},
			FailurePolicy:           P(ar.Fail),
			SideEffects:             P(ar.SideEffectClassNone),
			AdmissionReviewVersions: []string{"v1"},
		},
	},
}

var WebhookValidationKarpenterAWS = &ar.ValidatingWebhookConfiguration{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "admissionregistration.k8s.io/v1",
		Kind:       "ValidatingWebhookConfiguration",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:   "validation.webhook.karpenter.k8s.aws",
		Labels: commonLabels,
	},
	Webhooks: []ar.ValidatingWebhook{
		{
			AdmissionReviewVersions: []string{"v1"},
			ClientConfig: ar.WebhookClientConfig{
				Service: &ar.ServiceReference{
					Name:      "karpenter",
					Namespace: "karpenter",
				},
			},
			FailurePolicy: P(ar.FailurePolicyType("Fail")),
			Name:          "validation.webhook.karpenter.k8s.aws",
			Rules: []ar.RuleWithOperations{
				{
					Operations: []ar.OperationType{
						ar.OperationType("CREATE"),
						ar.OperationType("UPDATE"),
					},
					Rule: ar.Rule{
						APIGroups:   []string{"karpenter.k8s.aws"},
						APIVersions: []string{"v1alpha1"},
						Resources: []string{
							"awsnodetemplates",
							"awsnodetemplates/status",
						},
						Scope: P(ar.ScopeType("*")),
					},
				}, {
					Operations: []ar.OperationType{
						ar.OperationType("CREATE"),
						ar.OperationType("UPDATE"),
					},
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
			SideEffects: P(ar.SideEffectClass("None")),
		},
	},
}

var WebhookValidationKarpenterConfig = &ar.ValidatingWebhookConfiguration{
	TypeMeta: kubeutil.TypeValidatingWebhookConfigurationV1,
	ObjectMeta: metav1.ObjectMeta{
		Name:   "validation.webhook.config.karpenter.sh",
		Labels: commonLabels,
	},
	Webhooks: []ar.ValidatingWebhook{
		{
			Name:                    "validation.webhook.config.karpenter.sh",
			ClientConfig:            webHookClientConfig,
			ObjectSelector:          &metav1.LabelSelector{MatchLabels: map[string]string{"app.kubernetes.io/part-of": "karpenter"}},
			FailurePolicy:           P(ar.Fail),
			SideEffects:             P(ar.SideEffectClassNone),
			AdmissionReviewVersions: []string{"v1"},
		},
	},
}
