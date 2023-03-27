// Copyright (c) Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"golang.org/x/exp/slog"
)

// Bool returns a bool value
func Bool(b bool) BoolValue {
	return boolValue{
		value: cty.BoolVal(b),
	}
}

// BoolValue represents a bool value
type BoolValue interface {
	Value[BoolValue]

	AsString() StringValue
}

var _ BoolValue = (*boolValue)(nil)

type boolValue struct {
	value cty.Value
}

// AsString tries to convert a BoolValue to a StringValue
func (v boolValue) AsString() StringValue {
	val, err := convert.Convert(v.value, cty.String)
	if err != nil {
		// TODO: handle error
		slog.Error("converting number to bool", err)
	}
	return stringValue{
		value: val,
	}
}

func (v boolValue) InternalTraverse(hcl.Traverser) BoolValue {
	panic("cannot traverse a boolean")
}

func (v boolValue) InternalTokens() hclwrite.Tokens {
	return hclwrite.TokensForValue(v.value)
}

var _ BoolValue = (*boolRef)(nil)

type boolRef struct {
	ref ReferenceValue
}

func (r boolRef) InternalTokens() hclwrite.Tokens {
	return r.ref.InternalTokens()
}

func (r boolRef) InternalTraverse(step hcl.Traverser) BoolValue {
	return boolRef{
		ref: r.ref.InternalTraverse(step),
	}
}

// AsString converts a reference to a BoolValue to a StringValue reference
func (r boolRef) AsString() StringValue {
	return stringRef(r)
}
