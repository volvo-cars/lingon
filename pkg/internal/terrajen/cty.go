// Copyright (c) 2023 Volvo Car Corporation
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

func ctyTypeReturnValue(n *node, ct cty.Type, name string) *jen.Statement {
	// E.g.
	// i.InternalTraverse(hcl.TraverseAttr{Name: "issuer"}).Reference
	param := jen.Id(n.receiver).Dot(idStructReferenceValue).Dot(idFuncInternalTraverse).Call(hclTraverseAttr(jen.Lit(name)))
	switch {
	case ct == cty.Bool:
		return param.Dot(idFuncReferenceAsBool).Call()
	case ct == cty.String:
		return param.Dot(idFuncReferenceAsString).Call()
	case ct == cty.Number:
		return param.Dot(idFuncReferenceAsNumber).Call()
	case ct.IsMapType():
		return qualAsMapRefFunc().Call(
			ctyTypeReturnValue(
				n,
				ct.ElementType(),
				name,
			),
		)
	case ct.IsSetType():
		return qualAsSetRefFunc().Call(
			ctyTypeReturnValue(
				n,
				ct.ElementType(),
				name,
			),
		)
	case ct.IsListType():
		return qualAsListRefFunc().Call(
			ctyTypeReturnValue(
				n,
				ct.ElementType(),
				name,
			),
		)
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
