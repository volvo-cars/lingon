// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
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

// AsSetRef converts the given value to a set reference
func AsSetRef[T Value[T]](value T) SetValue[T] {
	return setRef[T]{
		value: value,
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

func (v setValue[T]) InternalTraverse(hcl.Traverser) SetValue[T] {
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

type setRef[T Value[T]] struct {
	value T
}

func (r setRef[T]) InternalTraverse(tr hcl.Traverser) SetValue[T] {
	return setRef[T]{
		value: r.value.InternalTraverse(tr),
	}
}

func (r setRef[T]) InternalTokens() hclwrite.Tokens {
	return r.value.InternalTokens()
}

func (r setRef[T]) Index(i int) T {
	return r.value.InternalTraverse(hcl.TraverseIndex{Key: cty.NumberIntVal(int64(i))})
}

func (r setRef[T]) Splat() T {
	return r.value.InternalTraverse(hcl.TraverseSplat{})
}
