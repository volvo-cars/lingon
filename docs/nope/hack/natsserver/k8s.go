package main

import (
	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NATSServer struct {
	kube.App
	Svc *corev1.Service
	SS  *appsv1.StatefulSet
}

func NewNATSServer() *NATSServer {
	svc := &corev1.Service{
		TypeMeta: kubeutil.TypeServiceV1,
		ObjectMeta: metav1.ObjectMeta{
			Name:      "natsserver",
			Namespace: "system", // same as operator, for now
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "natsserver",
					Protocol: corev1.ProtocolTCP,
					Port:     4222,
				},
			},
			Type: corev1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"app": "natsserver",
			},
		},
	}
	ss := &appsv1.StatefulSet{
		TypeMeta: kubeutil.TypeStatefulSetV1,
		ObjectMeta: metav1.ObjectMeta{
			Name:      "natsserver",
			Namespace: "system", // same as operator, for now
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: svc.Name,
			Replicas:    kubeutil.P[int32](1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "natsserver",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "natsserver",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "natsserver",
							Image:           "platypus.io/natsserver",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									Name:          "natsserver",
									ContainerPort: 4222,
									HostPort:      4222,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							// TODO: This should be written to a volume...
							Args: []string{"-out=/tmp/nats"},
						},
					},
				},
			},
		},
	}

	return &NATSServer{
		Svc: svc,
		SS:  ss,
	}
}
