// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

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
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"mvdan.cc/gofumpt/format"
)

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

// DiffLatest is meant to be used to show the difference between a kubernetes manifest
// in Go and the latest version of the same manifest in yaml.
//
//	out := filepath.Join(defOutDir, "update")
//	diff, err := kube.DiffLatest(
//		"tekton",
//		"tekton",
//		out,
//		nil,
//		tekton.New(),
//		[]string{"testdata/tekton-updated.yaml"},
//	)
func DiffLatest(
	appName, pkgName, outDir string,
	serializer runtime.Decoder,
	km Exporter,
	manifests []string,
) (string, error) {
	if err := upsertFolder(outDir); err != nil {
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
	err = export(km, tmpdir, false)
	if err != nil {
		return "", fmt.Errorf("export: %w", err)
	}

	// IMPORT BACK TO GO
	currManifests, err := ListYAMLFiles(tmpdir)
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
		filepath.Join(outDir, "diff.txt"),
		[]byte(d),
		0o644,
	); err != nil {
		return "", fmt.Errorf("write diff: %w", err)
	}

	if err := txtar.Write(arc, outDir); err != nil {
		return "", fmt.Errorf("write current: %w", err)
	}
	if err := txtar.Write(aru, outDir); err != nil {
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
	err := Import(
		WithOutputDirectory(outDir),
		WithManifestFiles(manifests),
		WithAppName(appName),
		WithRemoveAppName(true),
		WithPackageName(pkgName),
		WithSerializer(serializer),
		WithGroupByKind(true),
		WithWriter(&buf),
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

func defaultSerializer() runtime.Decoder {
	// NEEDED FOR CRDS
	//
	_ = apiextensions.AddToScheme(scheme.Scheme)
	return scheme.Codecs.UniversalDeserializer()
}

func upsertFolder(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		errmk := os.MkdirAll(dir, 0o755)
		if errmk != nil {
			return fmt.Errorf("create folder %s: %w", dir, err)
		}
		return nil
	}
	return nil
}
