// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kubeutil

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/txtar"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestTxtar2YAML(t *testing.T) {
	a, err := os.ReadFile("testdata/apps.txt")
	tu.AssertNoError(t, err, "reading txtar file")
	ar := txtar.Parse(a)
	ay, err := os.ReadFile("testdata/apps.yaml")
	tu.AssertNoError(t, err, "reading yaml file")
	ar2y := Txtar2YAML(ar)
	tu.AssertEqual(t, string(ay), string(ar2y))
}

func TestTxtar2JSON(t *testing.T) {
	a, err := os.ReadFile("testdata/apps.txt")
	tu.AssertNoError(t, err, "reading txtar file")
	ar := txtar.Parse(a)
	ar, err = TxtarYAML2TxtarJSON(ar)
	tu.AssertNoError(t, err, "converting to json")
	aj, err := os.ReadFile("testdata/apps.json")
	tu.AssertNoError(t, err, "reading json file")
	ar2j := Txtar2JSON(ar)
	tu.AssertEqual(t, string(aj), string(ar2j))
}
