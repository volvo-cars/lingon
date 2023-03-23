package terragen

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
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
	// OutDir is the filesystem location where the generated files will be created.
	// The path will be suffixed with /providers/<local-name-of-provider>.
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
	providers map[string]Provider,
	schemas *tfjson.ProviderSchemas,
) error {
	if args.OutDir == "" {
		return errors.New("outDir is empty")
	}
	if args.PkgPath == "" {
		return errors.New("pkgPath is empty")
	}

	for providerName, provider := range providers {
		slog.Info(
			"Generating Go wrapper",
			slog.String("provider", providerName),
			slog.String("source", provider.Source),
			slog.String("version", provider.Version),
		)
		tfArgs := terrajen.ProviderGenerator{
			GoProvidersPkgPath: path.Join(args.PkgPath, "providers"),
			GeneratedPackageLocation: filepath.Join(
				args.OutDir,
				"providers",
				providerName,
			),
			ProviderName:    providerName,
			ProviderSource:  provider.Source,
			ProviderVersion: provider.Version,
		}

		pSchema, ok := schemas.Schemas[provider.Source]
		if !ok {
			// Try adding registry.terraform.io/ prefix if not already added
			if !strings.HasPrefix(provider.Source, "registry.terraform.io/") {
				pSchema, ok = schemas.Schemas[fmt.Sprintf(
					"registry.terraform.io/%s",
					provider.Source,
				)]
			}
			// If still not ok, indicate an error
			if !ok {
				return fmt.Errorf(
					"provider source: %s: %w",
					provider.Source,
					ErrProviderSchemaNotFound,
				)
			}
		}
		err := generateProvider(args, tfArgs, pSchema)
		if err != nil {
			return err
		}
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
