// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
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

	return g.export(km)
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
		filenames := []string{}
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
		if _, err = g.o.ManifestWriter.Write(txtar.Format(g.ar)); err != nil {
			return fmt.Errorf("write: %w", err)
		}
		return nil
	}

	if err = os.MkdirAll(g.o.OutputDir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", g.o.OutputDir, err)
	}

	if err = txtar.Write(g.ar, "."); err != nil {
		return fmt.Errorf("txtar write: %w", err)
	}

	return err
}
