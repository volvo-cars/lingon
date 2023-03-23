package terrajen

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/zclconf/go-cty/cty"
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
				jen.Id("err").Op(":=").Qual("encoding/json", "NewDecoder").Call(jen.Id("av")).Dot(
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
				jen.Qual(pkgTerra, idFuncReference).Call(
					jen.Lit(s.Type), jen.Id(s.Receiver).Dot(idFieldName),
				),
			),
		)
}

func hclTraverseAttr(stmt *jen.Statement) *jen.Statement {
	return qualHCLTraverseAttr().Values(jen.Dict{jen.Id("Name"): stmt})
}

// funcReferenceFromCtyType takes a cty Type and returns the function call to create
// a reference of that cty Type
func funcReferenceFromCtyType(ct cty.Type, stmt *jen.Statement) *jen.Statement {
	if ct.IsCollectionType() {
		switch {
		case ct.IsMapType():
			return qualAsMapRefFunc().Call(
				funcReferenceFromCtyType(
					ct.ElementType(),
					stmt,
				),
			)
		case ct.IsListType():
			return qualAsListRefFunc().Call(
				funcReferenceFromCtyType(
					ct.ElementType(),
					stmt,
				),
			)
		case ct.IsSetType():
			return qualAsSetRefFunc().Call(
				funcReferenceFromCtyType(
					ct.ElementType(),
					stmt,
				),
			)
		default:
			panic(
				fmt.Sprintf(
					"unsupported collection cty type: %s",
					ct.FriendlyName(),
				),
			)
		}
	}
	switch ct {
	case cty.String:
		return jen.Qual(
			pkgTerra,
			idFuncReference,
		).Call(stmt).Dot(idFuncReferenceAsString).Call()
	case cty.Number:
		return jen.Qual(
			pkgTerra,
			idFuncReference,
		).Call(stmt).Dot(idFuncReferenceAsNumber).Call()
	case cty.Bool:
		return jen.Qual(
			pkgTerra,
			idFuncReference,
		).Call(stmt).Dot(idFuncReferenceAsBool).Call()
	default:
		panic(fmt.Sprintf("unsupported simple cty type: %s", ct.FriendlyName()))
	}
}
