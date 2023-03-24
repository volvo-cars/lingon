# Getting started with Lingon for Terraform

Lingon provides support for creating Terraform configurations via Go.
It does not provide a client library for calling the Terraform CLI, but one can be easily built using [terraform-exec](https://github.com/hashicorp/terraform-exec).

The following steps are needed to start using Lingon for Terraform:

1. Generate Go code for Terraform provider(s)
2. Create and Export a Terraform Stack
3. Run the Terraform CLI over the exported configurations

## Generating Go code for Terraform provider(s)

The [terragen](../../cmd/terragen) command is used to generate Go code for some Terraform providers.
First you need to decide which providers you want to use and provide them as arguments to the `terragen` command.
The generator requires three values for each provider:

1. Local name
2. Source
3. Version

These three values are used in Terraform's `required_providers` block.
For example, given the following `required_providers` block:

```terraform
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.60.0"
    }
  }
}
```

the argument to `terragen` for the provider would be `-provider aws=hashicorp/aws:4.60.0`.
You can provide multiple `-provider` arguments to the same call.
See the [Terraform documentation](https://developer.hashicorp.com/terraform/language/providers/requirements) for more information on what these values are.

Additionally, you need to provide an `out` location and the path to the `pkg` for the `out` directory.

We recommend creating a Go file with a `go:generate` directive to invoke the `terragen` command. E.g.

```go
//go:generate go run -mod=readonly github.com/volvo-cars/lingon/cmd/terragen -out ./gen -pkg mypkg/gen -provider local=hashicorp/aws:4.60.0 -force 
```

## Creating and Exporting Terraform Stacks

A Terraform "Stack" in Lingon is the Terraform configurations that make up a ["Root Module"](https://developer.hashicorp.com/terraform/language/modules#the-root-module).

A Stack is defined as a Go struct that, at a minimum, implements the `terra.Exporter` interface.
For convenience, Lingon provides the `terra.Stack` struct which can be embedded into a struct to implement this interface.
Creating a minimal stack and exporting it would look like the following:

```go
package terraform_test

import (
	"bytes"
	"fmt"

	"github.com/volvo-cars/lingon/pkg/terra"
	"golang.org/x/exp/slog"
)

type MinimalStack struct {
	terra.Stack
}

func ExampleMinimalStack() {
	// Initialise the minimal stack
	stack := MinimalStack{}
	// Export the stack to Terraform HCL
	var b bytes.Buffer
	if err := terra.ExportWriter(&stack, &b); err != nil {
		slog.Error("exporting stack", "err", err)
		return
	}
	fmt.Println(b.String())

	// Output:
	// terraform {
	// }
}
```

This is obviously not very useful, so let's add some more things to our stack.

### Defining a backend

A backend in Lingon is defined by creating a struct that implements the `terra.Backend` interface.
As the schema and list of [available backend types](https://developer.hashicorp.com/terraform/language/settings/backends/configuration)
cannot be automatically obtained, there is no utility to generate code for the backends.
If you have a good idea for this, please let us know!

```go
package terraform_test

import (
	"bytes"
	"fmt"

	"github.com/volvo-cars/lingon/pkg/terra"
	"golang.org/x/exp/slog"
)

var _ terra.Backend = (*BackendLocal)(nil)

// BackendLocal implements the Terraform local backend type.
// https://developer.hashicorp.com/terraform/language/settings/backends/local
type BackendLocal struct {
	Path string `hcl:"path,attr" validate:"required"`
}

// BackendType defines the type of the backend
func (b *BackendLocal) BackendType() string {
	return "local"
}

type BackendStack struct {
	terra.Stack

	Backend	*BackendLocal
}

func ExampleBackendLocal() {
	stack := BackendStack{
		Backend: &BackendLocal{Path: "terraform.tfstate"},
	}
	// Export the stack to Terraform HCL
	var b bytes.Buffer
	if err := terra.ExportWriter(&stack, &b); err != nil {
		slog.Error("exporting stack", "err", err)
		return
	}
	fmt.Println(b.String())

	// Output:
	// terraform {
	//   backend "local" {
	//     path = "terraform.tfstate"
	//   }
	// }
}
```

Note that the `hcl` struct tags are necessary on your custom backend for the HCL encoder to work.
You can define your backend to include only the minimal fields that you actually need and create a nice helper function for it.

### Defining providers, resources and data sources

TODO

### Referencing attributes

TODO

### Running the Terraform CLI

TODO

### Accessing the Terraform state

TODO

## Lingon's type system

TODO

## Known limitations

TODO: terraform functions, for_each, count
