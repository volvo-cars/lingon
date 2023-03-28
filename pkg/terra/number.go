// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"golang.org/x/exp/slog"
)

// Number returns a new number value
func Number(i int) NumberValue {
	return numberValue{
		value: cty.NumberIntVal(int64(i)),
	}
}

// NumberValue represents a number value
type NumberValue interface {
	Value[NumberValue]

	AsBool() BoolValue
	AsString() StringValue
}

var _ NumberValue = (*numberValue)(nil)

// numberValue is a concrete number, stored as a cty.Value
type numberValue struct {
	value cty.Value
}

func (v numberValue) AsBool() BoolValue {
	val, err := convert.Convert(v.value, cty.Bool)
	if err != nil {
		// TODO: handle error
		slog.Error("converting number to bool", err)
	}
	return boolValue{
		value: val,
	}
}

func (v numberValue) AsString() StringValue {
	val, err := convert.Convert(v.value, cty.String)
	if err != nil {
		// TODO: handle error
		slog.Error("converting number to string", err)
	}
	return stringValue{
		value: val,
	}
}

func (v numberValue) InternalWithRef(Reference) NumberValue {
	panic("cannot traverse a number")
}

func (v numberValue) InternalTokens() hclwrite.Tokens {
	return hclwrite.TokensForValue(v.value)
}

var _ NumberValue = (*numberRef)(nil)

// ReferenceNumber creates a number reference
func ReferenceNumber(ref Reference) NumberValue {
	return numberRef{
		ref: ref.copy(),
	}
}

// numberRef is a reference to a number in a Terraform configuration
type numberRef struct {
	ref Reference
}

func (r numberRef) InternalTokens() hclwrite.Tokens {
	return r.ref.InternalTokens()
}

func (r numberRef) InternalWithRef(ref Reference) NumberValue {
	return numberRef{
		ref: ref.copy(),
	}
}

func (r numberRef) AsBool() BoolValue {
	return boolRef(r)
}

func (r numberRef) AsString() StringValue {
	return stringRef(r)
}
