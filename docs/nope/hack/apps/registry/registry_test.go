package registry_test

import (
	"os"
	"testing"

	"github.com/volvo-cars/lingon/pkg/kube"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	"github.com/volvo-cars/nope/hack/apps/registry"
)

func TestRegistry(t *testing.T) {
	err := os.RemoveAll("out")
	tu.AssertNoError(t, err)
	reg := registry.NewRegistry()
	err = kube.Export(reg, kube.WithExportOutputDirectory("out"))
	tu.AssertNoError(t, err)
}
