// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

//go:build tools

// This file imports packages that are used when running go generate, or used
// during the development process but not otherwise depended on by built code.
package main

import (
	_ "github.com/anchore/syft/cmd/syft"                    // generate sbom
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint" // linting
	_ "github.com/google/go-licenses"                       // generate licenses
	_ "github.com/hashicorp/copywrite"                      // check license headers
	_ "golang.org/x/tools/cmd/goimports"                    // format code
	_ "golang.org/x/vuln/cmd/govulncheck"                   // check for vulnerabilities
	_ "gotest.tools/gotestsum"                              // run tests with formatted output
	_ "mvdan.cc/gofumpt"                                    // format code
	// having issues with this one's dependencies
	// _ "github.com/google/osv-scanner/cmd/osv-scanner"       // check for vulnerabilities
)
