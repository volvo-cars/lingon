// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	kyaml "sigs.k8s.io/yaml"
)

// encodeApp encodes kube.App to a map of YAML manifests.
// The keys are the struct field names.
func encodeApp(km Exporter) (map[string][]byte, error) {
	res := make(map[string][]byte)
	rv := reflect.ValueOf(km)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Type().Kind() != reflect.Struct {
		return nil, fmt.Errorf("cannot encode non-struct type: %v", rv)
	}

	if err := encodeStruct(rv, "", res); err != nil {
		return nil, err
	}

	return res, nil
}

func encodeStruct(
	rv reflect.Value,
	prefix string, // prefix is used for nested structs
	res map[string][]byte,
) error {
	if rv.Type().Kind() != reflect.Struct {
		return fmt.Errorf("cannot encode non-struct type: %v", rv)
	}

	for i := 0; i < rv.NumField(); i++ {
		sf := rv.Type().Field(i)
		fv := rv.Field(i)

		if sf.Anonymous {
			if err := encodeStruct(fv, sf.Name, res); err != nil {
				return err
			}
			continue
		}

		fieldVal := rv.FieldByName(sf.Name)
		switch v := fieldVal.Interface().(type) {
		case runtime.Object:
			if reflect.ValueOf(v).IsZero() {
				return fmt.Errorf(
					"%w: %q of type %q",
					ErrFieldMissing,
					sf.Name,
					sf.Type,
				)
			}
			r := rank(v)

			// It works by first marshalling to JSON, so no `yaml` tag necessary
			b, err := kyaml.Marshal(v)
			if err != nil {
				return fmt.Errorf(
					"error marshaling field %s: %w",
					sf.Name,
					err,
				)
			}

			res[r+"_"+prefix+sf.Name] = b

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
