package registry

import (
	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Registry struct {
	kube.App
	Svc *corev1.Service
	SS  *appsv1.StatefulSet
}

func NewRegistry() *Registry {
	svc := &corev1.Service{
		TypeMeta: kubeutil.TypeServiceV1,
		ObjectMeta: metav1.ObjectMeta{
			Name: "registry",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "registry",
					Protocol: corev1.ProtocolTCP,
					Port:     5000,
				},
			},
			Type: corev1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"app": "registry",
			},
		},
	}
	ss := &appsv1.StatefulSet{
		TypeMeta: kubeutil.TypeStatefulSetV1,
		ObjectMeta: metav1.ObjectMeta{
			Name: "registry",
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: svc.Name,
			Replicas:    kubeutil.P[int32](1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "registry",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "registry",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "registry",
							Image:           "registry:2",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									Name:          "registry",
									ContainerPort: 5000,
									HostPort:      5000,
									Protocol:      corev1.ProtocolTCP,
								},
							},
						},
					},
				},
			},
		},
	}

	return &Registry{
		Svc: svc,
		SS:  ss,
	}
}
