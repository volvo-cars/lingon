// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terrajen_test

import (
	"bytes"
	"fmt"
	"log/slog"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/volvo-cars/lingon/pkg/internal/terrajen"
	"github.com/zclconf/go-cty/cty"
)

var exampleProviderSchemas = tfjson.ProviderSchemas{
	Schemas: map[string]*tfjson.ProviderSchema{
		"registry.terraform.io/hashicorp/google": {
			// Schema for the provider configuration
			ConfigSchema: &tfjson.Schema{
				Block: &tfjson.SchemaBlock{
					Attributes: map[string]*tfjson.SchemaAttribute{
						"field": {
							AttributeType: cty.String,
							Required:      true,
						},
					},
				},
			},
			// Schema for the resource configurations
			ResourceSchemas: nil,
			// Schema for the data source configurations
			DataSourceSchemas: nil,
		},
	},
}

func Example() {
	// First you need to obtain a Terraform Provider Schema, which can be done using the
	// terragen package.
	// We will use some dummy data here
	ps := exampleProviderSchemas

	// Initialise the ProviderGenerator
	gen := terrajen.ProviderGenerator{
		// The module we are generating code into is called "mymodule" and the directory to generate
		// code into is "gen".
		// This allows the generated code to refer to itself using this pkg path as a prefix.
		GoProviderPkgPath: "mymodule/gen",
		// Filesystem location where to generate the files. This depends entirely on *where* you
		// run the generator from
		GeneratedPackageLocation: "./gen",
		// Name of the provider
		ProviderName: "google",
		// Source of the provider
		ProviderSource: "hashicorp/google",
		// Version of the provider
		ProviderVersion: "4.58.0",
	}

	// For each of the provider, resources and data sources in the provider schema, generate a
	// terrajen.Schema object using our ProviderGenerator
	prov := ps.Schemas["registry.terraform.io/hashicorp/google"]
	provSchema := gen.SchemaProvider(prov.ConfigSchema.Block)
	provFile := terrajen.ProviderFile(provSchema)
	// Typically we would save the file, e.g. provFile.Save("provider.go").
	// In this case we will write it to a buffer and print to stdout
	var pb bytes.Buffer
	if err := provFile.Render(&pb); err != nil {
		slog.Error("rendering provider file", err)
		return
	}
	fmt.Println(pb.String())

	// For each resource and data source (of which we have none), also do the same. E.g.
	// for name, res := range prov.ResourceSchemas {}
	// for name, data := range prov.DataSourceSchemas {}

	// Output:
	// // CODE GENERATED BY github.com/volvo-cars/lingon. DO NOT EDIT.
	//
	// package google
	//
	// import "github.com/volvo-cars/lingon/pkg/terra"
	//
	// func NewProvider(args ProviderArgs) *Provider {
	//	return &Provider{Args: args}
	// }
	//
	// var _ terra.Provider = (*Provider)(nil)
	//
	// type Provider struct {
	//	Args ProviderArgs
	// }
	//
	// // LocalName returns the provider local name for [Provider].
	// func (p *Provider) LocalName() string {
	//	return "google"
	// }
	//
	// // Source returns the provider source for [Provider].
	// func (p *Provider) Source() string {
	//	return "hashicorp/google"
	// }
	//
	// // Version returns the provider version for [Provider].
	// func (p *Provider) Version() string {
	//	return "4.58.0"
	// }
	//
	// // Configuration returns the configuration (args) for [Provider].
	// func (p *Provider) Configuration() interface{} {
	//	return p.Args
	// }
	//
	// // ProviderArgs contains the configurations for provider.
	// type ProviderArgs struct {
	//	// Field: string, required
	//	Field terra.StringValue `hcl:"field,attr" validate:"required"`
	// }
}
