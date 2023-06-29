package vmop

import (
	"os"
	"testing"

	"github.com/volvo-cars/lingon/pkg/kube"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestVMOp(t *testing.T) {
	tu.AssertNoError(t, os.RemoveAll("out"))
	tu.AssertNoError(t, kube.Export(New()))
}
