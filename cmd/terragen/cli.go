// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	"github.com/volvo-cars/lingon/pkg/terragen"
	"golang.org/x/exp/slog"
)

func main() {
	var (
		outDir      string
		tfOutDir    string
		pkgPath     string
		providerStr string
		force       bool
	)

	flag.StringVar(
		&tfOutDir,
		"tfout",
		filepath.Join(".lingon", "schemas"),
		"directory to generate Terraform provider schema",
	)
	flag.StringVar(&outDir, "out", "", "directory to generate Go files in")
	flag.StringVar(&pkgPath, "pkg", "", "Go pkg for the generated Go files")
	flag.StringVar(
		&providerStr,
		"provider",
		"",
		"provider to generate Go files for, e.g. aws=hashicorp/aws:4.49.0",
	)
	flag.BoolVar(
		&force,
		"force",
		false,
		"override any existing generated Go files",
	)
	flag.Parse()

	if outDir == "" {
		slog.Error("-out flag required")
		os.Exit(1)
	}
	if pkgPath == "" {
		slog.Error("-pkg flag required")
		os.Exit(1)
	}

	if providerStr == "" {
		slog.Error("-provider flag required")
		os.Exit(1)
	}

	provider, err := terragen.ParseProvider(providerStr)
	if err != nil {
		slog.Error("invalid provider", "err", err)
		os.Exit(1)
	}
	ctx := context.Background()
	slog.Info(
		"Generating Terraform provider schema",
		slog.String("provider", providerStr),
		slog.String("out", tfOutDir),
	)
	schemas, err := terragen.GenerateProviderSchema(ctx, provider)
	if err != nil {
		slog.Error("generating provider schema", "err", err)
		os.Exit(1)
	}

	slog.Info(
		"Generating Terraform Go wrappers",
		slog.String("provider", providerStr),
		slog.String("out", outDir),
		slog.String("pkg", pkgPath),
	)
	if err := terragen.GenerateGoCode(
		terragen.GenerateGoArgs{
			ProviderName:    provider.Name,
			ProviderSource:  provider.Source,
			ProviderVersion: provider.Version,
			OutDir:          outDir,
			PkgPath:         pkgPath,
			Force:           force,
		},
		schemas,
	); err != nil {
		slog.Error("generating Go wrapper", "err", err)
		os.Exit(1)
	}
}
