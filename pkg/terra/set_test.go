// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"fmt"
)

func ExampleSet_string() {
	s := Set(
		String("a"),
		String("b"),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: ["a", "b"]
}

func ExampleSet_number() {
	s := Set(
		Number(0),
		Number(1),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: [0, 1]
}

func ExampleSet_bool() {
	s := Set(
		Bool(false),
		Bool(true),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: [false, true]
}

func ExampleSet_ref() {
	s := Set(
		ReferenceString(newRef("a")),
		ReferenceString(newRef("b")),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: [a, b]
}

func ExampleSet_mixed() {
	s := Set(
		String("a"),
		ReferenceString(newRef("a")),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: ["a", a]
}

func ExampleSet_index() {
	// Create a reference set of string and Splat() it
	l := ReferenceSet[StringValue](
		newRef("a", "b", "c"),
	)
	index := l.Index(0)
	fmt.Println(string(index.InternalTokens().Bytes()))
	// Output: a.b.c[0]
}

func ExampleSet_splat() {
	// Create a reference set of string and Splat() it
	l := ReferenceSet[StringValue](
		newRef("a", "b", "c"),
	)
	splat := l.Splat()
	// Convert "splatted" set back to a Set
	var ls SetValue[StringValue] //nolint:gosimple
	ls = CastAsSet(splat)
	fmt.Println(string(ls.InternalTokens().Bytes()))
	// Output: a.b.c[*]
}

func ExampleSet_splatNested() {
	// Create a reference set of a set of string and Splat() it
	l := ReferenceSet[SetValue[StringValue]](
		newRef("a", "b", "c"),
	)
	splat := l.Splat()
	// Convert "splatted" set back to a Set of Set
	var ls SetValue[SetValue[StringValue]] //nolint:gosimple
	ls = CastAsSet(
		splat,
	)
	fmt.Println(string(ls.InternalTokens().Bytes()))
	// Output: a.b.c[*]
}
