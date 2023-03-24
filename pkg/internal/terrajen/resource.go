// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package terrajen

import (
	"github.com/dave/jennifer/jen"
	"github.com/volvo-cars/lingon/pkg/internal/str"
)

// ResourceFile generates a Go file for a Terraform resource configuration based on the given Schema
func ResourceFile(s *Schema) *jen.File {
	f := jen.NewFile(s.PackageName)
	f.ImportAlias(pkgHCL, pkgHCLAlias)
	f.ImportAlias(pkgTerra, pkgTerraAlias)
	f.HeaderComment(HeaderComment)
	f.Add(resourceNewFunc(s))
	f.Add(resourceStructCompileCheck(s))
	f.Add(resourceStruct(s))
	f.Add(argsStruct(s))
	f.Add(attributesStruct(s))
	f.Add(resourceStateStruct(s))

	return f
}

func resourceNewFunc(s *Schema) *jen.Statement {
	return jen.Func().Id(s.NewFuncName).Params(
		jen.Id("name").String(),
		jen.Id("args").Id(s.ArgumentStructName),
	).
		// Return
		Op("*").Id(s.StructName).
		// Block
		Block(
			jen.Return(
				jen.Op("&").Id(s.StructName).Values(
					jen.Dict{
						jen.Id(idFieldName): jen.Id("name"),
						jen.Id(idFieldArgs): jen.Id("args"),
					},
				),
			),
		)
}

func resourceStructCompileCheck(s *Schema) *jen.Statement {
	return jen.Var().Op("_").Qual(pkgTerra, "Resource").Op("=").
		Params(
			jen.Op("*").Id(s.StructName),
		).
		Params(jen.Nil()).
		Line()
}

func resourceStruct(s *Schema) *jen.Statement {
	stmt := jen.Type().Id(s.StructName).Struct(
		jen.Id(idFieldName).String(),
		jen.Id(idFieldArgs).Id(s.ArgumentStructName),
		jen.Id(idFieldState).Op("*").Id(s.StateStructName),
	)
	stmt.Line()
	stmt.Line()

	// Methods
	// Type
	stmt.Add(funcSchemaType(s, "Type"))
	stmt.Line()
	stmt.Line()
	// LocalName
	stmt.Add(funcLocalName(s))
	stmt.Line()
	stmt.Line()
	// Configuration
	stmt.Add(funcConfiguration(s))
	stmt.Line()
	stmt.Line()
	// Attributes
	stmt.Add(funcAttributes(s))
	stmt.Line()
	stmt.Line()
	// ImportState
	stmt.Add(funcResourceImportState(s))
	stmt.Line()
	stmt.Line()
	// State
	stmt.Add(funcResourceState(s))
	stmt.Line()
	stmt.Line()
	// StateMust
	stmt.Add(funcResourceStateMust(s))
	stmt.Line()
	stmt.Line()
	// DependOn
	stmt.Add(funcDependOn(s))
	stmt.Line()
	stmt.Line()

	return stmt
}

func resourceStateStruct(s *Schema) *jen.Statement {
	fields := make([]jen.Code, 0)
	for _, attr := range s.graph.attributes {
		pan := str.PascalCase(attr.name)
		stmt := jen.Id(pan)
		stmt.Add(ctyTypeToGoType(attr.ctyType, pan))
		stmt.Tag(
			map[string]string{
				tagJSON: attr.name,
			},
		)
		fields = append(fields, stmt)
	}

	for _, child := range s.graph.children {
		pbn := str.PascalCase(child.name)
		stmt := jen.Id(pbn)
		if len(child.nestingPath) == 0 {
			stmt.Op("*")
		} else {
			stmt.Index()
		}

		stmt.Qual(s.SubPkgQualPath(), str.PascalCase(child.name)+suffixState)
		stmt.Tag(
			map[string]string{
				tagJSON: child.name,
			},
		)
		fields = append(fields, stmt)
	}
	return jen.Type().Id(s.StateStructName).Struct(fields...)
}
