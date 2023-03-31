// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

func Bool(b bool) BoolValue {
	return BoolValue{
		isInit: true,
		isRef:  false,
		value:  cty.BoolVal(b),
	}
}

func ReferenceAsBool(ref Reference) BoolValue {
	return BoolValue{
		isInit: true,
		isRef:  true,
		ref:    ref,
	}
}

var _ Value[BoolValue] = (*BoolValue)(nil)

type BoolValue struct {
	isInit bool
	isRef  bool
	ref    Reference

	value cty.Value
}

func (v BoolValue) AsString() StringValue {
	if v.isRef {
		return ReferenceAsString(v.ref)
	}
	val, err := convert.Convert(v.value, cty.String)
	if err != nil {
		panic(fmt.Sprintf("converting bool to string: %s", err.Error()))
	}
	return StringValue{
		isInit: true,
		value:  val,
	}
}

func (v BoolValue) AsNumber() NumberValue {
	if v.isRef {
		return ReferenceAsNumber(v.ref)
	}
	val, err := convert.Convert(v.value, cty.Number)
	if err != nil {
		panic(fmt.Sprintf("converting bool to number: %s", err.Error()))
	}
	return NumberValue{
		isInit: true,
		value:  val,
	}
}

func (v BoolValue) InternalTokens() hclwrite.Tokens {
	if !v.isInit {
		return nil
	}
	if v.isRef {
		return v.ref.InternalTokens()
	}
	return hclwrite.TokensForValue(v.value)
}

func (v BoolValue) InternalRef() Reference {
	if !v.isRef {
		panic("BoolValue: cannot use value as reference")
	}
	return v.ref.copy()
}

func (v BoolValue) InternalWithRef(ref Reference) BoolValue {
	return ReferenceAsBool(ref)
}
