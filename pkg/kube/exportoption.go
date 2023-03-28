// Copyright (c) Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"io"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExportOption is used to configure conversion from Go code to kubernetes objects in YAML
type ExportOption func(*goky)

// option is used to configure the jamel, all fields have sane defaults
// Helpers function are provided to those field, see WithXXX functions
type exportOption struct {
	// AppName is the name of the application, used to name the generated struct
	// ex: "tekton"
	AppName string
	// OutputDir is the directory where the generated code will be written (default: out)
	// ex: "./tekton"
	OutputDir string
	// Writer is used to write the exported manifests in txtar format
	// for more info on txtar format see: https://pkg.go.dev/golang.org/x/tools/txtar
	// Note that we are using https://github.com/rogpeppe/go-internal/blob/master/txtar/ instead
	// ex: os.Stdout, bytes.Buffer
	ManifestWriter io.Writer
	// NameFileFunc formats the name of the file containing the kubernetes object
	NameFileFunc func(object metav1.ObjectMeta) string

	// SecretHook is used to process the secrets before they are exported.
	// The hook is called for each secret.
	// This is useful to redact the secrets in order not to save them in plain text.
	//
	SecretHook func(secret *corev1.Secret) error
	// GroupByKind flag groups the objects by kind
	GroupByKind bool
	// Kustomize flag adds a kustomization.yaml file to the output
	Kustomize bool
}

var exportDefaultOpts = exportOption{
	AppName:        "",
	ManifestWriter: os.Stdout,
	OutputDir:      "out",
	// NameFileFunc:   ExportNameFileFunc,
	GroupByKind: false,
}

// WithExportNameFileFunc sets the function to format the name of the file
// containing the kubernetes object
// default: ExportNameFileFunc
//
// Usage:
//
//	WithExportNameFileFunc(func(m metav1.Metadata) string {
//		return fmt.Sprintf("%s-%s.go", strings.ToLower(m.Kind),	m.Meta.Name)
//	})
func WithExportNameFileFunc(f func(object metav1.ObjectMeta) string) ExportOption {
	return func(g *goky) {
		g.o.NameFileFunc = f
	}
}

// WithExportGroupByKind groups the kubernetes objects by kind in the same file
//
// if there are 10 ConfigMaps and 5 Deployments, it will generate 2 files:
//   - configmaps.go
//   - deployments.go
//
// as opposed to one or 15 files.
func WithExportGroupByKind(b bool) ExportOption {
	return func(g *goky) {
		g.o.GroupByKind = b
	}
}

// WithExportExplodeManifests explodes the manifests into separate files
// organized by namespace to match closely the structure of the kubernetes cluster
func WithExportExplodeManifests(b bool) ExportOption {
	return func(g *goky) {
		g.explode = b
	}
}

// WithExportWriter writes the generated manifests to io.Writer
func WithExportWriter(w io.Writer) ExportOption {
	return func(g *goky) {
		g.useWriter = true
		g.o.ManifestWriter = w
	}
}

// WithExportKustomize adds a kustomization.yaml file to the output
func WithExportKustomize(b bool) ExportOption {
	return func(g *goky) {
		g.o.Kustomize = b
	}
}

// WithExportOutputDirectory sets the output directory for the generated manifests
func WithExportOutputDirectory(dir string) ExportOption {
	return func(g *goky) {
		g.o.OutputDir = dir
	}
}

// WithExportSecretHook is used to process the secrets before they are exported.
// The hook is called for each secret.
// This is useful to redact the secrets in order not to save them in plain text.
// Note the secrets will not be written to the output directory if this option is used.
func WithExportSecretHook(f func(s *corev1.Secret) error) ExportOption {
	return func(g *goky) {
		g.removeSecrets = true
		g.o.SecretHook = f
	}
}
