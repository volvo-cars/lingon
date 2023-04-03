// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terrajen

import (
	"go/token"
	"path/filepath"
	"strings"

	"github.com/volvo-cars/lingon/pkg/internal/str"

	tfjson "github.com/hashicorp/terraform-json"
)

// ProviderGenerator is created for each provider and is used to generate the schema
// for each resource and data object, and the provider configuration.
// The schemas are used by the generator to create the Go files and sub packages.
type ProviderGenerator struct {
	// GoProviderPkgPath is the Go pkg path to the generated provider directory.
	// E.g. github.com/volvo-cars/github.com/volvo-cars/lingon/gen/aws
	GoProviderPkgPath string
	// GeneratedPackageLocation is the directory on the filesystem where the generated
	// Go files will be created.
	// The GoProviderPkgPath path must match the location of the generated files
	// so that they can be imported correctly.
	// E.g. if we are in a Go module called "my-module" and we generate the files in a
	// "gen" directory within the root of "my-module", then GoProviderPkgPath is "my-module/gen"
	// and the GeneratedPackageLocation is "./gen" assuming we are running from the root of
	// "my-module"
	GeneratedPackageLocation string
	// ProviderName is the local name of the provider.
	// E.g. aws
	// https://developer.hashicorp.com/terraform/language/providers/requirements#local-names
	ProviderName string
	// ProviderSource is the source address of the provider.
	// E.g. registry.terraform.io/hashicorp/aws
	// https://developer.hashicorp.com/terraform/language/providers/requirements#source-addresses
	ProviderSource string
	// ProviderVersion is the version of thr provider.
	// E.g. 4.49.0
	ProviderVersion string
}

type SchemaType string

const (
	SchemaTypeProvider SchemaType = "provider"
	SchemaTypeResource SchemaType = "resource"
	SchemaTypeData     SchemaType = "data"
)

// SchemaProvider creates a schema for the provider config block for the provider
// represented by ProviderGenerator
func (a *ProviderGenerator) SchemaProvider(sb *tfjson.SchemaBlock) *Schema {
	return &Schema{
		SchemaType:               SchemaTypeProvider,
		GoProviderPkgPath:        a.GoProviderPkgPath,        // github.com/volvo-cars/github.com/volvo-cars/lingon/gen/aws
		GeneratedPackageLocation: a.GeneratedPackageLocation, // gen/aws
		ProviderName:             a.ProviderName,             // aws
		ProviderSource:           a.ProviderSource,           // registry.terraform.io/hashicorp/aws
		ProviderVersion:          a.ProviderVersion,          // 4.49.0
		PackageName:              a.ProviderName,             // aws
		ShortName:                "provider",
		Type:                     "provider",
		StructName:               "Provider",
		ArgumentStructName:       "ProviderArgs",
		Receiver:                 structReceiverFromName("provider"),

		NewFuncName:    "NewProvider",
		SubPackageName: "provider",
		FilePath: filepath.Join(
			a.GeneratedPackageLocation,
			"provider"+fileExtension,
		),
		graph: newGraph(sb),
	}
}

// SchemaResource creates a schema for the given resource for the provider
// represented by ProviderGenerator
func (a *ProviderGenerator) SchemaResource(
	name string,
	sb *tfjson.SchemaBlock,
) *Schema {
	shortName := providerShortName(name)
	spn := strings.ReplaceAll(shortName, "_", "")
	fp := filepath.Join(a.GeneratedPackageLocation, shortName+fileExtension)
	rs := &Schema{
		SchemaType:               SchemaTypeResource,
		GoProviderPkgPath:        a.GoProviderPkgPath,        // github.com/volvo-cars/github.com/volvo-cars/lingon/gen/aws
		GeneratedPackageLocation: a.GeneratedPackageLocation, // gen/aws
		ProviderName:             a.ProviderName,             // aws
		ProviderSource:           a.ProviderSource,           // hashicorp/aws
		ProviderVersion:          a.ProviderVersion,          // 4.49.0
		ShortName:                shortName,                  // aws_iam_role => iam_role
		PackageName:              a.ProviderName,             // aws
		Type:                     name,                       // aws

		StructName:           str.PascalCase(shortName),                   // iam_role => IamRole
		ArgumentStructName:   str.PascalCase(shortName) + suffixArgs,      // iam_role => IamRoleArgs
		AttributesStructName: str.CamelCase(shortName) + suffixAttributes, // iam_role => iamRoleAttributes
		StateStructName:      str.CamelCase(shortName) + suffixState,      // iam_role => IamRoleOut
		Receiver:             structReceiverFromName(shortName),           // iam_role => ir

		NewFuncName:    "New" + str.PascalCase(shortName),
		SubPackageName: spn, // iam_role => iamrole
		FilePath:       fp,
		graph:          newGraph(sb),
	}
	return rs
}

