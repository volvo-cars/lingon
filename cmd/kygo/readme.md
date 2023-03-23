# Kygo

Convert kubernetes YAML manifests to Go structs.

Why? Because we found it easier to manipulate manifests and automate with Go code than YAML.
Also, it interacts nicely with our other tool that manages our (multi) cloud infrastructure: [lingon](https://github.com/volvo-cars/lingon)

Why Go? Smarter people have a better to explain it: <https://github.com/bwplotka/mimic/blob/main/README.md#but-why-go>

## Usage

```
Usage of kygo:
  -app string
        specify the app name. This will be used as the package name and the prefix for the generated files. (default "myapp")
  -clean-name
        specify if the app name should be removed from the variable, struct and file name. (default true)
  -group
        specify if the output should be grouped by kind (default) or split by name. (default true)
  -in string
        specify the input directory of the yaml manifests, '-' for stdin (default "-")
  -out string
        specify the output directory for manifests. (default "out")
```

## Example

```shell
go build -o ./bin/kygo ./cmd/kygo
./bin/kygo -in=./pkg/kube/testdata/argocd.yaml -out=./out -app=argocd
ls -Rl1 out/
```

```
out/argocd
├── app.go
├── cluster-role-binding.go
├── cluster-role.go
├── config-map.go
├── custom-resource-definition.go
├── deployment.go
├── network-policy.go
├── role-binding.go
├── role.go
├── secret.go
├── service-account.go
├── service.go
└── stateful-set.go
```

## Deploying to kubernetes

Either:

* Use the `go-kart` library to generate the yaml from Go.
* Use the `k8s.io/client-go` library to directly apply to kubernetes.

### Using go-kart

```go
package main

import (
 "context"
 "os"
 "path/filepath"

 "github.com/XXX/YYY/myapp"
 "github.com/volvo-cars/lingon/pkg/cmdexec"
 "github.com/volvo-cars/lingon/pkg/kube"
)

func main() {
 app := myapp.New()
 manifestOut := filepath.Join(".k8s", "myapp")
 // create the output directory if it does not exist
 // and generate the yaml manifests in the .k8s/myapp/ directory
 if err := kube.Export(app, manifestOut); err != nil {
  panic(err)
 }
 
 // deploy the manifests to kubernetes, but it needs *kubectl* to be installed.
 ctx := context.Background()
 if err := cmdexec.KubectlOut(ctx, os.Stdout, os.Stderr, "apply", "-f", manifestOut); err != nil {
  panic(err)
 }
}
```

### Using client-go

Please refer to the [client-go repo](https://github.com/kubernetes/client-go).

There is an interesting issue on GitHub about "Add go generics support to client-go":
<https://github.com/kubernetes/kubernetes/issues/106846>
