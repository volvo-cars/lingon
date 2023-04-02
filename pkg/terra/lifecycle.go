// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	tkihcl "github.com/volvo-cars/lingon/pkg/internal/hcl"
)

// IgnoreChanges takes a list of object attributes to include in the
// `ignore_changes` list for the lifecycle of a resource.
func IgnoreChanges(attrs ...Referencer) LifecyleIgnoreChanges {
	refs := make(LifecyleIgnoreChanges, len(attrs))
	for i, attr := range attrs {
		ref, err := attr.InternalRef()
		if err != nil {
			panic(
				fmt.Sprintf(
					"IgnoreChanges: getting list of attributes: %s",
					err.Error(),
				),
			)
		}
		refs[i] = ref
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
func (l LifecyleIgnoreChanges) InternalTokens() (hclwrite.Tokens, error) {
	if len(l) == 0 {
		return nil, nil
	}
	elems := make([]hclwrite.Tokens, len(l))
	for i, ref := range l {
		tokens, err := tokensForSteps(ref.steps)
		if err != nil {
			return nil, fmt.Errorf(
				"creating tokens for lifecycle ignore_changes: %w",
				err,
			)
		}
		elems[i] = tokens
	}
	return hclwrite.TokensForTuple(elems), nil
}

// ReplaceTriggeredBy takes a list of object attributes to add to the
// `replace_triggered_by` list for the lifecycle of a resource.
func ReplaceTriggeredBy(attrs ...Referencer) LifecycleReplaceTriggeredBy {
	refs := make(LifecycleReplaceTriggeredBy, len(attrs))
	for i, attr := range attrs {
		ref, err := attr.InternalRef()
		if err != nil {
			panic(
				fmt.Sprintf(
					"ReplaceTriggeredBy: getting list of attributes: %s",
					err.Error(),
				),
			)
		}
		refs[i] = ref
	}
	return refs
}

var _ tkihcl.Tokenizer = (*LifecycleReplaceTriggeredBy)(nil)

// LifecycleReplaceTriggeredBy is a list of references to attributes that we
// want to trigger a replacement on if those attributes themselves are replaced.
type LifecycleReplaceTriggeredBy []Reference

// InternalTokens returns the HCL tokens for the `replace_triggered_by` list
func (r LifecycleReplaceTriggeredBy) InternalTokens() (hclwrite.Tokens, error) {
	if len(r) == 0 {
		return nil, nil
	}
	elems := make([]hclwrite.Tokens, len(r))
	for i, ref := range r {
		toks, err := ref.InternalTokens()
		if err != nil {
			return nil, fmt.Errorf(
				"creating tokens for lifecycle replace_triggered_by: %w", err,
			)
		}
		elems[i] = toks
	}
	return hclwrite.TokensForTuple(elems), nil
}

type Lifecycle struct {
	CreateBeforeDestroy BoolValue                   `hcl:"create_before_destroy,attr"`
	PreventDestroy      BoolValue                   `hcl:"prevent_destroy,attr"`
	IgnoreChanges       LifecyleIgnoreChanges       `hcl:"ignore_changes,attr"`
	ReplaceTriggeredBy  LifecycleReplaceTriggeredBy `hcl:"replace_triggered_by,attr"`
}
