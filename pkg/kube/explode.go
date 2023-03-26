// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/volvo-cars/lingon/pkg/kubeutil"
)

func Explode(r io.Reader, outDir string) error {
	content, err := splitManifest(r)
	if err != nil {
		return fmt.Errorf("explode: %w", err)
	}

	if outDir == "" {
		outDir = "out"
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
		dir := DirectoryName(outDir, m.Meta.Namespace, m.Kind)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		fn := fmt.Sprintf(
			"%d_%s.yaml",
			rankOfKind(m.Kind),
			basicName(m.Meta.Name, m.Kind),
		)
		outName := filepath.Join(dir, fn)
		if err := write(obj, outName); err != nil {
			return fmt.Errorf("explode: write %s: %w", outName, err)
		}
	}

	return nil
}
