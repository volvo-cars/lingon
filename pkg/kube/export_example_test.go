// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube_test

import (
	"bytes"
	"fmt"

	"github.com/rogpeppe/go-internal/txtar"
	"github.com/volvo-cars/lingon/pkg/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// validate the struct implements the interface
var _ kube.Exporter = (*MyK8sApp)(nil)

// MyK8sApp contains kubernetes manifests
type MyK8sApp struct {
	kube.App
	// Namespace is the namespace for the tekton-pipelines
	PipelinesNS *corev1.Namespace
}

// New returns a new MyK8sApp
func New() *MyK8sApp {
	return &MyK8sApp{
		PipelinesNS: &corev1.Namespace{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Namespace",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "tekton-pipelines",
				Labels: map[string]string{
					"app.kubernetes.io/name": "tekton-pipelines",
				},
			},
		},
	}
}

func ExampleExport() {
	tk := New()

	var buf bytes.Buffer
	_ = kube.Export(tk, kube.WithExportWriter(&buf))

	ar := txtar.Parse(buf.Bytes())

	if len(ar.Files) > 0 {
		ns := &corev1.Namespace{}
		_ = yaml.Unmarshal(ar.Files[0].Data, ns)
		// print line by line to avoid trailing whitespace
		fmt.Println("apiVersion:", ns.APIVersion)
		fmt.Println("kind:", ns.Kind)
		fmt.Println("metadata:")
		fmt.Println("  labels:")
		fmt.Println(
			"    app.kubernetes.io/name:",
			ns.Labels["app.kubernetes.io/name"],
		)
		fmt.Println("name:", ns.Name)
	}

	// Output:
	//
	// apiVersion: v1
	// kind: Namespace
	// metadata:
	//   labels:
	//     app.kubernetes.io/name: tekton-pipelines
	// name: tekton-pipelines
	//
}
