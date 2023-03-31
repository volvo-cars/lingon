// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package tekton

import (
	"context"
	"errors"
	"os"
	"os/exec"

	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	AppName            = "tekton-pipelines"
	WebhookName        = "webhook"
	WebhookFullName    = AppName + "-" + WebhookName
	ControllerName     = "controller"
	ControllerFullName = AppName + "-" + ControllerName
	ResolversName      = "resolvers"
	ResolversFullName  = AppName + "-" + ResolversName

	Version         = "v0.46.0"
	WebhookPort     = 8443
	WebhookImage    = "gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/webhook:v0.46.0@sha256:5dc383dc1bd71d81180e0e4da68be966ebf383cfd0ac9f53a72cff11463e7f59"
	ControllerImage = "gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/controller:v0.46.0@sha256:d67fb2fb69ec38571ce3f71ce09571154e4b5db9b4cf71d69c2cb32455a4f8b4"
	ResolversImage  = "gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/resolvers:v0.46.0@sha256:f57448b914c72c03cbf36228134cc9ed24e28fef6d2e0d6d72c34908f38d8742"
)

var (
	PipelinesNS = kubeutil.Namespace(
		AppName, kubeutil.MergeLabels(
			labelsPipelines, map[string]string{
				kubeutil.NSLabelPodSecurityEnforce: kubeutil.NSValuePodSecurityRestricted,
			},
		), nil,
	)

	ResolversNS = kubeutil.Namespace(
		ResolversFullName,
		kubeutil.MergeLabels(
			labelsResolvers, map[string]string{
				kubeutil.NSLabelPodSecurityEnforce: kubeutil.NSValuePodSecurityRestricted,
			},
		), nil,
	)
)

// validate the struct implements the interface
var _ kube.Exporter = (*Tekton)(nil)

// Tekton contains kubernetes manifests
type Tekton struct {
	kube.App
	CRD

	PipelinesNS *corev1.Namespace
	ResolversNS *corev1.Namespace

	AggregateEditCR                              *rbacv1.ClusterRole
	AggregateViewCR                              *rbacv1.ClusterRole
	BundleResolverConfigCM                       *corev1.ConfigMap
	ClusterResolverConfigCM                      *corev1.ConfigMap
	ConfigDefaultsCM                             *corev1.ConfigMap
	ConfigLeaderElectionCM                       *corev1.ConfigMap
	ConfigLeaderElectionResolversCM              *corev1.ConfigMap
	ConfigLoggingCM                              *corev1.ConfigMap
	ConfigLoggingResolversCM                     *corev1.ConfigMap
	ConfigObservabilityCM                        *corev1.ConfigMap
	ConfigObservabilityResolversCM               *corev1.ConfigMap
	ConfigRegistryCertCM                         *corev1.ConfigMap
	ConfigSpireCM                                *corev1.ConfigMap
	FeatureFlagsCM                               *corev1.ConfigMap
	GitResolverConfigCM                          *corev1.ConfigMap
	HubResolverConfigCM                          *corev1.ConfigMap
	PipelinesControllerClusterAccessCR           *rbacv1.ClusterRole
	PipelinesControllerClusterAccessCRB          *rbacv1.ClusterRoleBinding
	PipelinesControllerDeploy                    *appsv1.Deployment
	PipelinesControllerLeaderElectionRB          *rbacv1.RoleBinding
	PipelinesControllerRB                        *rbacv1.RoleBinding
	PipelinesControllerRole                      *rbacv1.Role
	PipelinesControllerSA                        *corev1.ServiceAccount
	PipelinesControllerSVC                       *corev1.Service
	PipelinesControllerTenantAccessCR            *rbacv1.ClusterRole
	PipelinesControllerTenantAccessCRB           *rbacv1.ClusterRoleBinding
	PipelinesInfoCM                              *corev1.ConfigMap
	PipelinesInfoRB                              *rbacv1.RoleBinding
	PipelinesInfoRole                            *rbacv1.Role
	PipelinesLeaderElectionRole                  *rbacv1.Role
	PipelinesRemoteResolversDeploy               *appsv1.Deployment
	PipelinesResolversCRB                        *rbacv1.ClusterRoleBinding
	PipelinesResolversNamespaceRbacRB            *rbacv1.RoleBinding
	PipelinesResolversNamespaceRbacRole          *rbacv1.Role
	PipelinesResolversResolutionRequestUpdatesCR *rbacv1.ClusterRole
	PipelinesResolversSA                         *corev1.ServiceAccount
	PipelinesWebhookClusterAccessCR              *rbacv1.ClusterRole
	PipelinesWebhookClusterAccessCRB             *rbacv1.ClusterRoleBinding
	PipelinesWebhookDeploy                       *appsv1.Deployment
	PipelinesWebhookHPA                          *autoscalingv2.HorizontalPodAutoscaler
	PipelinesWebhookLeaderelectionRB             *rbacv1.RoleBinding
	PipelinesWebhookRB                           *rbacv1.RoleBinding
	PipelinesWebhookRole                         *rbacv1.Role
	PipelinesWebhookSA                           *corev1.ServiceAccount
	PipelinesWebhookSVC                          *corev1.Service
	ResolversFeatureFlagsCM                      *corev1.ConfigMap
	WebhookCertsSecrets                          *corev1.Secret
	ValidatePipelineWebhook                      *admissionregistrationv1.ValidatingWebhookConfiguration
	ValidateConfigPipelineWebhook                *admissionregistrationv1.ValidatingWebhookConfiguration
	MutatePipelineWebhook                        *admissionregistrationv1.MutatingWebhookConfiguration
}

