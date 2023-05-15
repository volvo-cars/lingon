// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"testing"

	tu "github.com/volvo-cars/lingon/pkg/testutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestTypeMeta(t *testing.T) {
	tests := []struct {
		name string
		kind string
		want v1.TypeMeta
	}{
		{
			name: "deploy",
			kind: "Deployment",
			want: v1.TypeMeta{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := TypeMeta(tt.kind)
				tu.AssertEqual(t, tt.want, got)
			},
		)
	}
}
