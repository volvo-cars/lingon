# FAQ

- [FAQ](#faq)
  - [Why Go?](#why-go)
  - [Why did you build this project?](#why-did-you-build-this-project)
  - [Is the API stable?](#is-the-api-stable)
  - [How did you develop this project?](#how-did-you-develop-this-project)
    - [Proverbs](#proverbs)
  - [Why shelling out to kubectl instead of using the go-client?](#why-shelling-out-to-kubectl-instead-of-using-the-go-client)
  - [Why do you depend on k8s.io/client-go ?](#why-do-you-depend-on-k8sioclient-go-)

## Why Go?

See [Go](./go.md).

## Why did you build this project?

See [Rationale](./rationale.md).

## Is the API stable?

Consider this project as a beta release for now as we built this for us.
We haven't consider every use case and we might have missed some.
Also, there might be some bugs to iron out.

Once the API is stabilized, we will anounce it and will promise backward compatibility.
No API will be removed or changed after that. We will only add new APIs.
If we found an API that is not used, we will just deprecate it.
See [godoc](https://go.dev/blog/godoc), at the bottom of the page for
how Go handles deprecation. We want to adhere to the same rules as Go.

## How did you develop this project?

We experimented a lot and researched a lot of tools (see [comparison.md](./comparison.md))

We wanted to cater to the needs of platform teams in a big company.
We wanted to make sure we were not reinventing the wheel.

During our research, we even set a list of proverbs to adhere to.
It is looser than strict rules, but it is a good guideline.

> They are heavily inspired by the [Go proverbs](https://go-proverbs.github.io/).

### Proverbs

> Wait until it hurts a little

Sometimes doing thing that don't scale is the best way to get things done.
Think about what adds value and don't do things for the sake of doing things.

> If it hurts, do it more often

If something is painful, it's probably worth doing more often.
It's a good way to get better at it and will help automation later on.
Doing is a form of testing and we should test often.

> Set yourself up for automation

Unless it is something we do once every 3 months, don't automate, document it instead.

This is why we use Go, it's easy to read and write and it's easy to automate.
YAML is fine but it is prone to typos and it is harder to write tests for it.

> Do it manually before automating it

Too many times, the automation part is the cause of the problem.
If the automation fails, how to debug and fix it ? The answer is: manually.
Let's make sure the manual process is working before automating it.

> Don't split it up until it is too big

It is very hard and error prone to come up with the perfect architecture from the start.
It is much easier to split it up later, when it is too big.
Doing is learning about the architecture. Start small and iterate.

> Clear is better than clever

We read code more than we write it. Optimize for readability and maintainability.

> A little copying is better than a little dependency

We don't want to depend on a lot of external libraries as we want to build strong foundations.

> The bigger the interface the weaker the abstraction

We want to keep the API surface small and simple. Do one thing and one thing well.

> Documentation is for users

Users shouldn't need to know the internals of the system in order to use it.

> A user cannot do anything about errors in the system

See 4XX http status code vs 5XX http status code.
An error in the system is something that should be monitored, alerted and investigated by the team.
Whereas an error seen by the user is something that should help the user to fix the problem.

> Avoid polishing turds

Try to prove the value of experiments or research as quickly as possible, and avoid polishing them before you do.
Try to spend time working on things that will deliver value rather than perfecting things that
you know will soon change or be replaced.

## Why shelling out to kubectl instead of using the go-client?

Many projects are actually "shelling out" (use `exec.Command` in Go) in order to run a CLI,
instead of calling the APIs directly from Go.
That pattern occurs when it is difficult to work with APIs and the lack thereof.
While investigating, KinD is using this pattern to execute Docker.
We tried to make it work with the go-client but it was a lot of work and we decided to use the CLI instead.

## Why do you depend on k8s.io/client-go ?

We use the `"k8s.io/client-go/kubernetes/scheme"` package to register the CRD.
As many projects do have CRDs, we decided (relunctantly) to depend on the `client-go` package for now.
