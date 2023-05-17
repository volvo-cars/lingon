# Kubernetes manifests in Go

- [What is this?](#what-is-this)
- [But why?](#but-why)
- [Getting started](#getting-started)
- [Best practices](#best-practices)
- [CRDs](#crds)
- [Testing manifests](#testing-manifests)
- [CLI Utilities](#cli-utilities)
	- [Explode](#explode)
	- [Kygo](#kygo)
- [Packages](#packages)
	- [Kube](#kube)
	- [Kubeconfig](#kubeconfig)
	- [KubeUtil](#kubeutil)
	- [Testutils](#testutils)

## What is this?

This is a collection of libraries and helpers functions that helps to work with kubernetes manifests.

With this library you can:

- `Import` kubernetes manifests to Go structs ( YAML to Go )
  - see [Import Options](./options-import.md) for more information
- `Export` Go structs to kubernetes manifests ( Go to YAML )
  - see [Export Options](./options-export.md) for more information
- `Explode` kubernetes manifests to a single resource per file organized by namespaces (YAML to YAML)

## But why?

[Rationale.md](../rationale.md)

## Getting started

1. Get a kubernetes manifest
   - example:

```sh
wget https://storage.googleapis.com/tekton-releases/pipeline/previous/v0.45.0/release.yaml
```
  
2. Convert it to Go structs. Choose either:
   A) Most manifests don't include extensions CRD object, use the `kygo` CLI:
      ```sh
      go run cmd/kygo/ -in=./install.yaml -out=output -app=myapp -group`
      ```
   B) for **CRDs**, look at the [CRD example](./crd)
   C) for more options, see [Import Options](./options-import.md)

3. Modify the structs to your liking
   - see [best practices](docs/best-practices.md) for more information
   - example: [tekton](../platypus/pkg/platform/tekton/app.go)

4. Export:

{{ "Example_export" | example }}

5. Apply:

```sh
kubectl apply -f out/manifests/
```

done ✅.

Have a look at the [tests](../../pkg/kube/) and the [example](../kube/) for a full example.

What does the Go code looks like, see [tekton example](../platypus/pkg/platform/tekton/app.go)

## Best practices

> PLEASE READ [best practices](./best-practices.md) before using this library.

This project has been heavily inspired by :

- [Mimic](https://github.com/bwplotka/mimic) (you should definitely check it out, we copied the best practices from there)
- [NAML](https://github.com/krisnova/naml) (we found out about [valast](https://github.com/hexops/valast) from there)

Honorable mentions:

- [valast](https://github.com/hexops/valast) convert Go structs to its Go code representation
- [jennifer](https://github.com/dave/jennifer) generate Go code

## CRDs

How to convert CRDs from YAML to Go. By setting a kubernetes runtime scheme:

```go
"k8s.io/apimachinery/pkg/runtime"
"k8s.io/client-go/kubernetes/scheme"
```

Short example:

```go
func defaultSerializer() runtime.Decoder {
// Add CRDs to scheme
// This is needed to be able to import CRDs from kubernetes manifests files.
_ = apiextensions.AddToScheme(scheme.Scheme)
_ = apiextensionsv1.AddToScheme(scheme.Scheme)
_ = secretsstorev1.AddToScheme(scheme.Scheme)
_ = istionetworkingv1beta1.AddToScheme(scheme.Scheme)
return scheme.Codecs.UniversalDeserializer()
}

_ = kube.Import(
    kube.WithImportSerializer(defaultSerializer()),
    ...
)
```

Full example in the [crd folder](./crd)

## Testing manifests

see [testing manifests](./testing-manifests.md)

## [CLI Utilities](../../cmd/)

### [Explode](../../cmd/explode/)

Converts multi-kubernetes-object manifest into multiple files, organized by namespace.
A CLI was written to make it easier to use in the terminal instead of just a library.

### [Kygo](../../cmd/kygo/)

Converts kubernetes manifests to Go structs.

A CLI was written to make it easier to use in the terminal instead of just a library.
It does support CustomResourceDefinitions but not the custom resources themselves, although it is easy to add them manually.
An example of how to convert custom resource from YAML to Go can be found in the [CRD example](../crd/).

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