type CRD struct {
	ClusterTasksDevCRD                 *apiextensionsv1.CustomResourceDefinition
	CustomRunsDevCRD                   *apiextensionsv1.CustomResourceDefinition
	PipelineRunsDevCRD                 *apiextensionsv1.CustomResourceDefinition
	PipelinesDevCRD                    *apiextensionsv1.CustomResourceDefinition
	ResolutionRequestsResolutionDevCRD *apiextensionsv1.CustomResourceDefinition
	RunsDevCRD                         *apiextensionsv1.CustomResourceDefinition
	TaskRunsDevCRD                     *apiextensionsv1.CustomResourceDefinition
	TasksDevCRD                        *apiextensionsv1.CustomResourceDefinition
	VerificationPoliciesDevCRD         *apiextensionsv1.CustomResourceDefinition
}

var labelsPipelines = map[string]string{
	kubeutil.AppLabelInstance: "default",
	kubeutil.AppLabelPartOf:   AppName,
}

var labelsVersion = map[string]string{
	kubeutil.AppLabelVersion:      Version,
	"pipeline.tekton.dev/release": Version,
	// labels below are related to istio and should not be used for resource lookup
	"version": Version,
}

var labelsResolvers = map[string]string{
	kubeutil.AppLabelComponent: ResolversName,
	kubeutil.AppLabelInstance:  "default",
	kubeutil.AppLabelPartOf:    AppName,
}

var PipelinesResolversSA = kubeutil.ServiceAccount(
	ResolversFullName,
	ResolversNS.Name,
	labelsResolvers,
	nil,
)

var labelsController = map[string]string{
	kubeutil.AppLabelComponent: ControllerName,
	kubeutil.AppLabelInstance:  "default",
	kubeutil.AppLabelPartOf:    AppName,
}

var PipelinesControllerSA = kubeutil.ServiceAccount(
	ControllerFullName,
	PipelinesNS.Name,
	labelsController,
	nil,
)

var labelsWebhook = map[string]string{
	kubeutil.AppLabelComponent: WebhookName,
	kubeutil.AppLabelInstance:  "default",
	kubeutil.AppLabelPartOf:    AppName,
}

var PipelinesWebhookSA = kubeutil.ServiceAccount(
	WebhookFullName,
	PipelinesNS.Name,
	labelsWebhook,
	nil,
)

