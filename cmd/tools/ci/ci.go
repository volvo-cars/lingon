// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"golang.org/x/exp/slog"
)

const (
	// OSVScanner is the OSV Scanner to find vulnarabilities
	osvScannerRepo    = "github.com/google/osv-scanner/cmd/osv-scanner"
	osvScannerVersion = "@v1"

	// vuln is the GoVulnCheck to find vulnarabilities
	vulnRepo    = "golang.org/x/vuln/cmd/govulncheck"
	vulnVersion = "@latest"

	// syft is the Syft to generate SBOM
	syftRepo    = "github.com/anchore/syft/cmd/syft"
	syftVersion = "@v0.79.0"

	// goLicenses is Google's go-licenses to export all licenses
	goLicensesRepo    = "github.com/google/go-licenses"
	goLicensesVersion = "@v1.6.0"

	// goCILint is golangci/golangci-lint to lint code
	goCILintRepo    = "github.com/golangci/golangci-lint/cmd/golangci-lint"
	goCILintVersion = "@v1.52.2"

	// copyWriteCheck is hashicorp/copywrite to check license headers
	copyWriteRepo    = "github.com/hashicorp/copywrite"
	copyWriteVersion = "@v0.16.3"

	// goFumpt is mvdan.cc/gofumpt to format code
	goFumptRepo    = "mvdan.cc/gofumpt"
	goFumptVersion = "@v0.5.0"
)

const (
	dir    = "."
	recDir = "./..."
)

func main() {
	var fix, scan, verify bool
	flag.BoolVar(&fix, "fix", false, "Run all fixes")
	flag.BoolVar(&scan, "scan", false, "Run all scans")
	flag.BoolVar(&verify, "verify", false, "Run all verifications")
	flag.Parse()
	if fix {
		Fix()
	}
	if scan {
		Scan()
	}
	if verify {
		GoModVerify()
	}
}

type run func() error

func g0(args ...string) error {
	cmd := exec.Command("go", args...)
	slog.Info("exec", slog.String("cmd", cmd.String()))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go: %s", err)
	}
	return nil
}

func goRun(args ...string) error {
	return g0(append([]string{"run"}, args...)...)
}

func check(f run) {
	if err := f(); err != nil {
		panic(err)
	}
}

// Fix is used to run all fixes sequentially
func Fix() {
	slog.Info("running all fixes")
	defer slog.Info("DONE all fixes")
	check(Tidy)
	check(GoFumpt)
	check(CopyWriteFix)
	check(GoCILint)
}

func Scan() {
	slog.Info("running all scans")
	defer slog.Info("DONE all scans")
	check(Syft)
	check(GoVulnCheck)
	check(Notice)
	check(GoLicensesCheck)
	check(CopyWriteCheck)
	// check(OSVScanner) // because: https://osv.dev/GO-2022-0646 github.com/aws/aws-sdk-go │ 1.44.246 │ docs/go.mod
}

func GoModVerify() {
	check(CheckGoMod)
	check(Tidy)
	check(HasGitDiff)
}

// OSVScanner is the OSV Scanner to find vulnerabilities
func OSVScanner() error {
	slog.Info("running OSV Scanner")
	return goRun(osvScannerRepo+osvScannerVersion, "-r", dir)
}

// GoVulnCheck is the GoVulnCheck to find vulnerabilities
func GoVulnCheck() error {
	slog.Info("running govulncheck")
	defer slog.Info("DONE govulncheck")
	return goRun(vulnRepo+vulnVersion, recDir)
}

// Syft is used to generate SBOM
func Syft() error {
	slog.Info("running syft - generating SBOM")
	defer slog.Info("DONE syft - generating SBOM")
	cmd := exec.Command(
		"go", "run",
		syftRepo+syftVersion,
		"packages",
		"dir:"+dir,
		"-o=spdx-json",
		"--file=bin/sbom.json",
	)
	slog.Info("exec", slog.String("cmd", cmd.String()))
	cmd.Stdout = io.Discard
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "SYFT_QUIET=true")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go: %s", err)
	}

	return nil
}

// GoFumpt is used to format code
func GoFumpt() error {
	slog.Info("running gofumpt - formatting code")
	defer slog.Info("DONE gofumpt - formatting code")
	return goRun(goFumptRepo+goFumptVersion, "-w", "-extra", dir)
}

// Notice is used to generate a NOTICE file
func Notice() error {
	slog.Info("running go-licenses - generating report")
	defer slog.Info("DONE go-licenses - generating report")
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(
		"go", "run",
		goLicensesRepo+goLicensesVersion,
		"report",
		recDir,
		"--template=./cmd/tools/ci/licenses.tpl",
		"--ignore=github.com/volvo-cars/lingon",
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	slog.Info("exec", slog.String("cmd", cmd.String()))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"running go-licenses: %w:\n\n%s", err,
			stderr.String(),
		)
	}

	noticeFile, err := os.Create("NOTICE")
	if err != nil {
		return fmt.Errorf("creating NOTICE file: %w", err)
	}

	if err = GenerateNotice(noticeFile, &stdout); err != nil {
		return fmt.Errorf("generating NOTICE file: %w", err)
	}
	return nil
}

// GoLicensesCheck is used to check all licenses
func GoLicensesCheck() error {
	slog.Info("running go-licenses - exporting licenses")
	defer slog.Info("DONE go-licenses - exporting licenses")
	return goRun(
		goLicensesRepo+goLicensesVersion,
		"check", recDir,
	)
}

// CopyWriteCheck is hashicorp/copywrite to check license headers
func CopyWriteCheck() error {
	slog.Info("running copywrite - checking license headers")
	defer slog.Info("DONE copywrite - checking license headers")
	return goRun(
		copyWriteRepo+copyWriteVersion,
		"headers",
		"--dirPath", "./",
		"--config", "./.copywrite.hcl",
		"--plan",
	)
}

// CopyWriteFix is hashicorp/copywrite to fix license headers
func CopyWriteFix() error {
	slog.Info("running copywrite - fixing license headers")
	defer slog.Info("DONE copywrite - fixing license headers")
	return goRun(
		copyWriteRepo+copyWriteVersion,
		"headers",
		"--dirPath", "./",
		"--config", "./.copywrite.hcl",
	)
}

// GoCILint is used to lint code
func GoCILint() error {
	slog.Info("running golangci-lint - linting code")
	defer slog.Info("DONE golangci-lint - linting code")
	return goRun(goCILintRepo+goCILintVersion, "-v", "run", recDir)
}

func HasGitDiff() error {
	slog.Info("running git diff")
	defer slog.Info("DONE git diff")
	cmd := exec.Command("git", "--no-pager", "diff")
	slog.Info("exec", slog.String("cmd", cmd.String()))
	b, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}
	buf := bytes.NewBuffer(b)
	return fmt.Errorf("running git diff:\n\n%s", buf.String())
}

func CheckGoMod() error {
	return g0("mod", "verify")
}

// Tidy runs go mod tidy
func Tidy() error {
	return g0("mod", "tidy")
}
