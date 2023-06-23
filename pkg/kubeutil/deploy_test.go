// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMergeLabels(t *testing.T) {
	tests := []struct {
		name   string
		labels []map[string]string
		want   map[string]string
	}{
		{
			name: "merge",
			labels: []map[string]string{
				{"key1": "val1"},
				{"key2": "val2", "key3": "val3"},
			},
			want: map[string]string{
				"key1": "val1",
				"key2": "val2",
				"key3": "val3",
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if diff := cmp.Diff(
					tt.want,
					MergeLabels(tt.labels...),
				); diff != "" {
					t.Error(diff)
				}
			},
		)
	}
}

func TestNamespace(t *testing.T) {
	type args struct {
		name        string
		labels      map[string]string
		annotations map[string]string
	}
	tests := []struct {
		name string
		args args
		want *corev1.Namespace
	}{
		{
			name: "ns",
			args: args{
				name:        "testns",
				labels:      map[string]string{"mylabel": "labelvalue"},
				annotations: map[string]string{"annot": "tation"},
			},
			want: &corev1.Namespace{
				TypeMeta: TypeNamespaceV1,
				ObjectMeta: metav1.ObjectMeta{
					Name:        "testns",
					Labels:      map[string]string{"mylabel": "labelvalue"},
					Annotations: map[string]string{"annot": "tation"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := Namespace(
					tt.args.name,
					tt.args.labels,
					tt.args.annotations,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Namespace() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestSimpleDeployment(t *testing.T) {
	type args struct {
		name      string
		namespace string
		labels    map[string]string
		replicas  int32
		image     string
	}
	tests := []struct {
		name string
		args args
		want *appsv1.Deployment
	}{
		{
			name: "simple",
			args: args{
				name:      "depl",
				namespace: "default",
				labels:    map[string]string{"app": "depl"},
				replicas:  1,
				image:     "nginx",
			},
			want: &appsv1.Deployment{
				TypeMeta: TypeDeploymentV1,
				ObjectMeta: ObjectMeta(
					"depl", "default",
					map[string]string{"app": "depl"}, nil,
				),
				Spec: appsv1.DeploymentSpec{
					Replicas: P(int32(1)),
					Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "depl"}},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"app": "depl"},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{Name: "depl", Image: "nginx"},
							},
							ServiceAccountName: "depl",
						},
					},
					Strategy:                appsv1.DeploymentStrategy{},
					MinReadySeconds:         0,
					RevisionHistoryLimit:    nil,
					Paused:                  false,
					ProgressDeadlineSeconds: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := SimpleDeployment(
					tt.args.name,
					tt.args.namespace,
					tt.args.labels,
					tt.args.replicas,
					tt.args.image,
				)
				tu.AssertEqual(t, tt.want, got)
			},
		)
	}
}
