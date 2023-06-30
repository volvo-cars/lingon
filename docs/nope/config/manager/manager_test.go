package manager

import (
	"os"
	"testing"

	"github.com/volvo-cars/lingon/pkg/kube"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestManager(t *testing.T) {
	out := "out"
	err := os.RemoveAll(out)
	tu.AssertNoError(t, err)

	m := NewManager()
	err = kube.Export(m, kube.WithExportOutputDirectory(out))
	tu.AssertNoError(t, err)
}
