// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/rogpeppe/go-internal/txtar"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
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

func gatekeeperExportOptions(o exportOption) error {
	if o.Explode && o.SingleFile != "" {
		return fmt.Errorf(
			"WithExportExplodeManifests and WithExportAsSingleFile: %w",
			ErrIncompatibleOptions,
		)
	}

	if o.OutputJSON && o.Kustomize {
		return fmt.Errorf(
			"WithExportJSON and WithExportKustomize: %w",
			ErrIncompatibleOptions,
		)
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

	if g.o.Kustomize && !g.o.OutputJSON {
		// extract the filenames for kustomization.yaml
		var filenames []string
		for _, f := range g.ar.Files {
			filenames = append(filenames, f.Name)
		}

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
			// write to [io.Writer] as JSON array in a single file
			if g.o.OutputJSON {
				// The txtar is already in JSON format
				if _, err = g.o.ManifestWriter.Write(kubeutil.Txtar2JSON(g.ar)); err != nil {
					return fmt.Errorf("write: %w", err)
				}
				return nil
			}
			// write to [io.Writer] as YAML in a single file
			if _, err = g.o.ManifestWriter.Write(kubeutil.Txtar2YAML(g.ar)); err != nil {
				return fmt.Errorf("write: %w", err)
			}
			return nil
		}

		// write to [io.Writer] as multiple files (txtar format)
		if _, err = g.o.ManifestWriter.Write(txtar.Format(g.ar)); err != nil {
			return fmt.Errorf("write: %w", err)
		}
		return nil
	}

	// write to files on disk

	if err = os.MkdirAll(g.o.OutputDir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", g.o.OutputDir, err)
	}

	if g.o.SingleFile != "" {
		f := filepath.Join(g.o.OutputDir, g.o.SingleFile)
		// write to single file named g.o.SingleFile as JSON array
		if g.o.OutputJSON {
			if err = os.WriteFile(
				f,
				kubeutil.Txtar2JSON(g.ar),
				0o600,
			); err != nil {
				return fmt.Errorf("write file%s: %w", f, err)
			}
			return nil
		}
		// write to single file named g.o.SingleFile as YAML
		if err = os.WriteFile(f, kubeutil.Txtar2YAML(g.ar), 0o600); err != nil {
			return fmt.Errorf("write file%s: %w", f, err)
		}
		return nil
	}

	// write to multiple files (txtar format)
	if err = txtar.Write(g.ar, "."); err != nil {
		return fmt.Errorf("write txtar: %w", err)
	}

	return err
}
