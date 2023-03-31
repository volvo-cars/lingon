package main

import (
	ricobergerdev1alpha1 "github.com/ricoberger/vault-secrets-operator/api/v1alpha1"
	"github.com/volvo-cars/lingon/pkg/kube"
	istionetworkingv1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubectl/pkg/scheme"
	secretsstorev1 "sigs.k8s.io/secrets-store-csi-driver/apis/v1"
)

func main() {
	if err := kube.Import(
		kube.WithImportAppName("team"),
		kube.WithImportManifestFiles([]string{"./manifest.yaml"}),
		kube.WithImportOutputDirectory("./out"),
		kube.WithImportSerializer(defaultSerializer()),
	); err != nil {
		panic(err)
	}
}

func defaultSerializer() runtime.Decoder {
	// NEEDED FOR CRDS
	//
	if err := apiextensions.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}
	if err := apiextensionsv1.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}
	if err := secretsstorev1.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}
	if err := ricobergerdev1alpha1.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}
	if err := istionetworkingv1beta1.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}
	return scheme.Codecs.UniversalDeserializer()
}
