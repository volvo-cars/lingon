// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"errors"
	"fmt"

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
	ref, err := value.InternalRef()
	if err != nil {
		panic(
			fmt.Sprintf(
				"CastAsList: getting internal reference: %s",
				err.Error(),
			),
		)
	}
	return ReferenceAsList[T](ref)
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

func (v ListValue[T]) InternalTokens() (hclwrite.Tokens, error) {
	if !v.isInit {
		return nil, nil
	}
	if v.isRef {
		return v.ref.InternalTokens()
	}

	elems := make([]hclwrite.Tokens, len(v.values))
	for i, val := range v.values {
		toks, err := val.InternalTokens()
		if err != nil {
			return nil, fmt.Errorf("getting tokens: %w", err)
		}
		elems[i] = toks
	}
	return hclwrite.TokensForTuple(elems), nil
}

func (v ListValue[T]) InternalRef() (Reference, error) {
	if !v.isRef {
		return Reference{},
			errors.New("ListValue: cannot get reference from value")
	}
	return v.ref.copy(), nil
}

func (v ListValue[T]) InternalWithRef(ref Reference) ListValue[T] {
	return ReferenceAsList[T](ref)
}
