// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	tkihcl "github.com/volvo-cars/lingon/pkg/internal/hcl"
	"github.com/zclconf/go-cty/cty"
)

// Value represents the most generic type in terra's minimal type system.
// Every other type must implement this type.
// It is never expected that a user should have to interact with this type directly
type Value[T any] interface {
	tkihcl.Tokenizer
	// InternalWithRef returns a copy of this type T with the provided
	// reference. This is used internally to create a reference to attributes
	// within sets, lists,
	// maps and custom attributes types in the generated code.
	InternalWithRef(Reference) T
}

// ReferenceAttribute takes an address to a Terraform attribute.
// The address is represented as a list of strings that are converted into a
// references.
//
// Users are not expected to call this method. It is for the generated code to use when referencing
// a terraform object attribute
func ReferenceAttribute(address ...string) Reference {
	tr := make(hcl.Traversal, len(address))
	for i, s := range address {
		if i == 0 {
			tr[i] = hcl.TraverseRoot{Name: s}
		} else {
			tr[i] = hcl.TraverseAttr{Name: s}
		}
	}
	return Reference{
		tr: tr,
	}
}

// ReferenceSingle creates an instance of T with the given reference.
// It is a helper method for the generated code to use, to make it consistent
// with creating maps, sets, etc.
func ReferenceSingle[T Value[T]](ref Reference) T {
	var v T
	return v.InternalWithRef(ref)
}

// Reference represents a reference to some value in a Terraform configuration.
// The reference might include things like indices (i.e. [0]),
// nested objects or even the splat operator (i.e. [*])
type Reference struct {
	tr hcl.Traversal
}

func (r Reference) InternalTokens() hclwrite.Tokens {
	return hclwrite.TokensForTraversal(r.tr)
}

func (r Reference) InternalWithRef(ref Reference) Reference {
	return ref.copy()
}

// Append appends the given string to the reference
func (r Reference) Append(name string) Reference {
	cp := r.copy()
	return Reference{
		tr: append(cp.tr, hcl.TraverseAttr{Name: name}),
	}
}

// index adds a reference to an index (of a list/set) to the reference
func (r Reference) index(i int) Reference {
	cp := r.copy()
	return Reference{
		tr: append(cp.tr, hcl.TraverseIndex{Key: cty.NumberIntVal(int64(i))}),
	}
}

// key adds a reference to a key (of a map) to the reference
func (r Reference) key(k string) Reference {
	cp := r.copy()
	return Reference{
		tr: append(cp.tr, hcl.TraverseIndex{Key: cty.StringVal(k)}),
	}
}

// splat adds the Terraform splat operator to a reference
func (r Reference) splat() Reference {
	cp := r.copy()
	return Reference{
		tr: append(cp.tr, hcl.TraverseSplat{}),
	}
}

// copy makes a copy of the reference so that any slices can be safely passed
// around with modifying the original
func (r Reference) copy() Reference {
	tr := make(hcl.Traversal, len(r.tr))
	copy(tr, r.tr)
	return Reference{
		tr: tr,
	}
}
