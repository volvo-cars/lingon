// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terrajen

import (
	"github.com/dave/jennifer/jen"
)

func funcSchemaType(s *Schema, name string) *jen.Statement {
	return jen.Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id(name).Call().
		// Return type
		String().
		// Body
		Block(
			jen.Return(
				jen.Lit(s.Type),
			),
		)
}

func funcLocalName(s *Schema) *jen.Statement {
	return jen.Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id("LocalName").Call().
		// Return type
		String().
		// Body
		Block(
			jen.Return(
				jen.Id(s.Receiver).Dot(idFieldName),
			),
		)
}

func funcProviderLocalName(s *Schema) *jen.Statement {
	return jen.Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id("LocalName").Call().
		// Return type
		String().
		// Body
		Block(
			jen.Return(
				jen.Lit(s.ProviderName),
			),
		)
}

func funcProviderSource(s *Schema) *jen.Statement {
	return jen.Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id("Source").Call().
		// Return type
		String().
		// Body
		Block(
			jen.Return(
				jen.Lit(s.ProviderSource),
			),
		)
}

func funcProviderVersion(s *Schema) *jen.Statement {
	return jen.Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id("Version").Call().
		// Return type
		String().
		// Body
		Block(
			jen.Return(
				jen.Lit(s.ProviderVersion),
			),
		)
}

func funcConfiguration(s *Schema) *jen.Statement {
	return jen.Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id("Configuration").Call().
		// Return type
		Interface().
		// Body
		Block(
			jen.Return(
				jen.Id(s.Receiver).Dot(idFieldArgs),
			),
		)
}

func funcAttributes(s *Schema) *jen.Statement {
	return jen.Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id("Attributes").Call().
		//	Return
		Id(s.AttributesStructName).
		// Body
		Block(
			jen.Return(
				jen.Id(s.AttributesStructName).Values(
					jen.Dict{
						jen.Id("name"): jen.Id(s.Receiver).Dot(idFieldName),
					},
				),
			),
		)
}

func funcResourceImportState(s *Schema) *jen.Statement {
	return jen.Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id("ImportState").Call(jen.Id("av").Qual("io", "Reader")).
		// Return type
		Error().
		// Body
		Block(
			// Initialise the state
			jen.Id(s.Receiver).Dot(idFieldState).Op("=").Op("&").Id(s.StateStructName).Block(),
			jen.If(
				jen.Id("err").Op(":=").Qual(
					"encoding/json",
					"NewDecoder",
				).Call(jen.Id("av")).Dot(
					"Decode",
				).Call(jen.Id(s.Receiver).Dot(idFieldState)).Op(";").Id("err").Op("!=").Nil().Block(
					jen.Return(
						jen.Qual("fmt", "Errorf").Call(
							jen.Lit("decoding state into resource %s.%s: %w"),
							jen.Id(s.Receiver).Dot(idFuncType).Call(),
							jen.Id(s.Receiver).Dot(idFuncLocalName).Call(),
							jen.Id("err"),
						),
					),
				),
			),
			jen.Return(
				jen.Nil(),
			),
		)
}

func funcResourceState(s *Schema) *jen.Statement {
	return jen.Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id(idFuncState).Call().
		// Return type
		Params(jen.Op("*").Id(s.StateStructName), jen.Bool()).
		// Body
		Block(
			jen.Return(
				jen.Id(s.Receiver).Dot(idFieldState),
				jen.Id(s.Receiver).Dot(idFieldState).Op("!=").Nil(),
			),
		)
}

func funcResourceStateMust(s *Schema) *jen.Statement {
	return jen.Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id(idFuncStateMust).Call().
		// Return type
		Op("*").Id(s.StateStructName).
		// Body
		Block(
			jen.If(jen.Id(s.Receiver).Dot(idFieldState).Op("==").Nil()).Block(
				jen.Panic(
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("state is nil for resource %s.%s"),
						jen.Id(s.Receiver).Dot(idFuncType).Call(),
						jen.Id(s.Receiver).Dot(idFuncLocalName).Call(),
					),
				),
			),
			jen.Return(jen.Id(s.Receiver).Dot(idFieldState)),
		)
}

// funcDependOn, e.g.
//
//	func (irr *iamRoleResource) DependOn() terra.Value[terra.Reference] {
//		return terra.InternalRootRef("aws_iam_role", irr.Name)
//	}
func funcDependOn(s *Schema) *jen.Statement {
	return jen.Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id("DependOn").Call().
		// Return type
		Add(qualValue().Types(qualReferenceValue())).
		// Body
		Block(
			jen.Return(
				qualReferenceAttribute().Call(
					jen.Lit(s.Type), jen.Id(s.Receiver).Dot(idFieldName),
				),
			),
		)
}
