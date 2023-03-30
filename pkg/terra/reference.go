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
// It is never expected that a user should have to interact with this type
// directly.
//
// Value exposes some functions prefixed with "Internal". It is not possible
// to make these methods non-public, as generated code relies on this,
// but consumers of this package should **not** use these methods, please :)
type Value[T any] interface {
	tkihcl.Tokenizer
	// InternalWithRef returns a copy of this type T with the provided
	// reference. This is used internally to create a reference to attributes
	// within sets, lists,
	// maps and custom attributes types in the generated code.
	//
	// Internal: users should **not** use this!
	InternalWithRef(Reference) T
	// InternalRef returns the reference stored, if any.
	// If the Value T is not a reference, this method should panic to avoid
	// any nasty hidden errors (i.e.
	// silently converting a value to a reference).
	//
	// Internal: users should **not** use this!
	InternalRef() Reference
}

// ReferenceResource takes a resource and returns a Reference which is the
// address to that resource in the Terraform configuration.
func ReferenceResource(res Resource) Reference {
	return Reference{
		tr: hcl.Traversal{
			hcl.TraverseRoot{Name: res.Type()},
			hcl.TraverseAttr{Name: res.LocalName()},
		},
	}
}

// ReferenceDataResource takes a data resource and returns a Reference which
// is the address to that data resource in the Terraform configuration.
func ReferenceDataResource(data DataResource) Reference {
	return Reference{
		tr: hcl.Traversal{
			hcl.TraverseRoot{Name: "data"},
			hcl.TraverseAttr{Name: data.DataSource()},
			hcl.TraverseAttr{Name: data.LocalName()},
		},
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

// InternalTokens returns the tokens to represent this reference in Terraform
// configurations
func (r Reference) InternalTokens() hclwrite.Tokens {
	return hclwrite.TokensForTraversal(r.tr)
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
