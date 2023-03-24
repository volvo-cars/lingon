// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"io"
	"os"
	"strings"

	"github.com/volvo-cars/lingon/pkg/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

type ImportOption func(*jamel)

// Option is used to configure the jamel, all fields have sane defaults
// Helpers function are provided to those field, see WithXXX functions
type Option struct {
	// AppName is the name of the application, used to name the generated struct
	// ex: "karpenter"
	AppName string
	// OutputPkgName is the name of the package where the generated code will be written (default: same as AppName)
	// ex: "karpenter" but not "github.com/xxx/karpenter"
	OutputPkgName string
	// OutputDir is the directory where the generated code will be written (default: out)
	// ex: "./karpenter"
	OutputDir string
	// ManifestFiles is used to read the kubernetes objects from files, exclusive of ManifestReader
	// ex: []string{"./manifests/webapp1.yaml", "./manifests/webapp2.yaml"}
	ManifestFiles []string
	// ManifestReader is used to read the kubernetes objects from, exclusive of ManifestFiles
	// ex: os.Stdin
	ManifestReader io.Reader
	// GoCodeWriter is used to write the generated Go code in txtar format
	// for more info on txtar format see: https://pkg.go.dev/golang.org/x/tools/txtar
	// Note that we are using https://github.com/rogpeppe/go-internal/blob/master/txtar/ instead
	// ex: os.Stdout
	GoCodeWriter io.Writer
	// Serializer is used to decode the kubernetes objects
	// ex: scheme.Codecs.UniversalDeserializer()
	Serializer runtime.Decoder
	// NameFieldFunc formats the name of the field in the application struct
	NameFieldFunc func(object meta.Metadata) string
	// NameVarFunc formats the name of the variable containing the kubernetes object
	NameVarFunc func(object meta.Metadata) string
	// NameFileObjFunc formats the name of the file containing the kubernetes object
	NameFileObjFunc func(object meta.Metadata) string
	// RemoveAppName removes the app name from the object name
	RemoveAppName bool
	// GroupByKind groups the objects by kind
	GroupByKind bool
	// Implement the KubeApp interface
	AddMethods bool
	// RedactSecrets flag removes the value, but not the keys, of kubernetes secrets
	RedactSecrets bool
}

var defaultOpts = Option{
	AppName:         "myapp",
	OutputPkgName:   "myapp",
	ManifestFiles:   make([]string, 0),
	ManifestReader:  os.Stdin,
	GoCodeWriter:    os.Stdout,
	OutputDir:       "out",
	Serializer:      scheme.Codecs.UniversalDeserializer(), // no CRDs by default
	NameFieldFunc:   NameFieldFunc,
	NameVarFunc:     NameVarFunc,
	NameFileObjFunc: NameFileObjFunc,
	RemoveAppName:   false,
	GroupByKind:     false, // TODO: should default to true ?
	AddMethods:      true,
	RedactSecrets:   false,
}

func WithAppName(name string) ImportOption {
	return func(j *jamel) {
		j.o.AppName = name
	}
}

func WithPackageName(name string) ImportOption {
	if strings.Contains(name, "-") {
		panic("package name cannot contain a dash")
	}
	return func(j *jamel) {
		j.o.OutputPkgName = name
	}
}

func WithOutputDirectory(name string) ImportOption {
	return func(j *jamel) {
		j.o.OutputDir = name
	}
}

func WithManifestFiles(files []string) ImportOption {
	return func(j *jamel) {
		j.useReader = false
		for _, f := range files {
			if err := j.addManifest(f); err != nil {
				panic(err)
			}
		}
		j.o.ManifestFiles = files
	}
}

func WithRemoveAppName(b bool) ImportOption {
	return func(j *jamel) {
		j.o.RemoveAppName = b
	}
}

func WithGroupByKind(b bool) ImportOption {
	return func(j *jamel) {
		j.o.GroupByKind = b
	}
}

func WithReadStdIn() ImportOption {
	return func(j *jamel) {
		j.useReader = true
		j.o.ManifestReader = os.Stdin
		j.o.ManifestFiles = make([]string, 0)
	}
}

func WithReader(r io.Reader) ImportOption {
	return func(j *jamel) {
		j.useReader = true
		j.o.ManifestReader = r
	}
}

func WithWriter(w io.Writer) ImportOption {
	return func(j *jamel) {
		j.useWriter = true
		j.o.GoCodeWriter = w
	}
}

func WithRedactSecrets(b bool) ImportOption {
	return func(j *jamel) {
		j.o.RedactSecrets = b
	}
}

func WithMethods(b bool) ImportOption {
	return func(j *jamel) {
		j.o.AddMethods = b
	}
}

func WithSerializer(s runtime.Decoder) ImportOption {
	return func(j *jamel) {
		j.o.Serializer = s
	}
}

func WithNameFieldFunc(f func(object meta.Metadata) string) ImportOption {
	return func(j *jamel) {
		j.o.NameFieldFunc = f
	}
}

func WithNameVarFunc(f func(object meta.Metadata) string) ImportOption {
	return func(j *jamel) {
		j.o.NameVarFunc = f
	}
}

func WithNameFileObjFunc(f func(object meta.Metadata) string) ImportOption {
	return func(j *jamel) {
		j.o.NameFileObjFunc = f
	}
}
