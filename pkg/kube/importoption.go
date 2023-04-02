// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"io"
	"os"

	"github.com/volvo-cars/lingon/pkg/kubeutil"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

// ImportOption is used to configure conversion from kubernetes objects in YAML to Go code
// Helpers function are provided to those field, see WithExportXXX functions
type ImportOption func(*jamel)

// importOption is used to configure the jamel, all fields have sane defaults
// Helpers function are provided to those field, see WithImportXXX functions
type importOption struct {
	// AppName is the name of the application, used to name the generated struct
	// ex: "tekton"
	AppName string

	// OutputPkgName is the name of the package where the generated code will be written (default: same as AppName)
	// ex: "tekton" but not "github.com/xxx/tekton"
	OutputPkgName string

	// OutputDir is the directory where the generated code will be written (default: out)
	// ex: "./tekton"
	OutputDir string

	// ManifestFiles is used to read the kubernetes objects from files, exclusive of ManifestReader
	// ex: []string{"./manifests/webapp1.yaml", "./manifests/webapp2.yaml"}
	ManifestFiles []string

	// ManifestReader is used to read the kubernetes objects from, exclusive of ManifestFiles
	// ex: os.Stdout, bytes.Buffer
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
	NameFieldFunc func(object kubeutil.Metadata) string

	// NameVarFunc formats the name of the variable containing the kubernetes object
	NameVarFunc func(object kubeutil.Metadata) string

	// NameFileFunc formats the name of the file containing the kubernetes object
	NameFileFunc func(object kubeutil.Metadata) string

	// RemoveAppName flag removes the app name from the object name
	RemoveAppName bool

	// GroupByKind flag groups the objects by kind
	GroupByKind bool

	// AddMethods flag adds convenience methods to the generated code
	AddMethods bool

	// RedactSecrets flag removes the value, but not the keys, of kubernetes secrets
	RedactSecrets bool
}

var importDefaultOpts = importOption{
	AppName:        "app",
	OutputPkgName:  "",
	ManifestFiles:  make([]string, 0),
	ManifestReader: os.Stdin,
	GoCodeWriter:   os.Stdout,
	OutputDir:      "out",
	Serializer:     scheme.Codecs.UniversalDeserializer(), // no CRDs by default
	NameFieldFunc:  NameFieldFunc,
	NameVarFunc:    NameVarFunc,
	NameFileFunc:   NameFileFunc,
	RemoveAppName:  false,
	GroupByKind:    false, // FIXME: should default to true ?
	AddMethods:     true,
	RedactSecrets:  false,
}

// WithImportSerializer sets the serializer [runtime.Decoder] to decode the kubernetes objects
//
// Usage:
//
//	func defaultSerializer() runtime.Decoder {
//		// add the scheme of the kubernetes objects you want to import
//		// this is useful for CRDs to be converted in Go as well
//		_ = otherpackage.AddToScheme(scheme.Scheme)
//		// needed for `CustomResourceDefinition` objects
//		_ = apiextensions.AddToScheme(scheme.Scheme)
//		return scheme.Codecs.UniversalDeserializer()
//	}
//	_, _ := kube.Import(WithImportSerializer(defaultSerializer()))
func WithImportSerializer(s runtime.Decoder) ImportOption {
	return func(j *jamel) {
		j.o.Serializer = s
	}
}

//
// NAMES (field, var, file)
//

// WithImportAppName sets the application name for the generated code.
// This is used to name the generated struct.
// ex: "tekton"
//
// Default: "app"
//
// Note: the name can be used to name the package if none is defined,
// see WithImportPackageName
func WithImportAppName(name string) ImportOption {
	return func(j *jamel) {
		j.o.AppName = name
	}
}

// WithImportPackageName sets the package name for the generated code
// Note that the package name cannot contain a dash, it will panic otherwise.
//
// ex: "tekton" but not "github.com/xxx/tekton"
//
//	package tekton
//	...
func WithImportPackageName(name string) ImportOption {
	return func(j *jamel) {
		j.o.OutputPkgName = name
	}
}

// WithImportRemoveAppName tries to remove the name of the application from the object name.
// Default: false
func WithImportRemoveAppName(b bool) ImportOption {
	return func(j *jamel) {
		j.o.RemoveAppName = b
	}
}

// WithImportNameFieldFunc sets the function to format the name of the field
// in the application struct (containing [kube.App]).
//
// default: [NameFieldFunc]
//
// TIP: ALWAYS put the kind somewhere in the name to avoid collisions
//
//	type Tekton struct {
//		kube.App
//		// ...
//		ThisIsTheNameFieldCM  *corev1.ConfigMap
//	}
func WithImportNameFieldFunc(f func(object kubeutil.Metadata) string) ImportOption {
	return func(j *jamel) {
		j.o.NameFieldFunc = f
	}
}

