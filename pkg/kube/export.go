// Copyright (c) Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/veggiemonk/strcase"
)

var ErrFieldMissing = errors.New("missing")

// ExportWithKustomization exports Exporter
// containing kubernetes object to yaml files with kustomization.yaml.
func ExportWithKustomization(km Exporter, outDir string) error {
	return export(km, outDir, true)
}

// Export exports Exporter (struct containing kube.App) containing kubernetes object to YAML files.
// The YAML files are written to output directory.
// For example, see ExportWriter.
func Export(km Exporter, outDir string) error {
	return export(km, outDir, false)
}

// ExportWriter writes kubernetes object in YAML to io.Writer.
func ExportWriter(km Exporter, w io.Writer) error {
	manifests, err := encodeApp(km)
	if err != nil {
		return err
	}

	for _, k := range orderedKeys(manifests) {
		m := manifests[k]
		s := cleanManifest(m)
		if _, err = w.Write([]byte(s)); err != nil {
			return fmt.Errorf("export write: %w", err)
		}
		if _, err = w.Write([]byte("\n---\n")); err != nil {
			return fmt.Errorf("export write: %w", err)
		}
	}

	return nil
}

func export(km Exporter, destDir string, addKustomization bool) error {
	manifests, err := encodeApp(km)
	if err != nil {
		return err
	}

	for name, m := range manifests {
		if err = os.MkdirAll(destDir, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", destDir, err)
		}

		n := strcase.Snake(name) + ".yaml"

		f, err := os.Create(filepath.Join(destDir, n))
		if err != nil {
			return fmt.Errorf("create file %s: %w", name, err)
		}

		s := cleanManifest(m)
		if _, err = f.WriteString(s); err != nil {
			return fmt.Errorf("write file %s: %w", name, err)
		}
		if err = f.Close(); err != nil {
			return fmt.Errorf("close file %s: %w", name, err)
		}
	}

	if addKustomization {
		nn := make([]string, 0)
		for name := range manifests {
			nn = append(nn, name+".yaml")
		}
		if err := kustomization(destDir, nn...); err != nil {
			return fmt.Errorf("writing kustomize.yaml: %w", err)
		}
	}

	return err
}

var clm = strings.NewReplacer(
	"creationTimestamp: null", "",
	"status: {}", "",
	`status:
  loadBalancer: {}`, "",
	`status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: null
  storedVersions: null`, "",
)

func cleanManifest(m []byte) string {
	return strings.TrimSpace(clm.Replace(string(m)))
}

func kustomization(out string, files ...string) error {
	s := `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:`
	for _, name := range files {
		s = s + "\n  - " + name
	}

	f, err := os.Create(filepath.Join(out, "kustomization.yaml"))
	if err != nil {
		return err
	}

	_, err = f.WriteString(s)
	if err != nil {
		return err
	}

	return nil
}
