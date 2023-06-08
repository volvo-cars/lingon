// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package testutil_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rogpeppe/go-internal/txtar"
	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	"k8s.io/apimachinery/pkg/runtime"
	"mvdan.cc/gofumpt/format"
)

// WIP

// TODO: bells and whistles
// 1. have a utility to download the latest manifest releases yaml files
// --> How did we get the YAML to begin with? what is the starting point ?
// 2. export the existing app in Go to YAML files
// 3. Smart diff the two sets of YAML files
//  --> with small checks to what changes
//    * the numbers of configMaps, deployments, CRDs, etc
//    * in PodTemplates, the number of containers, the number of volumes, env vars, etc
//    * Services and Ingresses -> ports
//    * RBAC roles and bindings permissions

// diffLatest is meant to be used to show the difference between a kubernetes manifest
// in Go and the latest version of the same manifest in yaml.
//
//	out := filepath.Join(defaultImportOutputDir, "update")
//	diff, err := kube.DiffLatest(
//		"tekton",
//		"tekton",
//		out,
//		nil,
//		tekton.New(),
//		[]string{"testdata/tekton-updated.yaml"},
//	)
func diffLatest(
	appName, pkgName, destDir string,
	serializer runtime.Decoder,
	km kube.Exporter,
	manifests []string,
) (string, error) {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}

	tmpdir, err := os.MkdirTemp("", appName)
	if err != nil {
		return "", fmt.Errorf("mkdir temp: %w", err)
	}
	defer func() {
		if err != nil {
			fmt.Printf("temp dir: %s\n", tmpdir)
			return
		}
		err = errors.Join(os.RemoveAll(tmpdir), err)
		// TODO not sure about this error handling
	}()

	// EXPORT TO YAML
	err = kube.Export(km, kube.WithExportOutputDirectory(tmpdir))
	if err != nil {
		return "", fmt.Errorf("export: %w", err)
	}

	// TODO: remove all labels as an option

	// IMPORT BACK TO GO
	currManifests, err := kubeutil.ListYAMLFiles(tmpdir)
	if err != nil {
		return "", fmt.Errorf("list yaml files: %w", err)
	}

	if serializer == nil {
		serializer = defaultSerializer()
	}
	// IMPORT CURRENT MANIFESTS TO GO
	currentDir := "current"
	arc, err := importArchive(
		appName,
		pkgName,
		currentDir,
		serializer,
		currManifests,
	)

	latestDir := "latest"
	aru, err := importArchive(
		appName,
		pkgName,
		latestDir,
		serializer,
		manifests,
	)

	d := cmp.Diff(arc, aru, cmpopts.EquateEmpty())

	if err := os.WriteFile(
		filepath.Join(destDir, "diff.txt"),
		[]byte(d),
		0o644,
	); err != nil {
		return "", fmt.Errorf("write diff: %w", err)
	}

	if err := txtar.Write(arc, destDir); err != nil {
		return "", fmt.Errorf("write current: %w", err)
	}
	if err := txtar.Write(aru, destDir); err != nil {
		return "", fmt.Errorf("write latest: %w", err)
	}
	return d, err
}

func importArchive(
	appName string,
	pkgName string,
	outDir string,
	serializer runtime.Decoder,
	manifests []string,
) (*txtar.Archive, error) {
	// IMPORT YAML TO GO
	var buf bytes.Buffer
	err := kube.Import(
		kube.WithImportOutputDirectory(outDir),
		kube.WithImportManifestFiles(manifests),
		kube.WithImportAppName(appName),
		kube.WithImportRemoveAppName(true),
		kube.WithImportPackageName(pkgName),
		kube.WithImportSerializer(serializer),
		kube.WithImportGroupByKind(true),
		kube.WithImportWriter(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("import: %w", err)
	}

	// FORMAT SOURCE CODE with gofumpt -extra-rules
	ar := txtar.Parse(buf.Bytes())
	for _, f := range ar.Files {
		f.Data, err = format.Source(
			f.Data, format.Options{
				LangVersion: "1.20",
				// ModulePath:  "",
				ExtraRules: true,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("formating generated code: %w", err)
		}
	}
	sort.SliceStable(
		ar.Files, func(i, j int) bool {
			return ar.Files[i].Name < ar.Files[j].Name
		},
	)
	return ar, nil
}
