// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/volvo-cars/lingon/pkg/kube"
)

func main() {
	var in, out string
	var v bool
	flag.StringVar(
		&in,
		"in",
		"-",
		"input directory to (recursively) process yaml manifests (default: '-' for stdin)",
	)
	flag.StringVar(
		&out,
		"out",
		"out",
		"output directory to write split manifests to (default 'out')",
	)

	flag.BoolVar(&v, "v", false, "show version")

	flag.Parse()

	if v {
		printVersion()
		return
	}

	fmt.Println("in =", in, ", out =", out)

	if err := run(in, out); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "run: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("done")
}

func run(in, out string) error {
	// stdin
	if in == "-" {
		fmt.Println("reading from stdin")
		if err := kube.Explode(os.Stdin, out); err != nil {
			return fmt.Errorf("from stdin: %w", err)
		}
		return nil
	}

	// single file
	fi, err := os.Stat(in)
	if err != nil {
		return fmt.Errorf("checking %s: %w", in, err)
	}
	if !fi.IsDir() {
		fmt.Println("reading from file", in)
		if filepath.Ext(in) != ".yaml" && filepath.Ext(in) != ".yml" {
			return fmt.Errorf("file %s is not a yaml file", in)
		}
		fp, err := os.Open(in)
		if err != nil {
			return fmt.Errorf("file open %s: %w", in, err)
		}
		defer fp.Close()
		// explode them into individual files
		if err = kube.Explode(fp, out); err != nil {
			return fmt.Errorf("explode %s: %w", in, err)
		}
		return nil
	}

	// directory
	files, err := kube.ListYAMLFiles(in)
	if err != nil {
		return fmt.Errorf("list yaml files: %w", err)
	}

	fmt.Println("files:", len(files))
	for i, f := range files {
		fp, err := os.Open(f)
		if err != nil {
			return fmt.Errorf("file open %s: %w", f, err)
		}
		fmt.Printf("%02d: %s\n", i+1, f)
		// explode them into individual files
		if err = kube.Explode(fp, out); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "explode %s: %v'\n", f, err)
		}
		_ = fp.Close()
	}

	return nil
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
