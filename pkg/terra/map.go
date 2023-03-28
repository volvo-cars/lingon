// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"sort"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Map returns a map value
func Map[T Value[T]](value map[string]T) MapValue[T] {
	return mapValue[T]{
		values: value,
	}
}

type MapValue[T Value[T]] interface {
	Value[MapValue[T]]
	Key(string) T
}

var _ MapValue[StringValue] = (*mapValue[StringValue])(nil)

type mapValue[T Value[T]] struct {
	values map[string]T
}

func (v mapValue[T]) InternalWithRef(Reference) MapValue[T] {
	panic("cannot traverse a map")
}

func (v mapValue[T]) InternalTokens() hclwrite.Tokens {
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

func (v mapValue[T]) Key(s string) T {
	return v.values[s]
}

var _ MapValue[StringValue] = (*mapRef[StringValue])(nil)

// ReferenceMap creates a map reference
func ReferenceMap[T Value[T]](ref Reference) MapValue[T] {
	return mapRef[T]{
		ref: ref.copy(),
	}
}

type mapRef[T Value[T]] struct {
	ref Reference
}

func (r mapRef[T]) InternalWithRef(ref Reference) MapValue[T] {
	return mapRef[T]{
		ref: ref.copy(),
	}
}

func (r mapRef[T]) InternalTokens() hclwrite.Tokens {
	return r.ref.InternalTokens()
}

func (r mapRef[T]) Key(s string) T {
	var v T
	return v.InternalWithRef(r.ref.key(s))
}

func sortMapKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
