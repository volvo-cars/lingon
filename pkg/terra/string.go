// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"golang.org/x/exp/slog"
)

func String(s string) StringValue {
	return stringValue{
		value: cty.StringVal(s),
	}
}

type StringValue interface {
	Value[StringValue]

	AsNumber() NumberValue
	AsBool() BoolValue
}

var _ StringValue = (*stringValue)(nil)

type stringValue struct {
	value cty.Value
}

func (v stringValue) AsBool() BoolValue {
	val, err := convert.Convert(v.value, cty.Bool)
	if err != nil {
		// TODO: handle error
		slog.Error("converting number to bool", err)
	}
	return boolValue{
		value: val,
	}
}

func (v stringValue) AsNumber() NumberValue {
	val, err := convert.Convert(v.value, cty.Number)
	if err != nil {
		// TODO: handle error
		slog.Error("converting number to bool", err)
	}
	return numberValue{
		value: val,
	}
}

func (v stringValue) InternalTraverse(hcl.Traverser) StringValue {
	panic("cannot traverse a string")
}

func (v stringValue) InternalTokens() hclwrite.Tokens {
	return hclwrite.TokensForValue(v.value)
}

var _ StringValue = (*stringRef)(nil)

type stringRef struct {
	ref ReferenceValue
}

func (r stringRef) InternalTokens() hclwrite.Tokens {
	return r.ref.InternalTokens()
}

func (r stringRef) InternalTraverse(step hcl.Traverser) StringValue {
	return stringRef{
		ref: r.ref.InternalTraverse(step),
	}
}

func (r stringRef) AsBool() BoolValue {
	return boolRef(r)
}

func (r stringRef) AsNumber() NumberValue {
	return numberRef(r)
}
