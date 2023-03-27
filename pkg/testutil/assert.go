// Copyright (c) Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package testutil

import (
	"errors"
	"testing"
)

func AssertEqual[C comparable](t *testing.T, expected, actual C) {
	if diff := Diff(actual, expected); diff != "" {
		t.Error(Callers(), diff)
	}
}

func AssertNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Error(Callers(), msg, err)
	}
}

func equal[C comparable](a, b C) bool {
	return a == b
}

func Equal[C comparable](t *testing.T, expected, actual C) {
	t.Helper()
	if !equal(expected, actual) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func NotEqual[C comparable](t *testing.T, expected, actual C) {
	t.Helper()
	if equal(expected, actual) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func Nil(t *testing.T, obj any) {
	t.Helper()
	switch obj.(type) {
	case string:
		if obj != "" {
			t.Fatalf("expected empty string, got %q", obj)
		}
	default:
		if obj != nil {
			t.Fatalf("expected nil, got %v", obj)
		}
	}
}

func NotNil(t *testing.T, obj any) {
	t.Helper()
	switch obj.(type) {
	case string:
		if obj == "" {
			t.Fatalf("expected not nil, got empty string")
		}
	default:
		if obj == nil {
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

func Contains[C comparable](t *testing.T, haystack []C, needle C) {
	t.Helper()
	if !contains(haystack, needle) {
		t.Fatalf("expected %v to contain %v", haystack, needle)
	}
}

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
	if errors.Is(err, expected) {
		t.Fatalf("expected error %v, got %v", expected, err)
	}
}
