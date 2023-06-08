// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package infra

import (
	"github.com/volvo-cars/lingon/pkg/terra"
)

const KarpenterDiscoveryKey = "karpenter.sh/discovery"

var (
	S        = terra.String
	N        = terra.Number
	B        = terra.Bool
	Anywhere = S("0.0.0.0/0")
)

var TFBaseTags = map[string]string{
	TagManagedBy:  TagManagedByValue,
	TagEnv:        "Dev",
	TagPurpose:    "Experiment",
	TagProject:    "Platform",
	TagOwner:      "mlops",
	TagCostCenter: "mlops",
	"terraform":   "true",
}

func TFTags(app, role string) map[string]string {
	return map[string]string{
		TagName:    app,
		TagAppID:   app,
		TagAppRole: role,
	}
}
