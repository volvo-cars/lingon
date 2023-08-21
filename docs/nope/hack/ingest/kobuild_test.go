package ingest_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/volvo-cars/nope/hack/ingest"
	"github.com/volvo-cars/nope/hack/ingest/templates"
	"golang.org/x/tools/txtar"
)

func TestKoBuild(t *testing.T) {
	ar, err := packTxtar("testdata/event", []string{"event.proto"})
	if err != nil {
		t.Fatalf("creating txtar: %v", err)
	}
	err = ingest.KoBuild(ingest.KoBuildConfig{
		Workdir:        filepath.Join("out", t.Name()),
		GoServiceTxtar: templates.IngestionTxtar,
		// Workdir: t.TempDir(),
		Schema: ingest.Schema{
			Name:  "ingestion",
			Txtar: txtar.Format(ar),
		},
	})
	if err != nil {
		t.Fatalf("running KoBuild: %v", err)
	}
}

// packTxtar takes a list of files and creates a [txtar.Archive] for them.
func packTxtar(dir string, files []string) (*txtar.Archive, error) {
	var archive txtar.Archive
	for _, file := range files {
		path := filepath.Join(dir, file)
		contents, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading proto file %s: %w", path, err)
		}
		ttFile := txtar.File{
			Name: file,
			Data: contents,
		}
		archive.Files = append(archive.Files, ttFile)
	}
	return &archive, nil
}
