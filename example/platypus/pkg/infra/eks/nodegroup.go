package eks

import (
	"github.com/volvo-cars/lingon/example/platypus/gen/providers/aws"
	"github.com/volvo-cars/lingon/example/platypus/gen/providers/aws/dataiampolicydocument"
	"github.com/volvo-cars/lingon/example/platypus/gen/providers/aws/eksnodegroup"

	"github.com/volvo-cars/lingon/pkg/terra"
)

type ClusterNodes struct {
	NodeIAMPolicyDocument *aws.DataIamPolicyDocument   `validate:"required"`
	NodeIAMRole           *aws.IamRole                 `validate:"required"`
	WorkerNodePolicy      *aws.IamRolePolicyAttachment `validate:"required"`
	ECRReadOnlyPolicy     *aws.IamRolePolicyAttachment `validate:"required"`
	CNIPolicy             *aws.IamRolePolicyAttachment `validate:"required"`
	ManagedNodeGroup      *aws.EksNodeGroup            `validate:"required"`
}

func newEksClusterNodes(
	eksCluster *aws.EksCluster,
	subnetIDs terra.SetValue[terra.StringValue],
) ClusterNodes {
	iamPolicyDocument := aws.NewDataIamPolicyDocument(
		"node", aws.DataIamPolicyDocumentArgs{
			Statement: []dataiampolicydocument.Statement{
				{
					Sid:     S("EKSNodeAssumeRole"),
					Actions: terra.Set(S("sts:AssumeRole")),
					Principals: []dataiampolicydocument.Principals{
						{
							Type:        S("Service"),
							Identifiers: terra.Set(S("ec2.amazonaws.com")),
						},
					},
				},
			},
		},
	)
	iamRole := aws.NewIamRole(
		"node", aws.IamRoleArgs{
			AssumeRolePolicy: iamPolicyDocument.Attributes().Json(),
		},
	)

	workerNodePolicy := aws.NewIamRolePolicyAttachment(
		"worker_node_policy", aws.IamRolePolicyAttachmentArgs{
			PolicyArn: S("arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"),
			Role:      iamRole.Attributes().Name(),
		},
	)
	ecrReadOnlyPolicy := aws.NewIamRolePolicyAttachment(
		"ecr_ready_only_policy", aws.IamRolePolicyAttachmentArgs{
			PolicyArn: S("arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"),
			Role:      iamRole.Attributes().Name(),
		},
	)
	cniPolicy := aws.NewIamRolePolicyAttachment(
		"cni_policy", aws.IamRolePolicyAttachmentArgs{
			PolicyArn: S("arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"),
			Role:      iamRole.Attributes().Name(),
		},
	)

	return ClusterNodes{
		NodeIAMPolicyDocument: iamPolicyDocument,
		NodeIAMRole:           iamRole,
		WorkerNodePolicy:      workerNodePolicy,
		ECRReadOnlyPolicy:     ecrReadOnlyPolicy,
		CNIPolicy:             cniPolicy,
		ManagedNodeGroup: aws.NewEksNodeGroup(
			"node", aws.EksNodeGroupArgs{
				ClusterName:   eksCluster.Attributes().Name(),
				InstanceTypes: terra.List(S("t3.small")),
				NodeRoleArn:   iamRole.Attributes().Arn(),
				NodeGroupName: S("bootstrap"),
				SubnetIds:     subnetIDs,
				ScalingConfig: &eksnodegroup.ScalingConfig{
					DesiredSize: N(2),
					MinSize:     N(2),
					MaxSize:     N(2),
				},
				Taint: []eksnodegroup.Taint{
					{
						Effect: S("NO_EXECUTE"),
						Key:    S("node.cilium.io/agent-not-ready"),
						Value:  S("true"),
					},
					{
						Effect: S("NO_SCHEDULE"),
						Key:    S("dedicated"),
						Value:  S("karpenter"),
					},
				},
				UpdateConfig: &eksnodegroup.UpdateConfig{
					MaxUnavailablePercentage: N(33),
				},
			},
		),
	}
}
