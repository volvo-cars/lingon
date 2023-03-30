// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/volvo-cars/lingon/pkg/kubeutil"
)

// Explode reads a YAML manifests from [io.Reader] and writes the objects to files in dir.
// Each object is written to a file named after the object's kind and name.
// The file name is prefixed with a number that indicates the rank of the kind.
// The rank is used to sort the files in the directory and prioritize the order
// in which they are applied. For example, a namespace should be applied before
// any other object in the namespace.
func Explode(r io.Reader, dir string) error {
	content, err := splitManifest(r)
	if err != nil {
		return fmt.Errorf("explode: %w", err)
	}

	if dir == "" {
		dir = "out"
	}

	for _, obj := range content {
		if len(obj) == 0 {
			continue
		}

		// get name of the object in metadata
		m, err := kubeutil.ExtractMetadata([]byte(obj))
		if err != nil {
			return fmt.Errorf("extract metadata: %w", err)
		}
		dn := DirectoryName(m.Meta.Namespace, m.Kind)
		out := filepath.Join(dir, dn)
		if err := os.MkdirAll(out, 0o755); err != nil {
			return err
		}
		fn := fmt.Sprintf(
			"%d_%s.yaml",
			rankOfKind(m.Kind),
			basicName(m.Meta.Name, m.Kind),
		)
		outName := filepath.Join(out, fn)
		if err := write(obj, outName); err != nil {
			return fmt.Errorf("explode: write %s: %w", outName, err)
		}
	}

	return nil
}
