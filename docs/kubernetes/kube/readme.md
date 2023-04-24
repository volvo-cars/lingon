# Kube example


This example shows you how to convert a YAML manifest to Go structs and then export it back to YAML.

> Note that the secrets will be redacted.

## Usage

```shell
go mod download
go generate ./...
ls -l1 out/
```

The output should look like this:

```shell
$ ls -l1 out/tekton
app.go
cluster-role-binding.go
cluster-role.go
config-map.go
custom-resource-definition.go
deployment.go
horizontal-pod-autoscaler.go
mutating-webhook-configuration.go
namespace.go
role-binding.go
role.go
secret.go
service-account.go
service.go
validating-webhook-configuration.go
```

## Example

```go
tk := tekton.New()

out := filepath.Join("out", "export")
fmt.Printf("exporting to %s\n", out)

_ = os.RemoveAll(out)
defer os.RemoveAll(out)

err := kube.Export(tk, kube.WithExportOutputDirectory(out))
if err != nil {
	panic(err)
}

// or use io.Writer
var buf bytes.Buffer
_ = kube.Export(tk, kube.WithExportWriter(&buf))

ar := txtar.Parse(buf.Bytes())

fmt.Printf("\nexported %d manifests\n\n", len(ar.Files))
fmt.Println("\t>>> first manifest <<<")
if len(ar.Files) > 0 {
	// avoiding diff due to character invisible to the naked eye 😅
	l := strings.ReplaceAll(string(ar.Files[0].Data), "\n", "\n\t")
	fmt.Printf("\t%s\n", l)
}

// Output:
// exporting to out/export
//
// exported 65 manifests
//
//	>>> first manifest <<<
//	apiVersion: rbac.authorization.k8s.io/v1
//	kind: ClusterRole
//	metadata:
//	  labels:
//	    app.kubernetes.io/instance: default
//	    app.kubernetes.io/part-of: tekton-pipelines
//	    rbac.authorization.k8s.io/aggregate-to-admin: "true"
//	    rbac.authorization.k8s.io/aggregate-to-edit: "true"
//	  name: tekton-aggregate-edit
//	rules:
//	- apiGroups:
//	  - tekton.dev
//	  resources:
//	  - tasks
//	  - taskruns
//	  - pipelines
//	  - pipelineruns
//	  - pipelineresources
//	  - runs
//	  - customruns
//	  verbs:
//	  - create
//	  - delete
//	  - deletecollection
//	  - get
//	  - list
//	  - patch
//	  - update
//	  - watch
//
//
```
