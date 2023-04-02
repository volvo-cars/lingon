// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"fmt"
	"testing"

	tu "github.com/volvo-cars/lingon/pkg/testutil"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ExampleList_string() {
	s := List(
		String("a"),
		String("b"),
	)
	fmt.Println(exampleTokensOrError(s))
	// Output: ["a", "b"]
}

func ExampleList_number() {
	s := List(
		Number(0),
		Number(1),
	)
	fmt.Println(exampleTokensOrError(s))
	// Output: [0, 1]
}

func ExampleList_bool() {
	s := List(
		Bool(false),
		Bool(true),
	)
	fmt.Println(exampleTokensOrError(s))
	// Output: [false, true]
}

func ExampleList_ref() {
	// Create some dummy references
	refA := ReferenceAsString(ReferenceResource(&dummyResource{}))
	refB := ReferenceAsString(ReferenceDataResource(&dummyDataResource{}))

	s := List(refA, refB)
	fmt.Println(exampleTokensOrError(s))
	// Output: [dummy.dummy, data.dummy.dummy]
}

func ExampleList_mixed() {
	s := List(
		String("a"),
		Number(1).AsString(),
		ReferenceAsString(ReferenceResource(&dummyResource{})),
	)
	fmt.Println(exampleTokensOrError(s))
	// Output: ["a", "1", dummy.dummy]
}

func ExampleList_index() {
	// Create a reference list of string and Splat() it
	l := ReferenceAsList[StringValue](
		ReferenceResource(&dummyResource{}),
	)
	index := l.Index(0)
	fmt.Println(exampleTokensOrError(index))
	// Output: dummy.dummy[0]
}

func ExampleList_splat() {
	// Create a reference list of string and Splat() it
	l := ReferenceAsList[StringValue](
		ReferenceResource(&dummyResource{}),
	)
	splat := l.Splat()
	// Convert "splatted" list back to a List
	var ls ListValue[StringValue] //nolint:gosimple
	ls = CastAsList(splat)
	fmt.Println(exampleTokensOrError(ls))
	// Output: dummy.dummy[*]
}

func ExampleList_splatNested() {
	// Create a reference list of a list of string and Splat() it
	l := ReferenceAsList[ListValue[StringValue]](
		ReferenceResource(&dummyResource{}),
	)
	splat := l.Splat()
	// Convert "splatted" list back to a List of List
	var ls ListValue[ListValue[StringValue]] //nolint:gosimple
	ls = CastAsList(
		splat,
	)
	fmt.Println(exampleTokensOrError(ls))
	// Output: dummy.dummy[*]
}

var _ Value[Attrs] = (*Attrs)(nil)

// Attrs is a dummy implementation of an attribute that is generated for
// Terraform objects.
type Attrs struct {
	ref Reference
}

func (a Attrs) InternalTokens() (hclwrite.Tokens, error) {
	return a.ref.InternalTokens()
}

func (a Attrs) InternalRef() (Reference, error) {
	return a.ref.copy(), nil
}

func (a Attrs) InternalWithRef(ref Reference) Attrs {
	return Attrs{ref: ref}
}

func (a Attrs) Name() StringValue {
	return ReferenceAsString(a.ref.Append("name"))
}

func TestCustomTypes(t *testing.T) {
	l := ReferenceAsList[Attrs](ReferenceResource(&dummyResource{}))
	index := l.Index(0)
	name := index.Name()
	tu.AssertEqual(
		t, exampleTokensOrError(name), "dummy.dummy[0].name",
	)
	// Make sure index was not updated after updating name
	tu.AssertEqual(
		t, exampleTokensOrError(index), "dummy.dummy[0]",
	)
}
