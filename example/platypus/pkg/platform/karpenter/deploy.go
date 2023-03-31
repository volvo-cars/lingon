package karpenter

import (
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// P returns a pointer to the given value.
func P[T any](t T) *T {
	return &t
}

var Pdb = &policyv1.PodDisruptionBudget{
	TypeMeta: metav1.TypeMeta{
		Kind:       "PodDisruptionBudget",
		APIVersion: "policy/v1",
	},
	ObjectMeta: metav1.ObjectMeta{Name: "karpenter", Namespace: "karpenter"},
	Spec: policyv1.PodDisruptionBudgetSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app.kubernetes.io/instance": "karpenter",
				"app.kubernetes.io/name":     "karpenter",
			},
		}, MaxUnavailable: &intstr.IntOrString{IntVal: 1},
	},
}

var matchLabels = map[string]string{
	"app.kubernetes.io/instance": "karpenter",
	"app.kubernetes.io/name":     "karpenter",
}

var ContainerPorts = []corev1.ContainerPort{
	{
		Name:          "http-metrics",
		ContainerPort: 8080,
		Protocol:      corev1.Protocol("TCP"),
	},
	{
		Name:          "http",
		ContainerPort: 8081,
		Protocol:      corev1.Protocol("TCP"),
	},
	{
		Name:          "https-webhook",
		ContainerPort: 8443,
		Protocol:      corev1.Protocol("TCP"),
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
							Operator: corev1.NodeSelectorOperator("DoesNotExist"),
						},
					},
				},
			},
		},
	},
}

var Environment = []corev1.EnvVar{
	{
		Name:  "KUBERNETES_MIN_VERSION",
		Value: "1.19.0-0",
	},
	{
		Name:  "KARPENTER_SERVICE",
		Value: "karpenter",
	},
	{Name: "WEBHOOK_PORT", Value: "8443"},
	{Name: "METRICS_PORT", Value: "8080"},
	{
		Name:  "HEALTH_PROBE_PORT",
		Value: "8081",
	},
	{
		Name:      "SYSTEM_NAMESPACE",
		ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.namespace"}},
	},
	{
		Name: "MEMORY_LIMIT",
		ValueFrom: &corev1.EnvVarSource{
			ResourceFieldRef: &corev1.ResourceFieldSelector{
				ContainerName: "controller",
				Resource:      "limits.memory",
				Divisor:       resource.MustParse("0"),
			},
		},
	},
}

var Deploy = &appsv1.Deployment{
	TypeMeta: kubeutil.TypeDeploymentV1,
	ObjectMeta: metav1.ObjectMeta{
		Name:      "karpenter",
		Namespace: "karpenter",
		Labels:    commonLabels,
	},
	Spec: appsv1.DeploymentSpec{
		Replicas: P(int32(2)),
		Selector: &metav1.LabelSelector{
			MatchLabels: matchLabels,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: matchLabels,
			}, Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "controller",
						Image: "public.ecr.aws/karpenter/controller:v" + Version,
						Ports: ContainerPorts,
						Env:   Environment,
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceName("cpu"):    resource.MustParse("1"),
								corev1.ResourceName("memory"): resource.MustParse("1Gi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceName("cpu"):    resource.MustParse("1"),
								corev1.ResourceName("memory"): resource.MustParse("1Gi"),
							},
						},
						LivenessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/healthz",
									Port: intstr.IntOrString{
										Type:   intstr.Type(1),
										StrVal: "http",
									},
								},
							}, InitialDelaySeconds: 30, TimeoutSeconds: 30,
						},
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/readyz",
									Port: intstr.IntOrString{
										Type:   intstr.Type(1),
										StrVal: "http",
									},
								},
							}, TimeoutSeconds: 30,
						},
						ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
					},
				},
				DNSPolicy:          corev1.DNSPolicy("Default"),
				NodeSelector:       map[string]string{"kubernetes.io/os": "linux"},
				ServiceAccountName: "karpenter",
				SecurityContext:    &corev1.PodSecurityContext{FSGroup: P(int64(1000))},
				Affinity:           SetNodeAffinity,
				Tolerations: []corev1.Toleration{
					{
						Key:      "CriticalAddonsOnly",
						Operator: corev1.TolerationOperator("Exists"),
					},
				},
				PriorityClassName: "system-cluster-critical",
				TopologySpreadConstraints: []corev1.TopologySpreadConstraint{
					{
						MaxSkew:           1,
						TopologyKey:       "topology.kubernetes.io/zone",
						WhenUnsatisfiable: corev1.UnsatisfiableConstraintAction("ScheduleAnyway"),
						LabelSelector: &metav1.LabelSelector{
							MatchLabels: matchLabels,
						},
					},
				},
			},
		},
		Strategy:             appsv1.DeploymentStrategy{RollingUpdate: &appsv1.RollingUpdateDeployment{MaxUnavailable: &intstr.IntOrString{IntVal: 1}}},
		RevisionHistoryLimit: P(int32(10)),
	},
}
