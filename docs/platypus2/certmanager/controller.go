// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package certmanager

import (
	promoperatorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Controller struct {
	kube.App
	ControllerApproveIoCR                   *rbacv1.ClusterRole
	ControllerApproveIoCRB                  *rbacv1.ClusterRoleBinding
	ControllerCertificatesCR                *rbacv1.ClusterRole
	ControllerCertificatesCRB               *rbacv1.ClusterRoleBinding
	ControllerCertificateSigningRequestsCR  *rbacv1.ClusterRole
	ControllerCertificateSigningRequestsCRB *rbacv1.ClusterRoleBinding
	ControllerChallengesCR                  *rbacv1.ClusterRole
	ControllerChallengesCRB                 *rbacv1.ClusterRoleBinding
	ControllerClusterIssuersCR              *rbacv1.ClusterRole
	ControllerClusterIssuersCRB             *rbacv1.ClusterRoleBinding
	ControllerIngressShimCR                 *rbacv1.ClusterRole
	ControllerIngressShimCRB                *rbacv1.ClusterRoleBinding
	ControllerIssuersCR                     *rbacv1.ClusterRole
	ControllerIssuersCRB                    *rbacv1.ClusterRoleBinding
	ControllerOrdersCR                      *rbacv1.ClusterRole
	ControllerOrdersCRB                     *rbacv1.ClusterRoleBinding
	LeaderElectionRB                        *rbacv1.RoleBinding
	LeaderElectionRole                      *rbacv1.Role

	EditCR *rbacv1.ClusterRole
	ViewCR *rbacv1.ClusterRole

	ControllerDeploy *appsv1.Deployment
	PDB              *policyv1.PodDisruptionBudget
	SA               *corev1.ServiceAccount
	SVC              *corev1.Service
	ServiceMonitor   *promoperatorv1.ServiceMonitor
}

var ControllerDeploy = &appsv1.Deployment{
	TypeMeta:   ku.TypeDeploymentV1,
	ObjectMeta: CM.Controller.ObjectMeta(),
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(1)),
		Selector: &metav1.LabelSelector{MatchLabels: CM.Controller.MatchLabels()},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: CM.Controller.Labels()},
			Spec: corev1.PodSpec{
				ServiceAccountName: CM.Controller.ServiceAccount().Name,
				Containers: []corev1.Container{
					{
						Image:           CM.Controller.Img.URL(),
						ImagePullPolicy: corev1.PullIfNotPresent,
						Name:            CM.Controller.Name,
						Args: []string{
							"--v=2",
							"--cluster-resource-namespace=$(POD_NAMESPACE)",
							"--leader-election-namespace=" + ku.NSKubeSystem,
							"--acme-http01-solver-image=" + CM.ACMEImg.URL(),
							"--max-concurrent-challenges=60",
							"--enable-certificate-owner-ref=true",
						},
						Env: []corev1.EnvVar{
							ku.EnvVarDownAPI(
								"POD_NAMESPACE", "metadata.namespace",
							),
						},
						Ports: []corev1.ContainerPort{
							CM.ControllerPort.Container,
							{Name: "http-healthz", ContainerPort: int32(9403)},
						},
						Resources: ku.Resources(
							"10m", "32Mi", "100m", "64Mi",
						),
						SecurityContext: &corev1.SecurityContext{Capabilities: &corev1.Capabilities{Drop: []corev1.Capability{corev1.Capability("ALL")}}},
					},
				},
				NodeSelector: map[string]string{ku.LabelOSStable: "linux"},
				SecurityContext: &corev1.PodSecurityContext{
					RunAsNonRoot:   P(true),
					SeccompProfile: &corev1.SeccompProfile{Type: corev1.SeccompProfileType("RuntimeDefault")},
				},
			},
		},
	},
}

var ControllerSVC = &corev1.Service{
	TypeMeta:   ku.TypeServiceV1,
	ObjectMeta: CM.Controller.ObjectMeta(),
	Spec: corev1.ServiceSpec{
		Ports:    []corev1.ServicePort{CM.ControllerPort.Service},
		Selector: CM.Controller.MatchLabels(),
		Type:     corev1.ServiceTypeClusterIP,
	},
}

var ServiceMonitor = &v1.ServiceMonitor{
	ObjectMeta: CM.Controller.ObjectMeta(),
	Spec: v1.ServiceMonitorSpec{
		Endpoints: []v1.Endpoint{
			{
				Interval:      v1.Duration("60s"),
				Path:          "/metrics",
				ScrapeTimeout: v1.Duration("30s"),
				TargetPort:    P(intstr.FromString(CM.ControllerPort.Service.Name)),
			},
		},
		JobLabel: "cert-manager",
		Selector: metav1.LabelSelector{
			MatchLabels: CM.Controller.MatchLabels(),
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "monitoring.coreos.com/v1",
		Kind:       "ServiceMonitor",
	},
}

var PDB = &policyv1.PodDisruptionBudget{
	TypeMeta:   ku.TypePodDisruptionBudgetV1,
	ObjectMeta: CM.Controller.ObjectMeta(),
	Spec: policyv1.PodDisruptionBudgetSpec{
		MinAvailable: P(intstr.FromInt(1)),
		Selector:     &metav1.LabelSelector{MatchLabels: CM.Controller.MatchLabels()},
	},
}
