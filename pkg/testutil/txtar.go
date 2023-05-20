// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package testutil

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"

	"github.com/rogpeppe/go-internal/txtar"
)

func Folder2Txtar(folder string) (*txtar.Archive, error) {
	files, err := listFiles(folder)
	if err != nil {
		return nil, err
	}
	var ar txtar.Archive
	ar.Files = make([]txtar.File, 0, len(files))
	for _, f := range files {
		data, err := os.ReadFile(f)
		ar.Files = append(ar.Files, txtar.File{Name: f, Data: data})
		if err != nil {
			return nil, fmt.Errorf("read file: %w", err)
		}
	}
	return &ar, nil
}

func Filenames(ar *txtar.Archive) []string {
	filenames := []string{}
	for _, file := range ar.Files {
		filenames = append(filenames, file.Name)
	}
	return filenames
}

func DiffTxtarSort(got, want *txtar.Archive) string {
	sort.SliceStable(
		got.Files, func(i, j int) bool {
			return got.Files[i].Name < got.Files[j].Name
		},
	)
	sort.SliceStable(
		want.Files, func(i, j int) bool {
			return want.Files[i].Name < want.Files[j].Name
		},
	)
	return Diff(got, want)
}

func DiffTxtar(got, want *txtar.Archive) string {
	return Diff(string(txtar.Format(got)), string(txtar.Format(want)))
}

func listFiles(root string) ([]string, error) {
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

			if !info.IsDir() {
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

// VerifyGo parses a [txtar.Archive], if the files are not valid Go code
// and error is returned
func VerifyGo(ar *txtar.Archive) error {
	fset := token.NewFileSet()
	for _, file := range ar.Files {
		fset.AddFile(file.Name, fset.Base(), len(file.Data))
		fast, err := parser.ParseFile(
			fset,
			file.Name,
			file.Data,
			parser.AllErrors,
		)
		if err != nil {
			return err
		}
		// not really useful, but it's a start
		if hasBadNodes(fast) {
			return fmt.Errorf("invalid go file: %s", file.Name)
		}
	}

	return nil
}

func hasBadNodes(node ast.Node) bool {
	a := false
	ast.Inspect(
		node, func(n ast.Node) bool {
			if a {
				return false
			}
			switch n.(type) {
			case *ast.BadExpr, *ast.BadDecl, *ast.BadStmt:
				a = true
			}
			return true
		},
	)
	return a
}
