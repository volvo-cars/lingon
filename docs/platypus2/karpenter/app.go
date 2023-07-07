// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package karpenter

import (
	promoperatorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	"github.com/volvo-cars/lingoneks/meta"
	ar "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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

var KA = Core()

func Core() Meta {
	AppName := "karpenter"
	version := "0.29.0"
	metricsPort := 8000
	metricsPortName := "http-metrics"
	webhookPort := 8443
	webhookPortName := "https-webhook"

	m := meta.Metadata{
		Name:      AppName,
		Namespace: AppName,
		Instance:  AppName,
		Component: AppName,
		PartOf:    AppName,
		Version:   version,
		ManagedBy: "lingon",
		Img: meta.ContainerImg{
			Registry: "public.ecr.aws/karpenter",
			Image:    "controller",
			Sha:      "3009f10487d9338f77c325adee3c208513cd06c7f191653327ef3a44006bf9c8",
			Tag:      "v" + version,
		},
	}
	return Meta{
		Metadata: m,
		Probe: meta.NetPort{
			Container: corev1.ContainerPort{
				Name:          "http",
				ContainerPort: 8081,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		Metrics: meta.NetPort{
			Container: corev1.ContainerPort{
				Name:          metricsPortName,
				ContainerPort: int32(metricsPort),
				Protocol:      corev1.ProtocolTCP,
			},
			Service: corev1.ServicePort{
				Name:     metricsPortName,
				Port:     int32(metricsPort),
				Protocol: corev1.ProtocolTCP,
			},
		},
		Webhook: meta.NetPort{
			Container: corev1.ContainerPort{
				Name:          webhookPortName,
				ContainerPort: int32(webhookPort),
				Protocol:      corev1.ProtocolTCP,
			},
			Service: corev1.ServicePort{
				Name:       webhookPortName,
				Port:       int32(webhookPort),
				TargetPort: intstr.FromString(webhookPortName),
				Protocol:   corev1.ProtocolTCP,
			},
		},

		ConfigName:  m.Name + "-global-settings",
		KubeMinVer:  "1.19.0-0",
		ProfileName: "default",
	}
}

type Meta struct {
	meta.Metadata
	Probe   meta.NetPort
	Metrics meta.NetPort
	Webhook meta.NetPort

	ConfigName  string
	KubeMinVer  string
	ProfileName string
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
	Deploy     *appsv1.Deployment
	Svc        *corev1.Service
	Pdb        *policyv1.PodDisruptionBudget
	SvcMonitor *promoperatorv1.ServiceMonitor

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
	SA := ku.ServiceAccount(
		KA.Name, KA.Namespace, KA.Labels(),
		map[string]string{"eks.amazonaws.com/role-arn": opts.IAMRoleArn},
	)
	return &Karpenter{
		Ns: ku.Namespace(KA.Namespace, KA.Labels(), nil),

		CertSecret:    CertSecret,
		Settings:      GlobalSettings(opts),
		LoggingConfig: LoggingConfig,

		Deploy: ku.SetDeploySA(Deploy, SA.Name),
		Svc:    Svc,
		Pdb:    Pdb,
		SvcMonitor: &promoperatorv1.ServiceMonitor{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "monitoring.coreos.com/v1",
				Kind:       "ServiceMonitor",
			},
			ObjectMeta: KA.ObjectMeta(),
			// ObjectMeta: metav1.ObjectMeta{
			// 	Name:      KA.Name,
			// 	Namespace: monitoring.Namespace,
			// 	Labels:    KA.Labels(),
			// },
			Spec: promoperatorv1.ServiceMonitorSpec{
				JobLabel: ku.AppLabelName,
				Endpoints: []promoperatorv1.Endpoint{
					{Path: ku.PathMetrics, Port: KA.Metrics.Service.Name},
				},
				Selector:          metav1.LabelSelector{MatchLabels: KA.MatchLabels()},
				NamespaceSelector: promoperatorv1.NamespaceSelector{Any: true},
			},
		},

		SA:      SA,
		DNSRole: DnsRole,
		DNSRb:   DnsRoleBinding,
		Role:    Role,
		Rb: ku.BindRole(
			"karpenter-rb", SA, Role, KA.Labels(),
		),

		CR: CanUpdateWebhooks,
		CRB: ku.BindClusterRole(
			"karpenter-crb-hook", SA, CanUpdateWebhooks, KA.Labels(),
		),
		CoreCR: CoreCr,
		CoreCRB: ku.BindClusterRole(
			"karpenter-crb-core", SA, CoreCr, KA.Labels(),
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
