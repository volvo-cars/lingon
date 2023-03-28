// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package terra

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExport(t *testing.T) {
	type simpleStack struct {
		DummyBaseStack
		DummyRes  *dummyResource     `validate:"required"`
		DummyData *dummyDataResource `validate:"required"`
	}
	dr := &dummyResource{}
	ddr := &dummyDataResource{}
	st := simpleStack{
		DummyBaseStack: newDummyBaseStack(),
		DummyRes:       dr,
		DummyData:      ddr,
	}
	var b bytes.Buffer
	err := encodeStack(&st, &b)
	require.NoError(t, err)
	fmt.Println(b.String())
}
