// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube_test

import (
	"bytes"
	"fmt"

	"github.com/volvo-cars/lingon/pkg/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	fmt.Printf("%s\n", buf.String())

	// Output:
	//
	// -- out/0_pipelines_ns.yaml --
	// apiVersion: v1
	// kind: Namespace
	// metadata:
	//   labels:
	//     app.kubernetes.io/name: tekton-pipelines
	//   name: tekton-pipelines
	// spec: {}
}
