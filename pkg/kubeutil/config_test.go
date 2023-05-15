// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConfigAndMount_ConfigMap(t *testing.T) {
	type fields struct {
		ObjectMeta  metav1.ObjectMeta
		VolumeMount corev1.VolumeMount
		Data        map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   corev1.ConfigMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				m := ConfigAndMount{
					ObjectMeta:  tt.fields.ObjectMeta,
					VolumeMount: tt.fields.VolumeMount,
					Data:        tt.fields.Data,
				}
				if got := m.ConfigMap(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ConfigMap() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestConfigAndMount_HashEnv(t *testing.T) {
	type fields struct {
		ObjectMeta  metav1.ObjectMeta
		VolumeMount corev1.VolumeMount
		Data        map[string]string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   corev1.EnvVar
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				m := ConfigAndMount{
					ObjectMeta:  tt.fields.ObjectMeta,
					VolumeMount: tt.fields.VolumeMount,
					Data:        tt.fields.Data,
				}
				if got := m.HashEnv(tt.args.name); !reflect.DeepEqual(
					got,
					tt.want,
				) {
					t.Errorf("HashEnv() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestConfigAndMount_VolumeAndMount(t *testing.T) {
	type fields struct {
		ObjectMeta  metav1.ObjectMeta
		VolumeMount corev1.VolumeMount
		Data        map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   VolumeAndMount
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				m := ConfigAndMount{
					ObjectMeta:  tt.fields.ObjectMeta,
					VolumeMount: tt.fields.VolumeMount,
					Data:        tt.fields.Data,
				}
				if got := m.VolumeAndMount(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("VolumeAndMount() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestDataConfigMap(t *testing.T) {
	type args struct {
		name        string
		namespace   string
		labels      map[string]string
		annotations map[string]string
		data        map[string]string
	}
	tests := []struct {
		name string
		args args
		want *corev1.ConfigMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := DataConfigMap(
					tt.args.name,
					tt.args.namespace,
					tt.args.labels,
					tt.args.annotations,
					tt.args.data,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("DataConfigMap() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
