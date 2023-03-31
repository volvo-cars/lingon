package karpenter

import (
	"github.com/aws/karpenter-core/pkg/apis/v1alpha5"
	"github.com/aws/karpenter/pkg/apis/v1alpha1"
	"github.com/volvo-cars/lingon/pkg/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Provisioners struct {
	kube.App

	AWSNodeTemplate *v1alpha1.AWSNodeTemplate
	Default         *v1alpha5.Provisioner
}

type ProvisionersOpts struct {
	ClusterName       string
	AvailabilityZones [3]string
}

func NewProvisioners(
	opts ProvisionersOpts,
) Provisioners {
	var ttlSecondsAfterEmpty int64 = 30
	// Kill each node after one hour, testing this feature a bit
	var ttlSecondsUntilExpired int64 = 3600

	nodeTmpl := v1alpha1.AWSNodeTemplate{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSNodeTemplate",
			APIVersion: "karpenter.k8s.aws/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec: v1alpha1.AWSNodeTemplateSpec{
			AWS: v1alpha1.AWS{
				SubnetSelector: map[string]string{
					"karpenter.sh/discovery": opts.ClusterName,
				},
				SecurityGroupSelector: map[string]string{
					"karpenter.sh/discovery": opts.ClusterName,
				},
				Tags: map[string]string{
					"karpenter.sh/discovery": opts.ClusterName,
				},
			},
		},
	}

	defaultProvisioner := v1alpha5.Provisioner{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Provisioner",
			APIVersion: "karpenter.sh/v1alpha5",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec: v1alpha5.ProvisionerSpec{
			ProviderRef: &v1alpha5.MachineTemplateRef{
				Name: nodeTmpl.Name,
			},
			// TODO: cilium taints
			// StartupTaints: []corev1.Taint{
			//	{
			//		Key:    "node.cilium.io/agent-not-ready",
			//		Effect: "NoExecute",
			//	},
			// },
			Requirements: []corev1.NodeSelectorRequirement{
				{
					Key:      "karpenter.k8s.aws/instance-category",
					Operator: corev1.NodeSelectorOpIn,
					Values:   []string{"m"},
				},
				{
					Key:      "karpenter.k8s.aws/instance-cpu",
					Operator: corev1.NodeSelectorOpIn,
					Values:   []string{"4"},
				},
				{
					Key:      "topology.kubernetes.io/zone",
					Operator: corev1.NodeSelectorOpIn,
					// TODO: get values from terraform
					Values: opts.AvailabilityZones[:],
				},
				{
					Key:      "karpenter.sh/capacity-type",
					Operator: corev1.NodeSelectorOpIn,
					Values:   []string{"spot"},
				},
			},
			TTLSecondsAfterEmpty:   &ttlSecondsAfterEmpty,
			TTLSecondsUntilExpired: &ttlSecondsUntilExpired,
			// Limits:                 nil,
		},
	}
	return Provisioners{
		AWSNodeTemplate: &nodeTmpl,
		Default:         &defaultProvisioner,
	}
}
