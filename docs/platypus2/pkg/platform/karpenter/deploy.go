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

const containerName = "controller"

var Deploy = &appsv1.Deployment{
	TypeMeta: ku.TypeDeploymentV1,
	ObjectMeta: metav1.ObjectMeta{
		Name:      AppName,
		Namespace: Namespace,
		Labels:    commonLabels,
	},
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(2)),
		Selector: &metav1.LabelSelector{MatchLabels: matchLabels},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: matchLabels},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:            containerName,
						Image:           ImgController,
						Ports:           ContainerPorts,
						Env:             Environment,
						Resources:       ku.Resources("1", "1Gi", "1", "1Gi"),
						LivenessProbe:   probe("/healthz", 30),
						ReadinessProbe:  probe("/readyz", 0),
						ImagePullPolicy: corev1.PullIfNotPresent,
					},
				},
				DNSPolicy:          corev1.DNSDefault,
				NodeSelector:       map[string]string{ku.LabelOSStable: "linux"},
				ServiceAccountName: "karpenter",
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
						LabelSelector:     &metav1.LabelSelector{MatchLabels: matchLabels},
					},
				},
			},
		},
		Strategy:             appsv1.DeploymentStrategy{RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &intstr.IntOrString{IntVal: 1}}},
		RevisionHistoryLimit: P(int32(10)),
	},
}

var Pdb = &policyv1.PodDisruptionBudget{
	TypeMeta: metav1.TypeMeta{
		Kind:       "PodDisruptionBudget",
		APIVersion: "policy/v1",
	},
	ObjectMeta: metav1.ObjectMeta{Name: AppName, Namespace: Namespace},
	Spec: policyv1.PodDisruptionBudgetSpec{
		Selector:       &metav1.LabelSelector{MatchLabels: matchLabels},
		MaxUnavailable: P(intstr.FromInt(1)),
	},
}

var ContainerPorts = []corev1.ContainerPort{
	{
		Name:          PortNameMetrics,
		ContainerPort: PortMetrics,
	},
	{
		Name:          PortNameProbe,
		ContainerPort: PortProbe,
	},
	{
		Name:          PortNameWebhook,
		ContainerPort: PortWebhookCtnr,
	},
}

var kS = func(p int) string { return fmt.Sprintf("%d", p) }

var Environment = []corev1.EnvVar{
	{Name: "KUBERNETES_MIN_VERSION", Value: "1.19.0-0"},
	{Name: "KARPENTER_SERVICE", Value: Svc.Name},
	{Name: "WEBHOOK_PORT", Value: kS(PortWebhookCtnr)},
	{Name: "METRICS_PORT", Value: kS(PortMetrics)},
	{Name: "HEALTH_PROBE_PORT", Value: kS(PortProbe)},
	ku.EnvVarDownAPI("SYSTEM_NAMESPACE", "metadata.namespace"),
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
}

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
				LabelSelector: &metav1.LabelSelector{
					MatchLabels: matchLabels,
				},
				TopologyKey: ku.LabelHostname,
			},
		},
	},
}

func probe(path string, initDelaySec int) *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: path,
				Port: intstr.FromString(PortNameProbe),
			},
		},
		TimeoutSeconds:      int32(30),
		InitialDelaySeconds: int32(initDelaySec),
	}
}
