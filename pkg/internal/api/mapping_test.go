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
		name    string
		varQual [2]string
		want    string
	}{
		{
			name:    "pkg alias deploy",
			varQual: [2]string{"k8s.io/api/apps/v1", "Deployment"},
			want: `package mypackage

import appsv1 "k8s.io/api/apps/v1"

var bla = appsv1.Deployment{}
`,
		},
		{
			name:    "pkg alias cm",
			varQual: [2]string{"k8s.io/api/core/v1", "ConfigMap"},
			want: `package mypackage

import corev1 "k8s.io/api/core/v1"

var bla = corev1.ConfigMap{}
`,
		},
		{
			name: "pkg alias API service",
			varQual: [2]string{
				"k8s.io/kube-aggregator/pkg/apis/apiregistration/v1",
				"APIService",
			},
			want: `package mypackage

import apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

var bla = apiregistrationv1.APIService{}
`,
		},
		{
			name: "non existent package",
			varQual: [2]string{
				"github.com/bork/totallyborked/pkg/apis/v1",
				"Deployment",
			},
			want: `package mypackage

import v1 "github.com/bork/totallyborked/pkg/apis/v1"

var bla = v1.Deployment{}
`,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				f := jen.NewFile("mypackage")
				f.Add(
					jen.Var().Id("bla").Op("=").Qual(
						tt.varQual[0], tt.varQual[1],
					).Block(),
				)
				ImportKubernetesPkgAlias(f)
				s := fmt.Sprintf("%#v", f)
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
