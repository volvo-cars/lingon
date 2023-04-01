// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

/*
Terragen generates Go code for Terraform providers.
It accepts one Terraform provider and generates Go structs and
helper functions for the provider configuration,
resources and data sources for each provider.

Usage:

	gofmt [flags]

The flags are:

	-force
		override any existing generated Go files (required)
	-out string
		directory to generate Go files in (required)
	-pkg string
		Go pkg for the generated Go files (required)
	-provider value
		provider to generate Go files for (required), e.g. aws=hashicorp/aws:4.49.0
	-tfout string
		directory to generate Terraform providers schema (default ".lingon/schemas")
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

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
		v           bool
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
	flag.BoolVar(&v, "v", false, "show version")

	flag.Parse()

	if v {
		printVersion()
		return
	}
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

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func printVersion() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		_, _ = fmt.Fprintln(os.Stderr, "error reading build-info")
		os.Exit(1)
	}
	fmt.Printf("Build:\n%s\n", bi)
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Date: %s\n", date)
}
