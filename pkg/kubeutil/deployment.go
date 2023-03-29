// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SimpleDeployment(
	name, namespace string,
	labels map[string]string,
	replicas int32,
	image string,
) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta:   TypeDeploymentV1,
		ObjectMeta: ObjectMeta(name, namespace, labels, nil),
		Spec: appsv1.DeploymentSpec{
			Replicas: P(replicas),
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					ServiceAccountName: name,
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: image,
						},
					},
				},
			},
		},
	}
}

func SetDeploySA(deploy *appsv1.Deployment, saName string) *appsv1.Deployment {
	deploy.Spec.Template.Spec.ServiceAccountName = saName
	return deploy
}

func Resources(cpuWant, memWant, cpuMax, memMax string) corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(cpuWant),
			corev1.ResourceMemory: resource.MustParse(memWant),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(cpuMax),
			corev1.ResourceMemory: resource.MustParse(memMax),
		},
	}
}
