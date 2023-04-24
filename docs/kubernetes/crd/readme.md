# Example CRD

- [Import the manifest in Go](#import-the-manifest-in-go)
- [Export the Go struct to YAML](#export-the-go-struct-to-yaml)
- [How to find CRD types](#how-to-find-crd-types)
- [List of well-known CRDs types](#list-of-well-known-crds-types)

All can be done in one step:

```bash
go generate -v ./...
```

## Import the manifest in Go

```go
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
		kube.WithImportPackageName("team"),
		kube.WithImportManifestFiles([]string{"./manifest.yaml"}),
		kube.WithImportOutputDirectory("./out"),
		kube.WithImportSerializer(defaultSerializer()),
		// do not print verbose information
		kube.WithImportVerbose(false),
		// do not ignore errors
		kube.WithImportIgnoreErrors(false),
		// use the default logger even if unused with WithImportVerbose(false)
		kube.WithImportLogger(kube.Logger(os.Stderr)),
	); err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println("successfully imported CRDs from manifest.yaml")

	// Output:
	// successfully imported CRDs from manifest.yaml
}
```


## Export the Go struct to YAML


```go
package crd_test

import (
	"fmt"
	"os"
	"path/filepath"

	team "github.com/volvo-cars/lingon/docs/kubernetes/crd/out"
	"github.com/volvo-cars/lingon/pkg/kube"
)

var defaultOut = "out"

// Example_export to shows how to export to kubernetes manifests files in YAML.
func Example_export() {
	out := filepath.Join(defaultOut, "manifests")

	// Remove previously generated output directory
	_ = os.RemoveAll(out)

	tm := team.New()
	if err := kube.Export(
		tm,
		// directory where to export the manifests
		kube.WithExportOutputDirectory(out),

		// // Write the manifests to a bytes.Buffer.
		// // Note that it will be written as txtar format.
		// kube.WithExportWriter(buf),
		// // Add a kustomization.yaml file next to the manifests.
		// kube.WithExportKustomize(true),

		// // Write the manifest as a single file.
		// kube.WithExportAsSingleFile("manifest.yaml"),

		// // Write the manifests as JSON instead of YAML.
		// // Written as an JSON array if written as a single file.
		// kube.WithExportOutputJSON(true),

		// // Write to standard output instead of a file.
		// kube.WithExportStdOut(),

		// // Write the manifests as it would appear on the cluster.
		// // 1 resource per file and one folder per namespace.
		// kube.WithExportExplodeManifests(true),

		// // Define a hook to be called before exporting a secret.
		// // Note that this will remove the secret from the output.
		// kube.WithExportSecretHook(
		// 	func(s *corev1.Secret) error {
		// 		// do something with the secret
		// 		return nil
		// 	}),
	); err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println("successfully exported CRDs to manifests")

	// Output:
	// successfully exported CRDs to manifests
}
```


## How to find CRD types

This is more of a heuristic than a rule.

- go to the repo of the project
- search for `AddToScheme` function

Often it is in the `pkg/apis` or `apis` folder of the project.

## List of well-known CRDs types

> This is best-effort

```go
import(
    "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/runtime/serializer"
    apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
    capiahelm "sigs.k8s.io/cluster-api-addon-provider-helm/api/v1alpha1"
    certmanager "github.com/cert-manager/cert-manager/pkg/api"
    clientgoscheme "k8s.io/client-go/kubernetes/scheme"
    externalsecretsv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
    gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
    gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
    helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
    imageautov1 "github.com/fluxcd/image-automation-controller/api/v1beta1"
    imagereflectv1 "github.com/fluxcd/image-reflector-controller/api/v1beta2"
    istionetworkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
    istionetworkingv1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
    istiosecurityv1beta1 "istio.io/client-go/pkg/apis/security/v1beta1"
    karpenterapi "github.com/aws/karpenter/pkg/apis"
    karpenterv1alpha5 "github.com/aws/karpenter-core/pkg/apis"
    knativecachingalpha1 "knative.dev/caching/pkg/apis/caching/v1alpha1"
    knativeservingv1 "knative.dev/serving/pkg/apis/serving/v1"
    kservev1alpha1 "github.com/kserve/kserve/pkg/apis/serving/v1alpha1"
    kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta2"
    metacontrolleralpha "github.com/metacontroller/metacontroller/pkg/apis/metacontroller/v1alpha1"
    notificationv1b2 "github.com/fluxcd/notification-controller/api/v1beta2"
    otelv1alpha1 "github.com/open-telemetry/opentelemetry-operator/apis/v1alpha1"
    profilev1 "github.com/kubeflow/kubeflow/components/profile-controller/api/v1"
    profilev1beta1 "github.com/kubeflow/kubeflow/components/profile-controller/api/v1beta1"
    secretsstorev1 "sigs.k8s.io/secrets-store-csi-driver/apis/v1"
    servingv1alpha1 "github.com/kserve/modelmesh-serving/apis/serving/v1alpha1"
    slothv1alpha1 "github.com/slok/sloth/pkg/kubernetes/api/sloth/v1"
    sourcev1 "github.com/fluxcd/source-controller/api/v1"
    tektonpipelinesv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
    tektontriggersv1alpha1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
    utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)
```
