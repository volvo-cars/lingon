package terrajen

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/volvo-cars/lingon/pkg/internal/str"
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

	for _, child := range s.graph.children {
		tags := map[string]string{
			tagHCL: child.uniqueName + ",block",
		}
		stmt := jen.Comment(child.comment()).Line()
		stmt.Add(jen.Id(str.PascalCase(child.uniqueName)))
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
		stmt.Qual(s.SubPkgQualPath(), str.PascalCase(child.uniqueName))
		stmt.Tag(tags)
		fields = append(fields, stmt)
	}

	// Add additional Terraform fields, like depends_on
	if s.SchemaType == SchemaTypeResource {
		fields = append(
			fields,
			jen.Comment(
				fmt.Sprintf(
					"// DependsOn contains resources that %s depends on",
					s.StructName,
				),
			).
				Line().
				Id("DependsOn").
				Qual(pkgTerra, "Dependencies").
				Tag(
					map[string]string{
						tagHCL: "depends_on,attr",
					},
				),
		)
	}

	return jen.Type().Id(s.ArgumentStructName).Struct(fields...)
}

// attributesStruct takes a schema and generates the Attributes struct that is used by the user to creates references to
// attributes for the object that the schema represents (e.g. provider, resource, data resource)
func attributesStruct(s *Schema) *jen.Statement {
	var stmt jen.Statement

	attrStruct := jen.Type().Id(s.AttributesStructName).Struct(
		jen.Id("name").String(),
	)

	stmt.Add(attrStruct)
	stmt.Line()

	//
	// Methods
	//
	for _, attr := range s.graph.attributes {
		ct := attr.ctyType
		implFunc := jen.Func().Params(jen.Id(s.Receiver).Id(s.AttributesStructName)).Id(str.PascalCase(attr.name)).Call()
		implFunc.Add(ctyTypeReturnType(ct))

		refList := []jen.Code{
			jen.Lit(s.Type), jen.Id(s.Receiver).Dot("name"), jen.Lit(attr.name),
		}
		// If schema is a data resource we need to prefix the "data" qualifier to the reference
		if s.SchemaType == SchemaTypeData {
			refList = append([]jen.Code{jen.Lit("data")}, refList...)
		}
		implFunc.Block(
			jen.Return(
				funcReferenceFromCtyType(ct, jen.List(refList...)),
			),
		)
		stmt.Line()
		stmt.Add(implFunc)
		stmt.Line()
	}

	for _, child := range s.graph.children {
		structName := str.PascalCase(child.uniqueName) + suffixAttributes
		implFunc := jen.Func().
			// Receiver
			Params(jen.Id(s.Receiver).Id(s.AttributesStructName)).
			// Name
			Id(str.PascalCase(child.uniqueName)).Call().
			//	Return
			Add(
				jenNodeReturnType(
					child,
					jen.Qual(s.SubPkgQualPath(), structName),
				),
			)

		refList := []jen.Code{
			jen.Lit(s.Type),
			jen.Id(s.Receiver).Dot("name"),
			jen.Lit(child.name),
		}
		// If schema is a data resource we need to prefix the "data" qualifier to the reference
		if s.SchemaType == SchemaTypeData {
			refList = append([]jen.Code{jen.Lit("data")}, refList...)
		}
		structStmt := jen.Qual(s.SubPkgQualPath(), structName).Values(
			jen.Dict{
				jen.Id(idStructReferenceValue): qualInternalRootRef().Call(jen.List(refList...)),
			},
		)

		implFunc.Block(
			jen.Return(
				jenNodeReturnValue(child, structStmt),
			),
		)

		stmt.Line()
		stmt.Add(implFunc)
		stmt.Line()
	}

	return &stmt
}
