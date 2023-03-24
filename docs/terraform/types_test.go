// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terraform

import (
	aws "github.com/golingon/terraproviders/aws/4.60.0"
	"github.com/volvo-cars/lingon/pkg/terra"
)

func Example_typesVars() {
	var (
		S = terra.String
		B = terra.Bool
	)

	_ = aws.NewVpc(
		"vpc", aws.VpcArgs{
			CidrBlock:        S("10.0.0.0/16"),
			EnableDnsSupport: B(true),
		},
	)
}
