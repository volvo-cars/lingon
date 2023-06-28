package ingest_test

import (
	"path/filepath"
	"testing"

	tu "github.com/volvo-cars/lingon/pkg/testutil"
	"github.com/volvo-cars/nope/hack/ingest"
	"github.com/volvo-cars/nope/hack/ingest/templates"
	"golang.org/x/tools/txtar"
)

func TestKoBuild(t *testing.T) {
	ar, err := ingest.PackTxtar("testdata/event", []string{"event.proto"})
	tu.AssertNoError(t, err, "creating txtar")
	err = ingest.KoBuild(ingest.KoBuildConfig{
		Workdir:        filepath.Join("out", t.Name()),
		GoServiceTxtar: templates.IngestionTxtar,
		// Workdir: t.TempDir(),
		Schema: ingest.Schema{
			Name:  "ingestion",
			Txtar: txtar.Format(ar),
		},
	})
	tu.AssertNoError(t, err, "building ingestion service")
}
