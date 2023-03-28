// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// SetString returns a set value containing the given string values
func SetString(values ...string) SetValue[StringValue] {
	ss := make([]StringValue, len(values))
	for i, s := range values {
		ss[i] = String(s)
	}
	return setValue[StringValue]{
		values: ss,
	}
}

// Set returns a set value
func Set[T Value[T]](values ...T) SetValue[T] {
	return setValue[T]{
		values: values,
	}
}

type SetValue[T Value[T]] interface {
	Value[SetValue[T]]
	Index(int) T
	Splat() T
}

var _ SetValue[StringValue] = (*setValue[StringValue])(nil)

type setValue[T Value[T]] struct {
	values []T
}

func (v setValue[T]) InternalWithRef(Reference) SetValue[T] {
	panic("cannot traverse a set")
}

func (v setValue[T]) InternalTokens() hclwrite.Tokens {
	elems := make([]hclwrite.Tokens, len(v.values))

	for i, val := range v.values {
		elems[i] = val.InternalTokens()
	}
	return hclwrite.TokensForTuple(elems)
}

func (v setValue[T]) Index(i int) T {
	return v.values[i]
}

func (v setValue[T]) Splat() T {
	panic("cannot splat set of values")
}

var _ SetValue[StringValue] = (*setRef[StringValue])(nil)

// ReferenceSet creates a set reference
func ReferenceSet[T Value[T]](ref Reference) SetValue[T] {
	return setRef[T]{
		ref: ref.copy(),
	}
}

type setRef[T Value[T]] struct {
	ref Reference
}

func (r setRef[T]) InternalWithRef(ref Reference) SetValue[T] {
	return setRef[T]{
		ref: ref.copy(),
	}
}

func (r setRef[T]) InternalTokens() hclwrite.Tokens {
	return r.ref.InternalTokens()
}

func (r setRef[T]) Index(i int) T {
	var v T
	return v.InternalWithRef(r.ref.index(i))
}

func (r setRef[T]) Splat() T {
	var v T
	return v.InternalWithRef(r.ref.splat())
}
