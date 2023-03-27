# Development

## Documentation

To view the godoc, run :

```shell
godoc -http=:6060
# and open http://localhost:6060/pkg/github.com/volvo-cars/lingon/
```

## Run tests

```shell
go test -v ./...
# with coverage
go test -v -coverprofile=coverage.out ./...
```

## Automation

see [mage](https://magefile.org/)

```shell
# list all available targets
mage -l
```
