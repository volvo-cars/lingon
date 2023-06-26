// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package vmop

import (
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Deploy = &appsv1.Deployment{
	TypeMeta:   ku.TypeDeploymentV1,
	ObjectMeta: O.ObjectMeta(),
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(1)),
		Selector: &metav1.LabelSelector{MatchLabels: O.MatchLabels()},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: O.MatchLabels()},
			Spec: corev1.PodSpec{
				ServiceAccountName: SA.Name,

				Containers: []corev1.Container{
					{
						Image:           O.Img.URL(),
						ImagePullPolicy: corev1.PullIfNotPresent,
						Name:            O.Name,
						Args: []string{
							"--zap-log-level=info",
							"--enable-leader-election",
							// "--webhook.enable=true",
						},
						Command: []string{"manager"},
						Env: []corev1.EnvVar{
							{Name: "WATCH_NAMESPACE"},
							{Name: "OPERATOR_NAME", Value: O.Name},
							ku.EnvVarDownAPI("POD_NAME", "metadata.name"),
							// See https://github.com/VictoriaMetrics/operator/blob/master/vars.MD
							{
								Name:  "VM_VMSINGLEDEFAULT_VERSION",
								Value: "v" + O.VMVersion,
							},
							{
								Name:  "VM_USECUSTOMCONFIGRELOADER",
								Value: "true",
							},
							{
								Name:  "VM_PSPAUTOCREATEENABLED",
								Value: "false",
							},
							{
								// By default, the operator doesn't make converted objects disappear after original ones are deleted.
								// To change this behaviour, it adds `OwnerReferences` to converted objects.
								// Converted objects will be linked to the original ones
								// and will be deleted by kubernetes after the original ones are deleted.
								Name:  "VM_ENABLEDPROMETHEUSCONVERTEROWNERREFERENCES",
								Value: "true",
							},
						},
						Ports: []corev1.ContainerPort{
							O.Main.Container,
							O.Webhook.Container,
						},
						Resources: ku.Resources(
							"80m",
							"310Mi",
							"120m",
							"320Mi",
						),
						// TODO: certificate for webhook
						// VolumeMounts: []corev1.VolumeMount{
						// 	{
						// 		MountPath: "/tmp/k8s-webhook-server/serving-certs",
						// 		Name:      "cert",
						// 		ReadOnly:  true,
						// 	},
						// },
					},
				},
				// Volumes: []corev1.Volume{
				// 	{
				// 		Name: "cert",
				// 		VolumeSource: corev1.VolumeSource{
				// 			Secret: &corev1.SecretVolumeSource{
				// 				DefaultMode: P(int32(420)),
				// 				SecretName:  "vmop-victoria-metrics-operator-validation",
				// 			},
				// 		},
				// 	},
				// },
			},
		},
	},
}
