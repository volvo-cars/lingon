// Copyright (c) Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package hcl

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TerraformBlock struct {
	Backend           BackendBlock            `hcl:"backend,block"`
	RequiredProviders RequiredProvidersBlocks `hcl:"required_providers,block"`
}

type RequiredProvidersBlocks struct {
	Providers map[string]RequiredProvider `hcl:",remain"`
}

type RequiredProvider struct {
	Source  string `hcl:"source,attr" cty:"source"`
	Version string `hcl:"version,attr" cty:"version"`
}

type BackendBlock struct {
	Type string `hcl:",label"`

	CommonBlockConfig `hcl:",remain"`
}

type ProviderBlock struct {
	LocalName string `hcl:",label"`

	CommonBlockConfig `hcl:",remain"`
}

type DataResourceBlock struct {
	DataSource string `hcl:",label"`
	LocalName  string `hcl:",label"`

	CommonBlockConfig `hcl:",remain"`
}

type ResourceBlock struct {
	Type      string `hcl:",label"`
	LocalName string `hcl:",label"`

	CommonBlockConfig `hcl:",remain"`
}

type CommonBlockConfig struct {
	StringField string `hcl:"string_field,attr"`
	NumberField int    `hcl:"number_field,attr"`
}

// HCLFile is the test structure which we will use to encode HCL and also decode the resulting HCL
// into
type HCLFile struct {
	Terraform     TerraformBlock      `hcl:"terraform,block"`
	Providers     []ProviderBlock     `hcl:"provider,block"`
	DataResources []DataResourceBlock `hcl:"data,block"`
	Resources     []ResourceBlock     `hcl:"resource,block"`
}

// Create a common config which we reuse. We might want to make this a bit more diverse
// or create separate instance of this for greater coverage.
var cbcfg = CommonBlockConfig{
	StringField: "test",
	NumberField: 123,
}

// TestEncode tests the encoder.
// The approach is to
// 1. Create an expected HCL value
// 2. Encode HCL based on the expected value
// 3. Use the hcl library to decode the encoded HCL into a new instance of the expected value
// 4. Assert that they match
//
// This is a bit complex but much more convincing and easier to maintain in the long run than
// string-based comparisons.
func TestEncode(t *testing.T) {
	// Create the expected HCL values which we will use to create the encoder arguments
	expectedHCL := HCLFile{
		Terraform: TerraformBlock{
			Backend: BackendBlock{
				Type:              "dummy",
				CommonBlockConfig: cbcfg,
			},
			RequiredProviders: RequiredProvidersBlocks{
				Providers: map[string]RequiredProvider{
					"aws": {
						Source:  "hashicorp/aws",
						Version: "4.49.0",
					},
				},
			},
		},
		Providers: []ProviderBlock{
			{
				LocalName:         "aws",
				CommonBlockConfig: cbcfg,
			},
		},
		DataResources: []DataResourceBlock{
			{
				DataSource:        "test_data_source",
				LocalName:         "test",
				CommonBlockConfig: cbcfg,
			},
		},
		Resources: []ResourceBlock{
			{
				Type:              "test_resource",
				LocalName:         "test",
				CommonBlockConfig: cbcfg,
			},
		},
	}

	// Populate the encoder arguments using the expected HCL value
	args := EncodeArgs{
		Backend: &Backend{
			Type:          expectedHCL.Terraform.Backend.Type,
			Configuration: expectedHCL.Terraform.Backend.CommonBlockConfig,
		},
	}
	for name, prov := range expectedHCL.Terraform.RequiredProviders.Providers {
		args.Providers = append(
			args.Providers, Provider{
				LocalName:     name,
				Source:        prov.Source,
				Version:       prov.Version,
				Configuration: cbcfg,
			},
		)
	}
	for _, data := range expectedHCL.DataResources {
		args.DataResources = append(
			args.DataResources, DataResource{
				DataSource:    data.DataSource,
				LocalName:     data.LocalName,
				Configuration: data.CommonBlockConfig,
			},
		)
	}
	for _, res := range expectedHCL.Resources {
		args.Resources = append(
			args.Resources, Resource{
				Type:          res.Type,
				LocalName:     res.LocalName,
				Configuration: res.CommonBlockConfig,
			},
		)
	}
	// Run the encoder
	var b bytes.Buffer
	err := Encode(&b, args)
	require.NoError(t, err)

	fmt.Println(b.String())

	// Decode the encoded HCL into a new instance of our test structure and compare
	actualHCL := HCLFile{}
	err = hclsimple.Decode("test.hcl", b.Bytes(), nil, &actualHCL)
	require.NoError(t, err)
	assert.Equal(t, expectedHCL, actualHCL)
}

func TestEncodeRaw(t *testing.T) {
	expectedHCL := HCLFile{
		Terraform: TerraformBlock{
			Backend: BackendBlock{
				Type:              "sometype",
				CommonBlockConfig: cbcfg,
			},
			RequiredProviders: RequiredProvidersBlocks{
				Providers: map[string]RequiredProvider{
					"localprovider": {
						Source:  "somesource",
						Version: "someversion",
					},
					"another": {
						Source:  "anothersource",
						Version: "anotherversion",
					},
				},
			},
		},
		Providers: []ProviderBlock{
			{
				LocalName:         "localname",
				CommonBlockConfig: cbcfg,
			},
		},
		DataResources: []DataResourceBlock{
			{
				DataSource:        "some_data_source",
				LocalName:         "localname",
				CommonBlockConfig: cbcfg,
			},
		},
		Resources: []ResourceBlock{
			{
				Type:              "some_resource_type",
				LocalName:         "this",
				CommonBlockConfig: cbcfg,
			},
		},
	}
	var b bytes.Buffer
	err := EncodeRaw(&b, expectedHCL)
	require.NoError(t, err)

	fmt.Println(b.String())

	actualHCL := HCLFile{}
	err = hclsimple.Decode("test.hcl", b.Bytes(), nil, &actualHCL)
	require.NoError(t, err)

	assert.Equal(t, expectedHCL, actualHCL)
}
