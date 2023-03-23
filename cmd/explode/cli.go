package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/volvo-cars/lingon/pkg/kube"
	"golang.org/x/exp/slog"
)

func main() {
	var in, out string
	flag.StringVar(
		&in,
		"in",
		".",
		"specify the input directory of the yaml manifests, '-' for stdin",
	)
	flag.StringVar(
		&out,
		"out",
		"out",
		"specify the output directory for manifests",
	)

	flag.Parse()

	fmt.Println("in =", in, ", out =", out)

	if err := run(in, out); err != nil {
		slog.Error("run", err)
		os.Exit(1)
	}

	fmt.Println("done")
}

func run(in, out string) error {
	// stdin
	if in == "-" {
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
		slog.Error("list yaml files", err)
	}

	fmt.Println("files", files)
	for _, f := range files {
		fp, err := os.Open(f)
		if err != nil {
			return fmt.Errorf("file open %s: %w", f, err)
		}
		// explode them into individual files
		if err = kube.Explode(fp, out); err != nil {
			return fmt.Errorf("explode %s: %w", f, err)
		}
		_ = fp.Close()
	}

	return nil
}
