// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package karpenter

import "github.com/volvo-cars/lingoneks/pkg/platform/awsauth"

func AWSAuthMapRoles(nodeRoleARN, fargateRoleARN string) []*awsauth.RolesAuth {
	return []*awsauth.RolesAuth{
		{
			RoleARN:  nodeRoleARN,
			Username: "system:node:{{EC2PrivateDNSName}}",
			Groups: []string{
				"system:bootstrappers", "system:nodes",
			},
		},
		{
			RoleARN:  fargateRoleARN,
			Username: "system:node:{{SessionName}}",
			Groups: []string{
				"system:bootstrappers", "system:nodes", "system:node-proxier",
			},
		},
	}
}
