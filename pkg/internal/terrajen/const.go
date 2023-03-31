// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terrajen

const fileExtension = ".go"

const HeaderComment = `CODE GENERATED BY github.com/volvo-cars/lingon. DO NOT EDIT.`

const (
	pkgTerraAlias = "terra"
	pkgTerra      = "github.com/volvo-cars/lingon/pkg/terra"
	pkgHCLAlias   = "hcl"
	pkgHCL        = "github.com/hashicorp/hcl/v2"
	pkgHCLWrite   = "github.com/hashicorp/hcl/v2/hclwrite"
)

const (
	suffixArgs       = "Args"
	suffixAttributes = "Attributes"
	suffixState      = "State"
)

const (
	tagJSON     = "json"
	tagHCL      = "hcl"
	tagValidate = "validate"
)

const (
	idFuncInternalRef     = "InternalRef"
	idFuncInternalWithRef = "InternalWithRef"
	idFuncInternalTokens  = "InternalTokens"
)

const (
	idFieldName      = "Name"
	idFieldArgs      = "Args"
	idFieldState     = "state"
	idFieldLifecycle = "Lifecycle"
	idFieldDependsOn = "DependsOn"

	idFuncSource              = "Source"
	idFuncVersion             = "Version"
	idFuncType                = "Type"
	idFuncLocalName           = "LocalName"
	idFuncConfiguration       = "Configuration"
	idFuncDependOn            = "DependOn"
	idFuncDependencies        = "Dependencies"
	idFuncLifecycleManagement = "LifecycleManagement"
	idFuncAttributes          = "Attributes"
	idFuncImportState         = "ImportState"
	idFuncState               = "State"
	idFuncStateMust           = "StateMust"
)
