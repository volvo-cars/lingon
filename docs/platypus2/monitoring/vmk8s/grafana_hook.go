package vmk8s

import (
	"github.com/VictoriaMetrics/operator/api/victoriametrics/v1beta1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GrafanaTest struct {
	GrafanaTestCM   *corev1.ConfigMap
	GrafanaTestPO   *corev1.Pod
	GrafanaTestRB   *rbacv1.RoleBinding
	GrafanaTestRole *rbacv1.Role
	GrafanaTestSA   *corev1.ServiceAccount
}

func NewGrafanaTest() *GrafanaTest {
	return &GrafanaTest{
		GrafanaTestCM:   GrafanaTestCM,
		GrafanaTestPO:   GrafanaTestPO,
		GrafanaTestRB:   GrafanaTestRB,
		GrafanaTestRole: GrafanaTestRole,
		GrafanaTestSA:   GrafanaTestSA,
	}
}

var GrafanaTestCM = &corev1.ConfigMap{
	Data: map[string]string{
		"run.sh": `
@test "Test Health" {
  url="http://vmk8s-grafana/api/health"
  code=$(wget --server-response --spider --timeout 90 --tries 10 ${url} 2>&1 | awk '/^  HTTP/{print $2}')
  [ "$code" == "200" ]
}
`,
	},
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "test-success",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
		},
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name:      "vmk8s-grafana-test",
		Namespace: "monitoring",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ConfigMap",
	},
}

var GrafanaTestRole = &rbacv1.Role{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "test-success",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
		},
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name:      "vmk8s-grafana-test",
		Namespace: "monitoring",
	},
	Rules: []rbacv1.PolicyRule{
		{
			APIGroups:     []string{"policy"},
			ResourceNames: []string{"vmk8s-grafana-test"},
			Resources:     []string{"podsecuritypolicies"},
			Verbs:         []string{"use"},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "Role",
	},
}

var GrafanaTestRB = &rbacv1.RoleBinding{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "test-success",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
		},
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name:      "vmk8s-grafana-test",
		Namespace: "monitoring",
	},
	RoleRef: rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "Role",
		Name:     "vmk8s-grafana-test",
	},
	Subjects: []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      "vmk8s-grafana-test",
			Namespace: "monitoring",
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "rbac.authorization.k8s.io/v1",
		Kind:       "RoleBinding",
	},
}

var GrafanaTestSA = &corev1.ServiceAccount{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "test-success",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
		},
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name:      "vmk8s-grafana-test",
		Namespace: "monitoring",
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ServiceAccount",
	},
}

var GrafanaTestPO = &corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "test-success",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
		},
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "grafana",
			"app.kubernetes.io/version":    "9.3.0",
			"helm.sh/chart":                "grafana-6.44.11",
		},
		Name:      "vmk8s-grafana-test",
		Namespace: "monitoring",
	},
	Spec: corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Command: []string{
					"/opt/bats/bin/bats",
					"-t",
					"/tests/run.sh",
				},
				Image:           "bats/bats:v1.4.1",
				ImagePullPolicy: corev1.PullIfNotPresent,
				Name:            "vmk8s-test",
				VolumeMounts: []corev1.VolumeMount{
					{
						MountPath: "/tests",
						Name:      "tests",
						ReadOnly:  true,
					},
				},
			},
		},
		RestartPolicy:      corev1.RestartPolicy("Never"),
		ServiceAccountName: "vmk8s-grafana-test",
		Volumes: []corev1.Volume{
			{
				Name:         "tests",
				VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: "vmk8s-grafana-test"}}},
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Pod",
	},
}

var GrafanaScrape = &v1beta1.VMServiceScrape{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "vmk8s",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "victoria-metrics-k8s-stack",
			"app.kubernetes.io/version":    "v1.91.2",
			"helm.sh/chart":                "victoria-metrics-k8s-stack-0.16.3",
		},
		Name:      "vmk8s-victoria-metrics-k8s-stack-grafana",
		Namespace: "monitoring",
	},
	Spec: v1beta1.VMServiceScrapeSpec{
		Endpoints: []v1beta1.Endpoint{{Port: "service"}},
		Selector: metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app.kubernetes.io/instance": "vmk8s",
				"app.kubernetes.io/name":     "grafana",
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "operator.victoriametrics.com/v1beta1",
		Kind:       "VMServiceScrape",
	},
}
