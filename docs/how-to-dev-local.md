# How to Develop Locally

## Pre-requisites

* Git latest version installed
* Docker latest version installed*
* VS Code latest version installed*

*These are optional, but these instructions do not tell you how to set up without them.

Development setup (with isolation**) can be achieved faster using Docker. Should you have Docker installed, then
you are halfway there.

1. To begin, make sure Docker is running.
2. Clone this repo somewhere on your PC and open in VS Code. a command prompt \
   to that directory.
3. VS Code should prompt you to install lots of extensions, if you don't have
   them already, most important is the `Remote Container` development extension,
   answer or click `yes` should it ask you to install it.
4. Once installed, VS Code should ask you to start the development environment,
   say yes. After a minute or so (depending on your connection speed), the setup

**isolation in terms of development, means you have configured development in a way that will not disrupt other
environments you have on the same PC/Laptop/machine.

## Requiring module code in a local directory

see: [Requiring module code in a local directory]

To tell Go commands to use the local copy of the module's code, use the `replace` directive:
in your `go.mod` file to replace the module path given in a `require` directive. See the `go.mod` reference for more
about directives.

```shell
go mod edit -replace=github.com/kohirens/stdlib@v0.0.0-unpublished=../stdlib
go get -d github.com/kohirens/stdlib@v0.0.0-unpublished
```

Where

* `github.com/kohirens/stdlib` - points to where the real module lives.
* `@v0.0.0-unpublished` - is a made-up placeholder version, see [Module version numbering].
* `=../stdlib` - points to where the module lives locally.
* `get -d` - instructs Go to download but not install the package.

---
[Requiring module code in a local directory]: https://golang.org/doc/modules/managing-dependencies#local_directory
[Module version numbering]: https://golang.org/doc/modules/version-numbers
