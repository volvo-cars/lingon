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
	"strings"
	"time"

	"golang.org/x/exp/slog"
)

const (
	// OSVScanner is the OSV Scanner to find vulnerabilities
	osvScannerRepo    = "github.com/google/osv-scanner/cmd/osv-scanner"
	osvScannerVersion = "@v1.3.2"

	// goVuln to find vulnerabilities
	vulnRepo    = "golang.org/x/vuln/cmd/govulncheck"
	vulnVersion = "@latest"
	goVuln      = vulnRepo + vulnVersion

	// syft is for generating SBOM
	syftRepo    = "github.com/anchore/syft/cmd/syft"
	syftVersion = "@v0.83.0"

	// goLicenses is Google's go-licenses to export all licenses
	goLicensesRepo    = "github.com/google/go-licenses"
	goLicensesVersion = "@v1.6.0"
	goLicenses        = goLicensesRepo + goLicensesVersion

	// goCILint is for linting code
	goCILintRepo    = "github.com/golangci/golangci-lint/cmd/golangci-lint"
	goCILintVersion = "@v1.53.2"
	goCILint        = goCILintRepo + goCILintVersion

	// copyWriteCheck is hashicorp/copywrite to check license headers
	copyWriteRepo    = "github.com/hashicorp/copywrite"
	copyWriteVersion = "@v0.16.3"

	// goFumpt is mvdan.cc/gofumpt to format code
	goFumptRepo    = "mvdan.cc/gofumpt"
	goFumptVersion = "@v0.5.0"
	goFumpt        = goFumptRepo + goFumptVersion
)

const (
	curDir = "."
	recDir = "./..."
)

var mod = func() string {
	if isUpdatingGoModOnly() {
		slog.Info("updating go.mod")
		return "-mod=mod"
	} else {
		slog.Info("readonly go.mod")
		return "-mod=readonly"
	}
}()

func main() {
	var cover, lint, doc, examples, fix, nodiff, pr, scan, release, update bool
	flag.BoolVar(&cover, "cover", false, "tests with coverage")
	flag.BoolVar(
		&lint,
		"lint",
		false,
		"linting and formatting code (gofumpt, golangci-lint)",
	)
	flag.BoolVar(&doc, "doc", false, "generate all docs and readme")
	flag.BoolVar(
		&examples,
		"examples",
		false,
		"tests all docs examples /!\\ slow without build cache",
	)
	flag.BoolVar(
		&fix,
		"fix",
		false,
		"same as -lint + generating notice and licenses headers",
	)
	flag.BoolVar(&nodiff, "nodiff", false, "error if git diff is not empty")
	flag.BoolVar(&pr, "pr", false, "run pull request checks: -fix + go test")
	flag.BoolVar(&scan, "scan", false, "scan for vulnerabilities")
	flag.BoolVar(&release, "release", false, "create a new release")
	flag.BoolVar(&update, "update", false, "update dependencies")

	flag.Parse()

	if update {
		Update()
	}
	if cover {
		CoverP()
	}
	if lint {
		Lint()
	}
	if doc {
		Doc()
	}
	if examples {
		DocExamples()
	}
	if fix {
		Fix()
	}
	if pr {
		PullRequest()
	}
	if scan {
		Scan()
	}
	if release {
		Release()
	}
	if nodiff {
		// should be last
		HasGitDiff()
	}
}

func iferr(err error) {
	if err != nil {
		panic(err)
	}
}

func Release() {
	d := time.Now().UTC().Format("2006-01-02")
	ssha, err := shortSha()
	iferr(err)
	v := d + "-" + ssha
	iferr(TagRelease(v, "Release "+v))
}

func TagRelease(tag, msg string) error {
	cmd := exec.Command("git", "tag", "-a", tag, "-s", "-m", msg)
	slog.Info("exec", slog.String("cmd", cmd.String()))
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	cmdgp := exec.Command("git", "push", "--tags")
	slog.Info("exec", slog.String("cmd", cmdgp.String()))
	_, err = cmdgp.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func shortSha() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	slog.Info("exec", slog.String("cmd", cmd.String()))
	b, err := cmd.CombinedOutput()
	s := strings.ReplaceAll(string(b), "\n", "")
	return s, err
}

