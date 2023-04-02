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

	fmt.Println(exampleTokensOrError(s))
	// Output: ["a", "b"]
}

func ExampleSet_number() {
	s := Set(
		Number(0),
		Number(1),
	)

	fmt.Println(exampleTokensOrError(s))
	// Output: [0, 1]
}

func ExampleSet_bool() {
	s := Set(
		Bool(false),
		Bool(true),
	)

	fmt.Println(exampleTokensOrError(s))
	// Output: [false, true]
}

func ExampleSet_ref() {
	s := Set(
		ReferenceAsString(ReferenceResource(&dummyResource{})),
		ReferenceAsString(ReferenceDataResource(&dummyDataResource{})),
	)

	fmt.Println(exampleTokensOrError(s))
	// Output: [dummy.dummy, data.dummy.dummy]
}

func ExampleSet_mixed() {
	s := Set(
		String("a"),
		ReferenceAsString(ReferenceResource(&dummyResource{})),
	)

	fmt.Println(exampleTokensOrError(s))
	// Output: ["a", dummy.dummy]
}

func ExampleSet_index() {
	// Create a reference set of string and Splat() it
	l := ReferenceAsSet[StringValue](
		ReferenceResource(&dummyResource{}),
	)
	index := l.Index(0)
	fmt.Println(exampleTokensOrError(index))
	// Output: dummy.dummy[0]
}

func ExampleSet_splat() {
	// Create a reference set of string and Splat() it
	l := ReferenceAsSet[StringValue](
		ReferenceResource(&dummyResource{}),
	)
	splat := l.Splat()
	// Convert "splatted" set back to a Set
	var ls SetValue[StringValue] //nolint:gosimple
	ls = CastAsSet(splat)
	fmt.Println(exampleTokensOrError(ls))
	// Output: dummy.dummy[*]
}

func ExampleSet_splatNested() {
	// Create a reference set of a set of string and Splat() it
	l := ReferenceAsSet[SetValue[StringValue]](
		ReferenceResource(&dummyResource{}),
	)
	splat := l.Splat()
	// Convert "splatted" set back to a Set of Set
	var ls SetValue[SetValue[StringValue]] //nolint:gosimple
	ls = CastAsSet(
		splat,
	)
	fmt.Println(exampleTokensOrError(ls))
	// Output: dummy.dummy[*]
}
