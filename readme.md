[![Go Reference](https://pkg.go.dev/badge/github.com/volvo-cars/lingon.svg)](https://pkg.go.dev/github.com/volvo-cars/lingon)
[![GoReportCard example](https://goreportcard.com/badge/github.com/volvo-cars/lingon)](https://goreportcard.com/report/github.com/volvo-cars/lingon)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/github/go-mod/go-version/volvo-cars/lingon)](go.mod)
![Go Test Coverage](.github/coverage.svg)


# Lingon 🍒 - Libraries for building platforms with Go  <!-- omit in toc -->

- [What is this?](#what-is-this)
  - [It can do 4 things so far](#it-can-do-4-things-so-far)
- [Who is this for?](#who-is-this-for)
- [Project status](#project-status)
- [Required knowledge](#required-knowledge)
- [Getting started](#getting-started)
- [Examples](#examples)
- [Install binaries](#install-binaries)
  - [Homebrew](#homebrew)
  - [Github releases](#github-releases)
  - [Install from source](#install-from-source)
- [Motivation](#motivation)
- [Extraordinary use cases](#extraordinary-use-cases)
- [Why Go?](#why-go)
- [Similar projects](#similar-projects)
- [License](#license)

## What is this?

Write [Terraform](./docs/terraform/) (HCL) and [Kubernetes](./docs/kubernetes/) (YAML) in Go. see [Rationale](./docs/rationale.md) for more details.

> Lingon is not a platform, it is a thin wrapper around terraform and kubernetes API in a library
> meant to be consumed in a Go application that platform engineers write to manage their platforms.
> It is a tool to build and automate the creation and the management of platforms regardless of the target infrastructure and services.

### It can do 4 things so far

- import kubernetes YAML manifests to valid Go code (even CRDs)  `kube.Import`
- export Go code to kubernetes YAML manifests  `kube.Export`
- generate Go code from Terraform providers `terragen.GenerateProviderSchema` and `terragen.GenerateGoCode`
- export terraform Go code to valid Terraform HCL  `terra.Export`

The only dependencies you need are:

- Go
- Terraform CLI
- kubectl

## Who is this for?

Lingon is aimed at advanced platform teams who need to automate the lifecycle of their cloud infrastructure
and have suffered the pain of configuration languages and complexity of gluing tools together with more tools.

## Project status

This project is in beta.
The APIs are stable, but we do not promise backward compatibility at this point.
We will eventually promise backward compatibility when the project is more battle tested.

See [FAQ](./docs/faq.md) for more details.

## Required knowledge

This is not a tutorial on how to use Go, Terraform or Kubernetes.
Lingon doesn't try to hide the complexity of these technologies, it embraces it.

> Which is why you need to know how to use these technologies to use Lingon.

- [Go](https://golang.org/)
- [Terraform](https://www.terraform.io/)
- [Kubernetes](https://kubernetes.io/)

## Getting started

> Note that in the terraform case, the code generation is fast.
> Compiling all the generated resources will take a while.
> Thankfully, Go is fast at compiling and keeps a cache of compiled packages.
> Expect to wait a few minutes the first time you run `go build` after generating the code.

- [Terraform](./docs/terraform/)
- [Kubernetes](./docs/kubernetes/)

## Examples

- All the [Examples](./docs/) are in the [documentation](./docs).
- Convert kubernetes manifests from YAML containing CRDs to Go: [example_import_test.go](./docs/kubernetes/crd/example_import_test.go)
- A web app to showcase the conversion from kubernetes manifests from YAML to Go: <https://lingonweb.bisconti.cloud/>
- Export Go code to kubernetes manifests: [kube_test.go](./docs/kubernetes/kube/kube_test.go)
- Create the HCL code for Terraform: [aws_test.go](./docs/terraform/aws_test.go)
- An example is [Platypus](./docs/platypus/) which shows how
the [kubernetes](./docs/kubernetes/) and [terraform](./docs/terraform/) libraries can be used together.

## Install binaries

Lingon provides helper binaries.

- `explode` - explode a kubernetes manifests YAML file into multiple files organized by kind and namespace.
- `kygo` - convert kubernetes YAML manifests to Go code
- `terragen` - generate Go code from Terraform providers

### Homebrew

You can install the binaries with [Homebrew](https://brew.sh) :

```bash
brew tap golingon/homebrew-tap
brew install lingon
```

### Github releases

Or simply download the binaries from the [releases](https://github.com/volvo-cars/lingon/releases/latest) page.

### Install from source

```bash
go install github.com/volvo-cars/lingon/cmd/explode@latest
go install github.com/volvo-cars/lingon/cmd/kygo@latest
go install github.com/volvo-cars/lingon/cmd/terragen@latest 

```

## Motivation

See [Rationale](./docs/rationale.md) for more details.

## Extraordinary use cases

Lingon might be helpful if you need to:

- use the SDK of your cloud provider to access APIs (alpha, beta, deprecated) not included in a Terraform provider.
- authenticate to a multitude of providers or webhook with specific requirements (e.g. Azure SSO, AWS, Github, Slack, etc.)
- automate some parts of the infrastructure that are really hard to test (e.g. iptables, DNS, IAM, etc.)
- store the state of the infrastructure in a database for further analysis
- collect advanced metrics about the failures occurring during the deployment of the infrastructure
- enforce advanced rules on kubernetes manifests before deploying it (e.g. every service account must be related to a role and that role cannot have '*' in access rights, etc.)
- define CI/CD pipelines as imperative code, not declarative.
- execute smoke tests after deploying changes to the platform (HTTP, gRPC, DB connection, etc.)
- write unit tests for your infrastructure

## Why Go?

See [Why Go](./docs/go.md) for more details.

## Similar projects

See [Comparison](./docs/comparison.md) for more details.

## License

This code is released under the [Apache-2.0 License](./LICENSE).
