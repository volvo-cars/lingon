# Development

## FIRST run this

```shell
git clone https://github.com/volvo-cars/lingon.git
cd lingon

go mod tidy
go generate -v ./...
```

to download dependencies and to generate the code.

## Documentation

To view the godoc, run :

```shell
godoc -http=:6060
# and open http://localhost:6060/pkg/github.com/volvo-cars/lingon/
```

Or if you prefer to have the same experience as <https://pkg.go.dev>

```shell
go install golang.org/x/pkgsite/cmd/pkgsite@latest && pkgsite

# open http://localhost:8080/github.com/volvo-cars/lingon.
```

To run tests for the documentation and examples, run:

```shell
cd docs 
go generate -v ./... && go test -v ./...
```

⚠️ Running the tests for the docs will take a while (+15min on M1 pro) ⚠️

It will **peg your CPU to 100%** during that time.

When `Platypus` is built and tested, there are many (many!!) Go generated files to compile.
On the second run, the cache will be used, and it will be much faster.

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

## Release

The release process is automated with [goreleaser](https://goreleaser.com/).
A release is triggered by pushing a tag to the repo.

```shell
function version () {
 local shortsha=$(git rev-parse --short HEAD) # will output 91d9a52
 local shortdate=$(date "+%F")                # will output 2021-01-01
 echo "$shortdate-$shortsha"                  # will output 2021-01-01-91d9a52
}

git tag -a $(version) -s -m "Release $(version)" && git push --tags
```
