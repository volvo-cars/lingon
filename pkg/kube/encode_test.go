// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	"golang.org/x/tools/txtar"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IAM struct {
	App
	Sa   *corev1.ServiceAccount
	Crb  *rbacv1.ClusterRoleBinding
	Cr   *rbacv1.ClusterRole
	Depl *appsv1.Deployment
}

type EmbedStruct struct {
	App
	IAM
	Depl *appsv1.Deployment
}

type EmbedWithEmptyField struct {
	App
	IAM
	EmptyDepl *appsv1.Deployment
	ZDepl     *appsv1.Deployment
}

var name = "fyaml"

var labels = map[string]string{
	"app": name,
}

func TestEncode_EmptyField(t *testing.T) {
	sa := ku.SimpleSA(name, "defaultns")
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
	crb := ku.SimpleCRB(sa, cr)
	crb.Labels = labels
	ebd := &EmbedWithEmptyField{
		ZDepl: ku.SimpleDeployment(
			name, sa.Namespace, labels, int32(1), "nginx:latest",
		),
		IAM: IAM{
			Sa:  sa,
			Crb: crb,
			Cr:  cr,
			Depl: ku.SimpleDeployment(
				"another"+name, sa.Namespace, labels, int32(1), "nginx:latest",
			),
		},
	}

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
		dup: map[string]struct{}{},
	}

	rv := reflect.ValueOf(ebd)
	err := g.encodeStruct(rv, "")
	if !errors.Is(err, ErrFieldMissing) {
		t.Fatal(err)
	}
	tu.AssertErrorMsg(
		t,
		err,
		`"EmptyDepl" of type "*v1.Deployment" in "kube.EmbedWithEmptyField" is nil: missing`,
	)

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
	sa := ku.SimpleSA(name, "defaultns")
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
	crb := ku.SimpleCRB(sa, cr)
	crb.Labels = labels

	ebd := &EmbedStruct{
		Depl: ku.SimpleDeployment(
			name, sa.Namespace, labels, int32(1), "nginx:latest",
		),
		IAM: IAM{
			Sa:  sa,
			Crb: crb,
			Cr:  cr,
			Depl: ku.SimpleDeployment(
				"another"+name,
				sa.Namespace,
				labels,
				int32(1),
				"nginx:latest",
			),
		},
	}

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
		dup: map[string]struct{}{},
	}

	err := g.encodeStruct(reflect.ValueOf(ebd), "")
	tu.AssertNoError(t, err, "encodeStruct")
	want := []string{
		"out/1_iamcr.yaml",
		"out/1_iamsa.yaml",
		"out/2_iamcrb.yaml",
		"out/3_depl.yaml",
		"out/3_iamdepl.yaml",
	}

	tu.AssertEqualSlice(t, want, tu.Filenames(ar))
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
		App

		NS1 *corev1.Namespace
		NS2 *corev1.Namespace
	}{
		NS1: ku.Namespace("same-ns", nil, nil),
		NS2: ku.Namespace("same-ns", nil, nil),
	}
	g := goky{
		ar:  &txtar.Archive{},
		o:   exportOption{},
		dup: map[string]struct{}{},
	}

	ebd.Lingon()
	err := g.encodeStruct(reflect.ValueOf(ebd), "")
	if !errors.Is(err, ErrDuplicateDetected) {
		t.Error(tu.Callers(), "duplicate not detected")
	}
	tu.AssertErrorMsg(
		t,
		err,
		"v1, Kind=Namespace, NS=default, Name=same-ns: duplicate detected",
	)
}

func TestEncode_DuplicateNames_subfield(t *testing.T) {
	type SubDup struct {
		App
		NS *corev1.Namespace
	}
	ebd := struct {
		App
		NS1 *corev1.Namespace
		Sub *SubDup
	}{
		NS1: ku.Namespace("same-ns", nil, nil),
		Sub: &SubDup{
			NS: ku.Namespace("same-ns", nil, nil),
		},
	}
	g := goky{
		ar:  &txtar.Archive{},
		o:   exportOption{},
		dup: map[string]struct{}{},
	}

	err := g.encodeStruct(reflect.ValueOf(ebd), "")
	if !errors.Is(err, ErrDuplicateDetected) {
		t.Error(tu.Callers(), "duplicate not detected")
	}
	tu.AssertErrorMsg(
		t,
		err,
		"encoding field Sub: v1, Kind=Namespace, NS=default, Name=same-ns: duplicate detected",
	)
}

