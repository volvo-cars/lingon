// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"testing"

	"github.com/volvo-cars/lingon/pkg/kubeutil"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestBasicName(t *testing.T) {
	type TT struct {
		name string
		in   string
		kind string
		want string
	}

	tt := []TT{
		{
			name: "simple deployment",
			in:   "super-duper-app",
			kind: "Deployment",
			want: "super-duper-app_deploy",
		},
		{
			name: "deployment with suffix",
			in:   "super-duper-deployment",
			kind: "Deployment",
			want: "super-duper_deploy",
		},
		{
			name: "kind with dash",
			in:   "argo-cluster-role",
			kind: "ClusterRole",
			want: "argo_cr",
		},
		{
			name: "kind with dash in name and dash in kind",
			in:   "tensorboards-web-app-service-account",
			kind: "ServiceAccount",
			want: "tensorboards-web-app_sa",
		},
		{
			name: "kind with no dash in name",
			in:   "tensorboards",
			kind: "ServiceAccount",
			want: "tensorboards_sa",
		},
		{
			name: "short kind at the end",
			in:   "kubecost-cost-analyzer-psp",
			kind: "PodSecurityPolicy",
			want: "kubecost-cost-analyzer_psp",
		},
		{
			name: "short kind at the end with dash",
			in:   "argocd-notifications-cm",
			kind: "ConfigMap",
			want: "argocd-notifications_cm",
		},
	}

	assert := func(t *testing.T, tt TT) {
		got := basicName(tt.in, tt.kind)
		if diff := tu.Diff(got, tt.want); diff != "" {
			t.Error(tu.Callers(), diff)
		}
	}

	for _, tc := range tt {
		t.Run(
			tc.name, func(t *testing.T) {
				assert(t, tc)
			},
		)
	}
}

func TestRemoveAppNameWithFunc(t *testing.T) {
	type TT struct {
		tname     string
		in        string
		app       string
		wantVar   string
		wantField string
		wantFile  string
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
			app:       "super-duper",
			wantVar:   "AppDeploy",
			wantField: "AppDeploy",
			wantFile:  "app_deploy.go",
		},
		{
			tname: "argocd redis service account",
			in: `
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app.kubernetes.io/component: redis
    app.kubernetes.io/name: argocd-redis
    app.kubernetes.io/part-of: argocd
  name: argocd-redis
`,
			app:       "argocd",
			wantVar:   "RedisSA",
			wantField: "RedisSA",
			wantFile:  "redis_sa.go",
		},
		{
			tname: "karpenter cluster-role",
			in: `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: karpenter-admin
`,
			app:       "karpenter",
			wantVar:   "AdminCR",
			wantField: "AdminCR",
			wantFile:  "admin_cr.go",
		},
		{
			tname: "karpenter configmap",
			in: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: karpenter-global-settings
  namespace: karpenter
`,
			app:       "karpenter",
			wantVar:   "GlobalSettingsCM",
			wantField: "GlobalSettingsCM",
			wantFile:  "global-settings_cm.go",
		},
		{
			tname: "karpenter dns role",
			in: `
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: karpenter-dns
`,
			app:       "karpenter",
			wantVar:   "DnsRole",
			wantField: "DnsRole",
			wantFile:  "dns_role.go",
		},
		{
			tname: "argo configmap with dash",
			in: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-notifications-cm
`,
			app:       "argocd",
			wantVar:   "NotificationsCM",
			wantField: "NotificationsCM",
			wantFile:  "notifications_cm.go",
		},
	}

	for _, tc := range tt {
		m, err := kubeutil.ExtractMetadata([]byte(tc.in))
		if err != nil {
			t.Error("ExtractMetadata", tu.Callers(), err)
		}
		// Name of the variable
		t.Run(
			tc.tname+"-var", func(t *testing.T) {
				nameVar := NameVarFunc(*m)
				got := RemoveAppName(nameVar, tc.app)
				if diff := tu.Diff(got, tc.wantVar); diff != "" {
					t.Error("NameVarFunc", tu.Callers(), diff)
				}
			},
		)

		// Name of the field
		t.Run(
			tc.tname+"-field", func(t *testing.T) {
				nameField := NameFieldFunc(*m)
				got := RemoveAppName(nameField, tc.app)
				if diff := tu.Diff(got, tc.wantField); diff != "" {
					t.Error("NameFieldFunc", tu.Callers(), diff)
				}
			},
		)

		// Name of the go file
		t.Run(
			tc.tname+"-file", func(t *testing.T) {
				nameFile := NameFileFunc(*m)
				got := RemoveAppName(nameFile, tc.app)
				if diff := tu.Diff(got, tc.wantFile); diff != "" {
					t.Error("NameFileFunc", tu.Callers(), diff)
				}
			},
		)
	}
}
