// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

//go:build mage

package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"golang.org/x/exp/slog"
)

type Build mg.Namespace

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build.All

// in a helper file somewhere
var (
	g0      = sh.RunCmd("go") // go is a keyword :(
	goRun   = sh.RunCmd("go", "run")
	goBuild = sh.RunCmd("go", "build")
	goTest  = sh.RunCmd("go", "test")
)

// CI runs the CI build
func (Build) CI() {
	slog.Info("running CI build")
	mg.Deps(
		Clean,
		Run.GoFumpt,
		Test.Verbose,
	)
}

// GoBuild builds Go binaries
// This function is needed as mg.Deps() only accepts
//   - func()
//   - func() error
//   - func(context.Context)
//   - func(context.Context) error
func GoBuild(args ...string) func() error {
	return func() error {
		return goBuild(args...)
	}
}

// All runs all the build steps
func (Build) All() {
	slog.Info("running all build steps")
	mg.SerialDeps(Clean, Run.GoFumpt, Test.Verbose, Build.Local)
}

// Generate runs go generate ./...
func (Build) Generate() error {
	slog.Info("generating code")
	return g0("generate", "./...")
}

// Local builds the CLI binaries
func (Build) Local() error {
	slog.Info("building binaries")
	// os.Environ()
	mg.Deps(
		GoBuild("-o", "bin/kygo", "./cmd/kygo"),
		GoBuild("-o", "bin/explode", "./cmd/explode"),
		GoBuild("-o", "bin/terragen", "./cmd/terragen"),
	)
	return nil
}

// Removes built files
func Clean() error {
	slog.Info("cleaning up")
	if err := os.RemoveAll("./bin/"); err != nil {
		return err
	}
	return nil
}

// githash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
