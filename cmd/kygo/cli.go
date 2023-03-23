package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/volvo-cars/lingon/pkg/kube"
	"golang.org/x/exp/slog"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

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
	var serializer runtime.Decoder

	_ = apiextensions.AddToScheme(scheme.Scheme)
	slog.Info("CRD", slog.String("group name", apiextensions.GroupName))
	// ADD MORE CRDS HERE

	serializer = scheme.Codecs.UniversalDeserializer()

	opts := []kube.ImportOption{
		kube.WithAppName(appName),
		kube.WithPackageName(appName),
		kube.WithOutputDirectory(out),
		kube.WithSerializer(serializer),
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
