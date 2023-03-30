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
   - example: `wget https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml`
  
2. Convert it to Go structs
   - example `go run cmd/kygo/ -in=./install.yaml -out=output -app=myapp -group`
   - for CRDs, create a `main.go` file with the following content:

      ```go
      func main() {
         err := kube.Import(
            kube.WithImportAppName("my-app"),
            kube.WithImportPackageName("myapp"),
            kube.WithImportOutputDirectory("./output"),
            kube.WithImportManifestFiles([]string{"path/to/myapp.yaml"}),
            kube.WithImportSerializer(defaultSerializer()),
            kube.WithImportRemoveAppName(true),
            kube.WithImportGroupByKind(true),
            kube.WithImportAddMethods(true),
         )
         ...
      }

      func defaultSerializer() runtime.Decoder {
         // ADD MORE CRDS HERE

         _ = apiextensions.AddToScheme(kubescheme.Scheme)
         return kubescheme.Codecs.UniversalDeserializer()
      }
      ```

3. Modify the structs to your liking

4. Export:

   ```go
   // import "github.com/xxx/yyy/output/myapp
   //
   myApp := myapp.New() // function lives in "output/app.go"
   err := kube.Export(myApp, kube.WithExportOutputDirectory("./manifests"))
	if err != nil {
		return err
	}

   ```

5. Apply:

   ```shell
   kubectl apply -f <output folder>
   ```

done.

Have a look at the [tests](../../pkg/kube/) and the [example](../../example/kube/) for a full example.

What does the Go code looks like, see [tekton example](../../pkg/kube/testdata/go/tekton)

### Best practices

> PLEASE READ [best practices](docs/best-practices.md) before using this library.

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
An example of how to do it can be found in the [example](../../example/kube/).

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
