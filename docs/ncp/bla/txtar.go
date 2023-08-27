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

func UnpackTxtar(dir string, contents []byte) error {
	ar := txtar.Parse(contents)
	if len(ar.Files) == 0 {
		return fmt.Errorf("txtar archive contains no files")
	}

	// Write the files to disk
	for _, tFile := range ar.Files {
		path := filepath.Join(dir, tFile.Name)
		fileDir := filepath.Dir(path)
		if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
			return fmt.Errorf("creating directory for file %s: %w", path, err)
		}
		if err := os.WriteFile(path, tFile.Data, os.ModePerm); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
	}
	return nil
}
