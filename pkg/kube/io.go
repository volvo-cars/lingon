package kube

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ListGoFiles(root string) ([]string, error) {
	return listFiles(root, []string{".go"})
}

func ListYAMLFiles(root string) ([]string, error) {
	return listFiles(root, []string{".yaml", ".yml"})
}

func listFiles(root string, extensions []string) (
	[]string,
	error,
) {
	var files []string

	fi, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, errors.New("root is not a directory")
	}
	err = filepath.Walk(
		root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("walk %q %q, %w", path, info.Name(), err)
			}

			if !info.IsDir() && contains(
				filepath.Ext(filepath.Base(path)),
				extensions,
			) {
				files = append(files, path)
			}

			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("walk: %w", err)
	}
	return files, nil
}

func contains(e string, s []string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ReadManifest(filePath string) ([]string, error) {
	e := filepath.Ext(filePath)
	if e != ".yaml" && e != ".yml" {
		return nil, fmt.Errorf("not yaml file: %s", filePath)
	}
	yf, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read manifest %s: %w", filePath, err)
	}
	splitYaml, err := SplitManifest(bytes.NewReader(yf))
	if err != nil {
		return nil, fmt.Errorf("splitting manifest: %s: %w", filePath, err)
	}
	return splitYaml, nil
}

func SplitManifest(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var content []string
	var buf bytes.Buffer

	for scanner.Scan() {
		txt := scanner.Text()
		switch {
		// Skip comments
		case strings.HasPrefix(txt, "#"):
			continue
		// Split by '---'
		case strings.Contains(txt, "---"):
			if buf.Len() > 0 {
				content = append(content, buf.String())
				buf.Reset()
			}
		default:
			buf.WriteString(txt + "\n")
		}
	}

	s := buf.String()
	if len(s) > 0 { // if a manifest ends with '---'
		content = append(content, s)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("spliting manifests: %w", err)
	}
	return content, nil
}

func write(s, filename string) error {
	fp, err := os.Create(filename)
	if err != nil {
		var pe *os.PathError
		if errors.As(err, &pe) {
			return fmt.Errorf("path %q: %w", pe.Path, pe)
		}
		return err
	}
	defer fp.Close() //nolint:errcheck

	_, err = fp.WriteString(s)
	if err != nil {
		return err
	}
	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
