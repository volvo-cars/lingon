package certmanager

import (
	"github.com/volvo-cars/lingon/pkg/kube"
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CaInjector struct {
	kube.App

	CaInjectorDeploy             *appsv1.Deployment
	CaInjectorCR                 *rbacv1.ClusterRole
	CaInjectorCRB                *rbacv1.ClusterRoleBinding
	CaInjectorLeaderelectionRB   *rbacv1.RoleBinding
	CaInjectorLeaderelectionRole *rbacv1.Role
	CaInjectorSA                 *corev1.ServiceAccount
}

var CaInjectorDeploy = &appsv1.Deployment{
	TypeMeta:   ku.TypeDeploymentV1,
	ObjectMeta: CM.CAInj.ObjectMeta(),
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(1)),
		Selector: &metav1.LabelSelector{MatchLabels: CM.CAInj.MatchLabels()},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: CM.CAInj.Labels()},
			Spec: corev1.PodSpec{
				ServiceAccountName: CM.CAInj.ServiceAccount().Name,
				Containers: []corev1.Container{
					{
						Name:            CM.CAInj.Name,
						Image:           CM.CAInj.Img.URL(),
						ImagePullPolicy: corev1.PullIfNotPresent,
						Args: []string{
							"--v=2",
							"--leader-election-namespace=" + ku.NSKubeSystem,
						},
						Env: []corev1.EnvVar{
							ku.EnvVarDownAPI(
								"POD_NAMESPACE", "metadata.namespace",
							),
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

var CainjectorCR = &rbacv1.ClusterRole{
	TypeMeta:   ku.TypeClusterRoleV1,
	ObjectMeta: CM.CAInj.ObjectMetaNoNS(),
	// ObjectMeta: metav1.ObjectMeta{
	// 	Labels: map[string]string{
	// 		"app":                          "cainjector",
	// 		"app.kubernetes.io/component":  "cainjector",
	// 		"app.kubernetes.io/instance":   "cert-manager",
	// 		"app.kubernetes.io/managed-by": "Helm",
	// 		"app.kubernetes.io/name":       "cainjector",
	// 		"app.kubernetes.io/version":    "v1.12.2",
	// 		"helm.sh/chart":                "cert-manager-v1.12.2",
	// 	},
	// 	Name: "cert-manager-cainjector",
	// },
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{"cert-manager.io"},
			Resources: []string{"certificates"},
			Verbs:     []string{"get", "list", "watch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"secrets"},
			Verbs:     []string{"get", "list", "watch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"events"},
			Verbs:     []string{"get", "create", "update", "patch"},
		}, {
			APIGroups: []string{"admissionregistration.k8s.io"},
			Resources: []string{
				"validatingwebhookconfigurations",
				"mutatingwebhookconfigurations",
			},
			Verbs: []string{"get", "list", "watch", "update", "patch"},
		}, {
			APIGroups: []string{"apiregistration.k8s.io"},
			Resources: []string{"apiservices"},
			Verbs:     []string{"get", "list", "watch", "update", "patch"},
		}, {
			APIGroups: []string{"apiextensions.k8s.io"},
			Resources: []string{"customresourcedefinitions"},
			Verbs:     []string{"get", "list", "watch", "update", "patch"},
		},
	},
}

var CaInjectorCRB = ku.BindClusterRole(
	CM.CAInj.Name, CM.CAInjSA, CainjectorCR, CM.CAInj.Labels(),
)

//	var CaInjectorCRB = &rbacv1.ClusterRoleBinding{
//		ObjectMeta: metav1.ObjectMeta{
//			Labels: map[string]string{
//				"app":                          "cainjector",
//				"app.kubernetes.io/component":  "cainjector",
//				"app.kubernetes.io/instance":   "cert-manager",
//				"app.kubernetes.io/managed-by": "Helm",
//				"app.kubernetes.io/name":       "cainjector",
//				"app.kubernetes.io/version":    "v1.12.2",
//				"helm.sh/chart":                "cert-manager-v1.12.2",
//			},
//			Name: "cert-manager-cainjector",
//		},
//		RoleRef: rbacv1.RoleRef{
//			APIGroup: "rbac.authorization.k8s.io",
//			Kind:     "ClusterRole",
//			Name:     "cert-manager-cainjector",
//		},
//		Subjects: []rbacv1.Subject{
//			{
//				Kind:      "ServiceAccount",
//				Name:      "cert-manager-cainjector",
//				Namespace: "cert-manager",
//			},
//		},
//		TypeMeta: metav1.TypeMeta{
//			APIVersion: "rbac.authorization.k8s.io/v1",
//			Kind:       "ClusterRoleBinding",
//		},
//	}

var CaInjectorLeaderElectionRole = &rbacv1.Role{
	TypeMeta: ku.TypeRoleV1,
	ObjectMeta: metav1.ObjectMeta{
		Labels:    CM.CAInj.Labels(),
		Name:      CM.CAInj.Name + ":leaderelection",
		Namespace: ku.NSKubeSystem,
	},
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{"coordination.k8s.io"},
			ResourceNames: []string{
				"cert-manager-cainjector-leader-election",
				"cert-manager-cainjector-leader-election-core",
			},
			Resources: []string{"leases"},
			Verbs:     []string{"get", "update", "patch"},
		}, {
			APIGroups: []string{"coordination.k8s.io"},
			Resources: []string{"leases"},
			Verbs:     []string{"create"},
		},
	},
}

var CaInjectorLeaderElectionRB = ku.BindRole(
	CaInjectorLeaderElectionRole.Name,
	CM.CAInjSA,
	CaInjectorLeaderElectionRole,
	CM.CAInj.Labels(),
)

// var CaInjectorLeaderElectionRB = &rbacv1.RoleBinding{
// 	ObjectMeta: metav1.ObjectMeta{
// 		Labels: map[string]string{
// 			"app":                          "cainjector",
// 			"app.kubernetes.io/component":  "cainjector",
// 			"app.kubernetes.io/instance":   "cert-manager",
// 			"app.kubernetes.io/managed-by": "Helm",
// 			"app.kubernetes.io/name":       "cainjector",
// 			"app.kubernetes.io/version":    "v1.12.2",
// 			"helm.sh/chart":                "cert-manager-v1.12.2",
// 		},
// 		Name:      "cert-manager-cainjector:leaderelection",
// 		Namespace: "kube-system",
// 	},
// 	RoleRef: rbacv1.RoleRef{
// 		APIGroup: "rbac.authorization.k8s.io",
// 		Kind:     "Role",
// 		Name:     "cert-manager-cainjector:leaderelection",
// 	},
// 	Subjects: []rbacv1.Subject{
// 		{
// 			Kind:      "ServiceAccount",
// 			Name:      "cert-manager-cainjector",
// 			Namespace: "cert-manager",
// 		},
// 	},
// 	TypeMeta: metav1.TypeMeta{
// 		APIVersion: "rbac.authorization.k8s.io/v1",
// 		Kind:       "RoleBinding",
// 	},
// }
