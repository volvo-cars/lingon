package karpenter

import (
	"github.com/volvo-cars/lingon/example/platypus/pkg/platform/karpenter/crd"

	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	ar "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ kube.Exporter = (*Karpenter)(nil)

const (
	AppName   = "karpenter"
	Namespace = "karpenter"
	Version   = "0.25.0"
)

var commonLabels = map[string]string{
	kubeutil.AppLabelInstance:  AppName,
	kubeutil.AppLabelManagedBy: "lingon",
	kubeutil.AppLabelName:      AppName,
	kubeutil.AppLabelVersion:   Version,
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

	CustomResourceDefinitions

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

	WHMutation    *ar.MutatingWebhookConfiguration
	WHMutationAWS *ar.MutatingWebhookConfiguration
}

type CustomResourceDefinitions struct {
	AWSNodeTemplates *apiextensionsv1.CustomResourceDefinition
	Provisioner      *apiextensionsv1.CustomResourceDefinition
}

type Opts struct {
	ClusterName            string
	ClusterEndpoint        string
	IAMRoleArn             string
	DefaultInstanceProfile string
	InterruptQueue         string
}

func New(
	opts Opts,
) *Karpenter {
	sacc := &corev1.ServiceAccount{
		TypeMeta: kubeutil.TypeServiceAccountV1,
		ObjectMeta: kubeutil.ObjectMeta(
			AppName,
			Namespace,
			commonLabels,
			map[string]string{"eks.amazonaws.com/role-arn": opts.IAMRoleArn},
		),
	}

	return &Karpenter{
		CustomResourceDefinitions: CustomResourceDefinitions{
			AWSNodeTemplates: crd.AwsnodetemplatesKarpenterK8SAwsCRD,
			Provisioner:      crd.ProvisionersKarpenterShCRD,
		},

		Ns: &corev1.Namespace{
			TypeMeta: kubeutil.TypeNamespaceV1,
			ObjectMeta: metav1.ObjectMeta{
				Name:   Namespace,
				Labels: commonLabels,
			},
			Spec: corev1.NamespaceSpec{},
		},
		CertSecret:    CertSecret,
		Settings:      GlobalSettings(opts),
		LoggingConfig: LoggingConfig,

		Deploy: kubeutil.SetDeploySA(Deploy, sacc.Name),
		Svc:    Svc,
		Pdb:    Pdb,

		SA:      sacc,
		DNSRole: DnsRole,
		DNSRb:   DnsRoleBinding,
		Role:    Role,
		Rb: kubeutil.BindRole(
			"karpenter-rb",
			sacc,
			Role,
			commonLabels,
		),

		CR: CanUpdateWebhooks,
		CRB: kubeutil.BindClusterRole(
			"karpenter-crb-hook",
			sacc,
			CanUpdateWebhooks,
			commonLabels,
		),
		CoreCR: CoreCr,
		CoreCRB: kubeutil.BindClusterRole(
			"karpenter-crb-core",
			sacc,
			CoreCr,
			commonLabels,
		),
		AdminCR: AdminCr,

		WHValidation:       WebhookValidationKarpenter,
		WHValidationAWS:    WebhookValidationKarpenterAWS,
		WHValidationConfig: WebhookValidationKarpenterConfig,
		WHMutation:         WebhookMutatingKarpenter,
		WHMutationAWS:      WebhookMutatingKarpenterAws,
	}
}
