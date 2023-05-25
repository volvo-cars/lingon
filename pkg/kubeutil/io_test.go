// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	"testing"

	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestFileExists(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			name:     "exists",
			filename: "kubeutil_test.go",
			want:     true,
		},
		{
			name:     "does not exist",
			filename: "oops",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if FileExists(tt.filename) != tt.want {
					t.Failed()
				}
			},
		)
	}
}

func Test_contains(t *testing.T) {
	type args struct {
		e string
		s []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "in haystack",
			args: args{
				e: "hey",
				s: []string{"hey", "ho", "let's go"},
			},
			want: true,
		},
		{
			name: "not in haystack",
			args: args{
				e: "NIE",
				s: []string{"hey", "ho", "let's go"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := contains(tt.args.e, tt.args.s); got != tt.want {
					t.Errorf("contains() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestListGoFiles(t *testing.T) {
	tests := []struct {
		name    string
		root    string
		want    []string
		wantErr bool
	}{
		{
			name: "go files",
			root: ".",
			want: []string{
				"config.go",
				"config_test.go",
				"doc.go",
				"iam.go",
				"iambinding.go",
				"io.go",
				"io_test.go",
				"kubeutil.go",
				"kubeutil_test.go",
				"label.go",
				"manifest.go",
				"manifest_test.go",
				"meta.go",
				"meta_test.go",
				"name.go",
				"namespace.go",
				"p.go",
				"p_test.go",
				"pod.go",
				"secret.go",
				"txtar.go",
				"txtar_test.go",
				"typemeta_gen.go",
			},
		},
		{
			name: "file but none go",
			root: "testdata",
			want: []string{},
		},
		{
			name:    "path err",
			root:    "oops",
			wantErr: true,
		},
		{
			name:    "is a file",
			root:    "testdata/apps.yaml",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := ListGoFiles(tt.root)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"ListGoFiles() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
				tu.AssertEqualSlice(t, tt.want, got)
			},
		)
	}
}

func TestListYAMLFiles(t *testing.T) {
	tests := []struct {
		name    string
		root    string
		want    []string
		wantErr bool
	}{
		{
			name: "yaml files",
			root: "testdata",
			want: []string{
				"testdata/apps.yaml",
				"testdata/broken.yaml",
				"testdata/empty.yaml",
			},
		},
		{
			name: "file but none go",
			root: "../terra",
			want: []string{},
		},
		{
			name:    "path err",
			root:    "oops",
			wantErr: true,
		},
		{
			name:    "is a file",
			root:    "testdata/apps.yaml",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := ListYAMLFiles(tt.root)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"ListGoFiles() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
				tu.AssertEqualSlice(t, tt.want, got)
			},
		)
	}
}
