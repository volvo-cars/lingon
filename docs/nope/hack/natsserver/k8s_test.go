package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/volvo-cars/lingon/pkg/kube"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestRegistry(t *testing.T) {
	out := filepath.Join("out", "k8s")

	err := os.RemoveAll(out)
	tu.AssertNoError(t, err)
	reg := NewNATSServer()
	err = kube.Export(reg, kube.WithExportOutputDirectory(out))
	tu.AssertNoError(t, err)
}
