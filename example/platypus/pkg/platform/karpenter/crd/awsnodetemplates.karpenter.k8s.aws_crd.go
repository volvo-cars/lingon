// CODE GENERATED.

package crd

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var AwsnodetemplatesKarpenterK8SAwsCRD = &apiextensionsv1.CustomResourceDefinition{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{"controller-gen.kubebuilder.io/version": "v0.11.3"},
		Name:        "awsnodetemplates.karpenter.k8s.aws",
	},
	Spec: apiextensionsv1.CustomResourceDefinitionSpec{
		Group: "karpenter.k8s.aws",
		Names: apiextensionsv1.CustomResourceDefinitionNames{
			Categories: []string{"karpenter"},
			Kind:       "AWSNodeTemplate",
			ListKind:   "AWSNodeTemplateList",
			Plural:     "awsnodetemplates",
			Singular:   "awsnodetemplate",
		},
		Scope: apiextensionsv1.ResourceScope("Cluster"),
		Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
			apiextensionsv1.CustomResourceDefinitionVersion{
				Name: "v1alpha1",
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Description: "AWSNodeTemplate is the Schema for the AWSNodeTemplate API",
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"apiVersion": apiextensionsv1.JSONSchemaProps{
								Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
								Type:        "string",
							},
							"kind": apiextensionsv1.JSONSchemaProps{
								Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
								Type:        "string",
							},
							"metadata": apiextensionsv1.JSONSchemaProps{Type: "object"},
							"spec": apiextensionsv1.JSONSchemaProps{
								Description: "AWSNodeTemplateSpec is the top level specification for the AWS Karpenter Provider. This will contain configuration necessary to launch instances in AWS.",
								Properties: map[string]apiextensionsv1.JSONSchemaProps{
									"amiFamily": apiextensionsv1.JSONSchemaProps{
										Description: "AMIFamily is the AMI family that instances use.",
										Type:        "string",
									},
									"amiSelector": apiextensionsv1.JSONSchemaProps{
										AdditionalProperties: &apiextensionsv1.JSONSchemaPropsOrBool{
											Allows: true,
											Schema: &apiextensionsv1.JSONSchemaProps{Type: "string"},
										},
										Description: "AMISelector discovers AMIs to be used by Amazon EC2 tags.",
										Type:        "object",
									},
									"apiVersion": apiextensionsv1.JSONSchemaProps{
										Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
										Type:        "string",
									},
									"blockDeviceMappings": apiextensionsv1.JSONSchemaProps{
										Description: "BlockDeviceMappings to be applied to provisioned nodes.",
										Items: &apiextensionsv1.JSONSchemaPropsOrArray{
											Schema: &apiextensionsv1.JSONSchemaProps{
												Properties: map[string]apiextensionsv1.JSONSchemaProps{
													"deviceName": apiextensionsv1.JSONSchemaProps{
														Description: "The device name (for example, /dev/sdh or xvdh).",
														Type:        "string",
													},
													"ebs": apiextensionsv1.JSONSchemaProps{
														Description: "EBS contains parameters used to automatically set up EBS volumes when an instance is launched.",
														Properties: map[string]apiextensionsv1.JSONSchemaProps{
															"deleteOnTermination": apiextensionsv1.JSONSchemaProps{
																Description: "DeleteOnTermination indicates whether the EBS volume is deleted on instance termination.",
																Type:        "boolean",
															},
															"encrypted": apiextensionsv1.JSONSchemaProps{
																Description: "Encrypted indicates whether the EBS volume is encrypted. Encrypted volumes can only be attached to instances that support Amazon EBS encryption. If you are creating a volume from a snapshot, you can't specify an encryption value.",
																Type:        "boolean",
															},
															"iops": apiextensionsv1.JSONSchemaProps{
																Description: "IOPS is the number of I/O operations per second (IOPS). For gp3, io1, and io2 volumes, this represents the number of IOPS that are provisioned for the volume. For gp2 volumes, this represents the baseline performance of the volume and the rate at which the volume accumulates I/O credits for bursting. \n The following are the supported values for each volume type: \n * gp3: 3,000-16,000 IOPS \n * io1: 100-64,000 IOPS \n * io2: 100-64,000 IOPS \n For io1 and io2 volumes, we guarantee 64,000 IOPS only for Instances built on the Nitro System (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html#ec2-nitro-instances). Other instance families guarantee performance up to 32,000 IOPS. \n This parameter is supported for io1, io2, and gp3 volumes only. This parameter is not supported for gp2, st1, sc1, or standard volumes.",
																Format:      "int64",
																Type:        "integer",
															},
															"kmsKeyID": apiextensionsv1.JSONSchemaProps{
																Description: "KMSKeyID (ARN) of the symmetric Key Management Service (KMS) CMK used for encryption.",
																Type:        "string",
															},
															"snapshotID": apiextensionsv1.JSONSchemaProps{
																Description: "SnapshotID is the ID of an EBS snapshot",
																Type:        "string",
															},
															"throughput": apiextensionsv1.JSONSchemaProps{
																Description: "Throughput to provision for a gp3 volume, with a maximum of 1,000 MiB/s. Valid Range: Minimum value of 125. Maximum value of 1000.",
																Format:      "int64",
																Type:        "integer",
															},
															"volumeSize": apiextensionsv1.JSONSchemaProps{
																AnyOf: []apiextensionsv1.JSONSchemaProps{
																	apiextensionsv1.JSONSchemaProps{Type: "integer"},
																	apiextensionsv1.JSONSchemaProps{Type: "string"},
																},
																Description:  "VolumeSize in GiBs. You must specify either a snapshot ID or a volume size. The following are the supported volumes sizes for each volume type: \n * gp2 and gp3: 1-16,384 \n * io1 and io2: 4-16,384 \n * st1 and sc1: 125-16,384 \n * standard: 1-1,024",
																Pattern:      "^(\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))))?$",
																XIntOrString: true,
															},
															"volumeType": apiextensionsv1.JSONSchemaProps{
																Description: "VolumeType of the block device. For more information, see Amazon EBS volume types (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html) in the Amazon Elastic Compute Cloud User Guide.",
																Type:        "string",
															},
														},
														Type: "object",
													},
												},
												Type: "object",
											},
										},
										Type: "array",
									},
									"context": apiextensionsv1.JSONSchemaProps{
										Description: "Context is a Reserved field in EC2 APIs https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_CreateFleet.html",
										Type:        "string",
									},
									"detailedMonitoring": apiextensionsv1.JSONSchemaProps{
										Description: "DetailedMonitoring controls if detailed monitoring is enabled for instances that are launched",
										Type:        "boolean",
									},
									"instanceProfile": apiextensionsv1.JSONSchemaProps{
										Description: "InstanceProfile is the AWS identity that instances use.",
										Type:        "string",
									},
									"kind": apiextensionsv1.JSONSchemaProps{
										Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
										Type:        "string",
									},
									"launchTemplate": apiextensionsv1.JSONSchemaProps{
										Description: "LaunchTemplateName for the node. If not specified, a launch template will be generated. NOTE: This field is for specifying a custom launch template and is exposed in the Spec as `launchTemplate` for backwards compatibility.",
										Type:        "string",
									},
									"metadataOptions": apiextensionsv1.JSONSchemaProps{
										Description: "MetadataOptions for the generated launch template of provisioned nodes. \n This specifies the exposure of the Instance Metadata Service to provisioned EC2 nodes. For more information, see Instance Metadata and User Data (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html) in the Amazon Elastic Compute Cloud User Guide. \n Refer to recommended, security best practices (https://aws.github.io/aws-eks-best-practices/security/docs/iam/#restrict-access-to-the-instance-profile-assigned-to-the-worker-node) for limiting exposure of Instance Metadata and User Data to pods. If omitted, defaults to httpEndpoint enabled, with httpProtocolIPv6 disabled, with httpPutResponseLimit of 2, and with httpTokens required.",
										Properties: map[string]apiextensionsv1.JSONSchemaProps{
											"httpEndpoint": apiextensionsv1.JSONSchemaProps{
												Description: "HTTPEndpoint enables or disables the HTTP metadata endpoint on provisioned nodes. If metadata options is non-nil, but this parameter is not specified, the default state is \"enabled\". \n If you specify a value of \"disabled\", instance metadata will not be accessible on the node.",
												Type:        "string",
											},
											"httpProtocolIPv6": apiextensionsv1.JSONSchemaProps{
												Description: "HTTPProtocolIPv6 enables or disables the IPv6 endpoint for the instance metadata service on provisioned nodes. If metadata options is non-nil, but this parameter is not specified, the default state is \"disabled\".",
												Type:        "string",
											},
											"httpPutResponseHopLimit": apiextensionsv1.JSONSchemaProps{
												Description: "HTTPPutResponseHopLimit is the desired HTTP PUT response hop limit for instance metadata requests. The larger the number, the further instance metadata requests can travel. Possible values are integers from 1 to 64. If metadata options is non-nil, but this parameter is not specified, the default value is 1.",
												Format:      "int64",
												Type:        "integer",
											},
											"httpTokens": apiextensionsv1.JSONSchemaProps{
												Description: "HTTPTokens determines the state of token usage for instance metadata requests. If metadata options is non-nil, but this parameter is not specified, the default state is \"optional\". \n If the state is optional, one can choose to retrieve instance metadata with or without a signed token header on the request. If one retrieves the IAM role credentials without a token, the version 1.0 role credentials are returned. If one retrieves the IAM role credentials using a valid signed token, the version 2.0 role credentials are returned. \n If the state is \"required\", one must send a signed token header with any instance metadata retrieval requests. In this state, retrieving the IAM role credentials always returns the version 2.0 credentials; the version 1.0 credentials are not available.",
												Type:        "string",
											},
										},
										Type: "object",
									},
									"securityGroupSelector": apiextensionsv1.JSONSchemaProps{
										AdditionalProperties: &apiextensionsv1.JSONSchemaPropsOrBool{
											Allows: true,
											Schema: &apiextensionsv1.JSONSchemaProps{Type: "string"},
										},
										Description: "SecurityGroups specify the names of the security groups.",
										Type:        "object",
									},
									"subnetSelector": apiextensionsv1.JSONSchemaProps{
										AdditionalProperties: &apiextensionsv1.JSONSchemaPropsOrBool{
											Allows: true,
											Schema: &apiextensionsv1.JSONSchemaProps{Type: "string"},
										},
										Description: "SubnetSelector discovers subnets by tags. A value of \"\" is a wildcard.",
										Type:        "object",
									},
									"tags": apiextensionsv1.JSONSchemaProps{
										AdditionalProperties: &apiextensionsv1.JSONSchemaPropsOrBool{
											Allows: true,
											Schema: &apiextensionsv1.JSONSchemaProps{Type: "string"},
										},
										Description: "Tags to be applied on ec2 resources like instances and launch templates.",
										Type:        "object",
									},
									"userData": apiextensionsv1.JSONSchemaProps{
										Description: "UserData to be applied to the provisioned nodes. It must be in the appropriate format based on the AMIFamily in use. Karpenter will merge certain fields into this UserData to ensure nodes are being provisioned with the correct configuration.",
										Type:        "string",
									},
								},
								Type: "object",
							},
							"status": apiextensionsv1.JSONSchemaProps{
								Description: "AWSNodeTemplateStatus contains the resolved state of the AWSNodeTemplate",
								Properties: map[string]apiextensionsv1.JSONSchemaProps{
									"securityGroups": apiextensionsv1.JSONSchemaProps{
										Description: "SecurityGroups contains the current Security Groups values that are available to the cluster under the SecurityGroups selectors.",
										Items: &apiextensionsv1.JSONSchemaPropsOrArray{
											Schema: &apiextensionsv1.JSONSchemaProps{
												Description: "SecurityGroupStatus contains resolved SecurityGroup selector values utilized for node launch",
												Properties: map[string]apiextensionsv1.JSONSchemaProps{
													"id": apiextensionsv1.JSONSchemaProps{
														Description: "Id of the security group",
														Type:        "string",
													},
												},
												Type: "object",
											},
										},
										Type: "array",
									},
									"subnets": apiextensionsv1.JSONSchemaProps{
										Description: "Subnets contains the current Subnet values that are available to the cluster under the subnet selectors.",
										Items: &apiextensionsv1.JSONSchemaPropsOrArray{
											Schema: &apiextensionsv1.JSONSchemaProps{
												Description: "SubnetStatus contains resolved Subnet selector values utilized for node launch",
												Properties: map[string]apiextensionsv1.JSONSchemaProps{
													"id": apiextensionsv1.JSONSchemaProps{
														Description: "Id of the subnet",
														Type:        "string",
													},
													"zone": apiextensionsv1.JSONSchemaProps{
														Description: "The associated availability zone",
														Type:        "string",
													},
												},
												Type: "object",
											},
										},
										Type: "array",
									},
								},
								Type: "object",
							},
						},
						Type: "object",
					},
				},
				Served:       true,
				Storage:      true,
				Subresources: &apiextensionsv1.CustomResourceSubresources{},
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "apiextensions.k8s.io/v1",
		Kind:       "CustomResourceDefinition",
	},
}
