// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/txtar"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestManifestSplit(t *testing.T) {
	tt := []struct {
		name  string
		in    string
		want  string
		wantN int
	}{
		{
			name:  "empty",
			in:    "empty.txt",
			want:  "empty.yaml",
			wantN: 2,
		},
		{
			name:  "broken",
			in:    "broken.txt",
			want:  "broken.yaml",
			wantN: 5,
		},
	}

	for _, tc := range tt {
		t.Run(
			tc.name, func(t *testing.T) {
				gp := filepath.Join("testdata", tc.in)

				ar, err := txtar.ParseFile(gp)
				tu.AssertNoError(t, err, "parsing txtar file")

				ms, err := ManifestSplit(bytes.NewReader(Txtar2YAML(ar)))
				tu.AssertNoError(t, err, "manifest split")

				tu.AssertEqual(t, tc.wantN, len(ms))
				got := []byte(strings.Join(ms, "\n---\n"))
				tu.AssertNoError(t, err, "split manifest")

				want, err := os.ReadFile(filepath.Join("testdata", tc.want))
				tu.AssertNoError(t, err, "read want file")

				tu.AssertEqual(t, string(want), string(got))
			})
	}
}
