// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	tkihcl "github.com/volvo-cars/lingon/pkg/internal/hcl"
)

// Dependency represents a Terraform dependency using the depends_on meta-argument
type Dependency interface {
	DependOn() Value[Reference]
}

// DependsOn returns a list of dependencies
func DependsOn(dependencies ...Dependency) []Dependency {
	return dependencies
}

var _ tkihcl.Tokenizer = (*Dependencies)(nil)

type Dependencies []Dependency

func (d Dependencies) InternalTokens() hclwrite.Tokens {
	if len(d) == 0 {
		return nil
	}
	tokens := hclwrite.TokensForIdentifier("[")
	length := len(d)
	for i, dep := range d {
		tokens = append(tokens, dep.DependOn().InternalTokens()...)
		if i < (length - 1) {
			tokens = append(tokens, hclwrite.TokensForIdentifier(",")...)
		}
	}
	tokens = append(tokens, hclwrite.TokensForIdentifier("]")...)
	return tokens
}
