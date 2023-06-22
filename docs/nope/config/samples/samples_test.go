package samples_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	tu "github.com/volvo-cars/lingon/pkg/testutil"
	"github.com/volvo-cars/nope/config/samples"
)

func TestExport(t *testing.T) {

	outDir := "out"

	rmErr := os.RemoveAll(outDir)
	tu.AssertNoError(t, rmErr, "removing out directory")

	app := samples.NewApp()
	err := kube.Export(
		app,
		kube.WithExportOutputDirectory(outDir),
		kube.WithExportNameFileFunc(func(m *kubeutil.Metadata) string {
			return strings.ToLower(fmt.Sprintf("%s_%s.yaml", m.Kind, m.Meta.Name))
		}),
	)
	tu.AssertNoError(t, err, "exporting samples")
}
