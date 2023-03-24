// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package terrajen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProviderShortName(t *testing.T) {
	expectedName := "some_resource"
	trimmedName := providerShortName("aws_" + expectedName)
	assert.Equal(t, expectedName, trimmedName)
}
