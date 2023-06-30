package ingest

import (
	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IngestApp struct {
	kube.App
	Deploy *appsv1.Deployment
}

func NewIngestApp() *IngestApp {
	deployment := &appsv1.Deployment{
		TypeMeta: kubeutil.TypeStatefulSetV1,
		ObjectMeta: metav1.ObjectMeta{
			Name:      "schema-service-builder",
			Namespace: "system", // same as operator, for now
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: kubeutil.P[int32](1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "schema-service-builder",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "schema-service-builder",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "schema-service-builder",
							Image:           "platypus.io/schema-service-builder",
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
				},
			},
		},
	}

	return &IngestApp{
		Deploy: deployment,
	}
}
