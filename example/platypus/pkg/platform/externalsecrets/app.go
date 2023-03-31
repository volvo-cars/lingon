// Code generated by go-kart. EDIT AS MUCH AS YOU LIKE.

package externalsecrets

import (
	"github.com/volvo-cars/lingon/example/platypus/pkg/platform/externalsecrets/crd"
	"github.com/volvo-cars/lingon/example/platypus/pkg/platform/externalsecrets/webhook"

	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	arv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

const (
	AppName            = "external-secrets"
	webhookName        = AppName + "-webhook"
	controllerName     = AppName + "-controller"
	certControllerName = AppName + "-cert-controller"
	Namespace          = "external-secrets"
	Version            = "0.7.2"
	// FOR DEPLOYMENT
	containerImage         = "ghcr.io/external-secrets/external-secrets:v" + Version
	portMetric             = 8080
	healthzPort            = 8081
	webhookPort            = 10250
	webhookSecretMountPath = "/tmp/certs"
)

var _ kube.Exporter = (*ExternalSecrets)(nil)

type ExternalSecrets struct {
	kube.App
	CRD

	NS *corev1.Namespace

	SecretsWebhook *corev1.Secret

	CertControllerDeploy *appsv1.Deployment
	CertControllerCR     *rbacv1.ClusterRole
	CertControllerCRB    *rbacv1.ClusterRoleBinding
	CertControllerSA     *corev1.ServiceAccount

	ControllerCR  *rbacv1.ClusterRole
	ControllerCRB *rbacv1.ClusterRoleBinding

	EditCR *rbacv1.ClusterRole
	ViewCR *rbacv1.ClusterRole

	SA     *corev1.ServiceAccount
	Deploy *appsv1.Deployment

	WebhookSA     *corev1.ServiceAccount
	WebhookDeploy *appsv1.Deployment
	WebhookSVC    *corev1.Service

	ExternalSecretValidateWH *arv1.ValidatingWebhookConfiguration
	SecretStoreValidateWH    *arv1.ValidatingWebhookConfiguration

	LeaderElectionRole *rbacv1.Role
	LeaderElectionRB   *rbacv1.RoleBinding
}

type CRD struct {
	CrdACRAccessTokensGenerators        *apiextv1.CustomResourceDefinition
	CrdClusterExternalSecrets           *apiextv1.CustomResourceDefinition
	CrdClusterSecretStores              *apiextv1.CustomResourceDefinition
	CrdECRAuthorizationTokensGenerators *apiextv1.CustomResourceDefinition
	CrdExternalSecrets                  *apiextv1.CustomResourceDefinition
	CrdFakesGenerators                  *apiextv1.CustomResourceDefinition
	CrdGCRAccessTokensGenerators        *apiextv1.CustomResourceDefinition
	CrdPasswordsGenerators              *apiextv1.CustomResourceDefinition
	CrdPushSecrets                      *apiextv1.CustomResourceDefinition
	CrdSecretStores                     *apiextv1.CustomResourceDefinition
}

func New() *ExternalSecrets {
	sa := kubeutil.ServiceAccount(
		AppName,
		Namespace,
		ESLabels,
		nil,
	)
	certControllerSA := kubeutil.ServiceAccount(
		certControllerName,
		Namespace,
		CertControllerLabels,
		nil,
	)
	webhookSA := kubeutil.ServiceAccount(
		webhookName,
		Namespace,
		WebhookLabels,
		nil,
	)

	return &ExternalSecrets{
		NS: kubeutil.Namespace(Namespace, ESLabels, nil),

		CRD: CRD{
			CrdACRAccessTokensGenerators:        crd.AcraccesstokensGeneratorsCrd,
			CrdClusterExternalSecrets:           crd.ClusterexternalsecretsCrd,
			CrdClusterSecretStores:              crd.ClustersecretstoresCrd,
			CrdECRAuthorizationTokensGenerators: crd.EcrauthorizationtokensGeneratorsCrd,
			CrdExternalSecrets:                  crd.ExternalsecretsCrd,
			CrdFakesGenerators:                  crd.FakesGeneratorsCrd,
			CrdGCRAccessTokensGenerators:        crd.GcraccesstokensGeneratorsCrd,
			CrdPasswordsGenerators:              crd.PasswordsGeneratorsCrd,
			CrdPushSecrets:                      crd.PushsecretsCrd,
			CrdSecretStores:                     crd.SecretstoresCrd,
		},

		// CONTROLLER
		ControllerCR: ControllerCr,
		ControllerCRB: kubeutil.BindClusterRole(
			controllerName,
			sa,
			ControllerCr,
			ESLabels,
		),
		// CERT-CONTROLLER
		CertControllerDeploy: kubeutil.SetDeploySA(
			CertControllerDeploy,
			certControllerSA.Name,
		),
		CertControllerSA: certControllerSA,
		CertControllerCR: CertControllerCR,
		CertControllerCRB: kubeutil.BindClusterRole(
			certControllerName,
			sa,
			CertControllerCR,
			CertControllerLabels,
		),

		// Assumption: this is for applications to access the external secrets
		EditCR: EditCR,
		ViewCR: ViewCR,

		// Main External Secrets
		SA:     sa,
		Deploy: kubeutil.SetDeploySA(Deploy, sa.Name),

		// WEBHOOK
		ExternalSecretValidateWH: webhook.ExternalsecretValidateValidatingWH,
		SecretStoreValidateWH:    webhook.SecretstoreValidateValidatingWH,

		WebhookSA:     webhookSA,
		WebhookDeploy: kubeutil.SetDeploySA(WebhookDeploy, webhookSA.Name),
		WebhookSVC:    WebhookSVC,

		// Empty secret for webhook
		SecretsWebhook: &corev1.Secret{
			TypeMeta: kubeutil.TypeSecretV1,
			ObjectMeta: kubeutil.ObjectMeta(
				webhookName, // secret name used in CertControllerDeploy and WebhookDeploy
				Namespace,   // secret namespace used in CertControllerDeploy
				WebhookLabels,
				nil,
			),
			Data: nil, // okayyy ??
		},

		// LEADER ELECTION ?? not needed unless more than one instance
		LeaderElectionRB: kubeutil.BindRole(
			AppName+"-leaderelection",
			sa,
			LeaderElectionRole,
			ESLabels,
		),
		LeaderElectionRole: LeaderElectionRole,
	}
}

func P[T any](t T) *T {
	return &t
}
