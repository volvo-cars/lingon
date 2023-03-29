// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terrajen

import (
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
			for range child.nestingPath {
				stmt.Index()
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

	stmt := jen.Type().Id(structName).
		Struct(
			qualReferenceValue(),
		)
	stmt.Line()
	stmt.Line()

	// Methods
	// Override InternalTraverse, e.g.
	//
	// 	func (i OidcRef) InternalTraverse(step hcl.Traverser) OidcRef {
	// 		return OidcRef{
	// 			Reference: i.Reference.InternalTraverse(step),
	// 		}
	// 	}
	stepArg := "step"
	stmt.Add(
		jen.Func().
			// Receiver
			Params(jen.Id(n.receiver).Id(structName)).
			// Name
			Id(idFuncInternalTraverse).Call(
			jen.Id(stepArg).Qual(
				pkgHCL,
				"Traverser",
			),
		).
			// Return type
			Id(structName).
			// Body
			Block(
				jen.Return(
					jen.Id(structName).Values(
						jen.Dict{
							jen.Id(idStructReferenceValue): jen.Id(n.receiver).Dot(idStructReferenceValue).Dot(idFuncInternalTraverse).Call(jen.Id(stepArg)),
						},
					),
				),
			),
	)
	stmt.Line()
	stmt.Line()

	for _, attr := range n.attributes {
		// Want: return terra.AsList(InternalStrRef(i.Reference.InternalTraverse(hcl.TraverseAttr{Name: "issuer"})))
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
					jen.Return(ctyTypeReturnValue(n, attr.ctyType, attr.name)),
				),
		)
		stmt.Line()
		stmt.Line()
	}

	for _, child := range n.children {
		// Want: return terra.AsList(OidcRef(i.InternalTraverse(hcl.TraverseAttr{Name: "oidc"})))
		childStructName := child.attributesStructName()

		implFunc := jen.Func().
			// Receiver
			Params(jen.Id(n.receiver).Id(structName)).
			// Name
			Id(str.PascalCase(child.name)).Call().
			// Return type
			Add(jenNodeReturnType(child, jen.Id(childStructName)))

		// Want: ChildAttributes(i.InternalTraverse(hcl.TraverseAttr{Name: "child"}))
		childValue := jen.
			// Create instance of child struct
			Id(childStructName).Call(
			// Call InternalTraverse on self, and pass result to child struct
			jen.Id(n.receiver).Dot(idFuncInternalTraverse).Call(hclTraverseAttr(jen.Lit(child.name))),
		)
		implFunc.Block(
			jen.Return(jenNodeReturnValue(child, childValue)),
		)

		stmt.Add(implFunc)
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
