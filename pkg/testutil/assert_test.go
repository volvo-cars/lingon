// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package testutil

import (
	"io/fs"
	"os"
	"testing"
)

func TestAssertEqual(t *testing.T) {
	type testCase[C comparable] struct {
		name     string
		expected string
		actual   string
	}
	tests := []testCase[string]{
		{
			name:     "string",
			expected: "hello",
			actual:   "hello",
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				AssertEqual(t, tt.expected, tt.actual)
			},
		)
	}
}

func TestAssertEqualSlice(t *testing.T) {
	type testCase[C comparable] struct {
		name     string
		expected []C
		actual   []C
	}
	tests := []testCase[int]{
		{
			name:     "int",
			expected: []int{1, 2, 3},
			actual:   []int{1, 2, 3},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				AssertEqualSlice(t, tt.expected, tt.actual)
			},
		)
	}
}

func TestAssertError(t *testing.T) {
	type args struct {
		t   *testing.T
		err error
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				AssertErrorMsg(tt.args.t, tt.args.err, tt.args.msg)
			},
		)
	}
}

func TestAssertNoError(t *testing.T) {
	type args struct {
		t   *testing.T
		err error
		msg []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				AssertNoError(tt.args.t, tt.args.err, tt.args.msg...)
			},
		)
	}
}

func TestContains(t *testing.T) {
	type testCase struct {
		name     string
		haystack []string
		needle   string
	}
	tests := []testCase{
		{
			name:     "contains",
			haystack: []string{"hey", "ho"},
			needle:   "hey",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				Contains(t, tt.haystack, tt.needle)
			},
		)
	}
}

func TestErrorIs(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected error
	}{
		{
			name:     "oops",
			err:      fs.ErrNotExist,
			expected: os.ErrNotExist,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ErrorIs(t, tt.err, tt.expected)
			},
		)
	}
}

func TestFalse(t *testing.T) {
	tests := []struct {
		name      string
		condition bool
	}{
		{
			name:      "is false",
			condition: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				False(t, tt.condition, "should be false")
			},
		)
	}
}

func TestIsEqual(t *testing.T) {
	type testCase[C comparable] struct {
		name     string
		expected C
		actual   C
	}
	tests := []testCase[int]{
		{
			name:     "int",
			expected: 10,
			actual:   10,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				IsEqual(t, tt.expected, tt.actual)
			},
		)
	}
}

func TestIsNil(t *testing.T) {
	tests := []struct {
		name string
		obj  any
	}{
		{
			name: "is nil",
			obj:  nil,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				IsNil(t, tt.obj)
			},
		)
	}
}

func TestIsNotEqual(t *testing.T) {
	type testCase[C comparable] struct {
		name     string
		expected C
		actual   C
	}
	tests := []testCase[int]{
		{
			name:     "int",
			expected: 10,
			actual:   11,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				IsNotEqual(t, tt.expected, tt.actual)
			},
		)
	}
}

func TestIsNotNil(t *testing.T) {
	tests := []struct {
		name string
		obj  any
	}{
		{
			name: "not nil",
			obj:  []string{"hello"},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				IsNotNil(t, tt.obj)
			},
		)
	}
}

func TestNotContains(t *testing.T) {
	type testCase[C comparable] struct {
		name     string
		haystack []C
		needle   C
	}
	tests := []testCase[string]{
		{
			name:     "string",
			haystack: []string{"not", "here"},
			needle:   "where",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				NotContains(t, tt.haystack, tt.needle)
			},
		)
	}
}

func TestTrue(t *testing.T) {
	tests := []struct {
		name      string
		condition bool
	}{
		{
			name:      "is true",
			condition: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				True(t, tt.condition, "should be true")
			},
		)
	}
}
