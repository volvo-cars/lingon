// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package karpenter

import (
	"fmt"

	aws "github.com/golingon/terraproviders/aws/4.60.0"
	"github.com/golingon/terraproviders/aws/4.60.0/dataiampolicydocument"
	"github.com/golingon/terraproviders/aws/4.60.0/iamrole"

	"github.com/volvo-cars/lingon/pkg/terra"
)

var S = terra.String

type InfraOpts struct {
	Name             string
	ClusterName      string
	ClusterARN       string
	PrivateSubnetIDs [3]string
	OIDCProviderArn  string
	OIDCProviderURL  string
}

func NewInfra(
	opts InfraOpts,
) Infra {
	ip := newInstanceProfile(opts.ClusterName)
	return Infra{
		FargateProfile: newFargateProfile(
			opts,
		),
		InstanceProfile: ip,
		Controller: newController(
			opts,
			ip.IAMRole,
		),
	}
}

// Infra is all the cloud resources for Karpenter to run on a pre-existing EKS
// cluster with an OIDC provider
type Infra struct {
	FargateProfile
	InstanceProfile
	Controller
}

// Controller contains all the resources for the Karpenter controller.
// This includes the IAM role used to manage EC2 nodes (and more) via IRSA
// (IAM Roles for ServiceAccounts) and the SQS (Simple Queue Service) for being
// notified about nodes being terminated (typically spot instances).
type Controller struct {
	NodeTerminationQueue
	IAMRole
}

func newController(
	opts InfraOpts,
	ipRole *aws.IamRole,
) Controller {
	queue := newNodeTerminationQueue(opts)
	return Controller{
		IAMRole: newIAMRole(
			opts,
			ipRole,
			queue.SimpleQueue,
		),
		NodeTerminationQueue: queue,
	}
}

type IAMRole struct {
	AssumeRolePolicy *aws.DataIamPolicyDocument `validate:"required"`
	Role             *aws.IamRole               `validate:"required"`
	RolePolicy       *aws.DataIamPolicyDocument `validate:"required"`
}

