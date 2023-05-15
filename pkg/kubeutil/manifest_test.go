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

func TestManifestReadFile(t *testing.T) {
	tt := []struct {
		name    string
		in      string
		want    int
		wantErr bool
	}{
		{
			name: "empty",
			in:   "empty.yaml",
			want: 2,
		},
		{
			name: "broken",
			in:   "broken.yaml",
			want: 5,
		},
		{
			name: "apps",
			in:   "apps.yaml",
			want: 8,
		},
		{
			name:    "json",
			in:      "apps.json",
			wantErr: true,
		},
		{
			name:    "fake",
			in:      "fake.yaml",
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(
			tc.name, func(t *testing.T) {
				gp := filepath.Join("testdata", tc.in)

				rf, err := ManifestReadFile(gp)
				if (err != nil) != tc.wantErr {
					t.Error("expected error", err)
				}

				tu.AssertEqual(t, tc.want, len(rf))
			},
		)
	}
}

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
			},
		)
	}
}

func TestCleanUpYAML(t *testing.T) {
	tests := []struct {
		name    string
		in      []byte
		want    []byte
		wantErr bool
	}{
		{
			name: "metadata",
			in: []byte(`
metadata:
  name: "n"
  other: "to be removed"
`),
			want: []byte(`{"metadata":{"name":"n"}}`),
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := CleanUpYAML(tt.in)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"CleanUpYAML() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
				tu.AssertEqual(t, string(tt.want), string(got))
			},
		)
	}
}

func TestCleanUpJSON(t *testing.T) {
	tests := []struct {
		name    string
		in      []byte
		want    []byte
		wantErr bool
	}{
		{
			name: "metadata",
			in:   []byte(`{"metadata":{"name":"n","other": "remove me"}}`),
			want: []byte(`{"metadata":{"name":"n"}}`),
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := CleanUpJSON(tt.in)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"CleanUpJSON() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
				tu.AssertEqual(t, string(tt.want), string(got))
			},
		)
	}
}
