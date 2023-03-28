// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terrajen

import (
	"fmt"

	"github.com/volvo-cars/lingon/pkg/internal/str"

	"github.com/dave/jennifer/jen"
	"github.com/zclconf/go-cty/cty"
)

func ctyTypeReturnType(ct cty.Type) *jen.Statement {
	switch {
	case ct == cty.Bool:
		return qualBoolValue()
	case ct == cty.String:
		return qualStringValue()
	case ct == cty.Number:
		return qualNumberValue()
	case ct.IsMapType():
		return qualMapValue().Types(ctyTypeReturnType(ct.ElementType()))
	case ct.IsSetType():
		return qualSetValue().Types(ctyTypeReturnType(ct.ElementType()))
	case ct.IsListType():
		return qualListValue().Types(ctyTypeReturnType(ct.ElementType()))
	default:
		panic(fmt.Sprintf("unsupported AttributeType: %s", ct.GoString()))
	}
}

func funcReferenceByCtyType(ct cty.Type) *jen.Statement {
	switch {
	case ct == cty.Bool:
		return qualReferenceBool()
	case ct == cty.String:
		return qualReferenceString()
	case ct == cty.Number:
		return qualReferenceNumber()
	case ct.IsMapType():
		subType := ctyTypeReturnType(ct.ElementType())
		return qualReferenceMap().Types(subType)
	case ct.IsSetType():
		subType := ctyTypeReturnType(ct.ElementType())
		return qualReferenceSet().Types(subType)
	case ct.IsListType():
		subType := ctyTypeReturnType(ct.ElementType())
		return qualReferenceList().Types(subType)
	default:
		panic(fmt.Sprintf("unsupported AttributeType: %s", ct.GoString()))
	}
}

func ctyTypeToGoType(
	t cty.Type,
	attrName string,
) jen.Code {
	if t.IsObjectType() {
		return jen.StructFunc(
			func(g *jen.Group) {
				for k, v := range t.AttributeTypes() {
					g.Id(str.PascalCase(k)).Add(
						ctyTypeToGoType(
							v,
							str.PascalCase(attrName+k),
						),
					).Tag(
						map[string]string{
							tagJSON: k,
						},
					)
				}
			},
		)
	}
	if t.IsListType() {
		if et := t.ListElementType(); et != nil {
			c := ctyTypeToGoType(*et, attrName)
			return jen.Index().Add(c)
		}
		panic("unsupported list type")
	}
	if t.IsMapType() {
		if et := t.MapElementType(); et != nil {
			c := ctyTypeToGoType(*et, attrName)
			return jen.Map(jen.String()).Add(c)
		}
		panic("unsupported map type")
	}
	if t.IsSetType() {
		if et := t.SetElementType(); et != nil {
			c := ctyTypeToGoType(*et, attrName)
			return jen.Index().Add(c)
		}
		panic("unsupported set type")
	}
	switch t {
	case cty.String:
		return jen.String()
	case cty.Bool:
		return jen.Bool()
	case cty.Number:
		return jen.Float64()
	default:
		panic(fmt.Sprintf("unsupported type: %s", t.FriendlyName()))
	}
}
