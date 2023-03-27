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

3. Modify the structs to your liking

4. Export:

```go
// import "github.com/xxx/yyy/output/myapp
// manifests found in `output/` folder
myApp := myapp.New() // function lives in "output/app.go"
if err := kube.Export(myApp, output); err != nil {
   return err
}
```

5. Apply:

```shell
kubectl apply -f <output folder>
```

done.

Have a look at the [tests for more examples](../../pkg/kube/].

### Best practices

> PLEASE READ [best practices](docs/best-practices.md) before using this library.

This project has been heavily inspired by :

- [Mimic](https://github.com/bwplotka/mimic) (you should definitely check it out, we copied the best practices from there)
- [NAML](https://github.com/krisnova/naml) (we found out about [valast](https://github.com/hexops/valast) from there)

Honorable mentions:

- [valast](https://github.com/hexops/valast) convert Go structs to its Go code representation
- [jennifer](https://github.com/dave/jennifer) generate Go code

## CLI Utilities

### Explode

Converts multi-kubernetes-object manifest into multiple files, organized by namespace.
A CLI was written to make it easier to use in the terminal instead of just a library.

### Kygo

Converts kubernetes manifests to Go structs.
A CLI was written to make it easier to use in the terminal instead of just a library.
It does support CustomResourceDefinitions but not the custom resources themselves, although it is easy to add them manually.
An example of how to do it can be found in the [examples](../../examples/).

## Packages

Have a look at the godoc for more information.

### Kube

- `App` struct that is embedded to mark kubernetes applications
- `Export` kubernetes objects defined as Go struct to kubernetes manifests in YAML.
- `Explode` kubernetes manifests in YAML to multiple files, organized by namespace.
- `Import` kubernetes manifests in YAML to Go structs.

### Kubeconfig

Manipulate kubeconfig files **without** any dependencies on `go-client`.

### KubeUtil

Reusable functions used to create kubernetes objects in Go.

### Testutils

Reusable test functions.

## Why Go

- [But Why Go](https://github.com/bwplotka/mimic#but-why-go) from [Mimic](https://github.com/bwplotka/mimic)
- [Go for Cloud](https://rakyll.org/go-cloud/) by [rakyll](https://rakyll.org)
- [The yaml document from hell](https://ruudvanasseldonk.com/2023/01/11/the-yaml-document-from-hell) by [ruudvanasseldonk](https://ruudvanasseldonk.com)