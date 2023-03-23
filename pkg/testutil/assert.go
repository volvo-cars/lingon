package testutil

import "testing"

func AssertNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

func equal[T comparable](a, b T) bool {
	return a == b
}

func Equal[T comparable](t *testing.T, expected, actual T) {
	t.Helper()
	if !equal(expected, actual) {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}

func NotEqual[T comparable](t *testing.T, expected, actual T) {
	t.Helper()
	if equal(expected, actual) {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}

func Nil(t *testing.T, obj any) {
	t.Helper()
	if obj != nil {
		t.Fatalf("Expected nil, got %v", obj)
	}
}

func NotNil(t *testing.T, obj any) {
	t.Helper()
	switch obj.(type) {
	case string:
		if obj == "" {
			t.Fatalf("Expected not nil, got empty string")
		}
	default:
		if obj == nil {
			t.Fatalf("Expected not nil, got nil")
		}
	}
}

func contains[T comparable](haystack []T, needle T) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}

func Contains[T comparable](t *testing.T, haystack []T, needle T) {
	t.Helper()
	if !contains(haystack, needle) {
		t.Fatalf("Expected %v to contain %v", haystack, needle)
	}
}

func NotContains[T comparable](t *testing.T, haystack []T, needle T) {
	t.Helper()
	if contains(haystack, needle) {
		t.Fatalf("Expected %v to not contain %v", haystack, needle)
	}
}

func True(t *testing.T, condition bool, msg string) {
	t.Helper()
	if !condition {
		t.Fatalf("Expected true, got false: %s", msg)
	}
}

func False(t *testing.T, condition bool, msg string) {
	t.Helper()
	if condition {
		t.Fatalf("Expected false, got true: %s", msg)
	}
}