// SchemaData creates a schema for the given data object for the provider
// represented by ProviderGenerator
func (a *ProviderGenerator) SchemaData(
	name string,
	sb *tfjson.SchemaBlock,
) *Schema {
	shortName := providerShortName(name)
	spn := strings.ReplaceAll(shortName, "_", "")
	dataName := "data_" + shortName
	fp := filepath.Join(a.GeneratedPackageLocation, dataName+fileExtension)
	pn := str.PascalCase(shortName)

	ds := &Schema{
		SchemaType:               SchemaTypeData,
		GoProviderPkgPath:        a.GoProviderPkgPath,        // github.com/volvo-cars/github.com/volvo-cars/lingon/gen/aws
		GeneratedPackageLocation: a.GeneratedPackageLocation, // gen/aws
		ProviderName:             a.ProviderName,             // aws
		ProviderSource:           a.ProviderSource,           // hashicorp/aws
		ProviderVersion:          a.ProviderVersion,          // 4.49.0
		ShortName:                shortName,                  // aws_iam_role => iam_role
		PackageName:              a.ProviderName,             // aws
		Type:                     name,                       // aws_iam_role

		StructName:           "Data" + pn,                       // iam_role => DataIamRole
		ArgumentStructName:   "Data" + pn + suffixArgs,          // iam_role => DataIamRoleArgs
		AttributesStructName: "data" + pn + suffixAttributes,    // iam_role => dataIamRoleAttributes
		Receiver:             structReceiverFromName(shortName), // iam_role => ir

		NewFuncName:    "NewData" + pn, // iam_role => NewDataIamRole
		SubPackageName: "data" + spn,   // iam_role => dataiamrole
		FilePath:       fp,
		graph:          newGraph(sb),
	}

	return ds
}

// providerShortName takes a name like "aws_iam_role" and returns the name without
// the leading provider prefix, i.e. it returns "iam_role"
func providerShortName(name string) string {
	underscoreIndex := strings.Index(name, "_")
	if underscoreIndex == -1 {
		return name
	}
	return name[underscoreIndex+1:]
}

// structReceiverFromName calculates a suitable receiver from the name of the object.
// It gets the first character of each word separated by underscores, e.g. iam_role => ir
func structReceiverFromName(name string) string {
	ss := strings.Split(name, "_")
	var receiver strings.Builder
	for _, s := range ss {
		receiver.WriteString(s[0:1])
	}
	r := receiver.String()
	// Avoid using keywords for the receiver!
	if token.Lookup(r).IsKeyword() || r == "nil" {
		r = "_" + r
	}
	return r
}

// Schema is used to store all the relevant information required for the Go
// code generator.
// A schema can represent a resource, a data object or the provider configuration.
type Schema struct {
	SchemaType               SchemaType // resource / provider / data
	GoProviderPkgPath        string     // github.com/volvo-cars/github.com/volvo-cars/lingon/gen/providers
	GeneratedPackageLocation string     // gen/providers/aws
	ProviderName             string     // aws
	ProviderSource           string     // registry.terraform.io/hashicorp/aws
	ProviderVersion          string     // 4.49.0
	ShortName                string     // aws_iam_role => iam_role
	PackageName              string     // aws
	Type                     string     // aws_iam_role

	// Structs
	StructName           string // iam_role => IamRole
	ArgumentStructName   string // iam_role => IamRoleArgs
	AttributesStructName string // iam_role => iamRoleAttributes
	StateStructName      string // iam_role => iamRoleState

	Receiver string // iam_role => ir

	NewFuncName    string // iam_role => NewIamRole
	SubPackageName string // iam_role => iamrole
	FilePath       string // gen/providers/aws/ xxx
	graph          *graph
}

func (s *Schema) SubPkgQualPath() string {
	return s.GoProviderPkgPath + "/" + s.SubPackageName
}

func (s *Schema) SubPkgPath() string {
	return filepath.Join(
		s.GeneratedPackageLocation,
		s.SubPackageName,
		s.ShortName+fileExtension,
	)
}
