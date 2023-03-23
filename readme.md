# lingon - libraries for building platforms in Go

## What is this?

Lingon is a collection of libraries and tools for building platforms using Go.

The following technologies are supported:

1. Terraform
2. Kubernetes

## Who is this for?

Lingon is aimed at people managing cloud infrastructure who have suffered the pain of configuration languages and complexity of gluing tools together with more tools.

## Getting started

TODO: link to docs for each tool

## Motivation

Lingon was developed to achieve the following goals:

### Reduce cognitive load

Building a platform within a single context (i.e. Go) will reduce cognitive load by decreasing the number of tools and context switching in the process.

### Type safety

Detect misconfigurations in your text editor by using type-safe Go structs to exchange values across tool boundaries.
This "shifts left" the majority of errors that occur to the earliest possible point in time.

### Error handling

Go's error handling enables propagating meaningful errors to the user.
This significantly reduces the effort in finding the root cause of errors and provides a better developer experience.

### Limitless automation

TODO: something about less glue means more possibilities... Maybe mention "test first" approach?

## Why Go?

- [But Why Go](https://github.com/bwplotka/mimic#but-why-go) from [Mimic](https://github.com/bwplotka/mimic)
- [Go for Cloud](https://rakyll.org/go-cloud/) by [rakyll](https://rakyll.org)
- [The yaml document from hell](https://ruudvanasseldonk.com/2023/01/11/the-yaml-document-from-hell) by [ruudvanasseldonk](https://ruudvanasseldonk.com)

## Similar projects

TODO: write about similar projects and differences.
Primarily about being Go idiomatic and using structs to be "declarative".
Type-safety all the way down.
Don't involve nodejs/jsii.

- Pulumi: https://www.pulumi.com/
- CDK for AWS: https://aws.amazon.com/cdk/
- CDK for Terraform (cdktf): https://developer.hashicorp.com/terraform/cdktf
- CKD for Kubernetes (cdkk8s): https://cdk8s.io/

## Project status

TODO: something about the project status

## License

This code is released under the [Apache-2.0 License](./LICENSE).