[![Go Reference](https://pkg.go.dev/badge/github.com/volvo-cars/lingon.svg)](https://pkg.go.dev/github.com/volvo-cars/lingon)
[![GoReportCard example](https://goreportcard.com/badge/github.com/volvo-cars/lingon)](https://goreportcard.com/report/github.com/volvo-cars/lingon)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# Lingon 🍒 - Libraries for building platforms with Go  <!-- omit in toc -->

- [What is this?](#what-is-this)
- [Who is this for?](#who-is-this-for)
- [Project status](#project-status)
- [Getting started](#getting-started)
- [Examples](#examples)
- [Motivation](#motivation)
  - [Reduce cognitive load](#reduce-cognitive-load)
  - [Type safety](#type-safety)
  - [Error handling](#error-handling)
  - [Limitless automation](#limitless-automation)
- [Why Go?](#why-go)
- [Similar projects](#similar-projects)
- [License](#license)

## What is this?

Lingon is a collection of libraries and tools for building platforms using Go.

> In short, we write Terraform (HCL) and Kubernetes (YAML) in Go. see [Rationale](./docs/rationale.md) for more details.

The following technologies are currently supported:

- [Terraform](./docs/terraform/)
- [Kubernetes](./docs/kubernetes/)

The only dependencies you need are:

- Go
- Terraform CLI
- kubectl

## Who is this for?

Lingon is aimed at people who need to automate the lifecycle of their cloud infrastructure
and have suffered the pain of configuration languages and complexity of gluing tools together with more tools.

## Project status

This project is in beta.
The APIs are stable, but we do not promise backward compatibility at this point.
We will eventually promise backward compatibility when the project is more battle tested.

See [FAQ](./docs/faq.md) for more details.

## Getting started

> Note that in the terraform case, the code generation is fast.
> Compiling all the generated resources will take a while.
> Thankfully, Go is fast at compiling and keeps a cache of compiled packages.
> Expect to wait a few minutes the first time you run `go build` after generating the code.

- [Terraform](./docs/terraform/)
- [Kubernetes](./docs/kubernetes/)

## Examples

All the [Examples](./docs/) are in the [documentation](./docs).

A big example is [Platypus](./docs/platypus/) which shows how 
the [kubernetes](./docs/kubernetes/) and [terraform](./docs/terraform/) libraries can be used together.

## Motivation

See [Rationale](./docs/rationale.md) for more details.

Lingon was developed to achieve the following goals:

### Reduce cognitive load

Building a platform within a single context (i.e. Go) will reduce cognitive load by decreasing the number of tools and context switching in the process.
It provides a better developer experience with out-of-the-box IDE support and a single language to learn with smooth learning curve.

### Type safety

Detect misconfigurations in your text editor by using type-safe Go structs to exchange values across tool boundaries.
This "shifts left" the majority of errors that occur to the earliest possible point in time.

### Error handling

Go's error handling enables propagating meaningful errors to the user.
This significantly reduces the effort in finding the root cause of errors and provides a better developer experience.

### Limitless automation

We are only limited by what a programming language can do. 
We can reuse part of what we build in libraries without external tooling.
That is not possible with YAML as doesn't support "includes", therefore we need a tool for that.
Configuration languages are limited by the features they provide.
Gluing tools together with more tools and configuration to manage more tools and configuration is not a sustainable approach.
We do use a limited set of tools that we learn well and can extend, but we automate them and test them together using Go.

Note that we are in a particular situation where we need custom automation of the lifecycle of our cloud infrastructure.

## Why Go?

Because most outages are caused by a configuration error.

- [Why Go](./docs/go.md) from us, inspired by 👇
- [But Why Go](https://github.com/bwplotka/mimic#but-why-go) from [Mimic](https://github.com/bwplotka/mimic)
- [Not Another Markup Language](https://github.com/krisnova/naml) from [NAML](https://github.com/krisnova/naml)
- [Go for Cloud](https://rakyll.org/go-cloud/) by [rakyll](https://rakyll.org)
- [The yaml document from hell](https://ruudvanasseldonk.com/2023/01/11/the-yaml-document-from-hell) by [ruudvanasseldonk](https://ruudvanasseldonk.com)
- [noyaml website](https://noyaml.com)
- [YAML Considered Harmful - Philipp Krenn 🎥](https://youtu.be/WQurEEfSf8M)
- [Nightmares On Cloud Street 29/10/20 - Joe Beda 🎥](https://youtu.be/8PpgqEqkQWA)

## Similar projects

See [Comparison](./docs/comparison.md) for more details.

## License

This code is released under the [Apache-2.0 License](./LICENSE).
