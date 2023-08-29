// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terrajen

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/veggiemonk/strcase"
)

// argsStruct takes a schema and generates the Args struct that is used by the user to specify the arguments
// for the object that the schema represents (e.g. provider, resource, data resource)
func argsStruct(s *Schema) *jen.Statement {
	fields := make([]jen.Code, 0)
	for _, attr := range s.graph.attributes {
		if !attr.isArg {
			continue
		}
		stmt := jen.Comment(attr.comment()).Line()
		stmt.Add(jen.Id(strcase.Pascal(attr.name)))
		stmt.Add(ctyTypeReturnType(attr.ctyType))

		// Add tags
		tags := map[string]string{
			tagHCL: attr.name + ",attr",
		}
		if attr.isRequired {
			tags[tagValidate] = "required"
		}
		stmt.Tag(tags)
		fields = append(fields, stmt)
	}

	for _, child := range s.graph.children {
		tags := map[string]string{
			tagHCL: child.uniqueName + ",block",
		}
		stmt := jen.Comment(child.comment()).Line()
		stmt.Add(jen.Id(strcase.Pascal(child.uniqueName)))
		if len(child.nestingPath) == 0 || child.maxItems == 1 {
			stmt.Op("*")
			if child.isRequired {
				tags[tagValidate] = "required"
			}
		} else {
			for range child.nestingPath {
				stmt.Index()
			}
			tags[tagValidate] = nodeBlockListValidateTags(child)
		}
		stmt.Qual(s.SubPkgQualPath(), strcase.Pascal(child.uniqueName))
		stmt.Tag(tags)
		fields = append(fields, stmt)
	}

	return jen.Comment(
		fmt.Sprintf(
			"%s contains the configurations for %s.",
			s.ArgumentStructName,
			s.Type,
		),
	).
		Line().
		Type().
		Id(s.ArgumentStructName).
		Struct(fields...)
}

// attributesStruct takes a schema and generates the Attributes struct that is used by the user to creates references to
// attributes for the object that the schema represents (e.g. provider, resource, data resource)
func attributesStruct(s *Schema) *jen.Statement {
	var stmt jen.Statement

	// Attribute struct will have a field called "ref" containing a [terra.Reference]
	structFieldRefName := "ref"
	structFieldRef := jen.Id(s.Receiver).Dot(structFieldRefName).Clone

	attrStruct := jen.Type().Id(s.AttributesStructName).Struct(
		jen.Id("ref").Add(qualReferenceValue()),
	)

	stmt.Add(attrStruct)
	stmt.Line()
	stmt.Line()

	//
	// Methods
	//
	for _, attr := range s.graph.attributes {
		ct := attr.ctyType
		stmt.Add(
			jen.Comment(
				fmt.Sprintf(
					"%s returns a reference to field %s of %s.",
					strcase.Pascal(attr.name),
					attr.name,
					s.Type,
				),
			).
				Line().
				Func().
				// Receiver
				Params(jen.Id(s.Receiver).Id(s.AttributesStructName)).
				// Name
				Id(strcase.Pascal(attr.name)).Call().
				// Return type
				Add(ctyTypeReturnType(ct)).
				Block(
					jen.Return(
						funcReferenceByCtyType(ct).
							Call(
								structFieldRef().Dot("Append").Call(
									jen.Lit(attr.name),
								),
							),
					),
				),
		)
		stmt.Line()
		stmt.Line()
	}

	for _, child := range s.graph.children {
		structName := strcase.Pascal(child.uniqueName) + suffixAttributes
		qualStruct := jen.Qual(s.SubPkgQualPath(), structName).Clone
		stmt.Add(
			jen.Func().
				// Receiver
				Params(jen.Id(s.Receiver).Id(s.AttributesStructName)).
				// Name
				Id(strcase.Pascal(child.uniqueName)).Call().
				// Return type
				Add(
					returnTypeFromNestingPath(
						child.nestingPath,
						qualStruct(),
					),
				).
				Block(
					jen.Return(
						jenNodeReturnValue(child, qualStruct()).
							Call(
								structFieldRef().Dot("Append").Call(
									jen.Lit(child.name),
								),
							),
					),
				),
		)
		stmt.Line()
		stmt.Line()
	}

	return &stmt
}
