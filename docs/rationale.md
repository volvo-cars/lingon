# Rationale

Why we built this? Because we had no other choice, no tools supported a controlled and reproducible way to upgrade our infrastructure.

Note that most tools are fine, we just didn't want to manage YAML, HCL or simply put, text files anymore.
The declarative approach has worked well until it hasn't.
Above a certain threshold of complexity, code is necessary.
We didn't want another kubernetes cluster to manage all other kubernetes clusters.
We are software engineers and we wanted Go code, we wanted tests,
we wanted to abstract away when it makes sense and leave it bare when it makes sense.
Go is usually a second class citizen compared to TypeScript and Python
when it comes to infra-as-code tools such as Pulumi and CDK.
All the solutions we tried were too simple for our use case or too bloated.
Most of us worked with Terraform a lot and kubernetes a lot.
We were missing a way to combine the two and automate the glue with code.
Therefore, Lingon has evolved into what it is today.
We don't try to compete with other tools, we built it for us.
If it ends up being useful to someone, we are more than happy to help with that.

We didn't want to build yet another abstraction that tries to cover everyone's use cases (see <https://xkcd.com/927/>).
So we decided to build a library to facilitate the creation of **your** abstraction for your platform, regardless of the target infrastructure and services.

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

![kubernetes landscape](./assets/k8s.gif "kubernetes memes gif")

## Solution

When it comes to manipulating data, **general programming languages** such as Python, Go, Java, Rust, C#, ...
and software engineering best practices have been battle tested and proven.

Lingon has been built from our experience using a general programming language to configure
our applications and infrastructure in order to:

* **Reduce cognitive load** :Building a platform within a single context (i.e. Go) will reduce cognitive load
by decreasing the number of tools and context switching in the process.
It provides a better developer experience with out-of-the-box IDE support and a single language to learn with smooth learning curve.

* **Type safety**: Detect misconfigurations in your text editor at **compile time** by using type-safe Go structs
to exchange values across tool boundaries.
This "shifts left" the majority of errors that occur to the earliest possible point in time.

* **Error handling**: Go's error handling enables propagating meaningful errors to the user.
This significantly reduces the effort in finding the root cause of errors and provides a better developer experience.

* **Limitless automation**: no need to write bash scripts to glue tools together.
We have a general programming language at our disposal that enables us to automate and
**test** the most critical component before they reach production.
We are only limited by what a programming language can do.
We can reuse part of what we build in libraries without external tooling.
That is not possible with YAML as doesn't support "includes", therefore we need a tool for that.
Configuration languages are limited by the features they provide.
Gluing tools together with more tools and configuration to manage more tools and configuration is not a sustainable approach.
We do use a limited set of tools that we learn well and can extend, but we automate them and test them together using Go.

Note that we are in a particular situation where we need custom automation of the lifecycle of our cloud infrastructure.
Lingon is aimed at people managing cloud infrastructure who have suffered the pain of configuration languages
and the complexity of gluing tools together with yet another tool.

> Lingon was created to manage platforms living in various environments at scale.

Lingon is not a platform, it is a library meant to be consumed in a Go application that platform engineers write
to manage their platforms. It is a tool to build and automate the creation and the management of platforms
regardless of the target infrastructure and services.

## FAQ

See [FAQ](./faq.md)
