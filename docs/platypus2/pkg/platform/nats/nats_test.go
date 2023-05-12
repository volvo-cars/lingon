package nats

import (
	"os"
	"testing"

	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestMonitoring(t *testing.T) {
	_ = os.RemoveAll("out")

	n := New()
	if err := n.Export("out"); err != nil {
		tu.AssertNoError(t, err, "prometheus crd")
	}

}
