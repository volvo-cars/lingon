// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package kube_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const defaultExportOutputDir = "out/export/"

func TestExport(t *testing.T) {
	ebd := newEmbeddedStruct()
	out := filepath.Join(defaultExportOutputDir, "embeddedstruct")
	tu.AssertNoError(t, os.RemoveAll(out), "failed to remove out dir")
	defer os.RemoveAll(out)
	err := kube.Export(ebd, out)
	tu.AssertNoError(t, err, "failed to import")

	got, err := kube.ListYAMLFiles(out)
	tu.AssertNoError(t, err, "failed to list go files")
	sort.Strings(got)

	want := []string{
		"out/export/embeddedstruct/1_iamcr.yaml",
		"out/export/embeddedstruct/1_iamsa.yaml",
		"out/export/embeddedstruct/2_iamcrb.yaml",
		"out/export/embeddedstruct/3_depl.yaml",
		"out/export/embeddedstruct/3_iamdepl.yaml",
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestExportWithKustomization(t *testing.T) {
	ebd := newEmbeddedStruct()

	out := filepath.Join(defaultExportOutputDir, "embeddedstruct2")
	tu.AssertNoError(t, os.RemoveAll(out), "failed to remove out dir")
	defer os.RemoveAll(out)
	err := kube.ExportWithKustomization(ebd, out)
	tu.AssertNoError(t, err, "failed to import")

	got, err := kube.ListYAMLFiles(out)
	tu.AssertNoError(t, err, "failed to list go files")
	sort.Strings(got)

	want := []string{
		"out/export/embeddedstruct2/1_iamcr.yaml",
		"out/export/embeddedstruct2/1_iamsa.yaml",
		"out/export/embeddedstruct2/2_iamcrb.yaml",
		"out/export/embeddedstruct2/3_depl.yaml",
		"out/export/embeddedstruct2/3_iamdepl.yaml",
		"out/export/embeddedstruct2/kustomization.yaml",
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

type IAM struct {
	Sa   *corev1.ServiceAccount
	Crb  *rbacv1.ClusterRoleBinding
	Cr   *rbacv1.ClusterRole
	Depl *appsv1.Deployment
}

type EmbedStruct struct {
	kube.App

	IAM
	Depl *appsv1.Deployment
}

var name = "fyaml"

var labels = map[string]string{
	"app": name,
}

func newEmbeddedStruct() *EmbedStruct {
	sa := kubeutil.SimpleSA(name, "defaultns")
	sa.Labels = labels

	cr := &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{Name: sa.Name, Labels: labels},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
		},
	}
	crb := kubeutil.SimpleCRB(sa, cr)
	crb.Labels = labels

	iam := IAM{
		Sa:  sa,
		Crb: crb,
		Cr:  cr,
		Depl: kubeutil.SimpleDeployment(
			"another"+name,
			sa.Namespace,
			labels,
			int32(1),
			"nginx:latest",
		),
	}

	return &EmbedStruct{
		Depl: kubeutil.SimpleDeployment(
			name,
			sa.Namespace,
			labels,
			int32(1),
			"nginx:latest",
		),
		IAM: iam,
	}
}
