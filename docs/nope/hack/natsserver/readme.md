# NATS Server

This runs a minimal NATS server embedded in the small Go program.

What makes this special is that it generates an operator, system account and system user, and writes the operator NKey seed and user credentials.

Why? Because we need this for testing the operator.

For convenience, it outputs a `out/.env` file which can be sourced to get your environment setup for running the operator.

## WARNING

Run this from this directory because there is no path handling right now.

```bash
cd hack/natsserver
go run main.go
```

## TODO

1. Give this a better name
2. Handle paths better
3. Provide some args like port to run on
