// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package tekton

import (
	"fmt"

	"github.com/volvo-cars/lingon/pkg/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var PipelinesWebhookDeploy = &appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Labels: kubeutil.MergeLabels(
			map[string]string{kubeutil.AppLabelName: WebhookName},
			labelsWebhook,
			labelsVersion,
		),
		Name:      WebhookFullName,
		Namespace: PipelinesNS.Name,
	},
	Spec: appsv1.DeploymentSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: kubeutil.MergeLabels(
				map[string]string{kubeutil.AppLabelName: WebhookName},
				labelsWebhook,
			),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: kubeutil.MergeLabels(
					map[string]string{
						"app":                 WebhookFullName,
						kubeutil.AppLabelName: WebhookName,
					},
					labelsWebhook,
					labelsVersion,
				),
			},
			Spec: corev1.PodSpec{
				Affinity: &corev1.Affinity{
					NodeAffinity: &corev1.NodeAffinity{
						RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
							NodeSelectorTerms: []corev1.NodeSelectorTerm{
								kubeutil.NotInWindows,
							},
						},
					},
					PodAntiAffinity: &corev1.PodAntiAffinity{
						PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
							{
								PodAffinityTerm: corev1.PodAffinityTerm{
									LabelSelector: &metav1.LabelSelector{
										MatchLabels: kubeutil.MergeLabels(
											map[string]string{kubeutil.AppLabelName: WebhookName},
											labelsWebhook,
										),
									},
									TopologyKey: kubeutil.LabelHostname,
								},
								Weight: int32(100),
							},
						},
					},
				},
				Containers: []corev1.Container{
					{
						Env: []corev1.EnvVar{
							{
								Name:      "SYSTEM_NAMESPACE",
								ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.namespace"}},
							}, {
								Name:  "CONFIG_LOGGING_NAME",
								Value: ConfigLoggingCM.Name,
							}, {
								Name:  "CONFIG_OBSERVABILITY_NAME",
								Value: ConfigObservabilityCM.Name,
							}, {
								Name:  "CONFIG_LEADERELECTION_NAME",
								Value: ConfigLeaderElectionCM.Name,
							}, {
								Name:  "CONFIG_FEATURE_FLAGS_NAME",
								Value: FeatureFlagsCM.Name,
							}, {
								Name:  "WEBHOOK_PORT",
								Value: fmt.Sprintf("%d", WebhookPort),
							}, {
								Name:  "WEBHOOK_ADMISSION_CONTROLLER_NAME",
								Value: "webhook.pipeline.tekton.dev",
							}, {
								Name:  "WEBHOOK_SERVICE_NAME",
								Value: WebhookFullName,
							}, {
								Name:  "WEBHOOK_SECRET_NAME",
								Value: WebhookCertsSecrets.Name,
							}, {
								Name:  "METRICS_DOMAIN",
								Value: "tekton.dev/pipeline",
							},
						},
						Image: WebhookImage,
						LivenessProbe: &corev1.Probe{
							InitialDelaySeconds: int32(5),
							PeriodSeconds:       int32(10),
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/health",
									Port: intstr.IntOrString{
										StrVal: "probes",
										Type:   intstr.Type(int64(1)),
									},
									Scheme: corev1.URISchemeHTTP,
								},
							},
							TimeoutSeconds: int32(5),
						},
						Name: WebhookName,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: int32(9090),
								Name:          "metrics",
							}, {
								ContainerPort: int32(8008),
								Name:          "profiling",
							}, {
								ContainerPort: int32(WebhookPort),
								Name:          "https-webhook",
							}, {
								ContainerPort: int32(8080),
								Name:          "probes",
							},
						},
						ReadinessProbe: &corev1.Probe{
							InitialDelaySeconds: int32(5),
							PeriodSeconds:       int32(10),
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/readiness",
									Port: intstr.IntOrString{
										StrVal: "probes",
										Type:   intstr.Type(int64(1)),
									},
									Scheme: corev1.URISchemeHTTP,
								},
							},
							TimeoutSeconds: int32(5),
						},
						Resources: kubeutil.Resources(
							"100m",
							"100Mi",
							"500m",
							"500Mi",
						),
						SecurityContext: &corev1.SecurityContext{
							Capabilities:   &corev1.Capabilities{Drop: []corev1.Capability{corev1.Capability("ALL")}},
							RunAsGroup:     P(int64(65532)),
							RunAsNonRoot:   P(true),
							RunAsUser:      P(int64(65532)),
							SeccompProfile: &corev1.SeccompProfile{Type: corev1.SeccompProfileTypeRuntimeDefault},
						},
					},
				},
				ServiceAccountName: WebhookFullName,
			},
		},
	},
	TypeMeta: kubeutil.TypeDeploymentV1,
}