func TestEncode_SubMissing(t *testing.T) {
	type SubDup struct {
		App
		NS *corev1.Namespace
	}
	ebd := struct {
		App
		NS1 *corev1.Namespace
		Sub *SubDup
	}{
		NS1: ku.Namespace("same-ns", nil, nil),
	}
	g := goky{
		ar:  &txtar.Archive{},
		o:   exportOption{},
		dup: map[string]struct{}{},
	}

	err := g.encodeStruct(reflect.ValueOf(ebd), "")
	if !errors.Is(err, ErrFieldMissing) {
		t.Error(tu.Callers(), "field not missing")
	}
	tu.AssertErrorMsg(
		t,
		err,
		`"Sub" of type "*kube.SubDup" in "struct { kube.App; NS1 *v1.Namespace; Sub *kube.SubDup }" is nil: missing`,
	)
}

func TestEncode_SubMissing2(t *testing.T) {
	type SubDup struct {
		App
		NS *corev1.Namespace
	}
	ebd := struct {
		App
		NS1 *corev1.Namespace
		Sub *SubDup
	}{
		NS1: ku.Namespace("same-ns", nil, nil),
		Sub: &SubDup{
			// 	NS: ku.Namespace("same-ns", nil, nil),
		},
	}
	g := goky{
		ar:  &txtar.Archive{},
		o:   exportOption{},
		dup: map[string]struct{}{},
	}

	err := g.encodeStruct(reflect.ValueOf(ebd), "")
	if !errors.Is(err, ErrFieldMissing) {
		t.Error(tu.Callers(), "field not missing")
	}
	tu.AssertErrorMsg(
		t,
		err,
		`"Sub" of type "*kube.SubDup" in "struct { kube.App; NS1 *v1.Namespace; Sub *kube.SubDup }" is zero value: missing`,
	)
}

func TestEncode_SubZero(t *testing.T) {
	type SubDup struct {
		App
		NS corev1.Namespace
	}
	ebd := struct {
		App
		NS1 *corev1.Namespace
		Sub SubDup
	}{
		NS1: ku.Namespace("same-ns", nil, nil),
	}
	g := goky{
		ar:  &txtar.Archive{},
		o:   exportOption{},
		dup: map[string]struct{}{},
	}
	err := g.encodeStruct(reflect.ValueOf(ebd), "")
	if !errors.Is(err, ErrFieldMissing) {
		t.Error(tu.Callers(), "field not missing")
	}
	tu.AssertErrorMsg(
		t,
		err,
		`"Sub" of type "kube.SubDup" in "struct { kube.App; NS1 *v1.Namespace; Sub kube.SubDup }" is zero value: missing`,
	)
}

func TestEncode_SubZero2(t *testing.T) {
	type SubDup struct {
		App
		NS corev1.Namespace
	}
	ebd := struct {
		App
		NS1 *corev1.Namespace
		Sub SubDup
	}{
		NS1: ku.Namespace("same-ns", nil, nil),
		Sub: SubDup{},
	}
	g := goky{
		ar:  &txtar.Archive{},
		o:   exportOption{},
		dup: map[string]struct{}{},
	}

	err := g.encodeStruct(reflect.ValueOf(ebd), "")
	if !errors.Is(err, ErrFieldMissing) {
		t.Error(tu.Callers(), "field not missing")
	}
	tu.AssertErrorMsg(
		t,
		err,
		`"Sub" of type "kube.SubDup" in "struct { kube.App; NS1 *v1.Namespace; Sub kube.SubDup }" is zero value: missing`,
	)
}

func TestEncode_SubZero3(t *testing.T) {
	type SubDup struct {
		App
		NS *corev1.Namespace
	}
	ebd := struct {
		App
		NS1 *corev1.Namespace
		Sub *SubDup
	}{
		NS1: ku.Namespace("same-ns", nil, nil),
		Sub: &SubDup{
			NS: &corev1.Namespace{},
		},
	}
	g := goky{
		ar:  &txtar.Archive{},
		o:   exportOption{},
		dup: map[string]struct{}{},
	}

	err := g.encodeStruct(reflect.ValueOf(ebd), "")
	if !errors.Is(err, ErrFieldMissing) {
		t.Error(tu.Callers(), "field not missing")
	}
	tu.AssertErrorMsg(
		t,
		err,
		`encoding field Sub: "NS" of type "*v1.Namespace" in "kube.SubDup" is zero value: missing`,
	)
}

