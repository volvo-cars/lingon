// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package karpenter

import (
	"github.com/aws/karpenter-core/pkg/apis/v1alpha5"
	"github.com/aws/karpenter/pkg/apis/v1alpha1"
	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingoneks/infra"
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

func NewProvisioners(opts ProvisionersOpts) *Provisioners {
	// Scale down nodes after X seconds without workloads (excluding daemons).
	var ttlSecondsAfterEmpty int64 = 30
	// Kill each node after X seconds, testing this feature a bit.
	var ttlSecondsUntilExpired int64 = 3600

	nodeTmpl := v1alpha1.AWSNodeTemplate{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSNodeTemplate",
			APIVersion: "karpenter.k8s.aws/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{Name: KA.ProfileName},
		Spec: v1alpha1.AWSNodeTemplateSpec{
			AWS: v1alpha1.AWS{
				SubnetSelector: map[string]string{
					infra.KarpenterDiscoveryKey: opts.ClusterName,
				},
				SecurityGroupSelector: map[string]string{
					infra.KarpenterDiscoveryKey: opts.ClusterName,
				},
				Tags: map[string]string{
					infra.KarpenterDiscoveryKey: opts.ClusterName,
				},
			},
		},
	}

	defaultProvisioner := v1alpha5.Provisioner{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Provisioner",
			APIVersion: "karpenter.sh/v1alpha5",
		},
		ObjectMeta: metav1.ObjectMeta{Name: KA.ProfileName},
		Spec: v1alpha5.ProvisionerSpec{
			// Error from server (BadRequest): error when creating "STDIN":
			// admission webhook "validation.webhook.karpenter.sh"
			// denied the request: validation failed: expected exactly one,
			// got both: spec.consolidation.enabled, spec.ttlSecondsAfterEmpty
			//
			// Consolidation: &v1alpha5.Consolidation{
			// 	Enabled: P(true),
			// },
			ProviderRef: &v1alpha5.MachineTemplateRef{Name: nodeTmpl.Name},
			Requirements: []corev1.NodeSelectorRequirement{
				// see https://karpenter.sh/v0.29.0/concepts/provisioners/
				{
					Key:      "karpenter.k8s.aws/instance-category",
					Operator: corev1.NodeSelectorOpIn,
					Values:   []string{"m"},
				},
				{
					Key:      "karpenter.k8s.aws/instance-cpu",
					Operator: corev1.NodeSelectorOpIn,
					Values:   []string{"4", "8", "16"},
				},
				{
					Key:      "topology.kubernetes.io/zone",
					Operator: corev1.NodeSelectorOpIn,
					Values:   opts.AvailabilityZones[:],
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
	return &Provisioners{
		AWSNodeTemplate: &nodeTmpl,
		Default:         &defaultProvisioner,
	}
}
