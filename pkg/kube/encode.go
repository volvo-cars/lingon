// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"

	"github.com/rogpeppe/go-internal/txtar"
	"github.com/tidwall/sjson"
	"github.com/veggiemonk/strcase"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kyaml "sigs.k8s.io/yaml"
)

var (
	ErrFieldMissing      = errors.New("missing")
	ErrDuplicateDetected = errors.New("duplicate detected")
)

// encodeStruct encodes [kube.App] struct to a [txtar.Archive].
func (g *goky) encodeStruct(
	rv reflect.Value,
	prefix string, // prefix is used for nested structs
) error {
	if rv.Type().Kind() != reflect.Struct {
		return fmt.Errorf("cannot encode non-struct type: %v", rv)
	}

	for i := 0; i < rv.NumField(); i++ {
		sf := rv.Type().Field(i)
		fv := rv.Field(i)

		// embedded struct
		if sf.Anonymous {
			if err := g.encodeStruct(fv, sf.Name); err != nil {
				return err
			}
			continue
		}

		fieldVal := rv.FieldByName(sf.Name)
		switch t := fieldVal.Interface().(type) {
		case Exporter:
			if fv.Kind() == reflect.Ptr {
				fv = fv.Elem()
			}
			if err := g.encodeStruct(fv, sf.Name); err != nil {
				return err
			}
			continue
		case runtime.Object:
			v := reflect.ValueOf(t)
			// TODO: should we continue instead ? Should it be an option ?
			if v.IsZero() || v.IsNil() {
				return fmt.Errorf(
					"%w: %q of type %q in %q",
					ErrFieldMissing,
					sf.Name,
					sf.Type,
					rv.Type(),
				)
			}

			switch sec := t.(type) {
			case *corev1.Secret:
				if g.o.SecretHook != nil {
					if err := g.o.SecretHook(sec); err != nil {
						return err
					}
					// skip the secret if the hook is used
					continue
				}
			default:

			}

			kj, err := kyaml.Marshal(t)
			if err != nil {
				return fmt.Errorf("error marshaling field %s: %w", sf.Name, err)
			}
			// Extract metadata to get the name of the file
			m, err := kubeutil.ExtractMetadata(kj)
			if err != nil {
				return fmt.Errorf("extract metadata: %w", err)
			}

			id := m.String()
			if _, ok := g.kn[id]; ok {
				return fmt.Errorf("%w: %s", ErrDuplicateDetected, id)
			}
			g.kn[m.String()] = struct{}{}

			kj, err = kyaml.YAMLToJSON(kj)
			if err != nil {
				return fmt.Errorf("YAMLToJSON %s: %w", sf.Name, err)
			}
			// delete unwanted fields
			kj, err = sjson.DeleteBytes(kj, "metadata.creationTimestamp")
			if err != nil {
				return fmt.Errorf(
					"deleting creationTimestamp %s: %w", sf.Name, err,
				)
			}

			kj, err = sjson.DeleteBytes(kj, "status")
			if err != nil {
				return fmt.Errorf("deleting status %s: %w", sf.Name, err)
			}

			yb, err := kyaml.JSONToYAML(kj)
			if err != nil {
				return fmt.Errorf("error marshaling field %s: %w", sf.Name, err)
			}

			ext := "yaml"
			if g.o.OutputJSON {
				ext = "json"
				yb = kj
			}
			name := fmt.Sprintf(
				"%d_%s.%s",
				rankOfKind(m.Kind),
				strcase.Snake(prefix)+strcase.Snake(sf.Name),
				ext,
			)

			// compute the name of the file and directory
			if g.o.NameFileFunc != nil {
				name = g.o.NameFileFunc(m)
			}
			if g.o.Explode {
				dn := DirectoryName(m.Meta.Namespace, m.Kind)
				name = filepath.Join(dn, name)
			}
			if g.o.OutputDir != "" {
				name = filepath.Join(g.o.OutputDir, name)
			}

			g.ar.Files = append(g.ar.Files, txtar.File{Name: name, Data: yb})

		default:
			// Not sure if this should be an error, but rather be explicit at this point
			return fmt.Errorf(
				"unknown public field: %s, type: %s, kind: %s",
				sf.Name,
				sf.Type,
				fieldVal.Kind(),
			)
		}
	}
	// predictable output
	sort.SliceStable(
		g.ar.Files, func(i, j int) bool {
			return g.ar.Files[i].Name < g.ar.Files[j].Name
		},
	)

	return nil
}
