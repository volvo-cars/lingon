// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terrajen

import (
	"github.com/dave/jennifer/jen"
)

const (
	idStructReference        = "Reference"
	idFuncReferenceAttribute = "ReferenceAttribute"
)

var (
	qualValue          = jen.Qual(pkgTerra, "Value").Clone
	qualStringValue    = jen.Qual(pkgTerra, "StringValue").Clone
	qualNumberValue    = jen.Qual(pkgTerra, "NumberValue").Clone
	qualBoolValue      = jen.Qual(pkgTerra, "BoolValue").Clone
	qualListValue      = jen.Qual(pkgTerra, "ListValue").Clone
	qualSetValue       = jen.Qual(pkgTerra, "SetValue").Clone
	qualMapValue       = jen.Qual(pkgTerra, "MapValue").Clone
	qualReferenceValue = jen.Qual(pkgTerra, idStructReference).Clone

	qualReferenceAttribute = jen.Qual(pkgTerra, idFuncReferenceAttribute).Clone

	qualReferenceString = jen.Qual(pkgTerra, "ReferenceString").Clone
	qualReferenceNumber = jen.Qual(pkgTerra, "ReferenceNumber").Clone
	qualReferenceBool   = jen.Qual(pkgTerra, "ReferenceBool").Clone
	qualReferenceSingle = jen.Qual(pkgTerra, "ReferenceSingle").Clone
	qualReferenceMap    = jen.Qual(pkgTerra, "ReferenceMap").Clone
	qualReferenceSet    = jen.Qual(pkgTerra, "ReferenceSet").Clone
	qualReferenceList   = jen.Qual(pkgTerra, "ReferenceList").Clone

	qualHCLWriteTokens = jen.Qual(pkgHCLWrite, "Tokens").Clone
)
