// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// ListString returns a list value containing the given string values
func ListString(values ...string) ListValue[StringValue] {
	ss := make([]StringValue, len(values))
	for i, s := range values {
		ss[i] = String(s)
	}
	return List[StringValue](ss...)
}

// List returns a list value
func List[T Value[T]](values ...T) ListValue[T] {
	return ListValue[T]{
		isInit: true,
		isRef:  false,
		values: values,
	}
}

// CastAsList takes a value (as a reference) and wraps it in a ListValue
func CastAsList[T Value[T]](value T) ListValue[T] {
	return ReferenceAsList[T](value.InternalRef())
}

// ReferenceAsList creates a list reference
func ReferenceAsList[T Value[T]](ref Reference) ListValue[T] {
	return ListValue[T]{
		isInit: true,
		isRef:  true,
		ref:    ref.copy(),
	}
}

var _ Value[ListValue[StringValue]] = (*ListValue[StringValue])(nil)

type ListValue[T Value[T]] struct {
	isInit bool
	isRef  bool
	ref    Reference

	values []T
}

func (v ListValue[T]) Index(i int) T {
	if !v.isRef {
		panic("ListValue: cannot use Index on value")
	}
	var t T
	return t.InternalWithRef(v.ref.index(i))
}

func (v ListValue[T]) Splat() T {
	if !v.isRef {
		panic("ListValue: cannot use Splat on value")
	}
	var t T
	return t.InternalWithRef(v.ref.splat())
}

func (v ListValue[T]) InternalTokens() hclwrite.Tokens {
	if !v.isInit {
		return nil
	}
	if v.isRef {
		return v.ref.InternalTokens()
	}

	elems := make([]hclwrite.Tokens, len(v.values))
	for i, val := range v.values {
		elems[i] = val.InternalTokens()
	}
	return hclwrite.TokensForTuple(elems)
}

func (v ListValue[T]) InternalRef() Reference {
	if !v.isRef {
		panic("ListValue: cannot get reference from value")
	}
	return v.ref.copy()
}

func (v ListValue[T]) InternalWithRef(ref Reference) ListValue[T] {
	return ReferenceAsList[T](ref)
}