func newIAMRole(
	opts InfraOpts,
	ipRole *aws.IamRole,
	queue *aws.SqsQueue,
) IAMRole {
	assumeRolePolicy := aws.NewDataIamPolicyDocument(
		"karpenter_assume_role", aws.DataIamPolicyDocumentArgs{
			Statement: []dataiampolicydocument.Statement{
				{
					Actions: terra.Set(S("sts:AssumeRoleWithWebIdentity")),
					Principals: []dataiampolicydocument.Principals{
						{
							Type:        S("Federated"),
							Identifiers: terra.Set(S(opts.OIDCProviderArn)),
						},
					},
					Condition: []dataiampolicydocument.Condition{
						{
							Test:     S("StringEquals"),
							Variable: S(opts.OIDCProviderURL + ":sub"),
							Values: terra.ListString(
								fmt.Sprintf(
									"system:serviceaccount:%s:%s",
									Namespace,
									AppName,
								),
							),
						},
						{
							Test:     S("StringEquals"),
							Variable: S(opts.OIDCProviderURL + ":aud"),
							Values:   terra.ListString("sts.amazonaws.com"),
						},
					},
				},
			},
		},
	)
	policy := aws.NewDataIamPolicyDocument(
		"karpenter", aws.DataIamPolicyDocumentArgs{
			Statement: []dataiampolicydocument.Statement{
				{
					Actions: terra.SetString(
						"ec2:DescribeImages",
						"ec2:RunInstances",
						"ec2:DescribeSubnets",
						"ec2:DescribeSecurityGroups",
						"ec2:DescribeLaunchTemplates",
						"ec2:DescribeInstances",
						"ec2:DescribeInstanceTypes",
						"ec2:DescribeInstanceTypeOfferings",
						"ec2:DescribeAvailabilityZones",
						"ec2:DeleteLaunchTemplate",
						"ec2:CreateTags",
						"ec2:CreateLaunchTemplate",
						"ec2:CreateFleet",
						"ec2:DescribeSpotPriceHistory",
						"pricing:GetProducts",
						"ssm:GetParameter",
					),
					Effect:    S("Allow"),
					Resources: terra.SetString("*"),
				},
				{
					Actions: terra.SetString(
						"ec2:TerminateInstances",
						"ec2:DeleteLaunchTemplate",
					),
					Effect:    S("Allow"),
					Resources: terra.SetString("*"),
					Condition: []dataiampolicydocument.Condition{
						{
							Test:     S("StringEquals"),
							Variable: S("ec2:ResourceTag/karpenter.sh/discovery"),
							Values:   terra.ListString(opts.ClusterName),
						},
					},
				},
				{
					Actions: terra.SetString(
						"eks:DescribeCluster",
					),
					Effect:    S("Allow"),
					Resources: terra.SetString(opts.ClusterARN),
				},

				// The Karpenter IRSA role has to have permission to pass on the
				// InstanceProfile IAM Role
				{
					Actions: terra.SetString(
						"iam:PassRole",
					),
					Effect:    S("Allow"),
					Resources: terra.Set(ipRole.Attributes().Arn()),
				},

				// For AWS SQS spot interruption queue
				{
					Actions: terra.SetString(
						"sqs:DeleteMessage",
						"sqs:GetQueueUrl",
						"sqs:GetQueueAttributes",
						"sqs:ReceiveMessage",
					),
					Effect:    S("Allow"),
					Resources: terra.Set(queue.Attributes().Arn()),
				},
			},
		},
	)
	role := aws.NewIamRole(
		"karpenter", aws.IamRoleArgs{
			Name: terra.String(
				opts.Name + "-controller",
			),
			Description:      S("IAM Role for Karpenter Controller (pod) to assume"),
			AssumeRolePolicy: assumeRolePolicy.Attributes().Json(),

			InlinePolicy: []iamrole.InlinePolicy{
				{
					Name:   S(AppName),
					Policy: policy.Attributes().Json(),
				},
			},
		},
	)
	return IAMRole{
		AssumeRolePolicy: assumeRolePolicy,
		Role:             role,
		RolePolicy:       policy,
	}
}

func newNodeTerminationQueue(opts InfraOpts) NodeTerminationQueue {
	queue := aws.NewSqsQueue(
		"karpenter", aws.SqsQueueArgs{
			Name:                    S(opts.Name),
			MessageRetentionSeconds: terra.Number(300),
		},
	)
	policyDoc := aws.NewDataIamPolicyDocument(
		"node_termination_queue", aws.DataIamPolicyDocumentArgs{
			Statement: []dataiampolicydocument.Statement{
				{
					Sid:       S("SQSWrite"),
					Resources: terra.Set(queue.Attributes().Arn()),
					Actions:   terra.SetString("sqs:SendMessage"),
					Principals: []dataiampolicydocument.Principals{
						{
							Type: S("Service"),
							Identifiers: terra.SetString(
								"events.amazonaws.com",
								"sqs.amazonaws.com",
							),
						},
					},
				},
			},
		},
	)
	queuePolicy := aws.NewSqsQueuePolicy(
		"karpenter", aws.SqsQueuePolicyArgs{
			QueueUrl: queue.Attributes().Url(),
			Policy:   policyDoc.Attributes().Json(),
		},
	)
	return NodeTerminationQueue{
		SimpleQueue:         queue,
		QueuePolicy:         queuePolicy,
		QueuePolicyDocument: policyDoc,
	}
}

type NodeTerminationQueue struct {
	SimpleQueue         *aws.SqsQueue              `validate:"required"`
	QueuePolicy         *aws.SqsQueuePolicy        `validate:"required"`
	QueuePolicyDocument *aws.DataIamPolicyDocument `validate:"required"`
}