// New creates a new Tekton
func New() *Tekton {
	return &Tekton{
		CRD: CRD{
			ClusterTasksDevCRD:                 ClusterTasksDevCRD,
			CustomRunsDevCRD:                   CustomRunsDevCRD,
			PipelineRunsDevCRD:                 PipelineRunsDevCRD,
			PipelinesDevCRD:                    PipelinesDevCRD,
			ResolutionRequestsResolutionDevCRD: ResolutionRequestsCRD,
			RunsDevCRD:                         RunsDevCRD,
			TaskRunsDevCRD:                     TaskRunsDevCRD,
			TasksDevCRD:                        TasksDevCRD,
			VerificationPoliciesDevCRD:         VerificationPoliciesDevCRD,
		},

		PipelinesNS: PipelinesNS,
		ResolversNS: ResolversNS,

		AggregateEditCR: AggregateEditCR,
		AggregateViewCR: AggregateViewCR,
		PipelinesInfoCM: PipelinesInfoCM,
		PipelinesInfoRB: &rbacv1.RoleBinding{
			// DO NOT REPLACE WITH kubeutil.BindRole, see Subjects!
			ObjectMeta: metav1.ObjectMeta{
				Labels:    labelsPipelines,
				Name:      "tekton-pipelines-info",
				Namespace: PipelinesNS.Name,
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
				Name:     "tekton-pipelines-info",
			},
			Subjects: []rbacv1.Subject{
				{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     "Group",
					Name:     "system:authenticated",
				},
			},
			TypeMeta: kubeutil.TypeRoleV1,
		},
		PipelinesInfoRole:              PipelinesInfoRole,
		PipelinesRemoteResolversDeploy: PipelinesRemoteResolversDeploy,

		/*
			CONFIGS
		*/

		BundleResolverConfigCM:          BundleResolverConfigCM,
		ClusterResolverConfigCM:         ClusterResolverConfigCM,
		ConfigDefaultsCM:                ConfigDefaultsCM,
		ConfigLeaderElectionCM:          ConfigLeaderElectionCM,
		ConfigLeaderElectionResolversCM: ConfigLeaderElectionResolversCM,
		ConfigLoggingCM:                 ConfigLoggingCM,
		ConfigLoggingResolversCM:        ConfigLoggingResolversCM,
		ConfigObservabilityCM:           ConfigObservabilityCM,
		ConfigObservabilityResolversCM:  ConfigObservabilityResolversCM,
		ConfigRegistryCertCM:            ConfigRegistryCertCM,
		ConfigSpireCM:                   ConfigSpireCM,
		FeatureFlagsCM:                  FeatureFlagsCM,
		GitResolverConfigCM:             GitResolverConfigCM,
		HubResolverConfigCM:             HubresolverConfigCM,

		/*
			RESOLVERS
		*/
		PipelinesResolversSA:                         PipelinesResolversSA,
		PipelinesResolversResolutionRequestUpdatesCR: PipelinesResolversResolutionRequestUpdatesCR,
		PipelinesResolversCRB: kubeutil.BindClusterRole(
			ResolversFullName,
			PipelinesResolversSA,
			PipelinesResolversResolutionRequestUpdatesCR,
			labelsResolvers,
		),
		PipelinesResolversNamespaceRbacRole: PipelinesResolversNamespaceRbacRole,
		PipelinesResolversNamespaceRbacRB: kubeutil.BindRole(
			ResolversFullName+"-namespace-rbac",
			PipelinesResolversSA,
			PipelinesResolversNamespaceRbacRole,
			labelsResolvers,
		),
		ResolversFeatureFlagsCM: ResolversFeatureFlagsCM,

		/*
			CONTROLLER
		*/
		PipelinesControllerDeploy: PipelinesControllerDeploy,
		PipelinesControllerSVC:    PipelinesControllerSVC,

		PipelinesControllerSA:   PipelinesControllerSA,
		PipelinesControllerRole: PipelinesControllerRole,
		PipelinesControllerRB: kubeutil.BindRole(
			ControllerFullName,
			PipelinesControllerSA,
			PipelinesControllerRole,
			labelsController,
		),
		PipelinesControllerClusterAccessCR: PipelinesControllerClusterAccessCR,
		PipelinesControllerClusterAccessCRB: kubeutil.BindClusterRole(
			ControllerFullName+"-cluster-access",
			PipelinesControllerSA,
			PipelinesControllerClusterAccessCR,
			labelsController,
		),
		PipelinesControllerLeaderElectionRB: kubeutil.BindRole(
			ControllerFullName+"-leaderelection",
			PipelinesControllerSA,
			PipelinesLeaderElectionRole,
			labelsController,
		),
		PipelinesControllerTenantAccessCR: PipelinesControllerTenantAccessCR,
		PipelinesControllerTenantAccessCRB: kubeutil.BindClusterRole(
			ControllerFullName+"-tenant-access",
			PipelinesControllerSA,
			PipelinesControllerTenantAccessCR,
			labelsController,
		),

		PipelinesLeaderElectionRole: PipelinesLeaderElectionRole,

		/*
			WEBHOOK
		*/
		PipelinesWebhookDeploy: kubeutil.SetDeploySA(
			PipelinesWebhookDeploy,
			PipelinesWebhookSA.Name,
		),
		PipelinesWebhookHPA: PipelinesWebhookHPA,
		PipelinesWebhookSA:  PipelinesWebhookSA,
		PipelinesWebhookSVC: PipelinesWebhookSVC,

		PipelinesWebhookClusterAccessCR: PipelinesWebhookClusterAccessCR,
		PipelinesWebhookClusterAccessCRB: kubeutil.BindClusterRole(
			PipelinesWebhookClusterAccessCR.Name,
			PipelinesWebhookSA,
			PipelinesWebhookClusterAccessCR,
			labelsWebhook,
		),

		PipelinesWebhookRole: PipelinesWebhookRole,
		PipelinesWebhookRB: kubeutil.BindRole(
			WebhookFullName,
			PipelinesWebhookSA,
			PipelinesWebhookRole,
			labelsWebhook,
		),
		PipelinesWebhookLeaderelectionRB: kubeutil.BindRole(
			WebhookFullName+"-leaderelection",
			PipelinesWebhookSA,
			PipelinesLeaderElectionRole,
			labelsWebhook,
		),

		WebhookCertsSecrets:           WebhookCertsSecrets,
		ValidatePipelineWebhook:       ValidationWebhookPipelineDevValidatingwebhookconfigurations,
		ValidateConfigPipelineWebhook: ConfigWebhookPipelineDevValidatingwebhookconfigurations,
		MutatePipelineWebhook:         WebhookPipelineDevMutatingwebhookconfigurations,
	}
}

// Apply applies the kubernetes objects to the cluster
func (a *Tekton) Apply(ctx context.Context) error {
	return Apply(ctx, a)
}

// Export exports the kubernetes objects to YAML files in the given directory
func (a *Tekton) Export(dir string) error {
	return kube.Export(a, kube.WithExportOutputDirectory(dir))
}

// P converts T to *T, useful for basic types
func P[T any](t T) *T {
	return &t
}

// Apply applies the kubernetes objects contained in Exporter to the cluster
func Apply(ctx context.Context, km kube.Exporter) error {
	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "-")
	cmd.Env = os.Environ()        // inherit environment in case we need to use kubectl from a container
	stdin, err := cmd.StdinPipe() // pipe to pass data to kubectl
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go func() {
		defer func() {
			err = errors.Join(err, stdin.Close())
		}()
		if errEW := kube.Export(
			km,
			kube.WithExportWriter(stdin),
		); errEW != nil {
			err = errors.Join(err, errEW)
		}
	}()

	if errS := cmd.Start(); errS != nil {
		return errors.Join(err, errS)
	}

	// waits for the command to exit and waits for any copying
	// to stdin or copying from stdout or stderr to complete
	return errors.Join(err, cmd.Wait())
}
