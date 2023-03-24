// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package terrajen

import (
	"fmt"
	"strings"
)

func nodeBlockListValidateTags(n *node) string {
	validationTags := []string{
		fmt.Sprintf("min=%d", n.minItems),
	}
	if n.maxItems > 0 {
		validationTags = append(
			validationTags,
			fmt.Sprintf("max=%d", n.maxItems),
		)
	}
	return strings.Join(validationTags, ",")
}
