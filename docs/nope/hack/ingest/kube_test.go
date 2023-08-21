package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/volvo-cars/lingon/pkg/kube"
)

func TestKubeManifest(t *testing.T) {
	out := filepath.Join("out", "kube")
	err := os.RemoveAll(out)
	if err != nil {
		t.Fatalf("removing output directory: %v", err)
	}
	reg := NewIngestApp()
	err = kube.Export(reg, kube.WithExportOutputDirectory(out))
	if err != nil {
		t.Fatalf("exporting kube manifest: %v", err)
	}
}
