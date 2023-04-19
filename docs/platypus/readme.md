# Platypus  

## ⚠️ experimental ⚠️

This is a test case we use to test new concepts and get a feel for the APIs.

The Terraform code is in `pkg/infra` while the kubernetes manifestss are in `pkg/platform`.

The code entrypoint is in `cmd/platypus/cli.go`.

## Getting started

### Prerequisites

- [aws cli](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [terraform](https://learn.hashicorp.com/terraform/getting-started/install.html)

### Setup

Authenticate: Follow the instructions [here](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) to configure your AWS credentials.

  ```bash
  # if SSO
  aws sso login --profile=platypus-xxx
  # get the identity
  aws sts get-caller-identity --profile platypus-xxx
  ```

### Run

Terraform plan

  ```bash
  go run ./cmd/platypus/ --plan
  ```

Terraform apply

  ```bash
  go run ./cmd/platypus/ --apply
  ```
