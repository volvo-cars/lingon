package cilium

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var HubblePeerSvc = &v1.Service{
	TypeMeta: kubeutil.TypeMeta("Service"),
	ObjectMeta: metav1.ObjectMeta{
		Name:      "hubble-peer",
		Namespace: "kube-system",
		Labels:    map[string]string{"k8s-app": "cilium"},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{
			{
				Name:       "peer-service",
				Protocol:   v1.ProtocolTCP,
				Port:       443,
				TargetPort: intstr.IntOrString{IntVal: 4244},
			},
		},
		Selector:              map[string]string{"k8s-app": "cilium"},
		InternalTrafficPolicy: P(v1.ServiceInternalTrafficPolicyLocal),
	},
}
