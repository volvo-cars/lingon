// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/volvo-cars/lingon/pkg/kubeutil"
	"golang.org/x/tools/txtar"
)

type goky struct {
	ar            *txtar.Archive
	o             exportOption
	dup           map[string]struct{} // kind/name  for duplicate detection
	useWriter     bool
	useSingleFile bool
}

func Export(km Exporter, opts ...ExportOption) error {
	g := goky{
		ar:  &txtar.Archive{},
		o:   exportDefaultOpts,
		dup: make(map[string]struct{}, 0),
	}
	for _, o := range opts {
		o(&g)
	}

	if err := g.gatekeeperExportOptions(); err != nil {
		return fmt.Errorf("export: %w", err)
	}

	return g.export(km)
}

func (g *goky) gatekeeperExportOptions() error {
	if g.o.Explode && g.useSingleFile {
		return fmt.Errorf(
			"WithExportExplodeManifests and WithExportAsSingleFile: %w",
			ErrIncompatibleOptions,
		)
	}

	if g.o.OutputJSON && g.o.Kustomize {
		return fmt.Errorf(
			"WithExportJSON and WithExportKustomize: %w",
			ErrIncompatibleOptions,
		)
	}

	return nil
}

func (g *goky) export(km Exporter) error {
	var err error
	if km == nil {
		return fmt.Errorf("cannot export type %T: %v", km, km)
	}
	if err = g.encodeStruct(reflect.ValueOf(km), ""); err != nil {
		return fmt.Errorf("encoding: %w", err)
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
		if g.useSingleFile {
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

	if g.useSingleFile {
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
	if err = write(g.ar, "."); err != nil {
		return fmt.Errorf("write txtar: %w", err)
	}

	return err
}

// Following code taken from https://github.com/rogpeppe/go-internal
// as it is now deprecated in favor of
// https://github.com/golang/tools/blob/master/txtar/archive.go

// Write writes each File in an Archive to the given directory, returning any
// errors encountered. An error is also returned in the event a file would be
// written outside of dir.
func write(a *txtar.Archive, dir string) error {
	for _, f := range a.Files {
		fp := filepath.Clean(filepath.FromSlash(f.Name))
		if isAbs(fp) || strings.HasPrefix(fp, ".."+string(filepath.Separator)) {
			return fmt.Errorf("%q: outside parent directory", f.Name)
		}
		fp = filepath.Join(dir, fp)

		if err := os.MkdirAll(filepath.Dir(fp), 0o777); err != nil {
			return err
		}
		// Avoid overwriting existing files by using O_EXCL.
		out, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o666)
		if err != nil {
			return err
		}

		_, err = out.Write(f.Data)
		cerr := out.Close()
		if err != nil {
			return err
		}
		if cerr != nil {
			return cerr
		}
	}
	return nil
}

func isAbs(p string) bool {
	// Note: under Windows, filepath.IsAbs(`\foo`) returns false,
	// so we need to check for that case specifically.
	return filepath.IsAbs(p) || strings.HasPrefix(p, string(filepath.Separator))
}
