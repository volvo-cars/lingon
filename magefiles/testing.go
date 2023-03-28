// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Test mg.Namespace

// Summary runs tests with gotestsum
func (Test) Summary() error {
	return sh.Run("gotestsum", "--", "go test ./...")
}

// Run tests in normal mode
func (Test) Default() error {
	return goTest("./...")
}

// Verbose runs tests in verbose mode
func (Test) Verbose() error {
	return goTest("-v", "./...")
}

// Race runs the test suite with the data race detector enabled.
func (Test) Race() error {
	return goTest("-race", "-v", "./...")
}
