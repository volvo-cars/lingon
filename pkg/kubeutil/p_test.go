// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	"testing"

	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestP(t *testing.T) {
	type testCase[T any] struct {
		name string
		args T
		want *T
	}
	tests := []testCase[int]{
		{
			name: "int",
			args: 10,
			want: func(a int) *int { return &a }(10),
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := P(tt.args)
				tu.AssertEqual(t, tt.want, got)
			},
		)
	}
}
