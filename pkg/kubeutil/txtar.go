// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	"bytes"
	"fmt"

	"golang.org/x/tools/txtar"
	kyaml "sigs.k8s.io/yaml"
)

// Txtar2YAML converts a txtar.Archive to a YAML document.
// Each file in the txtar.Archive must be a YAML document.
// No conversion is done, the files are simply concatenated.
func Txtar2YAML(ar *txtar.Archive) []byte {
	var buf bytes.Buffer
	for _, f := range ar.Files {
		buf.WriteString("\n---\n")
		buf.WriteString("# " + f.Name + "\n")
		buf.Write(f.Data)
	}
	return buf.Bytes()
}

// Txtar2JSON converts a txtar.Archive to a JSON array of JSON objects.
// Each file in the txtar.Archive must be a JSON object.
// No conversion is done, the files are simply concatenated.
func Txtar2JSON(ar *txtar.Archive) []byte {
	var buf bytes.Buffer
	buf.WriteString("[")
	x := ",\n"
	stop := len(ar.Files) - 1
	for i, f := range ar.Files {
		buf.Write(f.Data)
		if i < stop {
			buf.WriteString(x)
		}
	}

	buf.WriteString("\n]")
	return buf.Bytes()
}

// TxtarYAML2TxtarJSON converts a [txtar.Archive] containing YAML files to JSON files.
func TxtarYAML2TxtarJSON(ar *txtar.Archive) (*txtar.Archive, error) {
	var jar txtar.Archive
	for _, f := range ar.Files {
		j, err := kyaml.YAMLToJSON(f.Data)
		if err != nil {
			return nil, fmt.Errorf("converting to json: %w", err)
		}
		jar.Files = append(
			jar.Files, txtar.File{
				Name: f.Name,
				Data: j,
			},
		)
	}
	return &jar, nil
}

// TxtarJSON2TxtarYAML converts a [txtar.Archive] containing JSON files to JSON files.
func TxtarJSON2TxtarYAML(ar *txtar.Archive) (*txtar.Archive, error) {
	var jar txtar.Archive
	for _, f := range ar.Files {
		j, err := kyaml.JSONToYAML(f.Data)
		if err != nil {
			return nil, fmt.Errorf("converting to yaml: %w", err)
		}
		jar.Files = append(
			jar.Files, txtar.File{
				Name: f.Name,
				Data: j,
			},
		)
	}
	return &jar, nil
}
