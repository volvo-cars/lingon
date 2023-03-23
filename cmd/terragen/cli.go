package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/volvo-cars/lingon/pkg/terragen"
	"golang.org/x/exp/slog"
)

func main() {
	var (
		outDir    string
		tfOutDir  string
		pkgPath   string
		providers providerFlags = map[string]terragen.Provider{}
		force     bool
	)

	flag.StringVar(
		&tfOutDir,
		"tfout",
		filepath.Join(".terriyaki", "schemas"),
		"directory to generate Terraform providers schema",
	)
	flag.StringVar(&outDir, "out", "", "directory to generate Go files in")
	flag.StringVar(&pkgPath, "pkg", "", "Go pkg for the generated Go files")
	flag.Var(
		&providers,
		"provider",
		"providers to generate Go files for, e.g. aws=hashicorp/aws:4.49.0",
	)
	flag.BoolVar(
		&force,
		"force",
		false,
		"override any existing generated Go files",
	)
	flag.Parse()

	if outDir == "" {
		slog.Error("-out flag required", nil)
		os.Exit(1)
	}
	if pkgPath == "" {
		slog.Error("-pkg flag required", nil)
		os.Exit(1)
	}

	if len(providers) == 0 {
		slog.Error("-provider flag required", nil)
		os.Exit(1)
	}

	ctx := context.Background()
	slog.Info(
		"Generating Terraform providers schema",
		slog.String("providers", providers.String()),
		slog.String("out", tfOutDir),
	)
	schemas, err := terragen.GenerateProvidersSchema(ctx, providers)
	if err != nil {
		slog.Error("generating providers schema", err)
		os.Exit(1)
	}

	slog.Info(
		"Generating Terraform Go wrappers",
		slog.String("providers", providers.String()),
		slog.String("out", outDir),
		slog.String("pkg", pkgPath),
	)
	if err := terragen.GenerateGoCode(
		terragen.GenerateGoArgs{
			OutDir:  outDir,
			PkgPath: pkgPath,
			Force:   force,
		},
		providers,
		schemas,
	); err != nil {
		slog.Error("generating Go wrapper", err)
		os.Exit(1)
	}
}

var _ flag.Value = (*providerFlags)(nil)

// providerFlags implements the flag.Value interface to decode Terraform providers in the form
// localName=source:version. E.g. aws=hashicorp/aws:4.49.0
type providerFlags map[string]terragen.Provider

func (f *providerFlags) String() string {
	provList := make([]string, 0, len(*f))
	for name, prov := range *f {
		provList = append(
			provList,
			fmt.Sprintf("%s=%s:%s", name, prov.Source, prov.Version),
		)
	}
	return strings.Join(provList, ",")
}

func (f *providerFlags) Set(value string) error {
	pMap := strings.SplitN(value, "=", 2)
	if len(pMap) == 1 {
		return fmt.Errorf("provider format incorrect: missing `=`")
	}
	localName := pMap[0]
	sourceVersion := strings.SplitN(pMap[1], ":", 2)
	if len(sourceVersion) == 1 {
		return fmt.Errorf("provider format incorrect: missing `:` in `source:version`")
	}
	source := sourceVersion[0]
	version := sourceVersion[1]
	if _, ok := (*f)[localName]; ok {
		return fmt.Errorf("duplicate provider local name provided")
	}
	// Add the provider to the map
	(*f)[localName] = terragen.Provider{
		Source:  source,
		Version: version,
	}
	return nil
}
