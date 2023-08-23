// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"

	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsbeta "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	kubescheme "k8s.io/client-go/kubernetes/scheme"
)

const crdMsg = "IF there is an issue with CRDs. Please visit this page to solve it https://github.com/volvo-cars/lingon/tree/main/docs/kubernetes/crd"

func main() {
	var in, out, appName, pkgName string
	var version, verbose, ignoreErr bool

	groupByKind := true
	removeAppName := true
	flag.StringVar(
		&in,
		"in",
		"-",
		"specify the input directory of the yaml manifests, '-' for stdin",
	)
	flag.StringVar(
		&out,
		"out",
		"out",
		"specify the output directory for manifests.",
	)
	flag.StringVar(
		&appName,
		"app",
		"myapp",
		"specify the app name. This will be used as the package name if none is specified.",
	)
	flag.StringVar(
		&pkgName,
		"pkg",
		"",
		"specify the package name. If none is specified the app name will be used. Cannot contain a dash.",
	)
	flag.BoolVar(
		&groupByKind,
		"group",
		true,
		"specify if the output should be grouped by kind (default) or split by name.",
	)
	flag.BoolVar(
		&removeAppName,
		"clean-name",
		true,
		"specify if the app name should be removed from the variable, struct and file name.",
	)
	flag.BoolVar(&version, "version", false, "show version")
	flag.BoolVar(&verbose, "v", false, "show logs")
	flag.BoolVar(
		&ignoreErr,
		"ignore-errors",
		false,
		"ignore errors, useful to generate as much as possible",
	)
	flag.Parse()

	if version {
		printVersion()
		return
	}

	if pkgName == "" {
		pkgName = strings.ReplaceAll(appName, "-", "")
	}

	slog.Info(
		"flags",
		slog.String("in", in),
		slog.String("out", out),
		slog.String("app", appName),
		slog.Bool("group", groupByKind),
		slog.Bool("clean-name", removeAppName),
		slog.Bool("verbose", verbose),
		slog.Bool("ignore-errors", ignoreErr),
	)

	if err := run(
		in,
		out,
		appName,
		pkgName,
		groupByKind,
		removeAppName,
		verbose,
		ignoreErr,
	); err != nil {
		slog.Error(
			"run",
			slog.Any("error", err),
			slog.String("CRD", crdMsg),
		)
		os.Exit(1)
	}

	slog.Info("done")
}

func defaultSerializer() runtime.Decoder {
	// ADD MORE CRDS HERE
	_ = apiextensions.AddToScheme(kubescheme.Scheme)
	_ = apiextensionsbeta.AddToScheme(kubescheme.Scheme)
	return kubescheme.Codecs.UniversalDeserializer()
}

func run(
	in, out, appName, pkgName string,
	groupByKind, removeAppName, verbose, ignoreErr bool,
) error {
	opts := []kube.ImportOption{
		kube.WithImportAppName(appName),
		kube.WithImportPackageName(pkgName),
		kube.WithImportOutputDirectory(out),
		kube.WithImportSerializer(defaultSerializer()),
	}
	opts = append(opts, kube.WithImportGroupByKind(groupByKind))
	opts = append(opts, kube.WithImportRemoveAppName(removeAppName))
	opts = append(opts, kube.WithImportVerbose(verbose))
	opts = append(opts, kube.WithImportIgnoreErrors(ignoreErr))

	// stdin
	if in == "-" {
		opts = append(opts, kube.WithImportReadStdIn())
		if err := kube.Import(opts...); err != nil {
			return fmt.Errorf("import: %w", err)
		}
		return nil
	}

	// single file
	fi, err := os.Stat(in)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		opts = append(opts, kube.WithImportManifestFiles([]string{in}))
		if err := kube.Import(opts...); err != nil {
			return fmt.Errorf("import: %w", err)
		}
		return nil
	}

	// directory
	files, err := kubeutil.ListYAMLFiles(in)
	if err != nil {
		slog.Error("list yaml files", err)
	}

	fmt.Printf("files:\n- %s\n", strings.Join(files, "\n- "))
	opts = append(opts, kube.WithImportManifestFiles(files))
	if err := kube.Import(opts...); err != nil {
		return fmt.Errorf("import: %w", err)
	}
	return nil
}

var (
	ver    = "dev"
	commit = "none"
	date   = "unknown"
)

func printVersion() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		_, _ = fmt.Fprintln(os.Stderr, "error reading build-info")
		os.Exit(1)
	}
	fmt.Printf("Build:\n%s\n", bi)
	fmt.Printf("Version: %s\n", ver)
	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Date: %s\n", date)
}
