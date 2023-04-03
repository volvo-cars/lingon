// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"

	"github.com/rogpeppe/go-internal/txtar"
)

type goky struct {
	ar        *txtar.Archive
	o         exportOption
	useWriter bool
}

func Export(km Exporter, opts ...ExportOption) error {
	g := goky{
		ar: &txtar.Archive{},
		o:  exportDefaultOpts,
	}
	for _, o := range opts {
		o(&g)
	}

	if err := gatekeeperExportOptions(g.o); err != nil {
		return fmt.Errorf("export: %w", err)
	}

	return g.export(km)
}

var errOptConflictOutFiles = errors.New("option conflict: WithExportOutputFiles not compatible with WithExportOutputDir")

func gatekeeperExportOptions(o exportOption) error {
	if o.Explode && o.SingleFile != "" {
		return errOptConflictOutFiles
	}
	return nil
}

func (g *goky) export(km Exporter) error {
	var err error
	rv := reflect.ValueOf(km)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Type().Kind() != reflect.Struct {
		return fmt.Errorf("cannot encode non-struct type: %v", rv)
	}

	if err = g.encodeStruct(rv, ""); err != nil {
		return err
	}

	if len(g.ar.Files) == 0 {
		return fmt.Errorf("no file to write")
	}

	if g.o.Kustomize {
		// extract the filenames for kustomization.yaml
		var filenames []string
		for _, f := range g.ar.Files {
			filenames = append(filenames, f.Name)
		}
		// predictable output
		sort.Strings(filenames)

		// add the kustomization.yaml file
		k := filepath.Join(g.o.OutputDir, "kustomization.yaml")
		s := `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:`
		for _, name := range filenames {
			s = s + "\n  - " + name
		}
		b := []byte(s + "\n")
		g.ar.Files = append(g.ar.Files, txtar.File{Name: k, Data: b})
	}

	if g.useWriter {
		if g.o.SingleFile != "" {
			if _, err = g.o.ManifestWriter.Write(Txtar2YAML(g.ar)); err != nil {
				return fmt.Errorf("write: %w", err)
			}
		} else {
			if _, err = g.o.ManifestWriter.Write(txtar.Format(g.ar)); err != nil {
				return fmt.Errorf("write: %w", err)
			}
		}
		return nil
	}

	if err = os.MkdirAll(g.o.OutputDir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", g.o.OutputDir, err)
	}

	if g.o.SingleFile != "" {
		f := filepath.Join(g.o.OutputDir, g.o.SingleFile)
		if err = os.WriteFile(f, Txtar2YAML(g.ar), 0o600); err != nil {
			return fmt.Errorf("write file%s: %w", f, err)
		}
		return nil
	} else {
		if err = txtar.Write(g.ar, "."); err != nil {
			return fmt.Errorf("write txtar: %w", err)
		}
	}

	return err
}
