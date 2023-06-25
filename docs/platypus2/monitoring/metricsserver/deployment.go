// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package metricsserver

import (
	"fmt"

	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var podAntiAff = &corev1.PodAntiAffinity{
	PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
		{
			PodAffinityTerm: corev1.PodAffinityTerm{
				LabelSelector: &metav1.LabelSelector{
					MatchLabels: M.MatchLabels(),
				},
				TopologyKey: corev1.LabelHostname,
			},
			Weight: int32(1),
		},
	},
}

var deployStrategy = appsv1.DeploymentStrategy{
	RollingUpdate: &appsv1.RollingUpdateDeployment{
		MaxUnavailable: P(intstr.FromInt(1)),
	},
	Type: appsv1.RollingUpdateDeploymentStrategyType,
}

var d = func(i int32) string { return fmt.Sprintf("%d", i) }

var Deploy = &appsv1.Deployment{
	TypeMeta:   ku.TypeDeploymentV1,
	ObjectMeta: M.ObjectMeta(),
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(1)),
		Selector: &metav1.LabelSelector{MatchLabels: M.MatchLabels()},
		Strategy: deployStrategy,
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: M.MatchLabels()},
			Spec: corev1.PodSpec{
				Affinity: &corev1.Affinity{
					PodAntiAffinity: podAntiAff,
				},
				Containers: []corev1.Container{
					{
						Name:            M.Name,
						Image:           M.Img.URL(),
						ImagePullPolicy: corev1.PullIfNotPresent,
						Args: []string{
							"--secure-port=" + d(M.P.Container.ContainerPort),
							"--cert-dir=/tmp",
							"--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname",
							"--kubelet-use-node-status-port",
							"--metric-resolution=15s",
							"--authorization-always-allow-paths=" + M.MetricsURL, // /metrics",
						},
						Ports: []corev1.ContainerPort{M.P.Container},
						LivenessProbe: &corev1.Probe{
							FailureThreshold: int32(3),
							PeriodSeconds:    int32(10),
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "/livez",
									Port:   intstr.FromString(M.P.Container.Name),
									Scheme: corev1.URISchemeHTTPS,
								},
							},
						},
						Resources: ku.Resources(
							"200m",
							"300Mi",
							"200m",
							"300Mi",
						),
						ReadinessProbe: &corev1.Probe{
							FailureThreshold:    int32(3),
							InitialDelaySeconds: int32(20),
							PeriodSeconds:       int32(10),
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "/readyz",
									Port:   intstr.FromString(M.P.Container.Name),
									Scheme: corev1.URISchemeHTTPS,
								},
							},
						},
						SecurityContext: &corev1.SecurityContext{
							Capabilities: &corev1.Capabilities{
								Drop: []corev1.Capability{corev1.Capability("ALL")},
							},
							ReadOnlyRootFilesystem: P(true),
							RunAsNonRoot:           P(true),
							RunAsUser:              P(int64(1000)),
							SeccompProfile: &corev1.SeccompProfile{
								Type: corev1.SeccompProfileTypeRuntimeDefault,
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{MountPath: "/tmp", Name: "tmp"},
						},
					},
				},
				PriorityClassName:  "system-cluster-critical",
				ServiceAccountName: SA.Name,
				Volumes: []corev1.Volume{
					{Name: "tmp", VolumeSource: corev1.VolumeSource{}},
				},
			},
		},
	},
}
