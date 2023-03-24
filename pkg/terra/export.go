// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/volvo-cars/lingon/pkg/internal/hcl"
)

const (
	tfSuffix = ".tf"
)

func ExportWriter(stack Exporter, w io.Writer) error {
	if err := encodeStack(stack, w); err != nil {
		return fmt.Errorf(
			"encoding stack: %w", err,
		)
	}
	return nil
}

// Export encodes Exporter to HCL and writes it to the given outDir
func Export(stack Exporter, outDir string) error {
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(
		filepath.Join(
			outDir,
			"main"+tfSuffix,
		),
	)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := encodeStack(stack, f); err != nil {
		return fmt.Errorf(
			"encoding stack: %w", err,
		)
	}

	return nil
}

func encodeStack(stack Exporter, w io.Writer) error {
	blocks, err := objectsFromStack(stack)
	if err != nil {
		return err
	}
	if err := validateStack(blocks); err != nil {
		return fmt.Errorf("validating stack: %w", err)
	}

	args := hcl.EncodeArgs{
		Providers:     make([]hcl.Provider, len(blocks.Providers)),
		DataResources: make([]hcl.DataResource, len(blocks.DataResources)),
		Resources:     make([]hcl.Resource, len(blocks.Resources)),
	}
	if blocks.Backend != nil {
		args.Backend = &hcl.Backend{
			Type:          blocks.Backend.BackendType(),
			Configuration: blocks.Backend,
		}
	}
	for i, prov := range blocks.Providers {
		args.Providers[i] = hcl.Provider{
			LocalName:     prov.LocalName(),
			Source:        prov.Source(),
			Version:       prov.Version(),
			Configuration: prov.Configuration(),
		}
	}
	for i, data := range blocks.DataResources {
		args.DataResources[i] = hcl.DataResource{
			DataSource:    data.DataSource(),
			LocalName:     data.LocalName(),
			Configuration: data.Configuration(),
		}
	}
	for i, res := range blocks.Resources {
		args.Resources[i] = hcl.Resource{
			Type:          res.Type(),
			LocalName:     res.LocalName(),
			Configuration: res.Configuration(),
		}
	}
	if err := hcl.Encode(w, args); err != nil {
		return err
	}
	return nil
}

// validateStack is an attempt to catch errors with the Terraform configuration before
// calling Terraform validate.
// The struct validate tags should handle most of the basic config validation but what we can
// check for additionally here are things like providers existing for resources.
// Future things to check for (TODO):
// 1. Each resource/data block's specific provider exists
func validateStack(sb *stackObjects) error {
	if (len(sb.Resources)+len(sb.DataResources)) > 0 && len(sb.Providers) == 0 {
		return ErrNoProviderBlock
	}
	return nil
}
