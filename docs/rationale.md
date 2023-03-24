# Rationale

Because we had no other choice.

## Problem definition

Building a platform to cater to a specific community of people with various skills and interests is challenging 
because we don't completely know what that community wants or needs. We can only make educated guesses and
then build a platform that is open to all and can be extended to cater to the needs of that community.

Let's assume we know **exactly** what the community wants and needs and we start building the platform. 
The number of tools and frameworks available is overwhelming. The CNCF has a [landscape](https://landscape.cncf.io/) of all the tools
and frameworks available. That landscape is a great resource to get an idea of what is available and what
is trending. The time it takes to evaluate the most popular projects would take as much time as building
the platform directly. So people just pick popular tools, read some blog posts about it, decide it is good
and start building. Popular tools feel like a safe choice and make it very easy to start running before you learn how to jog sustainably.

After the initial phase of building a scaffold of a platform, the pain starts.
Things like security, reliability, scalability, performance and observability
are time consuming (weeks, months) to get right. During that time, a [sunk cost bias](https://en.wikipedia.org/wiki/Sunk_cost)
kicks in and the investment in the chosen tools and frameworks is not questioned but rather defended.
However, in order to get the platform into a production-ready shape, more tools are needed.
The more tools that are added, the more complex the platform becomes. The more complex the platform becomes,
the more difficult it is to add features to it. Everything related to developer experience becomes 
a third class citizen and getting things to work becomes the only priority.

This has platform engineers jumping through hoops to get things to work, begging and pleading for help from
other gate-keeping teams which have their own schedule, which requires the engineers 
to find workarounds - mainly to write their own custom tools - and ultimately, the platform 
that promised productivity and flexibility for end-users (data scientists, data engineers, ML engineers, developers, ...)
becomes a bottleneck for the community.

The problem is not the tools and frameworks. The problem is the way we combine them to build a platform.
Each set of tools have their own way of doing things, their own security models, their own observability requirements 
and their own API surface. The problem is that we don't have a way to combine them in a way that is
_"automatable"_. Engineers have to manually glue them together and that is where the problem starts.

We found that defining resources in a declarative way, i.e. YAML, JSON, HCL, SQL, is great for simple individual purpose, 
such as deploying a single container, creating a database in the cloud, etc. But when it comes to building a platform, 
we need to be able to link and combine the behaviour of those units of resources in a way that fits our needs. 

> It is the same as playing an instrument instead of playing an orchestra.

Each tool and frameworks have their own way of doing things.

Amongst other things, we found out that engineers have no time to improve the platform. What are they busy with?
Mainly: support. As the glue between the tools and frameworks is not automated, end-users have to rely
on the platform documentation, often scarce and outdated, and support to get things to work.
Since the platform is complex, end-users don't want to spend weeks of their time to learn the ins and outs of the platform
to do their work. The platform team help onboarding and settings up the platform for the end users themselves,
instead of automating the process.

Technically, the problem by having to manipulate the descriptive languages (YAML, JSON, HCL, SQL) to configure tools 
using things like bash scripts, sometimes Python. How and where those scripts run is usually a hassle. 
Since those scripts are not interacting directly with the APIs of the tools and frameworks, 
but mainly through configuration files, errors are discovered at runtime, often reported by an end-user asking for support. 
Those scripts can also pose a security risk, as they are not audited and can sometimes be easily modified by end-users.


> Since those descriptive languages are mainly **data**, 
> we need a tool that can manipulate data in an automatic manner
> while respecting the format of the APIs from the tools.

Basically :

> We need to automate the glue between the tools and the infrastructure.

![kubernetes landscape](assets/kubernetes-landscape.png)

### References

* [noyaml website](https://noyaml.com/)
* [The yaml document from hell](https://ruudvanasseldonk.com/2023/01/11/the-yaml-document-from-hell) by [ruudvanasseldonk](https://ruudvanasseldonk.com)
* [Nightmares On Cloud Street 29/10/20 - Joe Beda 🎥](https://youtu.be/8PpgqEqkQWA)
* [YAML Considered Harmful - Philipp Krenn 🎥](https://youtu.be/WQurEEfSf8M)
* [You Broke Reddit: The Pi-Day Outage](https://www.reddit.com/r/RedditEng/comments/11xx5o0/you_broke_reddit_the_piday_outage/)
* [Collection of post-mortems](https://github.com/danluu/post-mortems#config-errors)
* [Kubernetes failure stories](https://k8s.af/)

## Solution

When it comes to manipulating data, **general programming languages** such as Python, Go, Java, Rust, C#, ...
and software engineering best practices have been battle tested and proven. 

Lingon has been built from our experience using a general programming language to configure 
our applications and infrastructure.

* **Reduce cognitive load**: Building a platform within a single context (i.e. Go) will reduce cognitive load 
by decreasing the number of tools and context switching in the process.
* **Type safety**: Detect misconfigurations in your text editor at **compile time** by using type-safe Go structs 
to exchange values across tool boundaries. 
This "shifts left" the majority of errors that occur to the earliest possible point in time.

* **Error handling**: Go's error handling enables propagating meaningful errors to the user.
This significantly reduces the effort in finding the root cause of errors and provides a better developer experience.

* **Limitless automation**: no need to write bash scripts to glue tools together. 
We have a general programming language at our disposal that enables us to automate and 
**test** the most critical component before they reach production.

Lingon is aimed at people managing cloud infrastructure who have suffered the pain of configuration languages 
and the complexity of gluing tools together with yet another tool.


> Lingon was created to manage platforms living in various environments at scale. 

Lingon is not a platform, it is a library meant to be consumed in a Go application that platform engineers write 
to manage their platforms. It is a tool to build and automate the creation and the management of platforms 
regardless of the target infrastructure and services.


## FAQ

### Why Go ?

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

* _Backwards compatibility_: Go is well-known for its backwards compatibility making it a great choice for long-lasting platforms reducing rework through breaking changes.


### Why not use just Terraform ?

We love terraform and we used it a lot. Which is why we found that there are some shortcoming that we wanted to address.
Once you start managing a large number of resources, Terraform becomes cumbersome to use.
State splitting is something that needs to be set right from the beginning, 
and it is difficult to change it later on.
Terragrunt is a great tool to help with that, but it is not a silver bullet as it has its own patched HCL syntax.

Notably, the lack of IDE support (autocompletion, type safety) and the long feedback loop.

Terraform would have been enough, but we wanted to go further and provide a better developer experience as well as 
manage infrastructure at scale (+10k of resources).

Note that Lingon is using Terraform under the hood. Think of it as a Terraform wrapper in Go. 

### Why not use just Pulumi ?

We looked at it, and it is a great tool. 
However, we found that the developer experience could be improved.
A lot of the arguments are raw strings and do not prevent us from writing invalid configurations.
We would put Pulumi in the category of "more flexible Terraform".

The context object puts us off a bit. It is a global variable that is used to pass values between resources.
Therefore, we would have to pass the key (just a string) around in order 
to avoid typos and ensure a value is available for the next resource.

Even though Pulumi uses general purpose programming languages, the support for Go is not very impressive.

### Why not use just CloudFormation ?

It is cloud specific and not portable. It is not a general purpose programming language.

### Why not use Helm ? 

Templating is great when the template is small and simple.

We have converted a few Helm charts to Lingon, and we found mistakes in the configuration which are obvious in Go but 
extremely hard to find in Helm. Even though Helm uses Go templates, and a lot can be done with it, 
it is not a general purpose programming language.

For each setting, there is a value associated with it. Unless the documentation is super clear on what the value does,
it is hard to understand what the value is for. Moreover, when there are too many values, the template becomes unreadable.

When in Go, the code is much more readable and maintainable. You can create your own abstraction.

### Why not use Kustomize ?

Kustomize is great for simple use cases where we need to patch a few values in a YAML file.

We've seen case where there were five layers of overrides, and it was hard to understand what the final value was.
Let alone finding an error when a manifest is updated somewhere in the chain. 

### Why not CDK for Terraform (cdktf)?

The developer experience for Go is not great and 
introduces [jsii](https://github.com/aws/jsii) which was a big turn-off for us.

Additionally, and much like Pulumi and AWS CDK, infrastructure is not defined declaratively, 
which is something we can achieve using Go structs.