// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/volvo-cars/lingon/pkg/kube"
	"golang.org/x/exp/slog"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsbeta "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	kubescheme "k8s.io/client-go/kubernetes/scheme"
)

func main() {
	var in, out, appName, pkgName string
	var v bool

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
	flag.BoolVar(&v, "v", false, "show version")

	flag.Parse()

	if v {
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
	)

	if err := run(
		in,
		out,
		appName,
		pkgName,
		groupByKind,
		removeAppName,
	); err != nil {
		slog.Error("run", slog.Any("error", err))
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
	groupByKind, removeAppName bool,
) error {
	opts := []kube.ImportOption{
		kube.WithImportAppName(appName),
		kube.WithImportPackageName(pkgName),
		kube.WithImportOutputDirectory(out),
		kube.WithImportSerializer(defaultSerializer()),
	}
	if groupByKind {
		opts = append(opts, kube.WithImportGroupByKind(true))
	}
	if removeAppName {
		opts = append(opts, kube.WithImportRemoveAppName(true))
	}

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
	files, err := kube.ListYAMLFiles(in)
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