func Cover() {
	coverOutput := "cover.out"
	coverMode := "count" // see `go help testflag` for more info
	iferr(
		Go(
			"test",
			recDir,
			"-coverprofile="+coverOutput,
			"-covermode="+coverMode,
		),
	)
	iferr(
		Go(
			"tool",
			"cover",
			"-func="+coverOutput,
			// "-html="+coverOutput,
			// "-o",
			// "cover.html",
		),
	)
	fmt.Println("✅ coverage generated: open cover.html to see results")
}

func Lint() {
	fmt.Println("🧹 code linting")
	iferr(Go("mod", "tidy"))
	iferr(Go("mod", "verify"))
	iferr(GoRun(goFumpt, "-w", "-extra", curDir))
	iferr(GoRun(goCILint, "-v", "run", recDir))
	fmt.Println("✅ code linted")
}

func Fix() {
	Lint()
	fmt.Println("📋 copywrite and licenses fix")
	iferr(CopyWriteFix())
	iferr(Notice())
	fmt.Println("✅ All fixes applied")
}

func PullRequest() {
	Fix()
	iferr(Go("test", "-v", recDir))
	fmt.Println("✅ pull request checks passed")
}

func Scan() {
	iferr(Sbom())
	iferr(GoRun(goVuln, recDir))
	iferr(OSVScanner())
	fmt.Println("✅ all scans completed")
}

func Doc() {
	fmt.Println("📝 generating docs")
	docRun("go", "generate", mod, recDir)
	docRun("go", "mod", "tidy")
	fmt.Println("✅ docs generated")
}

func DocExamples() {
	Doc()
	fmt.Println("📝 testing examples")
	docRun("go", "test", mod, "-v", recDir)
	fmt.Println("✅ docs generated and examples tested")
}

func Go(args ...string) error {
	cmd := exec.Command("go", args...)
	slog.Info("exec", slog.String("cmd", cmd.String()))
	defer slog.Info("done", slog.String("cmd", cmd.String()))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		_ = os.Stderr.Sync()
		_ = os.Stdout.Sync()
		return fmt.Errorf("go: %s", err)
	}
	return nil
}

func GoRun(args ...string) error {
	return Go(append([]string{"run", mod}, args...)...)
}

// CopyWriteFix is hashicorp/copywrite to fix license headers
func CopyWriteFix() error {
	return GoRun(
		copyWriteRepo+copyWriteVersion,
		"headers",
		"--dirPath", "./",
		"--config", "./.copywrite.hcl",
	)
}

// HasGitDiff displays the git diff and errors if there is a diff
func HasGitDiff() {
	cmd := exec.Command("git", "--no-pager", "diff")
	slog.Info("exec", slog.String("cmd", cmd.String()))
	b, err := cmd.CombinedOutput()
	iferr(err)
	if len(b) == 0 {
		return
	}
	buf := bytes.NewBuffer(b)
	fmt.Println(buf.String())
	panic("git diff is not empty")
}

// docRun runs a command in the docs directory
func docRun(args ...string) {
	cmd := exec.Command(args[0], args[1:]...) //nolint:gosec
	slog.Info("exec", slog.String("cmd", cmd.String()))
	defer slog.Info("exec", slog.String("cmd", args[0]+" done"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "./docs"
	if err := cmd.Run(); err != nil {
		_ = os.Stderr.Sync()
		_ = os.Stdout.Sync()
		panic(err)
	}
}

// OSVScanner is the OSV Scanner to find vulnerabilities
func OSVScanner() error {
	slog.Info("running OSV Scanner")
	defer slog.Info("DONE OSV Scanner")
	// return GoRun(osvScannerRepo+osvScannerVersion, "-r", curDir)
	// not scanning docs/go.mod because of github.com/aws/aws-sdk-go
	// and the osvScanner returns an error when a vulnerability is detected
	return GoRun(osvScannerRepo+osvScannerVersion, curDir)
}

// Sbom is used to generate SBOM
func Sbom() error {
	slog.Info("running syft - generating SBOM")
	defer slog.Info("DONE syft - generating SBOM")
	cmd := exec.Command(
		"go", "run",
		syftRepo+syftVersion,
		"packages",
		"dir:"+curDir,
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

// Notice is used to generate a NOTICE file
func Notice() error {
	slog.Info("running go-licenses - generating report")
	defer slog.Info("DONE go-licenses - generating report")
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(
		"go", "run",
		goLicenses,
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
