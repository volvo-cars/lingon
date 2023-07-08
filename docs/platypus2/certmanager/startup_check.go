package certmanager

import (
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StartUpAPICheck struct {
	StartupapicheckCreateCertRB   *rbacv1.RoleBinding
	StartupapicheckCreateCertRole *rbacv1.Role
	StartupapicheckJOBS           *batchv1.Job
	StartupapicheckSA             *corev1.ServiceAccount
}

func NewStartUpCheck() StartUpAPICheck {
	return StartUpAPICheck{
		StartupapicheckCreateCertRB:   StartupapicheckCreateCertRB,
		StartupapicheckCreateCertRole: StartupapicheckCreateCertRole,
		StartupapicheckJOBS:           StartupapicheckJOBS,
		StartupapicheckSA:             StartupapicheckSA,
	}
}

var StartupapicheckJOBS = &batchv1.Job{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "post-install",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
			"helm.sh/hook-weight":        "1",
		},
		Labels: map[string]string{
			"app":                          "startupapicheck",
			"app.kubernetes.io/component":  "startupapicheck",
			"app.kubernetes.io/instance":   "cert-manager",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "startupapicheck",
			"app.kubernetes.io/version":    "v1.12.2",
			"helm.sh/chart":                "cert-manager-v1.12.2",
		},
		Name:      "cert-manager-startupapicheck",
		Namespace: "cert-manager",
	},
	Spec: batchv1.JobSpec{
		BackoffLimit: P(int32(4)),
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app":                          "startupapicheck",
					"app.kubernetes.io/component":  "startupapicheck",
					"app.kubernetes.io/instance":   "cert-manager",
					"app.kubernetes.io/managed-by": "Helm",
					"app.kubernetes.io/name":       "startupapicheck",
					"app.kubernetes.io/version":    "v1.12.2",
					"helm.sh/chart":                "cert-manager-v1.12.2",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Args:            []string{"check", "api", "--wait=1m"},
						Image:           "quay.io/jetstack/cert-manager-ctl:v1.12.2",
						ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
						Name:            "cert-manager-startupapicheck",
						SecurityContext: &corev1.SecurityContext{Capabilities: &corev1.Capabilities{Drop: []corev1.Capability{corev1.Capability("ALL")}}},
					},
				},
				NodeSelector:  map[string]string{"kubernetes.io/os": "linux"},
				RestartPolicy: corev1.RestartPolicy("OnFailure"),
				SecurityContext: &corev1.PodSecurityContext{
					RunAsNonRoot:   P(true),
					SeccompProfile: &corev1.SeccompProfile{Type: corev1.SeccompProfileType("RuntimeDefault")},
				},
				ServiceAccountName: "cert-manager-startupapicheck",
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "batch/v1",
		Kind:       "Job",
	},
}

var StartupapicheckCreateCertRole = &rbacv1.Role{
	TypeMeta: ku.TypeRoleV1,
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "post-install",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
			"helm.sh/hook-weight":        "-5",
		},
		Labels: map[string]string{
			"app":                          "startupapicheck",
			"app.kubernetes.io/component":  "startupapicheck",
			"app.kubernetes.io/instance":   "cert-manager",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "startupapicheck",
			"app.kubernetes.io/version":    "v1.12.2",
			"helm.sh/chart":                "cert-manager-v1.12.2",
		},
		Name:      "cert-manager-startupapicheck:create-cert",
		Namespace: "cert-manager",
	},
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups: []string{"cert-manager.io"},
			Resources: []string{"certificates"},
			Verbs:     []string{"create"},
		},
	},
}

var StartupapicheckCreateCertRB = &rbacv1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "post-install",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
			"helm.sh/hook-weight":        "-5",
		},
		Labels: map[string]string{
			"app":                          "startupapicheck",
			"app.kubernetes.io/component":  "startupapicheck",
			"app.kubernetes.io/instance":   "cert-manager",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "startupapicheck",
			"app.kubernetes.io/version":    "v1.12.2",
			"helm.sh/chart":                "cert-manager-v1.12.2",
		},
		Name:      "cert-manager-startupapicheck:create-cert",
		Namespace: "cert-manager",
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "Role",
		Name:     "cert-manager-startupapicheck:create-cert",
	},
	Subjects: []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      "cert-manager-startupapicheck",
			Namespace: "cert-manager",
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "RoleBinding",
	},
}

var StartupapicheckSA = &corev1.ServiceAccount{
	AutomountServiceAccountToken: P(true),
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "post-install",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
			"helm.sh/hook-weight":        "-5",
		},
		Labels: map[string]string{
			"app":                          "startupapicheck",
			"app.kubernetes.io/component":  "startupapicheck",
			"app.kubernetes.io/instance":   "cert-manager",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "startupapicheck",
			"app.kubernetes.io/version":    "v1.12.2",
			"helm.sh/chart":                "cert-manager-v1.12.2",
		},
		Name:      "cert-manager-startupapicheck",
		Namespace: "cert-manager",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ServiceAccount",
	},
}
