// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terrajen

import (
	"github.com/dave/jennifer/jen"
)

// ProviderFile generates a Go file for a Terraform provider configuration based on the given Schema
func ProviderFile(s *Schema) *jen.File {
	f := jen.NewFile(s.ProviderName)
	f.HeaderComment(HeaderComment)
	f.Add(providerNewFunc(s))
	f.Add(providerStructCompileCheck(s))
	f.Add(providerStruct(s))
	f.Add(argsStruct(s))

	return f
}

func providerNewFunc(s *Schema) *jen.Statement {
	return jen.Func().Id(s.NewFuncName).Params(
		jen.Id("args").Id(s.ArgumentStructName),
	).
		// Return
		Op("*").Id(s.StructName).
		// Block
		Block(
			jen.Return(
				jen.Op("&").Id(s.StructName).Values(
					jen.Dict{
						jen.Id(idFieldArgs): jen.Id("args"),
					},
				),
			),
		)
}

func providerStructCompileCheck(s *Schema) *jen.Statement {
	return jen.Var().Op("_").Qual(pkgTerra, "Provider").Op("=").
		Params(
			jen.Op("*").Id(s.StructName),
		).
		Params(jen.Nil()).
		Line()
}

func providerStruct(s *Schema) *jen.Statement {
	stmt := jen.Type().Id(s.StructName).Struct(
		jen.Id(idFieldArgs).Id(s.ArgumentStructName),
		jen.Line(),
	)
	stmt.Line()
	stmt.Line()

	// LocalName
	stmt.Add(funcProviderLocalName(s))
	stmt.Line()
	stmt.Line()
	// Source
	stmt.Add(funcProviderSource(s))
	stmt.Line()
	stmt.Line()
	// Version
	stmt.Add(funcProviderVersion(s))
	stmt.Line()
	stmt.Line()
	// Configuration
	stmt.Add(funcConfiguration(s))
	stmt.Line()
	stmt.Line()

	return stmt
}
