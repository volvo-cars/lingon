// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package promstack

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var KubePromtheusStackKubeAdmissionCreateJOBS = &batchv1.Job{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "pre-install,pre-upgrade",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
		},
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-admission-create",
			"app.kubernetes.io/instance":   "kube-promtheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "47.0.0",
			"chart":                        "kube-prometheus-stack-47.0.0",
			"heritage":                     "Helm",
			"release":                      "kube-promtheus-stack",
		},
		Name:      "kube-promtheus-stack-kube-admission-create",
		Namespace: "monitoring",
	},
	Spec: batchv1.JobSpec{Template: corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app":                          "kube-prometheus-stack-admission-create",
				"app.kubernetes.io/instance":   "kube-promtheus-stack",
				"app.kubernetes.io/managed-by": "Helm",
				"app.kubernetes.io/part-of":    "kube-prometheus-stack",
				"app.kubernetes.io/version":    "47.0.0",
				"chart":                        "kube-prometheus-stack-47.0.0",
				"heritage":                     "Helm",
				"release":                      "kube-promtheus-stack",
			},
			Name: "kube-promtheus-stack-kube-admission-create",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Args:            []string{"create", "--host=kube-promtheus-stack-kube-operator,kube-promtheus-stack-kube-operator.monitoring.svc", "--namespace=monitoring", "--secret-name=kube-promtheus-stack-kube-admission"},
				Image:           "registry.k8s.io/ingress-nginx/kube-webhook-certgen:v20221220-controller-v1.5.1-58-g787ea74b6",
				ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
				Name:            "create",
				SecurityContext: &corev1.SecurityContext{
					Capabilities:           &corev1.Capabilities{Drop: []corev1.Capability{corev1.Capability("ALL")}},
					ReadOnlyRootFilesystem: P(true),
				},
			}},
			RestartPolicy: corev1.RestartPolicy("OnFailure"),
			SecurityContext: &corev1.PodSecurityContext{
				RunAsGroup:     P(int64(2000)),
				RunAsNonRoot:   P(true),
				RunAsUser:      P(int64(2000)),
				SeccompProfile: &corev1.SeccompProfile{Type: corev1.SeccompProfileType("RuntimeDefault")},
			},
			ServiceAccountName: "kube-promtheus-stack-kube-admission",
		},
	}},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "batch/v1",
		Kind:       "Job",
	},
}

var KubePromtheusStackKubeAdmissionPatchJOBS = &batchv1.Job{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "post-install,post-upgrade",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
		},
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-admission-patch",
			"app.kubernetes.io/instance":   "kube-promtheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "47.0.0",
			"chart":                        "kube-prometheus-stack-47.0.0",
			"heritage":                     "Helm",
			"release":                      "kube-promtheus-stack",
		},
		Name:      "kube-promtheus-stack-kube-admission-patch",
		Namespace: "monitoring",
	},
	Spec: batchv1.JobSpec{Template: corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app":                          "kube-prometheus-stack-admission-patch",
				"app.kubernetes.io/instance":   "kube-promtheus-stack",
				"app.kubernetes.io/managed-by": "Helm",
				"app.kubernetes.io/part-of":    "kube-prometheus-stack",
				"app.kubernetes.io/version":    "47.0.0",
				"chart":                        "kube-prometheus-stack-47.0.0",
				"heritage":                     "Helm",
				"release":                      "kube-promtheus-stack",
			},
			Name: "kube-promtheus-stack-kube-admission-patch",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Args:            []string{"patch", "--webhook-name=kube-promtheus-stack-kube-admission", "--namespace=monitoring", "--secret-name=kube-promtheus-stack-kube-admission", "--patch-failure-policy="},
				Image:           "registry.k8s.io/ingress-nginx/kube-webhook-certgen:v20221220-controller-v1.5.1-58-g787ea74b6",
				ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
				Name:            "patch",
				SecurityContext: &corev1.SecurityContext{
					Capabilities:           &corev1.Capabilities{Drop: []corev1.Capability{corev1.Capability("ALL")}},
					ReadOnlyRootFilesystem: P(true),
				},
			}},
			RestartPolicy: corev1.RestartPolicy("OnFailure"),
			SecurityContext: &corev1.PodSecurityContext{
				RunAsGroup:     P(int64(2000)),
				RunAsNonRoot:   P(true),
				RunAsUser:      P(int64(2000)),
				SeccompProfile: &corev1.SeccompProfile{Type: corev1.SeccompProfileType("RuntimeDefault")},
			},
			ServiceAccountName: "kube-promtheus-stack-kube-admission",
		},
	}},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "batch/v1",
		Kind:       "Job",
	},
}
