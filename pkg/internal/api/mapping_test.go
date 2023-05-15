// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"testing"

	"github.com/dave/jennifer/jen"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestImportKubernetesPkgAlias(t *testing.T) {
	tests := []struct {
		name string
		file *jen.File
		want string
	}{
		{
			name: "pkg alias",
			file: jen.NewFile("mypackage"),
			want: `package mypackage

import appsv1 "k8s.io/api/apps/v1"

var bla = appsv1.Deployment{}
`,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.file.Add(
					jen.Var().Id("bla").Op("=").Qual(
						"k8s.io/api/apps/v1",
						"Deployment",
					).Block(),
				)
				ImportKubernetesPkgAlias(tt.file)
				s := fmt.Sprintf("%#v", tt.file)
				tu.AssertEqual(t, s, tt.want)
			},
		)
	}
}

func TestPkgPathFromAPIVersion(t *testing.T) {
	tests := []struct {
		name       string
		apiVersion string
		want       string
		wantErr    bool
	}{
		{
			name:       "happy",
			apiVersion: "apps/v1",
			want:       "k8s.io/api/apps/v1",
		},
		{
			name:       "sad",
			apiVersion: "asjdhfajklsdhf",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := PkgPathFromAPIVersion(tt.apiVersion)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"PkgPathFromAPIVersion() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
				if got != tt.want {
					t.Errorf(
						"PkgPathFromAPIVersion() got = %v, want %v",
						got,
						tt.want,
					)
				}
			},
		)
	}
}
