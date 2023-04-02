// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package cilium

import (
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var opMeta = kubeutil.ObjectMeta(
	"cilium-operator",
	"kube-system",
	map[string]string{
		"io.cilium/app": "operator",
		"name":          "cilium-operator",
	},
	nil,
)

var Operator = &appsv1.Deployment{
	TypeMeta:   kubeutil.TypeDeploymentV1,
	ObjectMeta: opMeta,
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(2)),
		Selector: &metav1.LabelSelector{
			MatchLabels: opMeta.Labels,
		},
		Template: v1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: opMeta.Labels,
			},
			Spec: v1.PodSpec{
				Volumes: []v1.Volume{
					{
						Name: "cilium-config-path",
						VolumeSource: v1.VolumeSource{
							ConfigMap: &v1.ConfigMapVolumeSource{
								LocalObjectReference: v1.LocalObjectReference{Name: "cilium-config"},
							},
						},
					},
				},
				Containers: []v1.Container{
					{
						Name:    "cilium-operator",
						Image:   "quay.io/cilium/operator-generic:v1.12.4@sha256:071089ec5bca1f556afb8e541d9972a0dfb09d1e25504ae642ced021ecbedbd1",
						Command: []string{"cilium-operator-generic"},
						Args: []string{
							"--config-dir=/tmp/cilium/config-map",
							"--debug=$(CILIUM_DEBUG)",
						},
						Env: []v1.EnvVar{
							{
								Name: "K8S_NODE_NAME",
								ValueFrom: &v1.EnvVarSource{
									FieldRef: &v1.ObjectFieldSelector{
										APIVersion: "v1",
										FieldPath:  "spec.nodeName",
									},
								},
							},
							{
								Name: "CILIUM_K8S_NAMESPACE",
								ValueFrom: &v1.EnvVarSource{
									FieldRef: &v1.ObjectFieldSelector{
										APIVersion: "v1",
										FieldPath:  "metadata.namespace",
									},
								},
							},
							{
								Name: "CILIUM_DEBUG",
								ValueFrom: &v1.EnvVarSource{
									ConfigMapKeyRef: &v1.ConfigMapKeySelector{
										LocalObjectReference: v1.LocalObjectReference{Name: "cilium-config"},
										Key:                  "debug",
										Optional:             P(true),
									},
								},
							},
						},
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      "cilium-config-path",
								ReadOnly:  true,
								MountPath: "/tmp/cilium/config-map",
							},
						},
						LivenessProbe:            liveness,
						TerminationMessagePolicy: v1.TerminationMessagePolicy("FallbackToLogsOnError"),
						ImagePullPolicy:          v1.PullPolicy("IfNotPresent"),
					},
				},
				RestartPolicy:      v1.RestartPolicy("Always"),
				NodeSelector:       map[string]string{"kubernetes.io/os": "linux"},
				ServiceAccountName: "cilium-operator",
				HostNetwork:        true,
				Affinity: &v1.Affinity{
					PodAntiAffinity: &v1.PodAntiAffinity{
						RequiredDuringSchedulingIgnoredDuringExecution: []v1.PodAffinityTerm{
							{
								LabelSelector: &metav1.LabelSelector{MatchLabels: map[string]string{"io.cilium/app": "operator"}},
								TopologyKey:   "kubernetes.io/hostname",
							},
						},
					},
				},
				Tolerations:       []v1.Toleration{{Operator: v1.TolerationOperator("Exists")}},
				PriorityClassName: "system-cluster-critical",
			},
		},
		Strategy: appsv1.DeploymentStrategy{
			Type: appsv1.DeploymentStrategyType("RollingUpdate"),
			RollingUpdate: &appsv1.RollingUpdateDeployment{
				MaxUnavailable: &intstr.IntOrString{IntVal: 1},
				MaxSurge:       &intstr.IntOrString{IntVal: 1},
			},
		},
	},
}

var liveness = &v1.Probe{
	ProbeHandler: v1.ProbeHandler{
		HTTPGet: &v1.HTTPGetAction{
			Path:   "/healthz",
			Port:   intstr.IntOrString{IntVal: 9234},
			Host:   "127.0.0.1",
			Scheme: v1.URIScheme("HTTP"),
		},
	},
	InitialDelaySeconds: 60,
	TimeoutSeconds:      3,
	PeriodSeconds:       10,
}
