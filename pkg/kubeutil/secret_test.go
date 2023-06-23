// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestSecret(t *testing.T) {
	type args struct {
		name      string
		namespace string
		data      map[string][]byte
	}
	tests := []struct {
		name string
		args args
		want *corev1.Secret
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := Secret(
					tt.args.name,
					tt.args.namespace,
					tt.args.data,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Secret() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
