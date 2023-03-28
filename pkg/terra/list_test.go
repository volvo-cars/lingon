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

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: ["a", "b"]
}

func ExampleList_number() {
	s := List(
		Number(0),
		Number(1),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: [0, 1]
}

func ExampleList_bool() {
	s := List(
		Bool(false),
		Bool(true),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: [false, true]
}

func ExampleList_ref() {
	s := List(
		ReferenceString(ReferenceAttribute("a")),
		ReferenceString(ReferenceAttribute("b")),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: [a, b]
}

func ExampleList_mixed() {
	s := List(
		String("a"),
		Number(1).AsString(),
		ReferenceString(ReferenceAttribute("a")),
	)

	fmt.Println(string(s.InternalTokens().Bytes()))
	// Output: ["a", "1", a]
}

var _ Value[Attrs] = (*Attrs)(nil)

// Attrs is a dummy implementation of an attribute that is generated for
// Terraform objects.
type Attrs struct {
	ref Reference
}

func (a Attrs) InternalTokens() hclwrite.Tokens {
	return a.ref.InternalTokens()
}

func (a Attrs) InternalWithRef(ref Reference) Attrs {
	return Attrs{ref: ref}
}

func (a Attrs) Name() StringValue {
	return ReferenceString(a.ref.Append("name"))
}

func TestCustomTypes(t *testing.T) {
	l := ReferenceList[Attrs](ReferenceAttribute("bla_type", "name"))
	index := l.Index(0)
	name := index.Name()
	tu.AssertEqual(
		t, string(name.InternalTokens().Bytes()),
		"bla_type.name[0].name",
	)
	// Make sure index was not updated after updating name
	tu.AssertEqual(
		t, string(index.InternalTokens().Bytes()),
		"bla_type.name[0]",
	)
}
