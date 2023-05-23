package kubeutil

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Resources creates a ResourceRequirements struct.
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

func AntiAffinityHostnameByLabel(key, value string) *corev1.PodAntiAffinity {
	return &corev1.PodAntiAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
			{
				LabelSelector: &metav1.LabelSelector{
					MatchExpressions: []metav1.LabelSelectorRequirement{
						{
							Key:      key,
							Operator: metav1.LabelSelectorOpIn,
							Values:   []string{value},
						},
					},
				},
				TopologyKey: LabelHostname,
			},
		},
	}
}

func ProbeHTTP(path string, port int) corev1.ProbeHandler {
	return corev1.ProbeHandler{
		HTTPGet: &corev1.HTTPGetAction{
			Path: path,
			Port: intstr.FromInt(port),
		},
	}
}

func EnvVarDownAPI(name, fieldPath string) corev1.EnvVar {
	return corev1.EnvVar{
		Name:      name,
		ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: fieldPath}},
	}
}

func AnnotationPrometheus(path string, port int32) map[string]string {
	return map[string]string{
		"prometheus.io/path":   path,
		"prometheus.io/port":   fmt.Sprintf("%d", port),
		"prometheus.io/scrape": "true",
	}
}
