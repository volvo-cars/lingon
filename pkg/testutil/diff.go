// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package testutil

import (
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Diff compares two items and returns a human-readable diff string. If the
// items are equal, the string is empty.
func Diff[T any](got, want T, opts ...cmp.Option) string {
	// nolint: gocritic
	oo := append(
		opts,
		cmp.Exporter(func(reflect.Type) bool { return true }),
		cmpopts.EquateEmpty(),
	)

	diff := cmp.Diff(got, want, oo...)
	if diff != "" {
		return "\n-got +want\n" + diff
	}

	return ""
}

// Callers prints the stack trace of everything up til the line where Callers()
// was invoked.
func Callers() string {
	var pc [50]uintptr
	n := runtime.Callers(
		2,
		pc[:],
	) //nolint:gomnd    // skip runtime.Callers + Callers
	callsites := make([]string, 0, n)
	frames := runtime.CallersFrames(pc[:n])

	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		callsites = append(callsites, frame.File+":"+strconv.Itoa(frame.Line))
	}

	callsites = callsites[:len(callsites)-1] // skip testing.tRunner
	if len(callsites) == 1 {
		return ""
	}

	var b strings.Builder

	for i := len(callsites) - 1; i >= 0; i-- {
		if b.Len() > 0 {
			b.WriteString(" -> ")
		}

		b.WriteString(filepath.Base(callsites[i]))
	}

	return "\n" + b.String() + ":"
}

func AssertEqual[C comparable](t *testing.T, got, want C) {
	if diff := Diff(got, want); diff != "" {
		t.Error(Callers(), diff)
	}
}
