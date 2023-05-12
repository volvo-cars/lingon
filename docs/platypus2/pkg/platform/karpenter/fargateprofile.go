// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package karpenter

import (
	"fmt"

	aws "github.com/volvo-cars/lingoneks/providers/aws/4.66.1"
	"github.com/volvo-cars/lingoneks/providers/aws/4.66.1/dataiampolicydocument"
	"github.com/volvo-cars/lingoneks/providers/aws/4.66.1/eksfargateprofile"

	"github.com/volvo-cars/lingon/pkg/terra"
)

const (
	awsPolicyARNPrefix                  = "arn:aws:iam::aws:policy/"
	awsEKSWorkerNodePolicy              = "AmazonEKSWorkerNodePolicy"
	awsEC2ContainerRegistryReadOnly     = "AmazonEC2ContainerRegistryReadOnly"
	awsEKSFargatePodExecutionRolePolicy = "AmazonEKSFargatePodExecutionRolePolicy"
	awsEKSCNIPolicy                     = "AmazonEKS_CNI_Policy"
	awsSSMManagedInstanceCore           = "AmazonSSMManagedInstanceCore"
)

// FargateProfile is the AWS EKS Fargate profile for the Karpenter pods to
// run on
type FargateProfile struct {
	FargateProfile    *aws.EksFargateProfile         `validate:"required"`
	IAMRole           *aws.IamRole                   `validate:"required"`
	AssumeRole        *aws.DataIamPolicyDocument     `validate:"required"`
	PolicyAttachments []*aws.IamRolePolicyAttachment `validate:"required,dive,required"`
}

func newFargateProfile(opts InfraOpts) FargateProfile {
	arPolicy := aws.NewDataIamPolicyDocument(
		"fargate", aws.DataIamPolicyDocumentArgs{
			Statement: []dataiampolicydocument.Statement{
				{
					Effect:  S("Allow"),
					Actions: terra.SetString("sts:AssumeRole"),
					Principals: []dataiampolicydocument.Principals{
						{
							Type: S("Service"),
							Identifiers: terra.SetString(
								"eks-fargate-pods.amazonaws.com",
							),
						},
					},
				},
			},
		},
	)

	iamRole := aws.NewIamRole(
		"fargate", aws.IamRoleArgs{
			Name:             S(opts.Name + "-fargate"),
			Description:      S("IAM Role for Fargate profile for Karpenter pods to run"),
			AssumeRolePolicy: arPolicy.Attributes().Json(),
		},
	)

	policies := []string{
		awsEKSFargatePodExecutionRolePolicy,
		awsEKSCNIPolicy,
	}

	policyAttachments := make([]*aws.IamRolePolicyAttachment, len(policies))
	for i, policy := range policies {
		policyAttachments[i] = aws.NewIamRolePolicyAttachment(
			fmt.Sprintf("%s_attach_%s", "fargate", policy),
			aws.IamRolePolicyAttachmentArgs{
				PolicyArn: S(awsPolicyARNPrefix + policy),
				Role:      iamRole.Attributes().Name(),
			},
		)
	}

	fargateProfile := aws.NewEksFargateProfile(
		"karpenter", aws.EksFargateProfileArgs{
			ClusterName:         S(opts.ClusterName),
			FargateProfileName:  S("karpenter"),
			PodExecutionRoleArn: iamRole.Attributes().Arn(),
			SubnetIds:           terra.SetString(opts.PrivateSubnetIDs[:]...),
			Selector: []eksfargateprofile.Selector{
				{
					Namespace: S(Namespace),
				},
			},
		},
	)
	return FargateProfile{
		FargateProfile:    fargateProfile,
		IAMRole:           iamRole,
		AssumeRole:        arPolicy,
		PolicyAttachments: policyAttachments,
	}
}
