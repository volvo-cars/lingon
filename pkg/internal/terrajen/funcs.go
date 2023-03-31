// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terrajen

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

func funcSchemaType(s *Schema, name string) *jen.Statement {
	return jen.Comment(
		fmt.Sprintf(
			"%s returns the Terraform object type for [%s].",
			name,
			s.StructName,
		),
	).
		Line().
		Func().
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
	return jen.Comment(
		fmt.Sprintf(
			"%s returns the local name for [%s].",
			idFuncLocalName,
			s.StructName,
		),
	).
		Line().
		Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id(idFuncLocalName).Call().
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
	return jen.Comment(
		fmt.Sprintf(
			"%s returns the provider local name for [%s].",
			idFuncLocalName,
			s.StructName,
		),
	).
		Line().
		Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id(idFuncLocalName).Call().
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
	return jen.Comment(
		fmt.Sprintf(
			"%s returns the provider source for [%s].",
			idFuncSource,
			s.StructName,
		),
	).
		Line().
		Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id(idFuncSource).Call().
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
	return jen.Comment(
		fmt.Sprintf(
			"%s returns the provider version for [%s].",
			idFuncVersion,
			s.StructName,
		),
	).
		Line().
		Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id(idFuncVersion).Call().
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
	return jen.Comment(
		fmt.Sprintf(
			"%s returns the configuration (args) for [%s].",
			idFuncConfiguration,
			s.StructName,
		),
	).
		Line().
		Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id(idFuncConfiguration).Call().
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
	var createRefFunc *jen.Statement
	if s.SchemaType == SchemaTypeResource {
		createRefFunc = qualReferenceResource()
	} else {
		createRefFunc = qualReferenceDataResource()
	}
	return jen.Comment(
		fmt.Sprintf(
			"%s returns the attributes for [%s].",
			idFuncAttributes,
			s.StructName,
		),
	).
		Line().
		Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id(idFuncAttributes).Call().
		//	Return
		Id(s.AttributesStructName).
		// Body
		Block(
			jen.Return(
				jen.Id(s.AttributesStructName).Values(
					jen.Dict{
						jen.Id("ref"): createRefFunc.Call(jen.Id(s.Receiver)),
					},
				),
			),
		)
}

func funcResourceImportState(s *Schema) *jen.Statement {
	attributesArgs := jen.Id("av").Clone
	return jen.Comment(
		fmt.Sprintf(
			"%s imports the given attribute values into [%s]'s state.",
			idFuncImportState,
			s.StructName,
		),
	).
		Line().
		Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id(idFuncImportState).Call(attributesArgs().Qual("io", "Reader")).
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
				).Call(attributesArgs()).Dot(
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
	return jen.Comment(
		fmt.Sprintf(
			"%s returns the state and a bool indicating if [%s] has state.",
			idFuncState,
			s.StructName,
		),
	).
		Line().
		Func().
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
	return jen.Comment(
		fmt.Sprintf(
			"%s returns the state for [%s]. Panics if the state is nil.",
			idFuncStateMust,
			s.StructName,
		),
	).
		Line().
		Func().
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
	return jen.Comment(
		fmt.Sprintf(
			"%s is used for other resources to depend on [%s].",
			idFuncDependOn,
			s.StructName,
		),
	).
		Line().
		Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id(idFuncDependOn).Call().
		// Return type
		Add(qualReferenceValue()).
		// Body
		Block(
			jen.Return(
				qualReferenceResource().Call(jen.Id(s.Receiver)),
			),
		)
}

// funcDependencies, e.g.
//
//	func (irr *iamRoleResource) Dependencies() terra.Dependencies {
//		return irr.DependsOn
//	}
func funcDependencies(s *Schema) *jen.Statement {
	return jen.Comment(
		fmt.Sprintf(
			"%s returns the list of resources [%s] depends_on.",
			idFuncDependencies,
			s.StructName,
		),
	).
		Line().
		Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id(idFuncDependencies).Call().
		// Return type
		Add(qualTypeDependencies()).
		// Body
		Block(
			jen.Return(
				jen.Id(s.Receiver).Dot(idFieldDependsOn),
			),
		)
}

// funcLifecycleManagement, e.g.
//
//	func (irr *iamRoleResource) LifecycleManagement() *terra.Lifecycle {
//		return irr.Lifecycle
//	}
func funcLifecycleManagement(s *Schema) *jen.Statement {
	return jen.Comment(
		fmt.Sprintf(
			"%s returns the lifecycle block for [%s].",
			idFuncLifecycleManagement,
			s.StructName,
		),
	).
		Line().
		Func().
		// Receiver
		Params(jen.Id(s.Receiver).Op("*").Id(s.StructName)).
		// Name
		Id(idFuncLifecycleManagement).Call().
		// Return type
		Op("*").Add(qualStructLifecycle()).
		// Body
		Block(
			jen.Return(
				jen.Id(s.Receiver).Dot(idFieldLifecycle),
			),
		)
}
