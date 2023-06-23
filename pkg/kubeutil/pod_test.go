// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

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
