// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terrajen

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/volvo-cars/lingon/pkg/internal/str"
)

func SubPkgFile(s *Schema) (*jen.File, bool) {
	// Skip sub pkg if there are no blocks to render
	if s.graph.isEmpty() {
		return nil, false
	}
	f := jen.NewFile(s.SubPackageName)
	f.ImportAlias(pkgHCL, "hcl")
	f.HeaderComment(HeaderComment)
	for _, n := range s.graph.nodes {
		f.Add(subPkgArgStruct(n))
	}
	for _, n := range s.graph.nodes {
		f.Add(subPkgAttributeStruct(n))
	}
	for _, n := range s.graph.nodes {
		f.Add(subPkgStateStruct(n))
	}

	return f, true
}

func subPkgArgStruct(n *node) *jen.Statement {
	fields := make([]jen.Code, 0)

	for _, attr := range n.attributes {
		// Skip attributes that are not arguments
		if !attr.isArg {
			continue
		}
		stmt := jen.Comment(attr.comment()).Line()
		stmt.Add(jen.Id(str.PascalCase(attr.name)))
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

	for _, child := range n.children {
		stmt := jen.Comment(child.comment()).Line()
		stmt.Add(jen.Id(str.PascalCase(child.name)))
		tags := map[string]string{
			tagHCL: child.name + ",block",
		}
		if child.isSingularArg() {
			stmt.Op("*")
			if child.isRequired {
				tags[tagValidate] = "required"
			}
		} else {
			// For children the nesting type cannot be a map
			for _, path := range child.nestingPath {
				switch path {
				case nodeNestingModeList, nodeNestingModeSet:
					stmt.Index()
				default:
					panic(
						fmt.Sprintf(
							"unsupported nesting path %d for child",
							path,
						),
					)
				}
			}
			tags[tagValidate] = nodeBlockListValidateTags(child)
		}
		stmt.Id(str.PascalCase(child.uniqueName))
		stmt.Tag(tags)
		fields = append(fields, stmt)
	}

	stmt := jen.
		Type().Id(n.argsStructName()).
		Struct(fields...)
	stmt.Line()
	stmt.Line()

	return stmt
}

func subPkgAttributeStruct(n *node) *jen.Statement {
	structName := n.attributesStructName()

	structFieldRef := "ref"
	refArg := "ref"
	stmt := jen.Type().Id(structName).
		Struct(
			jen.Id(structFieldRef).Add(qualReferenceValue()),
		)
	stmt.Line()
	stmt.Line()

	// Methods
	// Override InternalRef, e.g.
	//
	// 	func (i OidcRef) InternalRef() (terra.Reference, error) {
	// 		return i.ref, nil
	// 	}
	stmt.Add(
		jen.Func().
			// Receiver
			Params(jen.Id(n.receiver).Id(structName)).
			// Name
			Id(idFuncInternalRef).Call().
			// Return type
			Params(qualReferenceValue(), jen.Error()).
			// Body
			Block(
				jen.Return(
					jen.Id(n.receiver).Dot(
						structFieldRef,
					),
					jen.Nil(),
				),
			),
	)
	stmt.Line()
	stmt.Line()

	// Override InternalWithRef, e.g.
	//
	// 	func (i OidcRef) InternalWithRef(ref terra.Reference) OidcRef {
	// 		return terra.ReferenceSingle[OidcRef](ref)
	// 	}
	stmt.Add(
		jen.Func().
			// Receiver
			Params(jen.Id(n.receiver).Id(structName)).
			// Name
			Id(idFuncInternalWithRef).Call(
			jen.Id(refArg).Add(qualReferenceValue()),
		).
			// Return type
			Id(structName).
			// Body
			Block(
				jen.Return(
					jen.Id(structName).Values(
						jen.Dict{
							jen.Id(structFieldRef): jen.Id(refArg),
						},
					),
				),
			),
	)
	stmt.Line()
	stmt.Line()
	// Override InternalTokens
	stmt.Add(
		jen.Func().
			// Receiver
			Params(jen.Id(n.receiver).Id(structName)).
			// Name
			Id(idFuncInternalTokens).
			Call().
			// Return type
			// Id(structName).
			Add(qualHCLWriteTokens()).
			// Body
			Block(
				jen.Return(
					jen.Id(n.receiver).Dot(
						structFieldRef,
					).Dot(idFuncInternalTokens).Call(),
				),
			),
	)
	stmt.Line()
	stmt.Line()

	for _, attr := range n.attributes {
		appendRef := jen.Id(n.receiver).
			Dot(refArg).
			Dot("Append").
			Call(jen.Lit(attr.name))
		stmt.Add(
			jen.Func().
				// Receiver
				Params(jen.Id(n.receiver).Id(structName)).
				// Name
				Id(str.PascalCase(attr.name)).Call().
				//	Return type
				Add(ctyTypeReturnType(attr.ctyType)).
				// Body
				Block(
					jen.Return(
						funcReferenceByCtyType(attr.ctyType).
							Call(appendRef),
					),
				),
		)
		stmt.Line()
		stmt.Line()
	}

	for _, child := range n.children {
		childStructName := child.attributesStructName()
		appendRef := jen.Id(n.receiver).
			Dot(refArg).
			Dot("Append").
			Call(jen.Lit(child.name))

		stmt.Add(
			jen.Func().
				// Receiver
				Params(jen.Id(n.receiver).Id(structName)).
				// Name
				Id(str.PascalCase(child.name)).Call().
				// Return type
				Add(
					returnTypeFromNestingPath(
						child.nestingPath,
						jen.Id(childStructName),
					),
				).Block(
				jen.Return(
					jenNodeReturnValue(child, jen.Id(childStructName)).
						Call(appendRef),
				),
			),
		)
		stmt.Line()
		stmt.Line()
	}

	return stmt
}

func subPkgStateStruct(n *node) *jen.Statement {
	fields := make([]jen.Code, 0)

	for _, attr := range n.attributes {
		pan := str.PascalCase(attr.name)
		stmt := jen.Id(pan)
		stmt.Add(ctyTypeToGoType(attr.ctyType, pan))
		// Add tags
		stmt.Tag(
			map[string]string{
				tagJSON: attr.name,
			},
		)
		fields = append(fields, stmt)
	}

	for _, child := range n.children {
		stmt := jen.Id(str.PascalCase(child.name))
		if child.isSingularState() {
			stmt.Op("*")
		} else {
			stmt.Index()
		}
		stmt.Id(child.stateStructName())
		stmt.Tag(
			map[string]string{
				tagJSON: child.name,
			},
		)
		fields = append(fields, stmt)
	}

	stmt := jen.
		Type().Id(n.stateStructName()).
		Struct(fields...)
	stmt.Line()
	stmt.Line()

	return stmt
}
