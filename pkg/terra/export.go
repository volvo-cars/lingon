// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/volvo-cars/lingon/pkg/internal/hcl"
)

const (
	tfSuffix = ".tf"
)

// ExportOption is used to configure the conversion from Go code to Terraform
// configurations.
// Use the helper functions WithExportXXX to configure the export.
type ExportOption func(*gotf)

// WithExportWriter writes the generated Terraform configuration to [io.Writer].
func WithExportWriter(w io.Writer) ExportOption {
	return func(g *gotf) {
		g.useWriter = true
		g.w = w
	}
}

// WithExportOutputDirectory writes the generated Terraform configuration to
// the given output directory.
func WithExportOutputDirectory(dir string) ExportOption {
	return func(g *gotf) {
		g.dir = dir
	}
}

type gotf struct {
	useWriter bool
	w         io.Writer

	dir string
}

// Export encodes [Exporter] to Terraform configurations
func Export(stack Exporter, opts ...ExportOption) error {
	var g gotf
	for _, o := range opts {
		o(&g)
	}

	if g.useWriter {
		if err := encodeStack(stack, g.w); err != nil {
			return fmt.Errorf(
				"encoding stack: %w", err,
			)
		}
		return nil
	}

	if g.dir == "" {
		return errors.New("output directory is empty")
	}

	if err := os.MkdirAll(g.dir, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(
		filepath.Join(
			g.dir,
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
			DependsOn:     res.Dependencies(),
			Lifecycle:     res.LifecycleManagement(),
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
