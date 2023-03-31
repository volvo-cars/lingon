// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"testing"

	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestReferenceCopy(t *testing.T) {
	ref := ReferenceResource(&dummyResource{})
	ref2 := ref.copy()
	tu.Equal(t, ref.underlyingType, ref2.underlyingType)
	tu.Equal(t, ref.res, ref2.res)
	tu.Equal(t, ref.data, ref2.data)
}

func TestReferenceTokens(t *testing.T) {
	ref := ReferenceResource(&dummyResource{})
	tu.Equal(t, "dummy.dummy", string(ref.InternalTokens().Bytes()))

	appendRef := ref.Append("abc")
	tu.Equal(t, "dummy.dummy.abc", string(appendRef.InternalTokens().Bytes()))

	keyRef := ref.key("a")
	tu.Equal(t, "dummy.dummy[\"a\"]", string(keyRef.InternalTokens().Bytes()))

	indexRef := ref.index(3)
	tu.Equal(t, "dummy.dummy[3]", string(indexRef.InternalTokens().Bytes()))

	splatRef := ref.splat()
	tu.Equal(t, "dummy.dummy[*]", string(splatRef.InternalTokens().Bytes()))

	// 	Check original has not been updated
	tu.Equal(t, "dummy.dummy", string(ref.InternalTokens().Bytes()))
}
