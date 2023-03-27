# Kygo

Convert kubernetes YAML manifests to Go structs.

Why? Because we found it easier to manipulate manifests and automate with Go code than YAML.

Why Go? Smarter people have a better to explain it: <https://github.com/bwplotka/mimic/blob/main/README.md#but-why-go>

## Usage

```
Usage of kygo:
  -app string
    	specify the app name. This will be used as the package name if none is specified. (default "myapp")
  -clean-name
    	specify if the app name should be removed from the variable, struct and file name. (default true)
  -group
    	specify if the output should be grouped by kind (default) or split by name. (default false)
  -in string
    	specify the input directory of the yaml manifests, '-' for stdin (default "-")
  -out string
    	specify the output directory for manifests. (default "out")
  -pkg string
    	specify the Go package name. Cannot contain a dash. If none is specified the app name will be used.
```

## Example

```shell
go build -o ./bin/kygo ./cmd/kygo
./bin/kygo -in=./pkg/kube/testdata/argocd.yaml -out=./out -app=argocd -group
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

* Use the `lingon` library to generate the yaml from Go.
* Use the `k8s.io/client-go` library to directly apply to kubernetes.

### Using lingon

```go
package main

import (
    "context"
    "path/filepath"

    "github.com/XXX/YYY/myapp"
    "github.com/volvo-cars/lingon/pkg/kube"
)

func main() {   
	app := myapp.New()
    manifestOut := filepath.Join("manifests", "myapp")
	
	// it will create the output directory if it does not exist
	// and generate the YAML manifests in the directory manifests/myapp/
	if err := kube.Export(app, manifestOut); err != nil {
  panic(err)
 }
    
    // OR 
	// apply the manifests to kubernetes directly to the cluster
	// it will pass the manifest output to  `kubectl apply -f -`
	if err := app.Apply(context.Background()); err != nil {
        panic(err)
    }
	
	// check if the manifests are applied correctly
	// ...
}
```

### Using client-go

Please refer to the [client-go repo](https://github.com/kubernetes/client-go).

There is an interesting issue on GitHub about "Add go generics support to client-go":
<https://github.com/kubernetes/kubernetes/issues/106846>