func TestEncode_FieldStringNested(t *testing.T) {
	type SubDup struct {
		App
		Str string
		NS  *corev1.Namespace
	}
	ebd := struct {
		App
		NS1 *corev1.Namespace
		Sub *SubDup
	}{
		NS1: ku.Namespace("same-ns", nil, nil),
		Sub: &SubDup{
			Str: "blabla",
			NS:  &corev1.Namespace{},
		},
	}
	g := goky{
		ar:  &txtar.Archive{},
		o:   exportOption{},
		dup: map[string]struct{}{},
	}

	err := g.encodeStruct(reflect.ValueOf(ebd), "")
	if errors.Is(err, ErrFieldMissing) {
		t.Error(tu.Callers(), "field missing")
	}
	tu.AssertErrorMsg(
		t,
		err,
		"encoding field Sub: unsupported type: Str, type: string, kind: string",
	)
}

func TestEncode_String(t *testing.T) {
	g := goky{
		ar:  &txtar.Archive{},
		o:   exportOption{},
		dup: map[string]struct{}{},
	}

	bla := "blow up"
	err := g.encodeStruct(reflect.ValueOf(bla), "")
	if errors.Is(err, ErrFieldMissing) {
		t.Error(tu.Callers(), "field missing")
	}
	tu.AssertErrorMsg(
		t,
		err,
		"cannot encode non-struct type: [string] blow up",
	)
}

func TestEncode_FieldStringInEmbed(t *testing.T) {
	type SubDup struct {
		App
		Str string
		NS  *corev1.Namespace
	}
	ebd := struct {
		App
		NS1 *corev1.Namespace
		SubDup
	}{
		NS1: ku.Namespace("same-ns", nil, nil),
		SubDup: SubDup{
			Str: "blabla",
			NS:  &corev1.Namespace{},
		},
	}
	g := goky{
		ar:  &txtar.Archive{},
		o:   exportOption{},
		dup: map[string]struct{}{},
	}

	err := g.encodeStruct(reflect.ValueOf(ebd), "")
	if errors.Is(err, ErrFieldMissing) {
		t.Error(tu.Callers(), "field missing")
	}
	tu.AssertErrorMsg(
		t,
		err,
		"encoding embedded SubDup: unsupported type: Str, type: string, kind: string",
	)
}

func TestEncode_Incompatible4(t *testing.T) {
	type SubDup struct {
		App
		NS *corev1.Namespace
	}
	ebd := struct {
		App
		NS1 *corev1.Namespace
		SubDup
	}{
		NS1: nil,
		SubDup: SubDup{
			NS: nil,
		},
	}
	g := goky{
		ar:  &txtar.Archive{},
		o:   exportOption{},
		dup: map[string]struct{}{},
	}

	err := g.encodeStruct(reflect.ValueOf(ebd), "")
	if !errors.Is(err, ErrFieldMissing) {
		t.Error(tu.Callers(), "field not missing")
	}
	tu.AssertErrorMsg(
		t,
		err,
		`"NS1" of type "*v1.Namespace" in "struct { kube.App; NS1 *v1.Namespace; kube.SubDup }" is nil: missing`,
	)
}

func TestEncode_Nil(t *testing.T) {
	g := goky{
		ar:  &txtar.Archive{},
		o:   exportOption{},
		dup: map[string]struct{}{},
	}
	err := g.encodeStruct(reflect.ValueOf(nil), "")
	tu.AssertErrorMsg(
		t, err, "probably a nil value: <invalid reflect.Value>",
	)
}

func TestEncode_NilStruct(t *testing.T) {
	g := goky{
		ar:  &txtar.Archive{},
		o:   exportOption{},
		dup: map[string]struct{}{},
	}
	var v *struct{}
	err := g.encodeStruct(reflect.ValueOf(v), "")
	tu.AssertErrorMsg(t, err, `"*struct {}" is nil: missing`)
}
