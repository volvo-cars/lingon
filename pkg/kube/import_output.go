// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/veggiemonk/strcase"
	"golang.org/x/exp/slog"
	"golang.org/x/tools/txtar"
	"mvdan.cc/gofumpt/format"
)

const (
	embeddedStructName = "App"
	kubeAppPkgPath     = "github.com/volvo-cars/lingon/pkg/kube"
)

func (j *jamel) render() error {
	if err := j.generateGo(); err != nil {
		if j.o.IgnoreErrors {
			j.l.Error("generate go", logErrIgnored, "err", err)
			return nil
		}
		return fmt.Errorf("generate go: %w", err)
	}

	// render all the kubernetes objects
	if j.o.GroupByKind {
		if err := j.renderFileByKind(); err != nil {
			return fmt.Errorf("by kind: %w", err)
		}
	} else {
		if err := j.renderFileByName(); err != nil {
			return fmt.Errorf("by name: %w", err)
		}
	}

	// app.go with kubeapp struct
	appFile := j.appFile()
	filename := filepath.Join(j.o.OutputDir, "app.go")
	if _, err := j.buf.Write([]byte("-- " + filename + " --\n")); err != nil {
		return err
	}
	if err := appFile.Render(&j.buf); err != nil {
		return fmt.Errorf("app.go: %w", err)
	}
	if j.useWriter {
		size := j.buf.Len()
		written, err := io.Copy(j.o.GoCodeWriter, &j.buf)
		if err != nil {
			return fmt.Errorf("writing output: %w", err)
		}
		if written != int64(size) {
			return fmt.Errorf(
				"not all bytes written: %d >< %d",
				written,
				j.buf.Len(),
			)
		}
		if j.o.Verbose {
			j.l.Info("output", slog.Int64("bytes written", written))
		}
	}

	return nil
}

// renderFileByKind renders all the kubernetes objects
// to each file containing all the objects of the same kind
func (j *jamel) renderFileByKind() error {
	kindFileMap, err := j.fileMap()
	if err != nil {
		return fmt.Errorf("filemap: %w", err)
	}

	for _, kind := range orderedKeys(kindFileMap) {
		file := kindFileMap[kind]

		// no choice in the filename
		filename := strcase.Kebab(kind) + ".go"
		if j.o.OutputDir != "" {
			filename = filepath.Join(j.o.OutputDir, filename)
		}
		if j.o.Verbose {
			j.l.Info("render", "filename", filename)
		}
		if _, err = j.buf.Write([]byte("-- " + filename + " --\n")); err != nil {
			if j.o.IgnoreErrors {
				j.l.Error("render", logErrIgnored, "err", err)
				continue
			}
			return err
		}

		if err = file.Render(&j.buf); err != nil {
			return fmt.Errorf("render: %w", err)
		}
	}
	return nil
}

func (j *jamel) renderFileByName() error {
	for _, nameVar := range orderedKeys(j.objectsCode) {
		objMeta, ok := j.objectsMeta[nameVar]
		if !ok {
			return fmt.Errorf("no object meta for %s", nameVar)
		}

		// rename the variable holding the kubernetes object
		nameVarObj := j.o.NameVarFunc(objMeta)
		if j.o.RemoveAppName {
			nameVarObj = RemoveAppName(nameVarObj, j.o.AppName)
		}

		stmt := j.objectsCode[nameVar]
		file := stmtKubeObjectFile(j.o.OutputPkgName, nameVarObj, stmt)

		// rename the file
		filename := j.o.NameFileFunc(objMeta)
		if j.o.RemoveAppName {
			filename = RemoveAppName(filename, j.o.AppName)
		}

		// set the correct path to the file
		if j.o.OutputDir != "" {
			filename = filepath.Join(j.o.OutputDir, filename)
		}
		if j.o.Verbose {
			j.l.Info("render", "filename", filename)
		}
		_, err := j.buf.Write([]byte("-- " + filename + " --\n"))
		if err != nil {
			return fmt.Errorf("write: %w", err)
		}

		err = file.Render(&j.buf)
		if err != nil {
			return fmt.Errorf("render: %w", err)
		}
	}
	return nil
}

// fileMap returns a map of all the kubernetes objects
// the key is the file name and the value is the jen.File meant to be rendered.
func (j *jamel) fileMap() (map[string]*jen.File, error) {
	filesCreated := make(map[string]struct{}, 0)
	kindFile := make(map[string]*jen.File, 0)

	// create a file for each kind
	for _, nameVar := range orderedKeys(j.objectsCode) {
		stmt := j.objectsCode[nameVar]
		objMeta, ok := j.objectsMeta[nameVar]
		if !ok {
			if j.o.IgnoreErrors {
				j.l.Error(
					"no object meta",
					logErrIgnored,
					"variable name",
					nameVar,
				)
				continue
			}
			return nil, fmt.Errorf("no object meta for %s", nameVar)
		}
		nameVarObj := j.o.NameVarFunc(objMeta)
		if j.o.RemoveAppName {
			nameVarObj = RemoveAppName(nameVarObj, j.o.AppName)
		}

		// if last letter of nameVar is a number, it is a duplicate
		// we add that number to the nameVarObj
		if lastChar := nameVar[len(nameVar)-1]; lastChar >= '0' && lastChar <= '9' {
			nameVarObj += string(lastChar)
		}

		// check if file exists for this kind
		if _, ok := filesCreated[objMeta.Kind]; ok {
			kindFile[objMeta.Kind].Line().
				Var().Id(nameVarObj).Op("=").Add(stmt)
			continue
		}
		// no file exists for this kind, create one
		filesCreated[objMeta.Kind] = struct{}{}
		kindFile[objMeta.Kind] = stmtKubeObjectFile(
			j.o.OutputPkgName,
			nameVarObj,
			stmt,
		)
	}

	return kindFile, nil
}

func (j *jamel) save() error {
	if err := j.render(); err != nil {
		return fmt.Errorf("render: %w", err)
	}

	if err := os.MkdirAll(j.o.OutputDir, 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	ar := txtar.Parse(j.buf.Bytes())
	if len(ar.Files) == 0 {
		return fmt.Errorf("no file to write")
	}
	if err := writeTxtar(ar); err != nil {
		return fmt.Errorf("write txtar: %w", err)
	}

	return nil
}

func writeTxtar(ar *txtar.Archive) error {
	var err error
	for _, f := range ar.Files {
		// format code with gofumpt extra rules as it is stricter
		// and will produce more predictable output
		f.Data, err = format.Source(
			f.Data, format.Options{LangVersion: "1.20", ExtraRules: true},
		)
		if err != nil {
			return fmt.Errorf("formating generated code: %w", err)
		}
	}

	// predictable output
	sort.SliceStable(
		ar.Files, func(i, j int) bool {
			return ar.Files[i].Name < ar.Files[j].Name
		},
	)

	// each files path is already prefixed with the output dir, using "." as
	// all the files will be written in relation to it.
	if err = write(ar, "."); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func orderedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
