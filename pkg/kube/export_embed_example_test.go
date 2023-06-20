package kube_test

import (
	"bytes"
	"fmt"

	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SubApp struct {
	kube.App

	CM *corev1.ConfigMap
}

func NewSubApp() *SubApp {
	return &SubApp{
		CM: &corev1.ConfigMap{
			TypeMeta: kubeutil.TypeConfigMapV1,
			ObjectMeta: metav1.ObjectMeta{
				Name:      "config-map-name",
				Namespace: "my-ns",
			},
			Data: map[string]string{"key": "value"},
		},
	}
}

type MainApp struct {
	kube.App

	SubApp *SubApp
	NS     *corev1.Namespace
}

func NewApp() *MainApp {
	return &MainApp{
		SubApp: NewSubApp(),
		NS: &corev1.Namespace{
			TypeMeta:   kubeutil.TypeNamespaceV1,
			ObjectMeta: metav1.ObjectMeta{Name: "my-ns"},
		},
	}
}

func ExampleExport_embedded() {
	app := NewApp()

	var buf bytes.Buffer
	_ = kube.Export(app, kube.WithExportWriter(&buf))

	fmt.Printf("%s\n", buf.String())

	// Output:
	//
	// -- out/0_ns.yaml --
	// apiVersion: v1
	// kind: Namespace
	// metadata:
	//   name: my-ns
	// spec: {}
	// -- out/2_sub_appcm.yaml --
	// apiVersion: v1
	// data:
	//   key: value
	// kind: ConfigMap
	// metadata:
	//   name: config-map-name
	//   namespace: my-ns
}
