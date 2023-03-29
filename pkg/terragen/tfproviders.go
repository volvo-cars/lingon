// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terragen

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/volvo-cars/lingon/pkg/internal/hcl"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

// TerraformVersions ...
type TerraformVersions struct {
	TerraformBlock TerraformBlock `hcl:"terraform,block"`
}

// TerraformBlock represents a terraform{} block in a Terraform stack
type TerraformBlock struct {
	RequiredProviders RequiredProviders `hcl:"required_providers,block"`
}

// RequiredProviders represents the map of required providers for a Terraform stack
type RequiredProviders struct {
	Providers map[string]Provider `hcl:",remain"`
}

// Provider represents a single element of a map of required providers
type Provider struct {
	Name    string // Lack of cty tag means it is ignored
	Source  string `cty:"source"`
	Version string `cty:"version"`
}

func GenerateProviderSchema(
	ctx context.Context,
	provider Provider,
) (*tfjson.ProviderSchemas, error) {
	versions := TerraformVersions{
		TerraformBlock: TerraformBlock{
			RequiredProviders: RequiredProviders{
				Providers: map[string]Provider{
					provider.Name: {
						Source:  provider.Source,
						Version: provider.Version,
					},
				},
			},
		},
	}
	workingDir := filepath.Join(
		".lingon", "schemas", provider.Name,
		provider.Version,
	)
	if err := os.MkdirAll(workingDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf(
			"creating schemas working directory: %s: %w",
			workingDir,
			err,
		)
	}

	// Write versions.tf file
	tfVersionsFile := filepath.Join(workingDir, "versions.tf")
	f, err := os.Create(tfVersionsFile)
	if err != nil {
		return nil, fmt.Errorf("creating file %s: %w", tfVersionsFile, err)
	}
	if err := hcl.EncodeRaw(f, versions); err != nil {
		return nil, fmt.Errorf("encoding file %s: %w", tfVersionsFile, err)
	}

	tf, err := tfexec.NewTerraform(workingDir, "terraform")
	if err != nil {
		return nil, fmt.Errorf("creating new terraform runtime: %w", err)
	}

	if err := tf.Init(ctx, tfexec.Upgrade(true)); err != nil {
		return nil, fmt.Errorf("running terraform init: %w", err)
	}

	providersSchema, err := tf.ProvidersSchema(ctx)
	if err != nil {
		return nil, fmt.Errorf("running terraform providers schema: %w", err)
	}

	return providersSchema, nil
}