var PipelinesControllerDeploy = &appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Labels: kubeutil.MergeLabels(
			map[string]string{kubeutil.AppLabelName: ControllerName},
			labelsController,
			labelsVersion,
		),
		Name:      ControllerFullName,
		Namespace: PipelinesNS.Name,
	},
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(1)),
		Selector: &metav1.LabelSelector{
			MatchLabels: kubeutil.MergeLabels(
				map[string]string{kubeutil.AppLabelName: ControllerName},
				labelsController,
			),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: kubeutil.MergeLabels(
					map[string]string{
						"app":                 ControllerFullName,
						kubeutil.AppLabelName: ControllerName,
					},
					labelsController,
					labelsVersion,
				),
			},
			Spec: corev1.PodSpec{
				Affinity: &corev1.Affinity{
					NodeAffinity: &corev1.NodeAffinity{
						RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
							NodeSelectorTerms: []corev1.NodeSelectorTerm{kubeutil.NotInWindows},
						},
					},
				},
				Containers: []corev1.Container{
					{
						Args: []string{
							"-entrypoint-image",
							"gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/entrypoint:v0.46.0@sha256:36114bab6037563667aa0620037e7a063ffe00f432866a293807f8029eddd645",
							"-nop-image",
							"gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/nop:v0.46.0@sha256:1b9ad2522b5a5ea0c51ac43e2838ea1535de9d9c82c7864ed9a88553db434a29",
							"-sidecarlogresults-image",
							"gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/sidecarlogresults:v0.46.0@sha256:4bc1d0dc796a2a85a72d431344b80a2ac93f259fdd199d17ebc6d31b52a571d6",
							"-workingdirinit-image",
							"gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/workingdirinit:v0.46.0@sha256:b066c05c1565675a573563557d2cd91bea48217091a3beda639f0dbdea5910bc",
							"-shell-image",
							"cgr.dev/chainguard/busybox@sha256:19f02276bf8dbdd62f069b922f10c65262cc34b710eea26ff928129a736be791",
							"-shell-image-win",
							"mcr.microsoft.com/powershell:nanoserver@sha256:b6d5ff841b78bdf2dfed7550000fd4f3437385b8fa686ec0f010be24777654d6",
						},
						Env: []corev1.EnvVar{
							{
								Name:      "SYSTEM_NAMESPACE",
								ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.namespace"}},
							}, {
								Name:  "CONFIG_DEFAULTS_NAME",
								Value: ConfigDefaultsCM.Name,
							}, {
								Name:  "CONFIG_LOGGING_NAME",
								Value: ConfigLoggingCM.Name,
							}, {
								Name:  "CONFIG_OBSERVABILITY_NAME",
								Value: ConfigObservabilityCM.Name,
							}, {
								Name:  "CONFIG_FEATURE_FLAGS_NAME",
								Value: FeatureFlagsCM.Name,
							}, {
								Name:  "CONFIG_LEADERELECTION_NAME",
								Value: ConfigLeaderElectionCM.Name,
							}, {
								Name:  "CONFIG_SPIRE",
								Value: ConfigSpireCM.Name,
							}, {
								Name:  "SSL_CERT_FILE",
								Value: "/etc/config-registry-cert/cert",
							}, {
								Name:  "SSL_CERT_DIR",
								Value: "/etc/ssl/certs",
							}, {
								Name:  "METRICS_DOMAIN",
								Value: "tekton.dev/pipeline",
							},
						},
						Image: ControllerImage,
						LivenessProbe: &corev1.Probe{
							InitialDelaySeconds: int32(5),
							PeriodSeconds:       int32(10),
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/health",
									Port: intstr.IntOrString{
										StrVal: "probes",
										Type:   intstr.Type(int64(1)),
									},
									Scheme: corev1.URISchemeHTTP,
								},
							},
							TimeoutSeconds: int32(5),
						},
						Name: ControllerFullName,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: int32(9090),
								Name:          "metrics",
							}, {
								ContainerPort: int32(8008),
								Name:          "profiling",
							}, {
								ContainerPort: int32(8080),
								Name:          "probes",
							},
						},
						ReadinessProbe: &corev1.Probe{
							InitialDelaySeconds: int32(5),
							PeriodSeconds:       int32(10),
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/readiness",
									Port: intstr.IntOrString{
										StrVal: "probes",
										Type:   intstr.Type(int64(1)),
									},
									Scheme: corev1.URISchemeHTTP,
								},
							},
							TimeoutSeconds: int32(5),
						},
						SecurityContext: &corev1.SecurityContext{
							Capabilities:   &corev1.Capabilities{Drop: []corev1.Capability{corev1.Capability("ALL")}},
							RunAsGroup:     P(int64(65532)),
							RunAsNonRoot:   P(true),
							RunAsUser:      P(int64(65532)),
							SeccompProfile: &corev1.SeccompProfile{Type: corev1.SeccompProfileTypeRuntimeDefault},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/etc/config-logging",
								Name:      ConfigLoggingCM.Name,
							}, {
								MountPath: "/etc/config-registry-cert",
								Name:      ConfigRegistryCertCM.Name,
							},
						},
					},
				},
				ServiceAccountName: PipelinesControllerSA.Name,
				Volumes: []corev1.Volume{
					{
						Name:         ConfigLoggingCM.Name,
						VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: ConfigLoggingCM.Name}}},
					}, {
						Name:         ConfigRegistryCertCM.Name,
						VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: ConfigRegistryCertCM.Name}}},
					},
				},
			},
		},
	},
	TypeMeta: kubeutil.TypeDeploymentV1,
}

