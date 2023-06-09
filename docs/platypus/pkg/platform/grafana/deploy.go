// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package grafana

import (
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Deployment(opts KubeOpts) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: kubeutil.TypeDeploymentV1,
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.Name,
			Namespace: opts.Namespace,
			Labels:    opts.CommonLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas:             P(int32(2)),
			RevisionHistoryLimit: P(int32(10)),
			Selector: &metav1.LabelSelector{
				MatchLabels: opts.CommonLabels,
			},
			Strategy: appsv1.DeploymentStrategy{Type: appsv1.RollingUpdateDeploymentStrategyType},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: opts.CommonLabels,
				},
				Spec: corev1.PodSpec{
					AutomountServiceAccountToken: P(true),
					Containers: []corev1.Container{
						{
							Env: []corev1.EnvVar{
								{
									Name:      "POD_IP",
									ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.podIP"}},
								},
								{
									Name: "GF_SECURITY_ADMIN_USER",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key:                  "admin-user",
											LocalObjectReference: corev1.LocalObjectReference{Name: "grafana"},
										},
									},
								},
								{
									Name: "GF_SECURITY_ADMIN_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key:                  "admin-password",
											LocalObjectReference: corev1.LocalObjectReference{Name: "grafana"},
										},
									},
								},
								{
									Name:  "GF_PATHS_DATA",
									Value: "/var/lib/grafana/",
								},
								{
									Name:  "GF_PATHS_LOGS",
									Value: "/var/log/grafana",
								},
								{
									Name:  "GF_PATHS_PLUGINS",
									Value: "/var/lib/grafana/plugins",
								},
								{
									Name:  "GF_PATHS_PROVISIONING",
									Value: "/etc/grafana/provisioning",
								},
								{
									Name:  "GF_ROOT_URL",
									Value: "",
								},
								{
									Name:  "GF_SERVER_DOMAIN",
									Value: "",
								},
								{
									Name:  "GF_DATABASE_TYPE",
									Value: "postgres",
								},
								{
									Name:  "GF_DATABASE_HOST",
									Value: opts.PostgresHost,
								},
								{
									Name:  "GF_DATABASE_NAME",
									Value: opts.PostgresDBName,
								},
								{
									Name:  "GF_DATABASE_USER",
									Value: opts.PostgresUser,
								},
								{
									Name:  "GF_DATABASE_PASSWORD",
									Value: opts.PostgresPassword,
								},
							},
							Image:           "grafana/grafana:" + Version,
							ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
							LivenessProbe: &corev1.Probe{
								FailureThreshold:    int32(10),
								InitialDelaySeconds: int32(60),
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/api/health",
										Port: intstr.IntOrString{IntVal: int32(3000)},
									},
								},
								TimeoutSeconds: int32(30),
							},
							Name: AppName,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: int32(3000),
									Name:          "grafana",
									Protocol:      corev1.ProtocolTCP,
								}, {
									ContainerPort: int32(9094),
									Name:          "gossip-tcp",
									Protocol:      corev1.ProtocolTCP,
								}, {
									ContainerPort: int32(9094),
									Name:          "gossip-udp",
									Protocol:      corev1.ProtocolUDP,
								},
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/api/health",
										Port: intstr.IntOrString{IntVal: int32(3000)},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: "/etc/grafana/grafana.ini",
									Name:      "config",
									SubPath:   "grafana.ini",
								}, {
									MountPath: "/var/lib/grafana",
									Name:      "storage",
								}, {
									MountPath: "/etc/grafana/provisioning/datasources/datasources.yaml",
									Name:      "config",
									SubPath:   "datasources.yaml",
								}, {
									MountPath: "/etc/grafana/provisioning/dashboards/dashboardproviders.yaml",
									Name:      "config",
									SubPath:   "dashboardproviders.yaml",
								},
							},
						},
					},
					EnableServiceLinks: P(true),
					InitContainers: []corev1.Container{
						{
							Args: []string{
								"-c",
								"mkdir -p /var/lib/grafana/dashboards/default && /bin/sh -x /etc/grafana/download_dashboards.sh",
							},
							Command:         []string{"/bin/sh"},
							Image:           "curlimages/curl:7.85.0",
							ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
							Name:            "download-dashboards",
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: "/etc/grafana/download_dashboards.sh",
									Name:      "config",
									SubPath:   "download_dashboards.sh",
								}, {
									MountPath: "/var/lib/grafana",
									Name:      "storage",
								},
							},
						},
					},
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup:    P(int64(472)),
						RunAsGroup: P(int64(472)),
						RunAsUser:  P(int64(472)),
					},
					ServiceAccountName: AppName,
					Volumes: []corev1.Volume{
						{
							Name:         "config",
							VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: AppName}}},
						}, {
							Name:         "dashboards-default",
							VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: "grafana-dashboards-default"}}},
						}, {
							Name:         "storage",
							VolumeSource: corev1.VolumeSource{},
						},
					},
				},
			},
		},
	}
}
