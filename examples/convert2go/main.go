package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	karpentercore "github.com/aws/karpenter-core/pkg/apis"
	karpenter "github.com/aws/karpenter/pkg/apis"
	"github.com/volvo-cars/lingon/pkg/kube"
	"golang.org/x/exp/slog"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

//go:generate go run main.go -in testdata/provisioner.yaml -out out/karpenter

func main() {
	var in, out, appName, pkgName string
	groupByKind := false
	removeAppName := false
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

	flag.Parse()

	slog.Info(
		"flags",
		slog.String("in", in),
		slog.String("out", out),
		slog.String("app", appName),
		slog.Bool("group", groupByKind),
		slog.Bool("clean-name", removeAppName),
	)

	if err := run(in, out, appName, groupByKind, removeAppName); err != nil {
		slog.Error("run", err)
		os.Exit(1)
	}

	slog.Info("done")
}

func run(in, out, appName string, groupByKind, removeAppName bool) error {
	_ = apiextensions.AddToScheme(scheme.Scheme)
	// ADD MORE CRDS HERE
	_ = karpenter.AddToScheme(scheme.Scheme)
	_ = karpentercore.AddToScheme(scheme.Scheme)

	decod := scheme.Codecs.UniversalDeserializer()

	opts := []kube.ImportOption{
		kube.WithAppName(appName),
		kube.WithPackageName(appName),
		kube.WithOutputDirectory(out),
		kube.WithSerializer(decod),
	}
	if groupByKind {
		opts = append(opts, kube.WithGroupByKind(true))
	}
	if removeAppName {
		opts = append(opts, kube.WithRemoveAppName(true))
	}

	// stdin
	if in == "-" {
		opts = append(opts, kube.WithReadStdIn())
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
		opts = append(opts, kube.WithManifestFiles([]string{in}))
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
	opts = append(opts, kube.WithManifestFiles(files))
	if err := kube.Import(opts...); err != nil {
		return fmt.Errorf("import: %w", err)
	}
	return nil
}