var PipelinesRemoteResolversDeploy = &appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Labels: kubeutil.MergeLabels(
			map[string]string{kubeutil.AppLabelName: ResolversName},
			labelsResolvers,
			labelsVersion,
		),
		Name:      ResolversFullName,
		Namespace: ResolversNS.Name,
	},
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(1)),
		Selector: &metav1.LabelSelector{
			MatchLabels: kubeutil.MergeLabels(
				map[string]string{kubeutil.AppLabelName: ResolversName},
				labelsResolvers,
			),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: kubeutil.MergeLabels(
					map[string]string{
						"app":                 ResolversFullName,
						kubeutil.AppLabelName: ResolversName,
					},
					labelsResolvers,
					labelsVersion,
				),
			},
			Spec: corev1.PodSpec{
				Affinity: &corev1.Affinity{
					PodAntiAffinity: &corev1.PodAntiAffinity{
						PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
							{
								PodAffinityTerm: corev1.PodAffinityTerm{
									LabelSelector: &metav1.LabelSelector{
										MatchLabels: kubeutil.MergeLabels(
											map[string]string{kubeutil.AppLabelName: ResolversName},
											labelsResolvers,
										),
									},
									TopologyKey: kubeutil.LabelHostname,
								},
								Weight: int32(100),
							},
						},
					},
				},
				Containers: []corev1.Container{
					{
						Env: []corev1.EnvVar{
							{
								Name:      "SYSTEM_NAMESPACE",
								ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.namespace"}},
							}, {
								Name:  "CONFIG_LOGGING_NAME",
								Value: ConfigLoggingResolversCM.Name,
							}, {
								Name:  "CONFIG_OBSERVABILITY_NAME",
								Value: ConfigObservabilityResolversCM.Name,
							}, {
								Name:  "CONFIG_FEATURE_FLAGS_NAME",
								Value: ResolversFeatureFlagsCM.Name,
							}, {
								Name:  "CONFIG_LEADERELECTION_NAME",
								Value: ConfigLeaderElectionResolversCM.Name,
							}, {
								Name:  "METRICS_DOMAIN",
								Value: "tekton.dev/resolution",
							}, {
								Name:  "ARTIFACT_HUB_API",
								Value: "https://artifacthub.io/",
							},
						},
						Image: ResolversImage,
						Name:  ResolversName,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: int32(9090),
								Name:          "metrics",
							},
						},
						Resources: kubeutil.Resources(
							"100m",
							"100Mi",
							"1000m",
							"4Gi",
						),
						SecurityContext: &corev1.SecurityContext{
							Capabilities:           &corev1.Capabilities{Drop: []corev1.Capability{corev1.Capability("ALL")}},
							ReadOnlyRootFilesystem: P(true),
							RunAsNonRoot:           P(true),
							SeccompProfile:         &corev1.SeccompProfile{Type: corev1.SeccompProfileTypeRuntimeDefault},
						},
					},
				},
				ServiceAccountName: PipelinesResolversSA.Name,
			},
		},
	},
	TypeMeta: kubeutil.TypeDeploymentV1,
}
