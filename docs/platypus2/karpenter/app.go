// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package karpenter

import (
	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	ar "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ kube.Exporter = (*Karpenter)(nil)

// KARPENTER
//
// to activate spot instance:
//		aws --profile=XXX iam create-service-linked-role --aws-service-name spot.amazonaws.com
//
// to see the logs:
// 		kubectl logs -f -n karpenter -l app.kubernetes.io/name=karpenter -c controller
//

const (
	AppName         = "karpenter"
	Namespace       = "karpenter"
	Version         = "0.27.5"
	ImgCtrlBase     = "public.ecr.aws/karpenter/controller"
	ImgSha          = "f9023101d05d0c0c6a5d67f19b8ecf754bf97cb4e94b41d9d80a75ee5be5150c"
	ImgController   = ImgCtrlBase + ":v" + Version + "@sha256:" + ImgSha
	ConfigName      = AppName + "-global-settings"
	PortNameMetrics = "http-metrics"
	PortMetrics     = 8080

	PortNameProbe = "http"
	PortProbe     = 8081

	PortNameWebhook = "https-webhook"
	PortWebhookSvc  = 443
	PortWebhookCtnr = 8443
)

var commonLabels = map[string]string{
	ku.AppLabelInstance:  AppName,
	ku.AppLabelManagedBy: "lingon",
	ku.AppLabelName:      AppName,
	ku.AppLabelVersion:   Version,
}

var matchLabels = map[string]string{
	ku.AppLabelInstance: AppName,
	ku.AppLabelName:     AppName,
}

func appendCommonLabels(items map[string]string) map[string]string {
	m := map[string]string{}
	for n, v := range commonLabels {
		m[n] = v
	}
	for n, v := range items {
		m[n] = v
	}
	return m
}

type Karpenter struct {
	kube.App

	Ns *corev1.Namespace
	// Configuration
	CertSecret *corev1.Secret
	Settings   *corev1.ConfigMap
	// LoggingConfig is not mounted but can be modified thanks to Role
	LoggingConfig *corev1.ConfigMap

	// Application
	Deploy *appsv1.Deployment
	Svc    *corev1.Service
	Pdb    *policyv1.PodDisruptionBudget

	// IAM
	SA *corev1.ServiceAccount

	DNSRole *rbacv1.Role
	DNSRb   *rbacv1.RoleBinding
	Role    *rbacv1.Role
	Rb      *rbacv1.RoleBinding

	// IAM cluster
	CR      *rbacv1.ClusterRole
	CRB     *rbacv1.ClusterRoleBinding
	CoreCR  *rbacv1.ClusterRole
	CoreCRB *rbacv1.ClusterRoleBinding
	AdminCR *rbacv1.ClusterRole
	// AdminCRB *rbacv1.ClusterRoleBinding // ???

	// Webhooks
	WHValidation       *ar.ValidatingWebhookConfiguration
	WHValidationAWS    *ar.ValidatingWebhookConfiguration
	WHValidationConfig *ar.ValidatingWebhookConfiguration

	// WHMutation    *ar.MutatingWebhookConfiguration
	WHMutationAWS *ar.MutatingWebhookConfiguration
}

type Opts struct {
	ClusterName            string
	ClusterEndpoint        string
	IAMRoleArn             string
	DefaultInstanceProfile string
	InterruptQueue         string
}

func New(opts Opts) *Karpenter {
	sacc := &corev1.ServiceAccount{
		TypeMeta: ku.TypeServiceAccountV1,
		ObjectMeta: ku.ObjectMeta(
			AppName,
			Namespace,
			commonLabels,
			map[string]string{"eks.amazonaws.com/role-arn": opts.IAMRoleArn},
		),
	}

	return &Karpenter{
		Ns: &corev1.Namespace{
			TypeMeta: ku.TypeNamespaceV1,
			ObjectMeta: metav1.ObjectMeta{
				Name:   Namespace,
				Labels: commonLabels,
			},
			Spec: corev1.NamespaceSpec{},
		},
		CertSecret:    CertSecret,
		Settings:      GlobalSettings(opts),
		LoggingConfig: LoggingConfig,

		Deploy: ku.SetDeploySA(Deploy, sacc.Name),
		Svc:    Svc,
		Pdb:    Pdb,

		SA:      sacc,
		DNSRole: DnsRole,
		DNSRb:   DnsRoleBinding,
		Role:    Role,
		Rb: ku.BindRole(
			"karpenter-rb",
			sacc,
			Role,
			commonLabels,
		),

		CR: CanUpdateWebhooks,
		CRB: ku.BindClusterRole(
			"karpenter-crb-hook",
			sacc,
			CanUpdateWebhooks,
			commonLabels,
		),
		CoreCR: CoreCr,
		CoreCRB: ku.BindClusterRole(
			"karpenter-crb-core",
			sacc,
			CoreCr,
			commonLabels,
		),
		AdminCR: AdminCr,

		WHValidation:       WebhookValidationKarpenter,
		WHValidationAWS:    WebhookValidationKarpenterAWS,
		WHValidationConfig: WebhookValidationKarpenterConfig,
		// WHMutation:         WebhookMutatingKarpenter, // removed when updating to 0.27.5
		WHMutationAWS: WebhookMutatingKarpenterAws,
	}
}

// P returns a pointer to the given value.
func P[T any](t T) *T {
	return &t
}
