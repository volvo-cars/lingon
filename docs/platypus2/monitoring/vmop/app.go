// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package vmop

import (
	"context"
	"errors"
	"os"
	"os/exec"

	"github.com/VictoriaMetrics/operator/api/victoriametrics/v1beta1"
	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	"github.com/volvo-cars/lingoneks/meta"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	// validate the struct implements the interface
	_ kube.Exporter = (*VMOperator)(nil)

	O = Core()

	SA = O.ServiceAccount()
)

func Core() Meta {
	n := "victoria-metrics-operator"
	ns := "monitoring"
	ver := "0.34.1"

	p := 8080
	pn := "http"

	return Meta{
		Metadata: meta.Metadata{
			Name:      n,
			Namespace: ns,
			Instance:  n + "-" + ns,
			Component: "operator",
			PartOf:    n,
			Version:   ver,
			ManagedBy: "lingon",
			Img: meta.ContainerImg{
				Registry: "",
				Image:    "victoriametrics/operator",
				Tag:      "v" + ver,
			},
		},

		VMVersion: "1.91.0",

		Main: meta.NetPort{
			Container: corev1.ContainerPort{
				Name:          pn,
				ContainerPort: int32(p),
				Protocol:      corev1.ProtocolTCP,
			},
			Service: corev1.ServicePort{
				Name:       pn,
				Port:       int32(p),
				TargetPort: intstr.FromString(pn),
				Protocol:   corev1.ProtocolTCP,
			},
		},
		Webhook: meta.NetPort{
			Container: corev1.ContainerPort{
				Name:          "webhook",
				ContainerPort: 9443,
				Protocol:      corev1.ProtocolTCP,
			},
			Service: corev1.ServicePort{
				Name:       "webhook",
				Port:       443,
				TargetPort: intstr.FromString("webhook"),
				Protocol:   corev1.ProtocolTCP,
			},
		},
	}
}

type Meta struct {
	meta.Metadata

	Webhook   meta.NetPort
	Main      meta.NetPort
	VMVersion string
}

// VMOperator contains kubernetes manifests
type VMOperator struct {
	kube.App

	NS        *corev1.Namespace
	CR        *rbacv1.ClusterRole
	CRB       *rbacv1.ClusterRoleBinding
	Deploy    *appsv1.Deployment
	PDB       *policyv1.PodDisruptionBudget
	RB        *rbacv1.RoleBinding
	Role      *rbacv1.Role
	SA        *corev1.ServiceAccount
	SVC       *corev1.Service
	SvcScrape *v1beta1.VMServiceScrape

	WHValidation *admissionregistrationv1.ValidatingWebhookConfiguration

	CRD
}

type CRD struct {
	VMAgentsCRD              *apiextensionsv1.CustomResourceDefinition
	VMAlertManagerConfigsCRD *apiextensionsv1.CustomResourceDefinition
	VMAlertManagersCRD       *apiextensionsv1.CustomResourceDefinition
	VMAlertsCRD              *apiextensionsv1.CustomResourceDefinition
	VMAuthsCRD               *apiextensionsv1.CustomResourceDefinition
	VMClustersCRD            *apiextensionsv1.CustomResourceDefinition
	VMNodeScrapesCRD         *apiextensionsv1.CustomResourceDefinition
	VMPodScrapesCRD          *apiextensionsv1.CustomResourceDefinition
	VMProbesCRD              *apiextensionsv1.CustomResourceDefinition
	VMRulesCRD               *apiextensionsv1.CustomResourceDefinition
	VMServiceScrapesCRD      *apiextensionsv1.CustomResourceDefinition
	VMSinglesCRD             *apiextensionsv1.CustomResourceDefinition
	VMStaticScrapesCRD       *apiextensionsv1.CustomResourceDefinition
	VMUsersCRD               *apiextensionsv1.CustomResourceDefinition
}

// New creates a new VMOperator
func New() *VMOperator {
	return &VMOperator{
		NS:           ku.Namespace(O.Namespace, nil, nil),
		WHValidation: WHValidation,
		CR:           CR,
		CRB:          ku.BindClusterRole(O.Name, SA, CR, O.Labels()),
		Deploy:       Deploy,
		RB:           ku.BindRole(O.Name, SA, Role, O.Labels()),
		Role:         Role,
		SA:           SA,
		PDB: &policyv1.PodDisruptionBudget{
			TypeMeta:   ku.TypePodDisruptionBudgetV1,
			ObjectMeta: O.ObjectMeta(),
			Spec: policyv1.PodDisruptionBudgetSpec{
				MinAvailable: P(intstr.FromInt(1)),
				Selector:     &metav1.LabelSelector{MatchLabels: O.MatchLabels()},
			},
		},

		SVC: &corev1.Service{
			TypeMeta:   ku.TypeServiceV1,
			ObjectMeta: O.ObjectMeta(),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					O.Main.Service,
					O.Webhook.Service,
				},
				Selector: O.MatchLabels(),
				Type:     corev1.ServiceTypeClusterIP,
			},
		},
		SvcScrape: &v1beta1.VMServiceScrape{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "operator.victoriametrics.com/v1beta1",
				Kind:       "VMServiceScrape",
			},
			ObjectMeta: O.ObjectMeta(),
			Spec: v1beta1.VMServiceScrapeSpec{
				Endpoints: []v1beta1.Endpoint{{Port: O.Main.Service.Name}},
				NamespaceSelector: v1beta1.NamespaceSelector{
					MatchNames: []string{O.Namespace},
				},
				Selector: metav1.LabelSelector{MatchLabels: O.MatchLabels()},
			},
		},

		CRD: CRD{
			VMAgentsCRD:              VMAgentsCRD,
			VMAlertManagerConfigsCRD: VMAlertManagerConfigsCRD,
			VMAlertManagersCRD:       VMAlertManagersCRD,
			VMAlertsCRD:              VMAlertsCRD,
			VMAuthsCRD:               VMAuthsCRD,
			VMClustersCRD:            VMClustersCRD,
			VMNodeScrapesCRD:         VMNodeScrapesCRD,
			VMPodScrapesCRD:          VMPodScrapesCRD,
			VMProbesCRD:              VMProbesCRD,
			VMRulesCRD:               VMRulesCRD,
			VMServiceScrapesCRD:      VMServiceScrapesCRD,
			VMSinglesCRD:             VMSinglesCRD,
			VMStaticScrapesCRD:       VMStaticScrapesCRD,
			VMUsersCRD:               VMUsersCRD,
		},
	}
}

// Apply applies the kubernetes objects to the cluster
func (a *VMOperator) Apply(ctx context.Context) error {
	return Apply(ctx, a)
}

// Export exports the kubernetes objects to YAML files in the given directory
func (a *VMOperator) Export(dir string) error {
	return kube.Export(a, kube.WithExportOutputDirectory(dir))
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
			kube.WithExportAsSingleFile("stdin"),
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

// P converts T to *T, useful for basic types
func P[T any](t T) *T {
	return &t
}
