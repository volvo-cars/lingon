// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	tkihcl "github.com/volvo-cars/lingon/pkg/internal/hcl"
)

// Value represents the most generic value in terra's minimal type system.
// Every other type must implement this type.
// It is never expected that a user should have to interact with this type directly
type Value[T any] interface {
	tkihcl.Tokenizer
	// InternalTraverse is for internal use
	InternalTraverse(hcl.Traverser) T
}

// Reference takes a list of steps as strings, that represent a reference to an object in
// a Terraform config (such as a resource or a data resource), and returns a ReferenceValue
// with those steps.
//
// Users are not expected to call this method. It is for the generated code to use when referencing
// a terraform object attribute
func Reference(steps ...string) ReferenceValue {
	tr := make(hcl.Traversal, len(steps))
	for i, s := range steps {
		if i == 0 {
			tr[i] = hcl.TraverseRoot{Name: s}
		} else {
			tr[i] = hcl.TraverseAttr{Name: s}
		}
	}
	return ReferenceValue{
		tr: tr,
	}
}

// ReferenceValue represents a reference to some value in a Terraform configuration.
// The reference might include things like indices (i.e. [0]), nested objects or even the splat
// operator (i.e. [*])
type ReferenceValue struct {
	tr hcl.Traversal
}

func (r ReferenceValue) InternalTokens() hclwrite.Tokens {
	return hclwrite.TokensForTraversal(r.tr)
}

func (r ReferenceValue) AsString() StringValue {
	return &stringRef{
		ref: r,
	}
}

func (r ReferenceValue) AsNumber() NumberValue {
	return &numberRef{
		ref: r,
	}
}

func (r ReferenceValue) AsBool() BoolValue {
	return &boolRef{
		ref: r,
	}
}

func (r ReferenceValue) InternalTraverse(step hcl.Traverser) ReferenceValue {
	return ReferenceValue{
		tr: append(r.tr, step),
	}
}
