// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package nats

import (
	"os"
	"testing"

	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestNatsManifestExport(t *testing.T) {
	_ = os.RemoveAll("out")
	n := New()
	if err := n.Export("out"); err != nil {
		tu.AssertNoError(t, err, "nats")
	}
}

// TODO: THIS IS INTEGRATION and needs KWOK

// func TestNatsDeploy(t *testing.T) {
// 	n := New()
// 	if err := n.Apply(context.Background()); err != nil {
// 		tu.AssertNoError(t, err, "nats")
// 	}
// }
