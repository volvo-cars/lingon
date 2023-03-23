# Documentation

To view the godoc, run :

```shell
godoc -http=:6060
# and open http://localhost:6060/pkg/github.com/volvo-cars/lingon/pkg/kube/
```

## Run tests

```shell
go test -v ./...
# with coverage
go test -v -coverprofile=coverage.out ./...
```

## build everything

see [build.go](../hack/build.go)

```shell
go generate ./...
```
