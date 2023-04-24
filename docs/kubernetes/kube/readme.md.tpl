# Kube example

- [Usage](#usage)
- [Example](#example)

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

{{ "Example" | example }}
