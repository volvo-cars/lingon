// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// ListString returns a list value containing the given string values
func ListString(values ...string) ListValue[StringValue] {
	ss := make([]StringValue, len(values))
	for i, s := range values {
		ss[i] = String(s)
	}
	return listValue[StringValue]{
		values: ss,
	}
}

// List returns a list value
func List[T Value[T]](values ...T) ListValue[T] {
	return listValue[T]{
		values: values,
	}
}

// AsListRef converts the given value to a list reference
func AsListRef[T Value[T]](value T) ListValue[T] {
	return listRef[T]{
		value: value,
	}
}

type ListValue[T Value[T]] interface {
	Value[ListValue[T]]
	Index(int) T
	Splat() T
}

var _ ListValue[StringValue] = (*listValue[StringValue])(nil)

type listValue[T Value[T]] struct {
	values []T
}

func (v listValue[T]) InternalCanTraverse() bool {
	return false
}

func (v listValue[T]) InternalTraverse(hcl.Traverser) ListValue[T] {
	panic("cannot traverse a list")
}

func (v listValue[T]) InternalTokens() hclwrite.Tokens {
	elems := make([]hclwrite.Tokens, len(v.values))

	for i, val := range v.values {
		elems[i] = val.InternalTokens()
	}
	return hclwrite.TokensForTuple(elems)
}

func (v listValue[T]) Index(i int) T {
	return v.values[i]
}

func (v listValue[T]) Splat() T {
	panic("cannot splat list of values")
}

var _ ListValue[StringValue] = (*listRef[StringValue])(nil)

type listRef[T Value[T]] struct {
	value T
}

func (r listRef[T]) InternalTraverse(tr hcl.Traverser) ListValue[T] {
	return listRef[T]{
		value: r.value.InternalTraverse(tr),
	}
}

func (r listRef[T]) InternalTokens() hclwrite.Tokens {
	return r.value.InternalTokens()
}

func (r listRef[T]) Index(i int) T {
	return r.value.InternalTraverse(hcl.TraverseIndex{Key: cty.NumberIntVal(int64(i))})
}

func (r listRef[T]) Splat() T {
	return r.value.InternalTraverse(hcl.TraverseSplat{})
}
