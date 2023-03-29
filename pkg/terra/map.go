// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// Map returns a map value
func Map[T Value[T]](value map[string]T) MapValue[T] {
	return mapValue[T]{
		values: value,
	}
}

// AsMapRef converts the given value to a map reference
func AsMapRef[T Value[T]](value T) MapValue[T] {
	return mapRef[T]{
		value: value,
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

func (v mapValue[T]) InternalTraverse(hcl.Traverser) MapValue[T] {
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

type mapRef[T Value[T]] struct {
	value T
}

func (r mapRef[T]) InternalTraverse(tr hcl.Traverser) MapValue[T] {
	return mapRef[T]{
		value: r.value.InternalTraverse(tr),
	}
}

func (r mapRef[T]) InternalTokens() hclwrite.Tokens {
	return r.value.InternalTokens()
}

func (r mapRef[T]) Key(s string) T {
	return r.value.InternalTraverse(hcl.TraverseIndex{Key: cty.StringVal(s)})
}

func sortMapKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
