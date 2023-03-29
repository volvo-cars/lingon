// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terragen

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/volvo-cars/lingon/pkg/internal/terrajen"

	tfjson "github.com/hashicorp/terraform-json"
	"golang.org/x/exp/slog"
)

var (
	ErrPackageLocationNotEmpty = errors.New("providers pkg location is not empty")
	ErrProviderSchemaNotFound  = errors.New("provider schema not found")
)

type GenerateGoArgs struct {
	ProviderName    string
	ProviderSource  string
	ProviderVersion string
	// OutDir is the filesystem location where the generated files will be created.
	OutDir string
	// PkgPath is the Go pkg path to the generated files location, specified by OutDir.
	// E.g. if OutDir is in a module called "my-module" in a directory called "gen",
	// then the PkgPath should be "my-module/gen".
	PkgPath string
	// Force enables overriding any existing generated files per-provider.
	Force bool
}

// GenerateGoCode generates Go code for creating Terraform objects for the given
// providers and their schemas.
func GenerateGoCode(
	args GenerateGoArgs,
	schemas *tfjson.ProviderSchemas,
) error {
	if args.OutDir == "" {
		return errors.New("outDir is empty")
	}
	if args.PkgPath == "" {
		return errors.New("pkgPath is empty")
	}

	slog.Info(
		"Generating Go wrapper",
		slog.String("provider", args.ProviderName),
		slog.String("source", args.ProviderSource),
		slog.String("version", args.ProviderVersion),
	)
	tfArgs := terrajen.ProviderGenerator{
		GoProviderPkgPath:        args.PkgPath,
		GeneratedPackageLocation: args.OutDir,
		ProviderName:             args.ProviderName,
		ProviderSource:           args.ProviderSource,
		ProviderVersion:          args.ProviderVersion,
	}

	pSchema, ok := schemas.Schemas[args.ProviderSource]
	if !ok {
		// Try adding registry.terraform.io/ prefix if not already added
		if !strings.HasPrefix(args.ProviderSource, "registry.terraform.io/") {
			pSchema, ok = schemas.Schemas[fmt.Sprintf(
				"registry.terraform.io/%s",
				args.ProviderSource,
			)]
		}
		// If still not ok, indicate an error
		if !ok {
			return fmt.Errorf(
				"provider source: %s: %w",
				args.ProviderSource,
				ErrProviderSchemaNotFound,
			)
		}
	}
	err := generateProvider(args, tfArgs, pSchema)
	if err != nil {
		return err
	}

	return nil
}

func generateProvider(
	genArgs GenerateGoArgs,
	provider terrajen.ProviderGenerator,
	providerSchema *tfjson.ProviderSchema,
) error {
	if err := createDirIfNotEmpty(
		provider.GeneratedPackageLocation,
		genArgs.Force,
	); err != nil {
		return fmt.Errorf(
			"creating providers pkg directory %s: %w",
			provider.GeneratedPackageLocation,
			err,
		)
	}
	//
	// Generate Provider
	//
	ps := provider.SchemaProvider(providerSchema.ConfigSchema.Block)
	f := terrajen.ProviderFile(ps)
	if err := f.Save(ps.FilePath); err != nil {
		terrajen.JenDebug(err)
		return fmt.Errorf("saving provider file %s: %w", ps.FilePath, err)
	}

	subPkgFile, ok := terrajen.SubPkgFile(ps)
	if ok {
		subPkgDir := filepath.Dir(ps.SubPkgPath())
		if err := os.MkdirAll(subPkgDir, os.ModePerm); err != nil {
			return fmt.Errorf(
				"creating sub package directory %s: %w",
				subPkgDir,
				err,
			)
		}
		if err := subPkgFile.Save(ps.SubPkgPath()); err != nil {
			terrajen.JenDebug(err)
			return fmt.Errorf(
				"saving sub package file %s: %w",
				ps.SubPkgPath(),
				err,
			)
		}
	}
	//
	// Generate Resources
	//
	for name, resource := range providerSchema.ResourceSchemas {
		rs := provider.SchemaResource(name, resource.Block)
		rsf := terrajen.ResourceFile(rs)
		if err := rsf.Save(rs.FilePath); err != nil {
			terrajen.JenDebug(err)
			return fmt.Errorf(
				"saving resource file %s: %w",
				rs.FilePath,
				err,
			)
		}

		rsSubPkgFile, ok := terrajen.SubPkgFile(rs)
		if !ok {
			continue
		}
		subPkgDir := filepath.Dir(rs.SubPkgPath())
		if err := os.MkdirAll(subPkgDir, os.ModePerm); err != nil {
			return fmt.Errorf(
				"creating sub package directory %s: %w",
				subPkgDir,
				err,
			)
		}
		if err := rsSubPkgFile.Save(rs.SubPkgPath()); err != nil {
			terrajen.JenDebug(err)
			return fmt.Errorf(
				"saving sub package file %s: %w",
				rs.SubPkgPath(),
				err,
			)
		}
	}

	//
	// Generate Data blocks
	//
	for name, data := range providerSchema.DataSourceSchemas {
		ds := provider.SchemaData(name, data.Block)
		df := terrajen.DataSourceFile(ds)
		if err := df.Save(ds.FilePath); err != nil {
			terrajen.JenDebug(err)
			return fmt.Errorf("saving data file %s: %w", ds.FilePath, err)
		}

		dataSubPkgFile, ok := terrajen.SubPkgFile(ds)
		if !ok {
			continue
		}
		subPkgDir := filepath.Dir(ds.SubPkgPath())
		if err := os.MkdirAll(subPkgDir, os.ModePerm); err != nil {
			return fmt.Errorf(
				"creating sub package directory %s: %w",
				subPkgDir,
				err,
			)
		}
		if err := dataSubPkgFile.Save(ds.SubPkgPath()); err != nil {
			terrajen.JenDebug(err)
			return fmt.Errorf(
				"saving sub package file %s: %w",
				ds.SubPkgPath(),
				err,
			)
		}
	}
	return nil
}

func createDirIfNotEmpty(path string, force bool) error {
	f, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		// Create the directory
		if err := os.MkdirAll(
			path,
			os.ModePerm,
		); err != nil {
			return err
		}
		return nil
	}

	_, readErr := f.Readdirnames(1)
	if readErr != nil {
		if readErr == io.EOF {
			return nil
		}
		return err
	}
	// The directory is not empty. If force flag is provided, clean the directory, else error
	if !force {
		return ErrPackageLocationNotEmpty
	}
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("cleaning directory: %w", err)
	}
	// Create the directory again now that it's gone
	return createDirIfNotEmpty(path, false)
}

// ParseProvider takes a provider as a string and returns a Provider object.
// An error is returned if the string could not be parsed.
// Example provider: aws=hashicorp/aws:4.60.0
func ParseProvider(s string) (Provider, error) {
	pMap := strings.SplitN(s, "=", 2)
	if len(pMap) == 1 {
		return Provider{}, fmt.Errorf("provider format incorrect: missing `=`")
	}
	p := Provider{
		Name: pMap[0],
	}
	sourceVersion := strings.SplitN(pMap[1], ":", 2)
	if len(sourceVersion) == 1 {
		return Provider{}, fmt.Errorf(
			"provider format incorrect: missing `:` in `source:version`",
		)
	}
	p.Source = sourceVersion[0]
	p.Version = sourceVersion[1]
	// Add the provider to the map
	return p, nil
}
