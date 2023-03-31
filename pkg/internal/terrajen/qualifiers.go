// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terrajen

import (
	"github.com/dave/jennifer/jen"
)

const (
	idStructReference           = "Reference"
	idFuncReferenceResource     = "ReferenceResource"
	idFuncReferenceDataResource = "ReferenceDataResource"
)

var (
	qualStringValue    = jen.Qual(pkgTerra, "StringValue").Clone
	qualNumberValue    = jen.Qual(pkgTerra, "NumberValue").Clone
	qualBoolValue      = jen.Qual(pkgTerra, "BoolValue").Clone
	qualListValue      = jen.Qual(pkgTerra, "ListValue").Clone
	qualSetValue       = jen.Qual(pkgTerra, "SetValue").Clone
	qualMapValue       = jen.Qual(pkgTerra, "MapValue").Clone
	qualReferenceValue = jen.Qual(pkgTerra, idStructReference).Clone

	qualReferenceResource = jen.Qual(
		pkgTerra,
		idFuncReferenceResource,
	).Clone
	qualReferenceDataResource = jen.Qual(
		pkgTerra,
		idFuncReferenceDataResource,
	).Clone

	qualReferenceAsString = jen.Qual(pkgTerra, "ReferenceAsString").Clone
	qualReferenceAsNumber = jen.Qual(pkgTerra, "ReferenceAsNumber").Clone
	qualReferenceAsBool   = jen.Qual(pkgTerra, "ReferenceAsBool").Clone
	qualReferenceAsSingle = jen.Qual(pkgTerra, "ReferenceAsSingle").Clone
	qualReferenceAsMap    = jen.Qual(pkgTerra, "ReferenceAsMap").Clone
	qualReferenceAsSet    = jen.Qual(pkgTerra, "ReferenceAsSet").Clone
	qualReferenceAsList   = jen.Qual(pkgTerra, "ReferenceAsList").Clone

	qualTypeDependencies = jen.Qual(pkgTerra, "Dependencies").Clone
	qualStructLifecycle  = jen.Qual(pkgTerra, "Lifecycle").Clone
	// qualFuncIgnoreChanges      = jen.Qual(pkgTerra, "IgnoreChanges").Clone
	// qualFuncReplaceTriggeredBy = jen.Qual(pkgTerra, "ReplaceTriggeredBy").Clone

	qualHCLWriteTokens = jen.Qual(pkgHCLWrite, "Tokens").Clone
)
