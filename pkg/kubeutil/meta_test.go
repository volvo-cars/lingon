// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	"reflect"
	"testing"

	tu "github.com/volvo-cars/lingon/pkg/testutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestObjectMeta(t *testing.T) {
	type args struct {
		name        string
		namespace   string
		labels      map[string]string
		annotations map[string]string
	}
	tests := []struct {
		name string
		args args
		want metav1.ObjectMeta
	}{
		{
			name: "meta",
			args: args{
				name:        "o",
				namespace:   "ns",
				labels:      nil,
				annotations: nil,
			},
			want: metav1.ObjectMeta{
				Name:      "o",
				Namespace: "ns",
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := ObjectMeta(
					tt.args.name,
					tt.args.namespace,
					tt.args.labels,
					tt.args.annotations,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ObjectMeta() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestMetadata_GVK(t *testing.T) {
	type fields struct {
		APIVersion string
		Kind       string
		Meta       Meta
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "gvk",
			fields: fields{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Meta: Meta{
					Name:        "TestDeploy",
					Namespace:   "TestNS",
					Labels:      nil,
					Annotations: nil,
				},
			},
			want: "apps/v1, Kind=Deployment",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				m := &Metadata{
					APIVersion: tt.fields.APIVersion,
					Kind:       tt.fields.Kind,
					Meta:       tt.fields.Meta,
				}
				tu.AssertEqual(t, tt.want, m.GVK())
			},
		)
	}
}

func TestExtractMetadata(t *testing.T) {
	type TT struct {
		tname string
		in    string
		want  *Metadata
	}

	tt := []TT{
		{
			tname: "simple deployment",
			in: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: super-duper-app
`,
			want: &Metadata{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Meta: Meta{
					Name:        "super-duper-app",
					Namespace:   "default",
					Labels:      map[string]string{},
					Annotations: map[string]string{},
				},
			},
		},
		{
			tname: "deployment with labels",
			in: `
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app.kubernetes.io/component: repo-server
    app.kubernetes.io/name: argocd-repo-server
    app.kubernetes.io/part-of: argocd
  name: argocd-repo-server
`,
			want: &Metadata{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Meta: Meta{
					Name:      "argocd-repo-server",
					Namespace: "default",
					Labels: map[string]string{
						"app.kubernetes.io/component": "repo-server",
						"app.kubernetes.io/name":      "argocd-repo-server",
						"app.kubernetes.io/part-of":   "argocd",
					},
					Annotations: map[string]string{},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(
			tc.tname, func(t *testing.T) {
				got, err := ExtractMetadata([]byte(tc.in))
				tu.AssertNoError(t, err)
				tu.AssertEqual(t, tc.want, got)
			},
		)
	}
}
