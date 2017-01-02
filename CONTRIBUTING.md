# Contributing

We love pull requests from everyone! We want to keep it as easy as possible to
contribute changes. There are a few guidelines that we need contributors to
follow so that we can keep on top of things.

## Getting Started

* Make sure you have a [GitHub account](https://github.com/signup/free).
* Submit a GitHub issue, assuming one does not already exist.
  * Clearly describe the issue including steps to reproduce when it is a bug.
  * Make sure you fill in the earliest version that you know has the issue.
* Fork the repository on GitHub.

## Making Changes

* Create a topic branch from the master branch.
* Make and commit your changes in the topic branch.

Commited code must pass:

* [gofmt](https://golang.org/cmd/gofmt)
* [go test](https://golang.org/cmd/go/#hdr-Test_packages) use make to run all non vendor tests

```
$ make
```

## Dependencies

* [Godep](https://github.com/tools/godep) is used for managing dependencies
* [go install](https://golang.org/cmd/go/#hdr-Compile_and_install_packages_and_dependencies) compiles dependencies to speed up ttests

## Submitting Changes

* Push your changes to a topic branch in your fork of the repository.
* Submit a pull request to the repository in the microscaling organization.
* Update your GitHub issue with a link to the pull request.
* At this point you're waiting on us. We like to at least comment on pull requests
within three business days (and, typically, one business day). We may suggest
some changes or improvements or alternatives.
* After feedback has been given we expect responses within two weeks. After two
  weeks we may close the pull request if it isn't showing any activity.
