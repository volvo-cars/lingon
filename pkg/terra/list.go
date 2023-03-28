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

func (v listValue[T]) InternalWithRef(Reference) ListValue[T] {
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

// ReferenceList creates a list reference
func ReferenceList[T Value[T]](ref Reference) ListValue[T] {
	return listRef[T]{
		ref: ref,
	}
}

type listRef[T Value[T]] struct {
	ref Reference
}

func (r listRef[T]) InternalWithRef(ref Reference) ListValue[T] {
	return listRef[T]{
		ref: ref.copy(),
	}
}

func (r listRef[T]) InternalTokens() hclwrite.Tokens {
	return r.ref.InternalTokens()
}

func (r listRef[T]) Index(i int) T {
	var v T
	return v.InternalWithRef(r.ref.index(i))
}

func (r listRef[T]) Splat() T {
	var v T
	return v.InternalWithRef(r.ref.splat())
}
