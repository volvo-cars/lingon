// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package cluster_eks

import (
	"fmt"

	aws "github.com/golingon/terraproviders/aws/4.60.0"
	"github.com/golingon/terraproviders/aws/4.60.0/dataiampolicydocument"
	"github.com/golingon/terraproviders/aws/4.60.0/ekscluster"
	tls "github.com/golingon/terraproviders/tls/4.0.4"
	"github.com/volvo-cars/lingon/pkg/terra"
)

var (
	S = terra.String
	N = terra.Number
)

var (
	arnClusterPolicy         = S("arn:aws:iam::aws:policy/AmazonEKSClusterPolicy")
	arnVPCResourceController = S("arn:aws:iam::aws:policy/AmazonEKSVPCResourceController")
	PORT_HTTPS               = N(443)
	PORT_DNS                 = N(53)
	PROTOCOL_TCP             = S("tcp")
	INGRESS                  = S("ingress")
	EGRESS                   = S("egress")
)

type ClusterOpts struct {
	Name             string
	Version          string
	VPCID            string
	PrivateSubnetIDs [3]string
}

type Cluster struct {
	EKSCluster           *aws.EksCluster              `validate:"required"`
	IAMPolicyDocument    *aws.DataIamPolicyDocument   `validate:"required"`
	IAMRole              *aws.IamRole                 `validate:"required"`
	IAMRoleClusterPolicy *aws.IamRolePolicyAttachment `validate:"required"`
	IAMRoleVPCController *aws.IamRolePolicyAttachment `validate:"required"`
	// SecurityGroup is the AWS security group for both the EKS control plane
	// and worker nodes
	SecurityGroup   *aws.SecurityGroup     `validate:"required"`
	IngressAllowAll *aws.SecurityGroupRule `validate:"required"`
	EgressAllowAll  *aws.SecurityGroupRule `validate:"required"`

	TLSCert         *tls.DataCertificate          `validate:"required"`
	IAMOIDCProvider *aws.IamOpenidConnectProvider `validate:"required"`
}

func NewEKSCluster(opts ClusterOpts) *Cluster {
	sg := aws.NewSecurityGroup(
		"eks", aws.SecurityGroupArgs{
			Name: S("eks-" + opts.Name),
			Description: S(
				fmt.Sprintf(
					"Main security group for EKS cluster %s", opts.Name,
				),
			),
			VpcId: S(opts.VPCID),
			Tags: terra.Map(
				map[string]terra.StringValue{
					"karpenter.sh/discovery": S("platypus-1"),
				},
			),
		},
	)

	sgAttrs := sg.Attributes()

	ingressAllowAll := aws.NewSecurityGroupRule(
		"eks", aws.SecurityGroupRuleArgs{
			SecurityGroupId:       sgAttrs.Id(),
			SourceSecurityGroupId: sgAttrs.Id(),
			Description:           S("Allow all for EKS control plane and managed worker nodes"),
			Protocol:              S("-1"),
			FromPort:              N(0),
			ToPort:                N(0),
			Type:                  INGRESS,
		},
	)
	egressAllowAll := aws.NewSecurityGroupRule(
		"node_egress_all", aws.SecurityGroupRuleArgs{
			SecurityGroupId: sgAttrs.Id(),
			Description:     S("Allow all egress"),
			Protocol:        S("-1"),
			FromPort:        N(0),
			ToPort:          N(0),
			Type:            EGRESS,
			CidrBlocks:      terra.List(S("0.0.0.0/0")),
		},
	)

	iamPolicyDocument := aws.NewDataIamPolicyDocument(
		"eks", aws.DataIamPolicyDocumentArgs{
			Statement: []dataiampolicydocument.Statement{
				{
					Sid:     S("EKSClusterAssumeRole"),
					Actions: terra.Set(S("sts:AssumeRole")),
					Principals: []dataiampolicydocument.Principals{
						{
							Type:        S("Service"),
							Identifiers: terra.Set(S("eks.amazonaws.com")),
						},
					},
				},
			},
		},
	)
	iamRole := aws.NewIamRole(
		"eks", aws.IamRoleArgs{
			Name:             S("eks-" + opts.Name),
			AssumeRolePolicy: iamPolicyDocument.Attributes().Json(),
		},
	)
	clusterPolicy := aws.NewIamRolePolicyAttachment(
		"cluster_policy", aws.IamRolePolicyAttachmentArgs{
			PolicyArn: arnClusterPolicy,
			Role:      iamRole.Attributes().Name(),
		},
	)
	vpcController := aws.NewIamRolePolicyAttachment(
		"vpc_controller", aws.IamRolePolicyAttachmentArgs{
			PolicyArn: arnVPCResourceController,
			Role:      iamRole.Attributes().Name(),
		},
	)

	eksCluster := aws.NewEksCluster(
		"eks", aws.EksClusterArgs{
			Name:    S(opts.Name),
			RoleArn: iamRole.Attributes().Arn(),
			VpcConfig: &ekscluster.VpcConfig{
				SecurityGroupIds: terra.Set(sgAttrs.Id()),
				SubnetIds:        terra.SetString(opts.PrivateSubnetIDs[:]...),
			},
			Version: S(opts.Version),
		},
	)
	eksCluster.DependsOn = terra.DependsOn(
		sg,
		iamRole,
		clusterPolicy,
		vpcController,
	)
	// How to add lifecycle to platform_version
	// eksCluster.Lifecycle = &terra.Lifecycle{
	// 	IgnoreChanges: terra.IgnoreChanges(
	// 		eksCluster.Attributes().PlatformVersion(),
	// 	),
	// }

	tlsCert := tls.NewDataCertificate(
		"eks", tls.DataCertificateArgs{
			Url: eksCluster.Attributes().Identity().Index(0).Oidc().Index(0).Issuer(),
		},
	)
	iamOIDCProvider := aws.NewIamOpenidConnectProvider(
		"eks", aws.IamOpenidConnectProviderArgs{
			ClientIdList: terra.List(terra.String("sts.amazonaws.com")),
			ThumbprintList: terra.CastAsList(
				tlsCert.Attributes().
					Certificates().
					Splat().Sha1Fingerprint(),
			),
			Url: eksCluster.Attributes().Identity().Index(0).Oidc().Index(0).Issuer(),
		},
	)

	return &Cluster{
		EKSCluster:           eksCluster,
		IAMPolicyDocument:    iamPolicyDocument,
		IAMRole:              iamRole,
		IAMRoleClusterPolicy: clusterPolicy,
		IAMRoleVPCController: vpcController,

		SecurityGroup:   sg,
		IngressAllowAll: ingressAllowAll,
		EgressAllowAll:  egressAllowAll,

		TLSCert:         tlsCert,
		IAMOIDCProvider: iamOIDCProvider,
	}
}
