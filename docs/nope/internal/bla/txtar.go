package bla

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/tools/txtar"
)

// PackTxtar takes a list of files and creates a [txtar.Archive] for them.
func PackTxtar(dir string, files []string) (*txtar.Archive, error) {
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
