// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	kyaml "sigs.k8s.io/yaml"
)

// ManifestReadFile reads a YAML file and splits it into a list of YAML documents.
func ManifestReadFile(filePath string) ([]string, error) {
	e := filepath.Ext(filePath)
	if e != ".yaml" && e != ".yml" {
		return nil, fmt.Errorf("not yaml file: %s", filePath)
	}
	yf, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read manifest %s: %w", filePath, err)
	}
	splitYaml, err := ManifestSplit(bytes.NewReader(yf))
	if err != nil {
		return nil, fmt.Errorf("splitting manifest: %s: %w", filePath, err)
	}
	return splitYaml, nil
}

// ManifestSplit splits a YAML manifest where each object is separated by '---'
// into a list of string containing YAML documents.
func ManifestSplit(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var content []string
	var buf bytes.Buffer

	for scanner.Scan() {
		txt := scanner.Text()
		tmp := strings.ReplaceAll(txt, " ", "")
		if len(tmp) == 0 {
			continue
		}
		switch {
		// Skip comments
		case strings.HasPrefix(txt, "#"):
			continue
		// Split by '---'
		case txt == "---":
			if buf.Len() > 0 {
				content = append(content, buf.String())
				buf.Reset()
			}
		default:
			buf.WriteString(txt + "\n")
		}
	}

	s := buf.String()
	if len(s) > 0 { // if a manifest ends with '---', don't add it
		content = append(content, s)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("spliting manifests: %w", err)
	}
	return content, nil
}

var unusedFields = []string{
	"status",
	// see https://github.com/tidwall/gjson/blob/master/SYNTAX.md
	`metadata.annotations.kubectl\.kubernetes\.io\/last-applied-configuration`,
}

const metadataFields = "{metadata.name,metadata.namespace,metadata.labels,metadata.annotations}"

func CleanUpYAML(in []byte) ([]byte, error) {
	j, err := kyaml.YAMLToJSON(in)
	if err != nil {
		return nil, err
	}
	return CleanUpJSON(j)
}

func CleanUpJSON(in []byte) ([]byte, error) {
	out := in
	var err error

	newMeta := gjson.GetBytes(out, metadataFields)
	out, err = sjson.SetBytes(out, "metadata", newMeta.Value())
	if err != nil {
		return in, fmt.Errorf("error setting new metadata : %v", err)
	}
	for _, field := range unusedFields {
		var errD error
		out, errD = sjson.DeleteBytes(out, field)
		if errD != nil {
			err = errors.Join(err, errD)
		}
	}
	return out, err
}
