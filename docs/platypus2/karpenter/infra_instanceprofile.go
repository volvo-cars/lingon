// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package karpenter

import (
	"fmt"

	aws "github.com/golingon/terraproviders/aws/5.13.1"
	"github.com/golingon/terraproviders/aws/5.13.1/dataiampolicydocument"

	"github.com/volvo-cars/lingon/pkg/terra"
)

// InstanceProfile is the AWS EC2 Instance Profile for the nodes provisioned by
// Karpenter to use.
type InstanceProfile struct {
	InstanceProfile   *aws.IamInstanceProfile        `validate:"required"`
	IAMRole           *aws.IamRole                   `validate:"required"`
	AssumeRole        *aws.DataIamPolicyDocument     `validate:"required"`
	PolicyAttachments []*aws.IamRolePolicyAttachment `validate:"required,dive,required"`
}

func newInstanceProfile() InstanceProfile {
	arPolicy := aws.NewDataIamPolicyDocument(
		"eks_node", aws.DataIamPolicyDocumentArgs{
			Statement: []dataiampolicydocument.Statement{
				{
					Sid:     S("EKSNodeAssumeRole"),
					Effect:  S("Allow"),
					Actions: terra.SetString("sts:AssumeRole"),
					Principals: []dataiampolicydocument.Principals{
						{
							Type: S("Service"),
							Identifiers: terra.SetString(
								"ec2.amazonaws.com",
							),
						},
					},
				},
			},
		},
	)

	iamRole := aws.NewIamRole(
		"eks_node", aws.IamRoleArgs{
			Name:             S("platypus-karpenter-node"),
			Description:      S("IAM Role for Karpenter's InstanceProfile to use when launching nodes"),
			AssumeRolePolicy: arPolicy.Attributes().Json(),
		},
	)

	policies := []string{
		awsEKSWorkerNodePolicy,
		awsEC2ContainerRegistryReadOnly,
		awsEKSCNIPolicy,
		awsSSMManagedInstanceCore,
	}

	policyAttachments := make([]*aws.IamRolePolicyAttachment, len(policies))
	for i, policy := range policies {
		policyAttachments[i] = aws.NewIamRolePolicyAttachment(
			fmt.Sprintf("eks_node_attach_%s", policy),
			aws.IamRolePolicyAttachmentArgs{
				PolicyArn: S(awsPolicyARNPrefix + policy),
				Role:      iamRole.Attributes().Name(),
			},
		)
	}

	instanceProfile := aws.NewIamInstanceProfile(
		KA.Name, aws.IamInstanceProfileArgs{
			Name: S("platypus-karpenter-instance-profile"),
			Role: iamRole.Attributes().Name(),
		},
	)

	return InstanceProfile{
		InstanceProfile:   instanceProfile,
		IAMRole:           iamRole,
		AssumeRole:        arPolicy,
		PolicyAttachments: policyAttachments,
	}
}
