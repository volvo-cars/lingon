// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

func Number(i int) NumberValue {
	return NumberValue{
		isInit: true,
		isRef:  false,
		value:  cty.NumberIntVal(int64(i)),
	}
}

func ReferenceAsNumber(ref Reference) NumberValue {
	return NumberValue{
		isInit: true,
		isRef:  true,
		ref:    ref,
	}
}

var _ Value[NumberValue] = (*NumberValue)(nil)

type NumberValue struct {
	isInit bool
	isRef  bool
	ref    Reference

	value cty.Value
}

func (v NumberValue) AsString() StringValue {
	if v.isRef {
		return ReferenceAsString(v.ref)
	}
	val, err := convert.Convert(v.value, cty.String)
	if err != nil {
		panic(fmt.Sprintf("converting number to string: %s", err.Error()))
	}
	return StringValue{
		isInit: true,
		value:  val,
	}
}

func (v NumberValue) AsBool() BoolValue {
	if v.isRef {
		return ReferenceAsBool(v.ref)
	}
	val, err := convert.Convert(v.value, cty.Bool)
	if err != nil {
		panic(fmt.Sprintf("converting number to bool: %s", err.Error()))
	}
	return BoolValue{
		isInit: true,
		value:  val,
	}
}

func (v NumberValue) InternalTokens() hclwrite.Tokens {
	if !v.isInit {
		return nil
	}
	if v.isRef {
		return v.ref.InternalTokens()
	}
	return hclwrite.TokensForValue(v.value)
}

func (v NumberValue) InternalRef() Reference {
	if !v.isRef {
		panic("NumberValue: cannot use value as reference")
	}
	return v.ref.copy()
}

func (v NumberValue) InternalWithRef(ref Reference) NumberValue {
	return ReferenceAsNumber(ref)
}
