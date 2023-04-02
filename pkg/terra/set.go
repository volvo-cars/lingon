// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"errors"
	"fmt"

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
	ref, err := value.InternalRef()
	if err != nil {
		panic(
			fmt.Sprintf(
				"CastAsSet: getting internal reference: %s",
				err.Error(),
			),
		)
	}
	return ReferenceAsSet[T](ref)
}

// ReferenceAsSet creates a list reference
func ReferenceAsSet[T Value[T]](ref Reference) SetValue[T] {
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

func (v SetValue[T]) InternalTokens() (hclwrite.Tokens, error) {
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
	if v.values == nil {
		return nil, nil
	}
	return hclwrite.TokensForTuple(elems), nil
}

func (v SetValue[T]) InternalRef() (Reference, error) {
	if !v.isRef {
		return Reference{},
			errors.New("SetValue: cannot get reference from value")
	}
	return v.ref.copy(), nil
}

func (v SetValue[T]) InternalWithRef(ref Reference) SetValue[T] {
	return ReferenceAsSet[T](ref)
}
