// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// SetString returns a list value containing the given string values
func SetString(values ...string) SetValue[StringValue] {
	ss := make([]StringValue, len(values))
	for i, s := range values {
		ss[i] = String(s)
	}
	return Set[StringValue](ss...)
}

// Set returns a list value
func Set[T Value[T]](values ...T) SetValue[T] {
	return SetValue[T]{
		isInit: true,
		isRef:  false,
		values: values,
	}
}

// CastAsSet takes a value (as a reference) and wraps it in a SetValue
func CastAsSet[T Value[T]](value T) SetValue[T] {
	return ReferenceSet[T](value.InternalRef())
}

// ReferenceSet creates a list reference
func ReferenceSet[T Value[T]](ref Reference) SetValue[T] {
	return SetValue[T]{
		isInit: true,
		isRef:  true,
		ref:    ref.copy(),
	}
}

var _ Value[SetValue[StringValue]] = (*SetValue[StringValue])(nil)

type SetValue[T Value[T]] struct {
	isInit bool
	isRef  bool
	ref    Reference

	values []T
}

func (v SetValue[T]) Index(i int) T {
	if !v.isRef {
		panic("SetValue: cannot use Index on value")
	}
	var t T
	return t.InternalWithRef(v.ref.index(i))
}

func (v SetValue[T]) Splat() T {
	if !v.isRef {
		panic("SetValue: cannot use Splat on value")
	}
	var t T
	return t.InternalWithRef(v.ref.splat())
}

func (v SetValue[T]) InternalTokens() hclwrite.Tokens {
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
	if v.values == nil {
		return nil
	}
	return hclwrite.TokensForTuple(elems)
}

func (v SetValue[T]) InternalRef() Reference {
	if !v.isRef {
		panic("SetValue: cannot get reference from value")
	}
	return v.ref
}

func (v SetValue[T]) InternalWithRef(ref Reference) SetValue[T] {
	return ReferenceSet[T](ref)
}
