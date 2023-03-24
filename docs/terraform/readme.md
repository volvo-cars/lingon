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
Invoke `terragen` multiple times to generate multiple providers.
See the [Terraform documentation](https://developer.hashicorp.com/terraform/language/providers/requirements) for more information on what these values are.

Additionally, you need to provide an `out` location and the path to the `pkg` for the `out` directory.

We recommend creating a Go file with a `go:generate` directive to invoke the `terragen` command. E.g.

```go
//go:generate go run -mod=readonly github.com/volvo-cars/lingon/cmd/terragen -out ./gen/aws -pkg mypkg/gen/aws -provider local=hashicorp/aws:4.60.0 -force 
```

## Creating and Exporting Terraform Stacks

A Terraform "Stack" in Lingon is the Terraform configuration that makes up a [Root Module](https://developer.hashicorp.com/terraform/language/modules#the-root-module).

A Stack is defined as a Go struct that implements the `terra.Exporter` interface.
For convenience, Lingon provides the `terra.Stack` struct which can be embedded into a struct to implement this interface.
Here is a minimal stack that we export to HCL:

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

func Example_minimalStack() {
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

Lingon uses Go reflection on the struct to identify all fields of a stack struct.
All fields need to be one of:

1. An exported (public) field implementing one of the Terraform object interfaces (such as `terra.Backend`, `terra.Provider`, `terra.Resource` or `terra.DataResource`)
2. A field with a struct tag `lingon:"-"` telling the encoder to ignore the field
3. An embedded struct, whose fields follow these same rules

We therefore recommend splitting complex objects up into multiple structs and embedding them into a parent struct.

### Defining a backend

A backend in Lingon is defined by creating a struct that implements the `terra.Backend` interface.
The schema and list of [available backend types](https://developer.hashicorp.com/terraform/language/settings/backends/configuration)
cannot be automatically obtained, there is no utility to generate code for the backends.
Hence, you are required to define your own backend struct that implements `terra.Backend`.
If you have a good idea for improving this process, please let us know!

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

func Example_backendLocal() {
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

For each Terraform provider that you pass to `terragen`, a Go struct is generated for each of the following Terraform objects:

1. Provider configuration
2. All resources defined by the provider schema
3. All data source defined by the provider schema

All the structs are generated in the same directory to make it easy to use auto-completion to find the object you are looking for.

Here is an example for configuring the `hashicorp/aws` provider and creating a VPC:

```go
type AWSStack struct {
	terra.Stack
	Provider	*aws.Provider	`validate:"required"`
}

// Initialise a stack with the AWS provider configuration
_ = AWSStack{
	Provider: aws.NewProvider(
		aws.ProviderArgs{
			Region: terra.String("eu-north-1"),
		},
	),
}
```

Let's add an example AWS VPC to this stack.

```go
type AWSStack struct {
	terra.Stack
	Provider	*aws.Provider	`validate:"required"`
	VPC		*aws.Vpc	`validate:"required"`
}

// Initialise a stack with the AWS provider configuration
stack := AWSStack{
	Provider: aws.NewProvider(
		aws.ProviderArgs{
			Region: terra.String("eu-north-1"),
		},
	),
	VPC: aws.NewVpc(
		"vpc", aws.VpcArgs{
			CidrBlock:		terra.String("10.0.0.0/16"),
			EnableDnsSupport:	terra.Bool(true),
		},
	),
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
//   required_providers {
//     aws = {
//       source  = "hashicorp/aws"
//       version = "4.60.0"
//     }
//   }
// }
//
// // Provider blocks
// provider "aws" {
//   region = "eu-north-1"
// }
//
// // Resource blocks
// resource "aws_vpc" "vpc" {
//   cidr_block         = "10.0.0.0/16"
//   enable_dns_support = true
// }
```

### Referencing attributes

When creating Terraform configurations it is often necessary to [reference](https://developer.hashicorp.com/terraform/language/expressions/references) other resources or data sources.
Without this capability we found that Terraform stacks will become very small and we would end up with "stack sprawl".
Hence, we decided to add the ability to reference resource and data source attributes within a Terraform stack, with type-safety.

Let's add a subnet to our VPC we created earlier, which requires us to use the VPC ID before:

```go
type AWSStack struct {
	terra.Stack
	Provider	*aws.Provider	`validate:"required"`
	VPC		*aws.Vpc	`validate:"required"`
	Subnet		*aws.Subnet	`validate:"required"`
}

vpc := aws.NewVpc(
	"vpc", aws.VpcArgs{
		CidrBlock:		terra.String("10.0.0.0/16"),
		EnableDnsSupport:	terra.Bool(true),
	},
)
subnet := aws.NewSubnet(
	"subnet", aws.SubnetArgs{
		// Reference the VPC's ID, which will translate into a reference
		// in the Terraform configuration
		VpcId: vpc.Attributes().Id(),
	},
)

// Initialise a stack with the AWS provider configuration
stack := AWSStack{
	Provider: aws.NewProvider(
		aws.ProviderArgs{
			Region: terra.String("eu-north-1"),
		},
	),
	VPC:	vpc,
	Subnet:	subnet,
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
//   required_providers {
//     aws = {
//       source  = "hashicorp/aws"
//       version = "4.60.0"
//     }
//   }
// }
//
// // Provider blocks
// provider "aws" {
//   region = "eu-north-1"
// }
//
// // Resource blocks
// resource "aws_vpc" "vpc" {
//   cidr_block         = "10.0.0.0/16"
//   enable_dns_support = true
// }
//
// resource "aws_subnet" "subnet" {
//   vpc_id = aws_vpc.vpc.id
// }
```

### Running the Terraform CLI

Lingon does not provide a Terraform CLI client.
We feel it would be too opinionated and context-specific that it is easier for users to build their own on top of the
[terraform-exec](https://github.com/hashicorp/terraform-exec) library.

### Accessing the Terraform state

Lingon includes a simple utility for importing a Terraform state back into a Terraform Stack, and into the relevant resources within that stack.
Here is an example using the AWS VPC.

```go
type AWSStack struct {
	terra.Stack
	Provider	*aws.Provider	`validate:"required"`
	VPC		*aws.Vpc	`validate:"required"`
}

// Initialise a stack with the AWS provider configuration
stack := AWSStack{
	Provider: aws.NewProvider(
		aws.ProviderArgs{
			Region: terra.String("eu-north-1"),
		},
	),
	VPC: aws.NewVpc(
		"vpc", aws.VpcArgs{
			CidrBlock:		terra.String("10.0.0.0/16"),
			EnableDnsSupport:	terra.Bool(true),
		},
	),
}

// At this point, you would invoke the Terraform CLI, and at a minimum
// run the `terraform show` command to get the state back.
// The state can then be decoded back into our stack.
// For this test, we will create some dummy state data
// (don't do this at home!)

state := tfjson.State{
	Values: &tfjson.StateValues{
		RootModule: &tfjson.StateModule{
			Resources: []*tfjson.StateResource{
				{
					Type:	"aws_vpc",
					Name:	"vpc",
					AttributeValues: map[string]interface{}{
						"id": "12345",
					},
				},
			},
		},
	},
}
ok, err := terra.StackImportState(&stack, &state)
if err != nil {
	slog.Error("importing stack state", "error", err)
}
if !ok {
	// This means the stack includes resources that did not have values
	// in the Terraform state.
	// This is happens if you have not applied your Terraform configuration.
	slog.Info("stack state is not complete")
}
// Access the VPC state. If you know the state is complete, you can also use
// the StateMust() function
vpcState, ok := stack.VPC.State()
if !ok {
	slog.Info("vpc does not have state")
	return
}
fmt.Println(vpcState.Id)
// Output: 12345
```

In traditional Terraform setups, one would use [output values](https://developer.hashicorp.com/terraform/language/values/outputs) to get values from a module or from state.
We decided not to support outputs as we found having access to the entire state with Go type-safety to be so much better.

## Type system

Lingon's [terra](../../pkg/terra) package includes a very minimal type system.
The reason for creating this was to allow for [references](https://developer.hashicorp.com/terraform/language/expressions/references) within Terraform configuration.
See examples in the package for more details on how the type system is used.

One tip is to create some package-level vars to avoid bloat in the code, e.g.

```go
var (
	S	= terra.String
	B	= terra.Bool
)

_ = aws.NewVpc(
	"vpc", aws.VpcArgs{
		CidrBlock:		S("10.0.0.0/16"),
		EnableDnsSupport:	B(true),
	},
)
```

We also highly recommend not passing `terra` values between stacks or Go modules.
The API for a Terraform stack should only pass native Go values to avoid potentially passing references to Terraform attributes in separate stacks.

## Known limitations

There are many features of Terraform that are not supported by Lingon.
This is partly by design, as we believe Lingon should not be a direct port your Terraform configurations from HCL to Go without taking full advantage of what it offers.

### Not supported

Using [modules](https://developer.hashicorp.com/terraform/language/modules/syntax) is not supported.
We re-wrote the Terraform modules we needed as Go packages with a `NewXYZ` function to create a struct which we embed in our stack.
Other blocks like `locals`, `variable`, `output` are not supported.
We have not come across a case where we needed them.

Meta-arguments such as `for_each` or `count` are not supported.
Instead, you can use for-loops in Go.
If you need to `for_each` over values from the state, we suggest splitting your stack up into two.
Then in your second stack you can use the state values from the first.

There is no support for the long list of [Terraform functions](https://developer.hashicorp.com/terraform/language/functions).
If you need to do complex manipulation of data we suggest doing it in native Go.

If you have a feature request for missing functionality, please raise an issue.
