# PLATYPUS 2

## Requirements 

### Generate providers types

Using `terragen` :

```shell 
# install terragen
go install github.com/volvo-cars/lingon/cmd/terragen@latest

# generate providers
go generate -x ./...
```

### Valid credentials

```shell
# for sso
aws sso login --profile=XXXX

# for access key
export AWS_ACCESS_KEY_ID="anaccesskey"
export AWS_SECRET_ACCESS_KEY="asecretkey"
export AWS_REGION="eu-north-1"
```

## Deploying

```shell

# /!\ except the first compilation to peg the CPU to 100%
# after that, the build cache is hot and compilation is fast
go run ./cmd/platypus/*.go -plan

# it will only show the plan for the VPC and the rest of the state
# depends on it. 

# this will deploy everything
go run ./cmd/platypus/*.go -apply

```

## WIP

Getting new manifests from hell charts to Go

```shell
# execute script (will be converted later in a CLI)
./docs/platypus2/scripts/mleh.sh
```
