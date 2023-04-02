// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package crd_test

import (
	"fmt"
	"os"

	"github.com/volvo-cars/lingon/pkg/kube"
	istionetworkingv1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	secretsstorev1 "sigs.k8s.io/secrets-store-csi-driver/apis/v1"
)

func defaultSerializer() runtime.Decoder {
	// Add CRDs to scheme
	// This is needed to be able to import CRDs from kubernetes manifests files.
	_ = apiextensions.AddToScheme(scheme.Scheme)
	_ = apiextensionsv1.AddToScheme(scheme.Scheme)
	_ = secretsstorev1.AddToScheme(scheme.Scheme)
	_ = istionetworkingv1beta1.AddToScheme(scheme.Scheme)
	return scheme.Codecs.UniversalDeserializer()
}

// Example_import to shows how to import CRDs from kubernetes manifests files.
func Example_import() {
	// Remove previously generated output directory
	_ = os.RemoveAll("./out")

	if err := kube.Import(
		kube.WithImportAppName("team"),
		kube.WithImportManifestFiles([]string{"./manifest.yaml"}),
		kube.WithImportOutputDirectory("./out"),
		kube.WithImportSerializer(defaultSerializer()),
	); err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println("successfully imported CRDs from manifest.yaml")

	// Output:
	// successfully imported CRDs from manifest.yaml
}
