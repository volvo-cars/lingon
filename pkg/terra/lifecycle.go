// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	tkihcl "github.com/volvo-cars/lingon/pkg/internal/hcl"
)

// IgnoreChanges takes a list of object attributes to include in the
// `ignore_changes` list for the lifecycle of a resource.
func IgnoreChanges(attrs ...Referencer) LifecyleIgnoreChanges {
	refs := make(LifecyleIgnoreChanges, len(attrs))
	for i, attr := range attrs {
		// Make sure we get a copy of the reference
		refs[i] = attr.InternalRef()
	}
	return refs
}

var _ tkihcl.Tokenizer = (*LifecyleIgnoreChanges)(nil)

// LifecyleIgnoreChanges is a list of references to attributes that we want to
// ignore in the `lifecycle` block
type LifecyleIgnoreChanges []Reference

// InternalTokens only returns the relative address of a reference.
// This is due to the specification for the `ignore_changes` list inside the
// `lifecycle` block. E.g.
//
//	resource "aws_instance" "example" {
//	  # ...
//	   lifecycle {
//	    ignore_changes = [
//	      # Use a relative reference to the tags,
//	      # i.e. not `aws_instance.example.tags`
//	      tags,
//	    ]
//	  }
//	}
func (l LifecyleIgnoreChanges) InternalTokens() hclwrite.Tokens {
	if len(l) == 0 {
		return nil
	}
	elems := make([]hclwrite.Tokens, len(l))
	for i, ref := range l {
		// Ensure the first element in the traversal is a root traversal.
		tr, err := reRootTraversal(ref.tr)
		if err != nil {
			panic(
				fmt.Sprintf(
					"LifeCycleIngoreChanges: cannot creat tokens"+
						" for traversal: %s", err.Error(),
				),
			)
		}
		elems[i] = hclwrite.TokensForTraversal(tr)
	}
	return hclwrite.TokensForTuple(elems)
}

// ReplaceTriggeredBy takes a list of object attributes to add to the
// `replace_triggered_by` list for the lifecycle of a resource.
func ReplaceTriggeredBy(attrs ...Referencer) LifecycleReplaceTriggeredBy {
	refs := make(LifecycleReplaceTriggeredBy, len(attrs))
	for i, attr := range attrs {
		refs[i] = attr.InternalRef()
	}
	return refs
}

var _ tkihcl.Tokenizer = (*LifecycleReplaceTriggeredBy)(nil)

// LifecycleReplaceTriggeredBy is a list of references to attributes that we
// want to trigger a replacement on if those attributes themselves are replaced.
type LifecycleReplaceTriggeredBy []Reference

// InternalTokens returns the HCL tokens for the `replace_triggered_by` list
func (r LifecycleReplaceTriggeredBy) InternalTokens() hclwrite.Tokens {
	if len(r) == 0 {
		return nil
	}
	elems := make([]hclwrite.Tokens, len(r))
	for i, ref := range r {
		elems[i] = ref.InternalTokens()
	}
	return hclwrite.TokensForTuple(elems)
}

type Lifecycle struct {
	CreateBeforeDestroy BoolValue                   `hcl:"create_before_destroy,attr"`
	PreventDestroy      BoolValue                   `hcl:"prevent_destroy,attr"`
	IgnoreChanges       LifecyleIgnoreChanges       `hcl:"ignore_changes,attr"`
	ReplaceTriggeredBy  LifecycleReplaceTriggeredBy `hcl:"replace_triggered_by,attr"`
}

// reRootTraversal takes a hcl.Traversal that may not have a root traverser
// as it's first element, and converts it into a root.
// This is needed because when turning a hcl.Traversal into hclwrite.Tokens
// a "." is prefixed to the tokens if the first step in the hcl.Traversal is not
// a hcl.TraverseRoot.
func reRootTraversal(tr hcl.Traversal) (hcl.Traversal, error) {
	if len(tr) == 0 {
		return nil, errors.New("cannot re-root a traversal with no steps in it")
	}
	switch t := tr[0].(type) {
	case hcl.TraverseRoot:
		return tr, nil
	case hcl.TraverseAttr:
		// Convert the attribute into a root
		tr[0] = hcl.TraverseRoot{Name: t.Name}
		return tr, nil
	case hcl.TraverseIndex, hcl.TraverseSplat:
		return nil, errors.New(
			"cannot re-root a traversal with first" +
				" element index or splat",
		)
	default:
		return nil, errors.New("unknown traverser")
	}
}
