// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	tkihcl "github.com/volvo-cars/lingon/pkg/internal/hcl"
	"github.com/zclconf/go-cty/cty"
)

type Referencer interface {
	// InternalRef returns a copy of the reference stored, if any.
	// If the Value T is not a reference, this method should panic to avoid
	// any nasty hidden errors (i.e.
	// silently converting a value to a reference).
	//
	// Internal: users should **not** use this!
	InternalRef() Reference
}

// Value represents the most generic type in terra's minimal type system.
// Every other type must implement this type.
// It is never expected that a user should have to interact with this type
// directly.
//
// Value exposes some functions prefixed with "Internal". It is not possible
// to make these methods non-public, as generated code relies on this,
// but consumers of this package should **not** use these methods, please :)
type Value[T any] interface {
	tkihcl.Tokenizer
	Referencer
	// InternalWithRef returns a copy of this type T with the provided
	// reference. This is used internally to create a reference to attributes
	// within sets, lists,
	// maps and custom attributes types in the generated code.
	//
	// Internal: users should **not** use this!
	InternalWithRef(Reference) T
}

// ReferenceResource takes a resource and returns a Reference which is the
// address to that resource in the Terraform configuration.
func ReferenceResource(res Resource) Reference {
	return Reference{
		underlyingType: referenceResource,
		res:            res,
	}
}

// ReferenceDataResource takes a data resource and returns a Reference which
// is the address to that data resource in the Terraform configuration.
func ReferenceDataResource(data DataResource) Reference {
	return Reference{
		underlyingType: referenceDataResource,
		data:           data,
	}
}

// ReferenceAsSingle creates an instance of T with the given reference.
// It is a helper method for the generated code to use, to make it consistent
// with creating maps, sets, etc.
func ReferenceAsSingle[T Value[T]](ref Reference) T {
	var v T
	return v.InternalWithRef(ref)
}

type referenceUnderlyingType int

const (
	referenceResource     referenceUnderlyingType = 1
	referenceDataResource referenceUnderlyingType = 2
)

var _ tkihcl.Tokenizer = (*Reference)(nil)

// Reference represents a reference to some value in a Terraform configuration.
// The reference might include things like indices (i.e. [0]),
// nested objects or even the splat operator (i.e. [*])
type Reference struct {
	underlyingType referenceUnderlyingType
	res            Resource
	data           DataResource
	tr             hcl.Traversal
}

// InternalTokens returns the tokens to represent this reference in Terraform
// configurations
func (r Reference) InternalTokens() hclwrite.Tokens {
	var tr hcl.Traversal
	switch r.underlyingType {
	case referenceResource:
		tr = hcl.Traversal{
			hcl.TraverseRoot{Name: r.res.Type()},
			hcl.TraverseAttr{Name: r.res.LocalName()},
		}
	case referenceDataResource:
		tr = hcl.Traversal{
			hcl.TraverseRoot{Name: "data"},
			hcl.TraverseAttr{Name: r.data.DataSource()},
			hcl.TraverseAttr{Name: r.data.LocalName()},
		}
	default:
		panic("unknown underlying type for reference")
	}

	tr = append(tr, r.tr...)

	return hclwrite.TokensForTraversal(tr)
}

// Append appends the given string to the reference
func (r Reference) Append(name string) Reference {
	cp := r.copy()
	cp.tr = append(cp.tr, hcl.TraverseAttr{Name: name})
	return cp
}

// index adds a reference to an index (of a list/set) to the reference
func (r Reference) index(i int) Reference {
	cp := r.copy()
	cp.tr = append(cp.tr, hcl.TraverseIndex{Key: cty.NumberIntVal(int64(i))})
	return cp
}

// key adds a reference to a key (of a map) to the reference
func (r Reference) key(k string) Reference {
	cp := r.copy()
	cp.tr = append(cp.tr, hcl.TraverseIndex{Key: cty.StringVal(k)})
	return cp
}

// splat adds the Terraform splat operator to a reference
func (r Reference) splat() Reference {
	cp := r.copy()
	cp.tr = append(cp.tr, hcl.TraverseSplat{})
	return cp
}

// copy makes a copy of the reference so that any slices can be safely passed
// around with modifying the original
func (r Reference) copy() Reference {
	tr := make(hcl.Traversal, len(r.tr))
	copy(tr, r.tr)
	return Reference{
		underlyingType: r.underlyingType,
		res:            r.res,
		data:           r.data,
		tr:             tr,
	}
}
