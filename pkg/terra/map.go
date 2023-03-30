// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"sort"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// MapString returns a map value containing the given string values
func MapString(value map[string]string) MapValue[StringValue] {
	ms := make(map[string]StringValue, len(value))
	for key, val := range value {
		ms[key] = String(val)
	}
	return Map(ms)
}

// Map returns a map value
func Map[T Value[T]](value map[string]T) MapValue[T] {
	return MapValue[T]{
		isInit: true,
		isRef:  false,
		values: value,
	}
}

// CastAsMap takes a value (as a reference) and wraps it in a MapValue
func CastAsMap[T Value[T]](value T) MapValue[T] {
	return ReferenceMap[T](value.InternalRef())
}

// ReferenceMap returns a map value
func ReferenceMap[T Value[T]](ref Reference) MapValue[T] {
	return MapValue[T]{
		isInit: true,
		isRef:  true,
		ref:    ref.copy(),
	}
}

var _ Value[MapValue[StringValue]] = (*MapValue[StringValue])(nil)

type MapValue[T Value[T]] struct {
	isInit bool
	isRef  bool
	ref    Reference

	values map[string]T
}

func (v MapValue[T]) Key(s string) T {
	if !v.isRef {
		panic("MapValue: cannot use Key on value")
	}
	var t T
	return t.InternalWithRef(v.ref.key(s))
}

func (v MapValue[T]) InternalTokens() hclwrite.Tokens {
	if !v.isInit {
		return nil
	}
	if v.isRef {
		return v.ref.InternalTokens()
	}

	elems := make([]hclwrite.ObjectAttrTokens, len(v.values))
	i := 0
	for _, key := range sortMapKeys(v.values) {
		elems[i] = hclwrite.ObjectAttrTokens{
			Name:  hclwrite.TokensForIdentifier("\"" + key + "\""),
			Value: v.values[key].InternalTokens(),
		}
		i++
	}
	return hclwrite.TokensForObject(elems)
}

func (v MapValue[T]) InternalRef() Reference {
	if !v.isRef {
		panic("MapValue: cannot get reference from value")
	}
	return v.ref
}

func (v MapValue[T]) InternalWithRef(ref Reference) MapValue[T] {
	return ReferenceMap[T](ref)
}

func sortMapKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
