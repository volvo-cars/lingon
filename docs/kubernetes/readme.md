# Kubernetes manifests in Go

## What is this?

This is a collection of libraries and helpers functions that helps to work with kubernetes manifests.

With this library you can:

- `Import` kubernetes manifests to Go structs ( YAML to Go )
- `Export` Go structs to kubernetes manifests ( Go to YAML )
- `Explode` kubernetes manifests to a single resource per file organized by namespaces (YAML to YAML)

## But why?

[Rationale.md](../rationale.md)

### Getting started

1. Get a kubernetes manifest
   - example:
```
wget https://storage.googleapis.com/tekton-releases/pipeline/previous/v0.45.0/release.yaml
```
  
2. Convert it to Go structs
   - Most manifests don't include extensions CRD object, use the `kygo` CLI:
      ```sh
      go run cmd/kygo/ -in=./install.yaml -out=output -app=myapp -group`
      ```
   - for specific CRDs, create a `main.go` file with the following content:

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
```

3. Modify the structs to your liking
   - see [best practices](docs/best-practices.md) for more information
   - example: [tekton](../platypus/pkg/platform/tekton/app.go)

4. Export:


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
	if err := kube.Export(tm, kube.WithExportOutputDirectory(out)); err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println("successfully exported CRDs to manifests")

	// Output:
	// successfully exported CRDs to manifests
}
```

5. Apply:

```sh
kubectl apply -f out/manifests/
```

done ✅.

Have a look at the [tests](../../pkg/kube/) and the [example](../kube/) for a full example.

What does the Go code looks like, see [tekton example](../platypus/pkg/platform/tekton/app.go)

### Best practices

> PLEASE READ [best practices](./best-practices.md) before using this library.

This project has been heavily inspired by :

- [Mimic](https://github.com/bwplotka/mimic) (you should definitely check it out, we copied the best practices from there)
- [NAML](https://github.com/krisnova/naml) (we found out about [valast](https://github.com/hexops/valast) from there)

Honorable mentions:

- [valast](https://github.com/hexops/valast) convert Go structs to its Go code representation
- [jennifer](https://github.com/dave/jennifer) generate Go code

## [CLI Utilities](../../cmd/)

### [Explode](../../cmd/explode/)

Converts multi-kubernetes-object manifest into multiple files, organized by namespace.
A CLI was written to make it easier to use in the terminal instead of just a library.

### [Kygo](../../cmd/kygo/)

Converts kubernetes manifests to Go structs.

A CLI was written to make it easier to use in the terminal instead of just a library.
It does support CustomResourceDefinitions but not the custom resources themselves, although it is easy to add them manually.
An example of how to do it can be found in the [example](../example/kube/).

## [Packages](../../pkg/)

Have a look at the [godoc](https://pkg.go.dev/github.com/volvo-cars/lingon) for more information.

### [Kube](../../pkg/kube/)

- `App` struct that is embedded to mark kubernetes applications
- `Export` kubernetes objects defined as Go struct to kubernetes manifests in YAML.
- `Explode` kubernetes manifests in YAML to multiple files, organized by namespace.
- `Import` kubernetes manifests in YAML to Go structs.

### [Kubeconfig](../../pkg/kubeconfig/)

Manipulate kubeconfig files **without** any dependencies on `go-client`.

### [KubeUtil](../../pkg/kubeutil/)

Reusable functions used to create kubernetes objects in Go.

### [Testutils](../../pkg/testutils/)

Reusable test functions.
