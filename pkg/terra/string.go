// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

func String(s string) StringValue {
	return StringValue{
		isInit: true,
		isRef:  false,
		value:  cty.StringVal(s),
	}
}

func ReferenceAsString(ref Reference) StringValue {
	return StringValue{
		isInit: true,
		isRef:  true,
		ref:    ref,
	}
}

var _ Value[StringValue] = (*StringValue)(nil)

type StringValue struct {
	isInit bool
	isRef  bool
	ref    Reference

	value cty.Value
}

func (v StringValue) AsBool() BoolValue {
	if v.isRef {
		return ReferenceAsBool(v.ref)
	}
	val, err := convert.Convert(v.value, cty.Bool)
	if err != nil {
		panic(fmt.Sprintf("converting string to bool: %s", err.Error()))
	}
	return BoolValue{
		value: val,
	}
}

func (v StringValue) AsNumber() NumberValue {
	if v.isRef {
		return ReferenceAsNumber(v.ref)
	}
	val, err := convert.Convert(v.value, cty.Number)
	if err != nil {
		panic(fmt.Sprintf("converting string to bool: %s", err.Error()))
	}
	return NumberValue{
		isInit: true,
		value:  val,
	}
}

func (v StringValue) InternalTokens() (hclwrite.Tokens, error) {
	if !v.isInit {
		return nil, nil
	}
	if v.isRef {
		return v.ref.InternalTokens()
	}
	return hclwrite.TokensForValue(v.value), nil
}

func (v StringValue) InternalRef() (Reference, error) {
	if !v.isRef {
		return Reference{},
			errors.New("StringValue: cannot use value as reference")
	}
	return v.ref.copy(), nil
}

func (v StringValue) InternalWithRef(ref Reference) StringValue {
	return ReferenceAsString(ref)
}
