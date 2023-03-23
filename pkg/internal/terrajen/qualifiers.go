package terrajen

import (
	"github.com/dave/jennifer/jen"
)

const (
	idStructReferenceValue = "ReferenceValue"
	idFuncReference        = "Reference"
)

var (
	qualValue          = jen.Qual(pkgTerra, "Value").Clone
	qualStringValue    = jen.Qual(pkgTerra, "StringValue").Clone
	qualNumberValue    = jen.Qual(pkgTerra, "NumberValue").Clone
	qualBoolValue      = jen.Qual(pkgTerra, "BoolValue").Clone
	qualListValue      = jen.Qual(pkgTerra, "ListValue").Clone
	qualSetValue       = jen.Qual(pkgTerra, "SetValue").Clone
	qualMapValue       = jen.Qual(pkgTerra, "MapValue").Clone
	qualReferenceValue = jen.Qual(pkgTerra, idStructReferenceValue).Clone

	qualInternalRootRef = jen.Qual(pkgTerra, idFuncReference).Clone

	qualAsMapRefFunc  = jen.Qual(pkgTerra, "AsMapRef").Clone
	qualAsSetRefFunc  = jen.Qual(pkgTerra, "AsSetRef").Clone
	qualAsListRefFunc = jen.Qual(pkgTerra, "AsListRef").Clone

	qualHCLTraverseAttr = jen.Qual(pkgHCL, "TraverseAttr").Clone
)
