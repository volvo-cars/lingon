
## Go (Golang)


## Why Go

* _Configuration as code_ as in programming language code, not JSON, YAML or HCL.

* _Go is a strongly typed language_. IDEs provide a great developer experience with autocompletion and type safety.

* _Tests are a first class citizen_: it makes trivial to write a test to ensure correctness of the configuration

* _Most of the infrastructure tools are written in Go_: i.e. Kubernetes, Prometheus, Istio, Tekton, ...
  it is easy to interact with them and leverage their APIs/structs.

* _Dependency management_: Go modules provide a simple way
  to manage dependencies and reuse code with the correct version to ensure compatibility

* _Documentation_: Go recommends _godoc_ formatting as it can leverage native comments for each struct's fields
  in order to provide comments/examples in the IDE of the developer, which increase productive and correctness.

* _Quick feedback loop_: Catch most mistakes and incompatibilities in Go compile time.
  Go has very fast compilation time, which feels like you are running a script.

* _Limit the number of languages used in the organization to a minimum_ : Go is one of the cleanest,
  simplest and developer friendly languages that exists.

* _Backwards compatibility_: Go is well-known for its backwards compatibility making
  it a great choice for long-lasting platforms reducing rework through breaking changes.

### References

- [But Why Go](https://github.com/bwplotka/mimic#but-why-go) from [Mimic](https://github.com/bwplotka/mimic)
- [Go for Cloud](https://rakyll.org/go-cloud/) by [rakyll](https://rakyll.org)
- [The yaml document from hell](https://ruudvanasseldonk.com/2023/01/11/the-yaml-document-from-hell) by [ruudvanasseldonk](https://ruudvanasseldonk.com)
- [noyaml](https://noyaml.com)

## Go resources

- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [Google Style guide](https://google.github.io/styleguide/go/guide)
- [Google best practices](https://google.github.io/styleguide/go/best-practices)
- [Context about decisions](https://google.github.io/styleguide/go/decisions)

