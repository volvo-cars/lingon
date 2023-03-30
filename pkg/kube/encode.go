// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/rogpeppe/go-internal/txtar"
	"github.com/veggiemonk/strcase"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kyaml "sigs.k8s.io/yaml"
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

		if sf.Anonymous {
			if err := g.encodeStruct(fv, sf.Name); err != nil {
				return err
			}
			continue
		}

		fieldVal := rv.FieldByName(sf.Name)
		switch t := fieldVal.Interface().(type) {
		case runtime.Object:
			v := reflect.ValueOf(t)
			if v.IsZero() || v.IsNil() {
				return fmt.Errorf(
					"%w: %q of type %q",
					ErrFieldMissing,
					sf.Name,
					sf.Type,
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
			// It works by first marshalling to JSON, so no `yaml` tag necessary
			b, err := kyaml.Marshal(t)
			if err != nil {
				return fmt.Errorf("error marshaling field %s: %w", sf.Name, err)
			}

			// Extract metadata to get the name of the file
			m, err := kubeutil.ExtractMetadata(b)
			if err != nil {
				return fmt.Errorf("extract metadata: %w", err)
			}

			name := fmt.Sprintf(
				"%d_%s.yaml",
				rankOfKind(m.Kind),
				strcase.Snake(prefix)+strcase.Snake(sf.Name),
			)

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

			g.ar.Files = append(g.ar.Files, txtar.File{Name: name, Data: b})

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

	return nil
}
