// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
				if err != nil {
					t.Errorf("err: %s", err)
				}
				if !cmp.Equal(got, tc.want) {
					t.Errorf(cmp.Diff(got, tc.want))
				}
			},
		)
	}
}
