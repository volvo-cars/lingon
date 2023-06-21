// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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

func TestResources(t *testing.T) {
	type args struct {
		cpuWant string
		memWant string
		cpuMax  string
		memMax  string
	}
	tests := []struct {
		name string
		args args
		want corev1.ResourceRequirements
	}{
		// TODO: Add test cases.
		{
			name: "ram cpu",
			args: args{
				cpuWant: "2",
				memWant: "2Gi",
				cpuMax:  "4",
				memMax:  "4Gi",
			},
			want: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("2"),
					corev1.ResourceMemory: resource.MustParse("2Gi"),
				},
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("4"),
					corev1.ResourceMemory: resource.MustParse("4Gi"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := Resources(
					tt.args.cpuWant,
					tt.args.memWant,
					tt.args.cpuMax,
					tt.args.memMax,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Resources() = %v, want %v", got, tt.want)
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := SimpleDeployment(
					tt.args.name,
					tt.args.namespace,
					tt.args.labels,
					tt.args.replicas,
					tt.args.image,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("SimpleDeployment() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
