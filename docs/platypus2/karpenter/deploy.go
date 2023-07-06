// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package karpenter

import (
	"fmt"

	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var Svc = &corev1.Service{
	TypeMeta:   ku.TypeServiceV1,
	ObjectMeta: KA.ObjectMeta(),
	Spec: corev1.ServiceSpec{
		Type:     corev1.ServiceTypeClusterIP,
		Selector: KA.MatchLabels(),
		Ports: []corev1.ServicePort{
			KA.Metrics.Service,
			KA.Webhook.Service,
		},
	},
}

const containerName = "controller"

var Deploy = &appsv1.Deployment{
	TypeMeta:   ku.TypeDeploymentV1,
	ObjectMeta: KA.ObjectMeta(),
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(2)),
		Selector: &metav1.LabelSelector{MatchLabels: KA.MatchLabels()},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: KA.MatchLabels()},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  containerName,
						Image: KA.Img.URL(),
						Ports: []corev1.ContainerPort{
							KA.Probe.Container,
							KA.Metrics.Container,
							KA.Webhook.Container,
						},
						Env: []corev1.EnvVar{
							{
								Name:  "KUBERNETES_MIN_VERSION",
								Value: KA.KubeMinVer,
							},
							{Name: "KARPENTER_SERVICE", Value: Svc.Name},
							{
								Name:  "WEBHOOK_PORT",
								Value: d32(KA.Webhook.Container.ContainerPort),
							},
							{
								Name:  "METRICS_PORT",
								Value: d32(KA.Metrics.Container.ContainerPort),
							},
							{
								Name:  "HEALTH_PROBE_PORT",
								Value: d32(KA.Probe.Container.ContainerPort),
							},
							ku.EnvVarDownAPI(
								"SYSTEM_NAMESPACE",
								"metadata.namespace",
							),
							{
								Name: "MEMORY_LIMIT",
								ValueFrom: &corev1.EnvVarSource{
									ResourceFieldRef: &corev1.ResourceFieldSelector{
										ContainerName: containerName,
										Resource:      "limits.memory",
										Divisor:       resource.MustParse("0"),
									},
								},
							},
						},
						Resources:       ku.Resources("1", "1Gi", "1", "1Gi"),
						LivenessProbe:   probe("/healthz", 30),
						ReadinessProbe:  probe("/readyz", 0),
						ImagePullPolicy: corev1.PullIfNotPresent,
					},
				},
				DNSPolicy:          corev1.DNSDefault,
				NodeSelector:       map[string]string{ku.LabelOSStable: "linux"},
				ServiceAccountName: "TO_BE_SET_IN_NEW",
				SecurityContext:    &corev1.PodSecurityContext{FSGroup: P(int64(1000))},
				Affinity:           SetNodeAffinity,
				Tolerations: []corev1.Toleration{
					{
						Key:      "CriticalAddonsOnly",
						Operator: corev1.TolerationOpExists,
					},
				},
				PriorityClassName: "system-cluster-critical",
				TopologySpreadConstraints: []corev1.TopologySpreadConstraint{
					{
						MaxSkew:           1,
						TopologyKey:       ku.LabelTopologyZone,
						WhenUnsatisfiable: corev1.UnsatisfiableConstraintAction("ScheduleAnyway"),
						LabelSelector:     &metav1.LabelSelector{MatchLabels: KA.MatchLabels()},
					},
				},
			},
		},
		Strategy: appsv1.DeploymentStrategy{
			RollingUpdate: &appsv1.RollingUpdateDeployment{
				MaxUnavailable: P(intstr.FromInt(1)),
			},
		},
	},
}

var Pdb = &policyv1.PodDisruptionBudget{
	TypeMeta: metav1.TypeMeta{
		Kind:       "PodDisruptionBudget",
		APIVersion: "policy/v1",
	},
	ObjectMeta: KA.ObjectMeta(),
	Spec: policyv1.PodDisruptionBudgetSpec{
		Selector:       &metav1.LabelSelector{MatchLabels: KA.MatchLabels()},
		MaxUnavailable: P(intstr.FromInt(1)),
	},
}

var d32 = func(p int32) string { return fmt.Sprintf("%d", p) }

var SetNodeAffinity = &corev1.Affinity{
	NodeAffinity: &corev1.NodeAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
			NodeSelectorTerms: []corev1.NodeSelectorTerm{
				{
					MatchExpressions: []corev1.NodeSelectorRequirement{
						{
							Key:      "karpenter.sh/provisioner-name",
							Operator: corev1.NodeSelectorOpDoesNotExist,
						},
					},
				},
			},
		},
	},
	PodAntiAffinity: &corev1.PodAntiAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
			{
				LabelSelector: &metav1.LabelSelector{MatchLabels: KA.MatchLabels()},
				TopologyKey:   ku.LabelHostname,
			},
		},
	},
}

func probe(path string, initDelaySec int) *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: path,
				Port: intstr.FromString(KA.Probe.Container.Name),
			},
		},
		TimeoutSeconds:      int32(30),
		InitialDelaySeconds: int32(initDelaySec),
	}
}
