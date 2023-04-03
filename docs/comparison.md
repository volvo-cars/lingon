# Comparison <!-- omit in toc -->

There are no shortage of tools for [Kubernetes](https://docs.google.com/spreadsheets/d/1FCgqz1Ci7_VCz_wdh8vBitZ3giBtac_H8SBw4uxnrsE/edit#gid=0) and to manage infrastructure.

Most of the tools are good tools, they just don't work for us or are too complicated and requires a lot of investment to get started.

> We may sound negative, but it is not our intention.
> We are just trying to be honest about our experience.

- [Other programming languages](#other-programming-languages)
- [Mimic](#mimic)
- [NAML](#naml)
- [Kustomize](#kustomize)
- [Helm](#helm)
- [Terraform](#terraform)
- [Pulumi](#pulumi)
- [CUE](#cue)
- [CDTK and CDK8s](#cdtk-and-cdk8s)
- [Cluster API](#cluster-api)
- [CloudFormation](#cloudformation)
- [Kpt](#kpt)
- [Jsonnet](#jsonnet)
- [Ksonnet](#ksonnet)

## Other programming languages

Since both terraform and kubernetes are written in Go, it felt natural to use Go and its dependency management system.
We have tried to use other languages, but we found that Go was the best fit for our use case.
We like the opinions that Go has and that let us focus on the problem at hand.

- Python dependencies are just hell to manage, and we have seen a lot of issues with them. Also, it is not as strongly type as Go.
- Javascript/Typescript have a steeper learning curve. We found that we would need to pull in a lot of dependencies to get the same result.
- Rust is a great language but the learning curve is too steep for us to teach it.
- Java/C# are great languages but their runtime dependency prevent us to easily share the code with others.

## Mimic

[Mimic](https://github.com/bwplotka/mimic) was a great inspiration for Lingon and we have a lot of respect for its author (Bartlomiej Plotka).
We even read his book `Efficient Go`, great book, highly recommended.

We haven't seen much activity on the project lately, and we wanted to go further.
Also, it is not possible to import kubernetes manifests to Go structs which makes the migration from YAML to Go
a manual tedious and time-consuming process.

## NAML

[NAML](https://github.com/krisnova/naml) was a great inspiration for Lingon and we have a lot of respect for its author (Kris Nova).

We haven't seen much activity on the project lately, and we wanted to go further.

We started to play with it as it is possible to import kubernetes manifests to Go.
We found out that we needed more control over the generated code and that we wanted to be able to easily move from YAML to Go and back
in order to make the migration from YAML to Go, and the deployment, as smooth as possible since we already have a lot (a lot!!!) of YAML manifests.

The support for CustomResourceDefinitions is not great without digging through the code, and we also have a lot of CRDs in our clusters.
Comparing with Lingon, just adding the types to the schema is enough to convert a manifest to Go.

## Kustomize

[Kustomize](https://kubectl.docs.kubernetes.io/) is great for simple use cases where we need to patch a few values in a YAML file.

The problem with Kustomize is the possible overlays on top of overlays on top of overlays. A problem that CUE forbids for that reason.
We've seen case where there were five layers of overrides, and it was hard to understand what the final value was.
Let alone finding an error when a manifest is updated somewhere in the chain.
When the resulting manifest is big enough, finding the right place to do a change feels like
changing the data on disk with a magnetic needle and steady hands.
Any update in the layer below what you control will break something in the upper layer,
the error message will be a bit cryptic.
It is not impossible to fix but not enjoyable as well.
That process has to be repeated for each update.

Kustomize alone is not enough to manage infrastructure at scale.

It is not a programming language, and it is not possible to write unit tests for it.
So the automation is limited to what it can do.

Finally, it does not solve the cloud resource management part.

## Helm

[Helm](https://helm.sh/) is great when the template is small and simple.

We have converted a few Helm charts to Lingon, and we found mistakes in the configuration which are obvious in Go but
extremely hard to find and debug in Helm. Even though Helm uses Go templates, and a lot can be done with it,
it is not a general purpose programming language.

For each setting, there is a value associated with it.
Unless the documentation is super clear on what the value does,
it is hard to understand what the value is for.
Moreover, when there are too many values, the template becomes unreadable.
Have a look at a [deployment](https://github.com/prometheus-community/helm-charts/blob/main/charts/prometheus/templates/deploy.yaml)
and try to understand what it does.

When in Go, the code is much more readable and maintainable.
You can create your own abstraction.

Templating only works for abstracting away some straightforward,
low level details that are common or should be enforced.  
The debugging experience is painful and requires a lot of knowledge and focus.

We haven't seen it work properly beyond simple string replacement or some naming conventions.

Front-end developers know that all too well and keep on finding new abstraction
all the time with every possible ways to generate HTML & CSS.

Finally, it is only for kubernetes and requires context switching with other tools.

## Terraform

We love [Terraform](https://www.terraform.io/) and we used it a lot.
Which is why we found that there are some shortcoming that we wanted to address.

- Once you start managing a large number of resources, Terraform becomes cumbersome to use.
- State splitting is something that needs to be set right from the beginning, and it is difficult to change it later on.
- `terragrunt` is a great tool to help with that, but it is not a silver bullet as it has its own patched HCL syntax.
- Notably, the lack of IDE support (autocompletion, type safety) and the long feedback loop.

Terraform would have been enough, but we wanted to go further and provide a better developer experience as well as
manage infrastructure at scale (+10k of resources).

Note that Lingon is using Terraform under the hood. Think of it as a Terraform wrapper in Go.

Managing the lifecycle of resources in Terraform requires us to execute a manual process that is impossible to safeguard with tools.
The process for an upgrade, for instance, consists of:

- copy-pasting the resources code
- make the required changes
- deploying the new resource
- migrating other resources dependent on the old one to the new one
- manually verifying that everything works as expected
- optionally falling back to the old resource in case of an issue
- once fully migrated, removing the old resources

This process is tedious and time-consuming which leads to human error and the necessity for checklists.
The goal of Lingon is to automate that process while reusing all the software engineering practices such as
linting, unit testing, e2e testing, continuous integration, deployment manifests, smoke testing after deploying and so on.

All of that is somewhat possible with Terraform, but it requires a lot of investment,
and we still have context switching between Terraform and writing code.

So why not write the code directly and avoid context switching altogether?

## Pulumi

We looked at [Pulumi](https://www.pulumi.com/), and it is a great tool.
However, we found that the developer experience could be improved.
A lot of the arguments are raw strings and do not prevent us from writing invalid configurations.
We would put Pulumi in the category of "more flexible Terraform".

Pulumi does what Lingon is doing but forces us to use unnecessary abstraction and
store relevant information in the context `*pulumi.Context`.
The context object puts us off a bit. It is a global variable that is used to pass values between resources.
Therefore, we would have to pass the key (just a string) around in order
to avoid typos and ensure a value is available for the next resource.

Even though Pulumi uses general purpose programming languages, the support for Go is not very impressive.

The open source and hosted part is unclear, it requires a token and an account.
We could not find how to avoid that from their website which doesn't build trust.

## CUE

We love [CUE](https://cuelang.org/), we've been trying it since 2019, we even [submitted it for review to TGIK](https://github.com/vmware-archive/tgik/issues/211).
Unfortunately, as much as we wanted to make it work, we felt the pain points were too big to ignore:

- It is not possible to import kubernetes manifests to Go structs (yet).
- The APIs are not stable (yet).
- The documentation is severely lacking as the APIs are still in flux.
- The code examples are hard to understand and do not provide a lot of context.
- It is yet another language to learn.
- The error messages are cryptic and hard to understand.
- The debugging experience is not great.
- The IDE support is not great.
- The community is small.
- The support for Terraform is not there (yet).

Also, on the CUE website, there is a great [comparison with other configuration languages](https://cuelang.org/docs/usecases/configuration/#inheritance-based-configuration-languages).

## CDTK and CDK8s

With [CDKTF](https://developer.hashicorp.com/terraform/cdktf) and [CDK8s](https://cdk8s.io/), the developer experience for Go is not great and
introduces [jsii](https://github.com/aws/jsii) which was a big turn-off for us as it needs NodeJS.

Additionally, and much like Pulumi and AWS CDK, infrastructure is not defined declaratively,
which is something we can achieve using Go structs.

## Cluster API

[Cluster API](https://cluster-api.sigs.k8s.io/) is a great tool to manage infrastructure and kubernetes clusters, but it is not a general purpose programming language.
So the automation is limited to what it can do.

It requires a running kubernetes cluster to manage other kubernetes clusters and admin level access to everything.

The debugging experience is really hard as the logs are in the kubernetes cluster,
and it requires a lot more knowledge and tooling.

Finally, it does not solve the kubernetes manifest management part.

## CloudFormation

[CloudFormation](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-whatis-concepts.html) is cloud specific and not portable. It is not a general purpose programming language.

## Kpt

[Kpt](https://github.com/GoogleContainerTools/kpt) requires a kubernetes cluster.
The installation instructions alone have many of steps, and it is not clear why something is required and what is optional.

The support for Go is great, and it comes with a lot of built-in functions.
However, it comes with its own abstraction layer, which needs to be studied.
We saw that the cost of getting started is pretty high. We wanted to avoid that as it would be a blocker for adoption.

It uses a declarative approach to manage infrastructure at scale.
Declarative language are fine for simple use cases, and they do a lot of automation under the hood.
It is that automation that makes the debugging experience hard.

## Jsonnet

[Jsonnet](https://jsonnet.org/) is a reat language, we used it extensively at a previous company.

In short: tooling, small community, developer experience, debugging experience, IDE support, documentation, error messages, etc.

## Ksonnet

[Ksonnet](https://github.com/ksonnet/ksonnet) is unmaintained and the project is archived on GitHub.
