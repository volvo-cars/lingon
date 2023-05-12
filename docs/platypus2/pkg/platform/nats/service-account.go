// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package nats

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var SA = &corev1.ServiceAccount{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "nats",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "nats",
			"app.kubernetes.io/version":    "2.9.16",
			"helm.sh/chart":                "nats-0.19.13",
		},
		Name:      "nats",
		Namespace: "nats",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ServiceAccount",
	},
}
