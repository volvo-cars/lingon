// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package metricsserver

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
)

var V1Beta1MetricsK8SIoApiservices = &apiregistrationv1.APIService{
	ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "metrics-server",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/name":       "metrics-server",
			"app.kubernetes.io/version":    "0.6.3",
			"helm.sh/chart":                "metrics-server-3.10.0",
		},
		Name: "v1beta1.metrics.k8s.io",
	},
	Spec: apiregistrationv1.APIServiceSpec{
		Group:                 "metrics.k8s.io",
		GroupPriorityMinimum:  int32(100),
		InsecureSkipTLSVerify: true,
		Service: &apiregistrationv1.ServiceReference{
			Name:      "metrics-server",
			Namespace: namespace,
			Port:      P(int32(443)),
		},
		Version:         "v1beta1",
		VersionPriority: int32(100),
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "apiregistration.k8s.io/v1",
		Kind:       "APIService",
	},
}
