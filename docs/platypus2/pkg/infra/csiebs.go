// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package infra

import (
	"github.com/volvo-cars/lingon/pkg/terra"
	aws "github.com/volvo-cars/lingoneks/providers/aws/4.66.1"
	"github.com/volvo-cars/lingoneks/providers/aws/4.66.1/dataiampolicydocument"
	"github.com/volvo-cars/lingoneks/providers/aws/4.66.1/iamrole"
)

type CSI struct {
	CSIDriver *aws.EksAddon `validate:"required"`
	IAMRole   `validate:"required"`
}

type CSIOpts struct {
	ClusterName     string
	OIDCProviderArn string
	OIDCProviderURL string
}

type IAMRole struct {
	AssumeRolePolicy *aws.DataIamPolicyDocument   `validate:"required"`
	Role             *aws.IamRole                 `validate:"required"`
	RolePolicy       *aws.DataIamPolicyDocument   `validate:"required"`
	PolicyAttach     *aws.IamRolePolicyAttachment `validate:"required"`
}

func NewCSIEBS(opts CSIOpts) *CSI {
	ir := newIAMRole(opts)
	return &CSI{
		CSIDriver: aws.NewEksAddon(
			opts.ClusterName+"-csiebs", aws.EksAddonArgs{
				AddonName:             S("aws-ebs-csi-driver"),
				AddonVersion:          S("v1.19.0-eksbuild.1"),
				ClusterName:           S(opts.ClusterName),
				ServiceAccountRoleArn: ir.Role.Attributes().Arn(),
				ResolveConflicts:      S("OVERWRITE"),
			},
		),
		IAMRole: *ir,
	}
}

func newIAMRole(opts CSIOpts) *IAMRole {
	assumeRolePolicy := aws.NewDataIamPolicyDocument(
		"csi_assume_role", aws.DataIamPolicyDocumentArgs{
			Statement: []dataiampolicydocument.Statement{
				{
					Actions: terra.Set(S("sts:AssumeRoleWithWebIdentity")),
					Effect:  S("Allow"),

					Condition: []dataiampolicydocument.Condition{
						{
							Test:     S("StringEquals"),
							Variable: S(opts.OIDCProviderURL + ":sub"),
							Values:   terra.List(S("system:serviceaccount:kube-system:ebs-csi-controller-sa")),
						},
						{
							Test:     S("StringEquals"),
							Variable: S(opts.OIDCProviderURL + ":aud"),
							Values:   terra.ListString("sts.amazonaws.com"),
						},
					},
					Principals: []dataiampolicydocument.Principals{
						{
							Type:        S("Federated"),
							Identifiers: terra.Set(S(opts.OIDCProviderArn)),
						},
					},
				},
			},
		},
	)

	// small utility function to avoid repeting fields in the policy
	cond := func(action, v, val string) dataiampolicydocument.Statement {
		return dataiampolicydocument.Statement{
			Effect:    S("Allow"),
			Actions:   terra.SetString(action),
			Resources: terra.SetString("*"),
			Condition: []dataiampolicydocument.Condition{
				{
					Test:     S("StringLike"),
					Variable: S(v),
					Values:   terra.ListString(val),
				},
			},
		}
	}

	// converted from
	// https://github.com/kubernetes-sigs/aws-ebs-csi-driver/blob/master/docs/example-iam-policy.json
	//
	policy := aws.NewDataIamPolicyDocument(
		"csiebs", aws.DataIamPolicyDocumentArgs{
			Statement: []dataiampolicydocument.Statement{
				{
					Effect: S("Allow"),
					Actions: terra.SetString(
						"ec2:CreateSnapshot",
						"ec2:AttachVolume",
						"ec2:DetachVolume",
						"ec2:ModifyVolume",
						"ec2:DescribeAvailabilityZones",
						"ec2:DescribeInstances",
						"ec2:DescribeSnapshots",
						"ec2:DescribeTags",
						"ec2:DescribeVolumes",
						"ec2:DescribeVolumesModifications",
					),
					Resources: terra.SetString("*"),
				},
				{
					Effect:  S("Allow"),
					Actions: terra.SetString("ec2:CreateTags"),
					Resources: terra.SetString(
						"arn:aws:ec2:*:*:volume/*",
						"arn:aws:ec2:*:*:snapshot/*",
					),
					Condition: []dataiampolicydocument.Condition{
						{
							Test:     S("StringEquals"),
							Variable: S("ec2:CreateAction"),
							Values: terra.ListString(
								"CreateVolume",
								"CreateSnapshot",
							),
						},
					},
				},
				{
					Effect:  S("Allow"),
					Actions: terra.SetString("ec2:DeleteTags"),
					Resources: terra.SetString(
						"arn:aws:ec2:*:*:volume/*",
						"arn:aws:ec2:*:*:snapshot/*",
					),
				},
				cond(
					"ec2:CreateVolume",
					"aws:RequestTag/ebs.csi.aws.com/cluster", "true",
				),
				cond(
					"ec2:CreateVolume",
					"aws:RequestTag/CSIVolumeName", "*",
				),
				cond(
					"ec2:DeleteVolume",
					"ec2:ResourceTag/ebs.csi.aws.com/cluster", "true",
				),
				cond(
					"ec2:DeleteVolume",
					"ec2:ResourceTag/CSIVolumeName", "*",
				),
				cond(
					"ec2:DeleteVolume",
					"ec2:ResourceTag/kubernetes.io/created-for/pvc/name", "*",
				),
				cond(
					"ec2:DeleteSnapshot",
					"ec2:ResourceTag/CSIVolumeSnapshotName", "*",
				),
				cond(
					"ec2:DeleteSnapshot",
					"ec2:ResourceTag/ebs.csi.aws.com/cluster", "true",
				),
			},
		},
	)

	csiRole := aws.NewIamRole(
		"csiebs_role", aws.IamRoleArgs{
			Name:             S(opts.ClusterName + "-csi"),
			Description:      S("IAM Role for CSI EBS driver"),
			AssumeRolePolicy: assumeRolePolicy.Attributes().Json(),

			InlinePolicy: []iamrole.InlinePolicy{
				{
					Name:   S("csi-ebs-driver"),
					Policy: policy.Attributes().Json(),
				},
			},
		},
	)
	pa := aws.NewIamRolePolicyAttachment(
		"csiebs_attach_AmazonEBSCSIDriverPolicy",
		aws.IamRolePolicyAttachmentArgs{
			PolicyArn: S("arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy"),
			Role:      csiRole.Attributes().Name(),
		},
	)

	return &IAMRole{
		AssumeRolePolicy: assumeRolePolicy,
		Role:             csiRole,
		RolePolicy:       policy,
		PolicyAttach:     pa,
	}
}
