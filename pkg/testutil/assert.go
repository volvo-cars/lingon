// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package testutil

import (
	"errors"
	"testing"
)

// AssertEqual checks if the expected and actual are equal. Errors if not.
func AssertEqual[C comparable](t *testing.T, expected, actual C) {
	t.Helper()
	if diff := Diff(actual, expected); diff != "" {
		t.Fatal(Callers(), diff)
	}
}

// AssertEqualSlice checks if the expected and actual slices are equal. Errors if not.
func AssertEqualSlice[C comparable](t *testing.T, expected, actual []C) {
	t.Helper()
	if diff := Diff(actual, expected); diff != "" {
		t.Fatal(Callers(), diff)
	}
}

// AssertNoError checks if the error is nil.
func AssertNoError(t *testing.T, err error, msg ...string) {
	t.Helper()
	if err != nil {
		t.Fatal(Callers(), err, msg)
	}
}

// AssertErrorMsg checks if the error is not nil.
func AssertErrorMsg(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		t.Fatal(Callers(), msg, errors.New("error was expected but is nil"))
	}
	if diff := Diff(err.Error(), msg); diff != "" {
		t.Fatal(Callers(), diff)
	}
}

func equal[C comparable](a, b C) bool {
	return a == b
}

// IsEqual checks if the expected and actual are equal. Errors if not.
func IsEqual[C comparable](t *testing.T, expected, actual C) {
	t.Helper()
	if !equal(expected, actual) {
		t.Fatal(Callers(), Diff(actual, expected))
	}
}

// IsNotEqual checks if the expected and actual are not equal. Errors if not.
func IsNotEqual[C comparable](t *testing.T, expected, actual C) {
	t.Helper()
	if equal(expected, actual) {
		t.Fatal(Callers(), Diff(actual, expected))
	}
}

// IsNil checks if the obj is nil. Errors if not.
func IsNil(t *testing.T, obj any) {
	t.Helper()
	switch o := obj.(type) {
	case bool:
		t.Fatalf("expected nil, got %v", o)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
		uintptr, float32, float64, complex64, complex128:
		t.Fatalf("expected nil, got %v", o)
	case string:
		t.Fatalf("expected nil, got %v", o)
	default:
		if obj != nil {
			t.Fatalf("expected nil, got %v", obj)
		}
	}
}

// IsNotNil checks if the obj is not nil. Errors if it is.
func IsNotNil(t *testing.T, obj any) {
	t.Helper()
	switch o := obj.(type) {
	case bool:
		t.Fatalf("expected not nil, got %v", o)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
		uintptr, float32, float64, complex64, complex128:
		t.Fatalf("expected not nil, got %v", o)
	case string:
		t.Fatalf("expected not nil, got %v", o)
	default:
		if o == nil {
			t.Fatalf("expected not nil, got nil")
		}
	}
}

func contains[C comparable](haystack []C, needle C) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}

// Contains checks if the needle is in the haystack. Errors if not.
func Contains[C comparable](t *testing.T, haystack []C, needle C) {
	t.Helper()
	if !contains(haystack, needle) {
		t.Fatalf("expected %v to contain %v", haystack, needle)
	}
}

// NotContains checks if the needle is not in the haystack. Errors if it is.
func NotContains[C comparable](t *testing.T, haystack []C, needle C) {
	t.Helper()
	if contains(haystack, needle) {
		t.Fatalf("expected %v to not contain %v", haystack, needle)
	}
}

func True(t *testing.T, condition bool, msg string) {
	t.Helper()
	if !condition {
		t.Fatalf("expected true, got false: %s", msg)
	}
}

func False(t *testing.T, condition bool, msg string) {
	t.Helper()
	if condition {
		t.Fatalf("expected false, got true: %s", msg)
	}
}

func ErrorIs(t *testing.T, err, expected error) {
	t.Helper()
	if !errors.Is(err, expected) {
		t.Fatalf(`expected error "%v", got "%v"`, expected, err)
	}
}
