// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	tkihcl "github.com/volvo-cars/lingon/pkg/internal/hcl"
	"github.com/zclconf/go-cty/cty"
)

type Referencer interface {
	// InternalRef returns a copy of the reference stored, if any.
	// If the Value T is not a reference, this method should return an error to
	// avoid any nasty hidden errors (i.e.
	// silently converting a value to a reference).
	//
	// Internal: users should **not** use this!
	InternalRef() (Reference, error)
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
// nested objects or even the splat operator (i.e. [*]).
//
// A reference can be created by passing a [Resource] to
// [ReferenceResource] or passing a [DataResource] to
// [ReferenceDataResource].
type Reference struct {
	underlyingType referenceUnderlyingType
	res            Resource
	data           DataResource

	steps []referenceStep
}

// referenceStepType represents the type of step in a reference.
type referenceStepType int

const (
	referenceStepAttribute = 1
	referenceStepIndex     = 2
	referenceStepKey       = 3
	referenceStepSplat     = 4
)

// referenceStep is a part of a reference to some value in a Terraform
// configuration. A [Reference] might refer to a field of an attribute,
// and that reference is made up of multiple referenceStep.
type referenceStep struct {
	stepType  referenceStepType
	attribute string
	index     int
	key       string
}

// InternalTokens returns the tokens to represent this reference in Terraform
// configurations
func (r Reference) InternalTokens() hclwrite.Tokens {
	var fullSteps []referenceStep
	switch r.underlyingType {
	case referenceResource:
		fullSteps = []referenceStep{
			{
				stepType:  referenceStepAttribute,
				attribute: r.res.Type(),
			},
			{
				stepType:  referenceStepAttribute,
				attribute: r.res.LocalName(),
			},
		}
	case referenceDataResource:
		fullSteps = []referenceStep{
			{
				stepType:  referenceStepAttribute,
				attribute: "data",
			},
			{
				stepType:  referenceStepAttribute,
				attribute: r.data.DataSource(),
			},
			{
				stepType:  referenceStepAttribute,
				attribute: r.data.LocalName(),
			},
		}
	default:
		panic("unknown underlying type for reference")
	}
	fullSteps = append(fullSteps, r.steps...)
	tokens, err := tokensForSteps(fullSteps)
	if err != nil {
		panic(fmt.Sprintf("creating tokens for reference steps: %s", err))
	}
	return tokens
}

// Append appends the given string to the reference
func (r Reference) Append(name string) Reference {
	cp := r.copy()
	cp.steps = append(
		cp.steps, referenceStep{
			stepType:  referenceStepAttribute,
			attribute: name,
		},
	)
	return cp
}

// index adds a reference to an index (of a list/set) to the reference
func (r Reference) index(i int) Reference {
	cp := r.copy()
	cp.steps = append(
		cp.steps, referenceStep{
			stepType: referenceStepIndex,
			index:    i,
		},
	)
	return cp
}

// key adds a reference to a key (of a map) to the reference
func (r Reference) key(k string) Reference {
	cp := r.copy()
	cp.steps = append(
		cp.steps, referenceStep{
			stepType: referenceStepKey,
			key:      k,
		},
	)
	return cp
}

// splat adds the Terraform splat operator to a reference
func (r Reference) splat() Reference {
	cp := r.copy()
	cp.steps = append(
		cp.steps, referenceStep{
			stepType: referenceStepSplat,
		},
	)
	return cp
}

// copy makes a copy of the reference so that any slices can be safely passed
// around with modifying the original
func (r Reference) copy() Reference {
	steps := make([]referenceStep, len(r.steps))
	copy(steps, r.steps)
	return Reference{
		underlyingType: r.underlyingType,
		res:            r.res,
		data:           r.data,
		steps:          steps,
	}
}

func tokensForSteps(steps []referenceStep) (hclwrite.Tokens, error) {
	var tokens hclwrite.Tokens
	for i, step := range steps {
		switch step.stepType {
		case referenceStepAttribute:
			// If not the first step, add the "." separator
			if i > 0 {
				tokens = append(
					tokens,
					&hclwrite.Token{
						Type:  hclsyntax.TokenDot,
						Bytes: []byte{'.'},
					},
				)
			}
			tokens = append(
				tokens,
				hclwrite.TokensForIdentifier(step.attribute)...,
			)
		case referenceStepIndex:
			tokens = append(
				tokens,
				&hclwrite.Token{
					Type:  hclsyntax.TokenOBrack,
					Bytes: []byte{'['},
				},
			)
			tokens = append(
				tokens,
				hclwrite.TokensForValue(
					cty.NumberIntVal(int64(step.index)),
				)...,
			)
			tokens = append(
				tokens,
				&hclwrite.Token{
					Type:  hclsyntax.TokenCBrack,
					Bytes: []byte{']'},
				},
			)
		case referenceStepKey:
			tokens = append(
				tokens,
				&hclwrite.Token{
					Type:  hclsyntax.TokenOBrack,
					Bytes: []byte{'['},
				},
			)
			tokens = append(
				tokens,
				hclwrite.TokensForValue(
					cty.StringVal(step.key),
				)...,
			)
			tokens = append(
				tokens, &hclwrite.Token{
					Type:  hclsyntax.TokenCBrack,
					Bytes: []byte{']'},
				},
			)
		case referenceStepSplat:
			tokens = append(
				tokens,
				&hclwrite.Token{
					Type:  hclsyntax.TokenOBrack,
					Bytes: []byte{'['},
				},
				&hclwrite.Token{
					Type:  hclsyntax.TokenStar,
					Bytes: []byte{'*'},
				},
				&hclwrite.Token{
					Type:  hclsyntax.TokenCBrack,
					Bytes: []byte{']'},
				},
			)
		default:
			return nil, fmt.Errorf(
				"unknown reference step type: %d",
				step.stepType,
			)
		}
	}
	return tokens, nil
}
