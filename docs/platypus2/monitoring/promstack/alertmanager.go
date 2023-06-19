// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package promstack

import (
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var KubePromtheusStackKubeAlertmanagerAlertmanager = &v1.Alertmanager{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-alertmanager",
			"app.kubernetes.io/instance":   "kube-promtheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "46.8.0",
			"chart":                        "kube-prometheus-stack-46.8.0",
			"heritage":                     "Helm",
			"release":                      "kube-promtheus-stack",
		},
		Name:      "kube-promtheus-stack-kube-alertmanager",
		Namespace: "monitoring",
	},
	Spec: v1.AlertmanagerSpec{
		ExternalURL: "http://kube-promtheus-stack-kube-alertmanager.monitoring:9093",
		Image:       P("quay.io/prometheus/alertmanager:v0.25.0"),
		LogFormat:   "logfmt",
		LogLevel:    "info",
		PortName:    "http-web",
		Replicas:    P(int32(1)),
		Retention:   v1.GoDuration("120h"),
		RoutePrefix: "/",
		SecurityContext: &corev1.PodSecurityContext{
			FSGroup:        P(int64(2000)),
			RunAsGroup:     P(int64(2000)),
			RunAsNonRoot:   P(true),
			RunAsUser:      P(int64(1000)),
			SeccompProfile: &corev1.SeccompProfile{Type: corev1.SeccompProfileType("RuntimeDefault")},
		},
		ServiceAccountName: "kube-promtheus-stack-kube-alertmanager",
		Version:            "v0.25.0",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "monitoring.coreos.com/v1",
		Kind:       "Alertmanager",
	},
}
