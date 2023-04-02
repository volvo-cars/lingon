// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package updater

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetLatestVersion(t *testing.T) {
	lv, err := GetLatestVersion("grafana/grafana", ">9.3.6, <9.4")
	require.NoError(t, err)
	fmt.Println(lv)
}
