// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/rogpeppe/go-internal/txtar"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestEncode_EmptyField(t *testing.T) {
	ebd := newEmbedWithEmptyField()

	ar := &txtar.Archive{}

	g := goky{
		useWriter: false,
		ar:        ar,
		o: exportOption{
			OutputDir:      "out",
			ManifestWriter: nil,
			NameFileFunc:   nil,
			SecretHook:     nil,
			Kustomize:      false,
			Explode:        false,
		},
		kn: map[string]struct{}{},
	}

	rv := reflect.ValueOf(ebd)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Type().Kind() != reflect.Struct {
		t.Errorf("cannot encode non-struct type: %v", rv)
	}
	err := g.encodeStruct(rv, "")
	if !errors.Is(err, ErrFieldMissing) {
		t.Fatal(err)
	}

	want := []string{
		"out/1_iamcr.yaml",
		"out/1_iamsa.yaml",
		"out/2_iamcrb.yaml",
		"out/3_iamdepl.yaml",
	}

	if diff := tu.Diff(want, tu.Filenames(ar)); diff != "" {
		t.Error(tu.Callers(), diff)
	}
}

func TestEncode_EmbeddedStruct(t *testing.T) {
	ebd := newEmbeddedStruct()

	ar := &txtar.Archive{}

	g := goky{
		useWriter: false,
		ar:        ar,
		o: exportOption{
			OutputDir:      "out",
			ManifestWriter: nil,
			NameFileFunc:   nil,
			SecretHook:     nil,
			Kustomize:      false,
			Explode:        false,
		},
		kn: map[string]struct{}{},
	}

	rv := reflect.ValueOf(ebd)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Type().Kind() != reflect.Struct {
		t.Errorf("cannot encode non-struct type: %v", rv)
	}
	err := g.encodeStruct(rv, "")
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		"out/1_iamcr.yaml",
		"out/1_iamsa.yaml",
		"out/2_iamcrb.yaml",
		"out/3_depl.yaml",
		"out/3_iamdepl.yaml",
	}

	if diff := tu.Diff(want, tu.Filenames(ar)); diff != "" {
		t.Error(tu.Callers(), diff)
	}

	golden := filepath.Join("testdata", "golden", "encode.txt")
	expected, err := txtar.ParseFile(golden)
	tu.AssertNoError(t, err)

	if diff := tu.DiffTxtar(ar, expected); diff != "" {
		t.Error(tu.Callers(), diff)
	}
}

func TestEncode_DuplicateNames(t *testing.T) {
	ebd := struct {
		NS1 *corev1.Namespace
		NS2 *corev1.Namespace
	}{
		NS1: kubeutil.Namespace("same-ns", nil, nil),
		NS2: kubeutil.Namespace("same-ns", nil, nil),
	}
	g := goky{
		ar: &txtar.Archive{},
		o:  exportOption{},
		kn: map[string]struct{}{},
	}

	rv := reflect.ValueOf(ebd)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Type().Kind() != reflect.Struct {
		t.Errorf("cannot encode non-struct type: %v", rv)
	}
	err := g.encodeStruct(rv, "")
	if !errors.Is(err, ErrDuplicateDetected) {
		t.Error(tu.Callers(), "duplicate not detected")
	}
}

type IAM struct {
	App
	Sa   *corev1.ServiceAccount
	Crb  *rbacv1.ClusterRoleBinding
	Cr   *rbacv1.ClusterRole
	Depl *appsv1.Deployment
}

type EmbedStruct struct {
	IAM
	Depl *appsv1.Deployment
}

type EmbedWithEmptyField struct {
	IAM
	EmptyDepl *appsv1.Deployment
	ZDepl     *appsv1.Deployment
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

func newEmbedWithEmptyField() *EmbedWithEmptyField {
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

	return &EmbedWithEmptyField{
		ZDepl: kubeutil.SimpleDeployment(
			name,
			sa.Namespace,
			labels,
			int32(1),
			"nginx:latest",
		),
		IAM: iam,
	}
}