// WithImportNameVarFunc sets the function to format the name of the variable
// containing the kubernetes object.
//
// default: [NameVarFunc]
//
//	var ThisIsTheNameOfTheVar = &appsv1.Deployment{...}
//
//	 // ...
//
//	func New() *Tekton {
//		return &Tekton{
//			NameField:          ThisIsTheNameOfTheVar,
//			...
//		}
//	}
func WithImportNameVarFunc(f func(object kubeutil.Metadata) string) ImportOption {
	return func(j *jamel) {
		j.o.NameVarFunc = f
	}
}

// WithImportNameFileFunc sets the function to format the name of the file
// containing the kubernetes object.
//
// default: [NameFileFunc]
//
// Usage:
//
//	WithImportNameFileFunc(func(m kubeutil.Metadata) string {
//		return fmt.Sprintf("%s-%s.go", strings.ToLower(m.Kind),	m.Meta.Name)
//	})
func WithImportNameFileFunc(f func(object kubeutil.Metadata) string) ImportOption {
	return func(j *jamel) {
		j.o.NameFileFunc = f
	}
}

//
//  INPUT (files, reader)
//

// WithImportManifestFiles sets the manifest files to read the kubernetes objects from.
func WithImportManifestFiles(files []string) ImportOption {
	return func(j *jamel) {
		j.useReader = false
		j.o.ManifestFiles = files
	}
}

// WithImportReadStdIn reads the kubernetes objects from [os.Stdin].
func WithImportReadStdIn() ImportOption {
	return func(j *jamel) {
		j.useReader = true
		j.o.ManifestReader = os.Stdin
	}
}

// WithImportReader reads the kubernetes manifest (YAML) from a [io.Reader]
// Note that this is exclusive with [WithImportManifestFiles]
//
// If you want to read from [os.Stdin] use [WithImportReadStdIn].
func WithImportReader(r io.Reader) ImportOption {
	return func(j *jamel) {
		j.useReader = true
		j.o.ManifestReader = r
	}
}

//
//   OUTPUT
//   - writer
//   - directory
//   - group files by kind
//   - redact secrets
//   - methods
//

// WithImportGroupByKind groups the kubernetes objects by kind in the same file
//
// if there are 10 ConfigMaps and 5 Secrets, it will generate 2 files:
//   - configmaps.go
//   - secrets.go
//
// as opposed to 15 files.
//
// Default: false
func WithImportGroupByKind(b bool) ImportOption {
	return func(j *jamel) {
		j.o.GroupByKind = b
	}
}

// WithImportWriter writes the generated Go code to [io.Writer].
// Note that the format is txtar, for more info on [golang.org/x/tools/txtar.Archive] format
// see: https://pkg.go.dev/golang.org/x/tools/txtar
//
// A txtar archive is zero or more comment lines and then a sequence of file entries.
// Each file entry begins with a file marker line of the form "-- FILENAME --" and
// is followed by zero or more file content lines making up the file data.
// The comment or file content ends at the next file marker line.
// The file marker line must begin with the three-byte sequence "-- " and
// end with the three-byte sequence " --", but the enclosed file name can be
// surrounding by additional white space, all of which is stripped.
//
// If the txtar file is missing a trailing newline on the final line,
// parsers should consider a final newline to be present anyway.
//
// There are no possible syntax errors in a txtar archive.
func WithImportWriter(w io.Writer) ImportOption {
	return func(j *jamel) {
		j.useWriter = true
		j.o.GoCodeWriter = w
	}
}

// WithImportOutputDirectory sets the output directory for the generated code.
// Default: "./out"
func WithImportOutputDirectory(name string) ImportOption {
	return func(j *jamel) {
		j.o.OutputDir = name
	}
}

// WithImportRedactSecrets removes the value, but not the keys, of kubernetes secrets.
// Default: true
func WithImportRedactSecrets(b bool) ImportOption {
	return func(j *jamel) {
		j.o.RedactSecrets = b
	}
}

// WithImportAddMethods adds convenience methods to the generated code.
//
// Default: true
//
//	// Apply applies the kubernetes objects to the cluster
//	func (a *Tekton) Apply(ctx context.Context) error
//
//	// Export exports the kubernetes objects to YAML files in the given directory
//	func (a *Tekton) Export(dir string) error
//
//	// Apply applies the kubernetes objects contained in [Exporter] to the cluster
//	func Apply(ctx context.Context, km kube.Exporter) error
func WithImportAddMethods(b bool) ImportOption {
	return func(j *jamel) {
		j.o.AddMethods = b
	}
}
