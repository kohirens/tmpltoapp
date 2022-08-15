# How to Develop Locally

## Pre-requisites

* Git latest version installed
* Docker latest version installed*
* VS Code latest version installed* (optional)

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

### Testing

This tool makes use of Git commands which requires a repository. For
testing we don't want to use a real repository. However, we do wish to test
and validate the commands are performing as expected. One
approach is to make Git bundles (archives) and un-bundle them during test runs.
It's one cool way to make test fixture repos that work very well in test cases.

**How to make a fixture repository for a test**

1. Make a new folder in `testdata`, you can follow the existing naming
   convention of `repo-xx`, where `xx` is a number with a leading 0.
2. cd into that new directory and run `git init` to initialize it.
3. Now just add files and commit them in this directory. Be careful. Make Sure
   you do **NOT** commit any of these test repository files to the main project.
4. Once you get the test repo to a state that you want, it's time to bundle it
   up using the Git bundle command:
   ```
   git bundle create <bundle-filename> <branch> --tags
   ```
   NOTE: For all branches and tags
   ```
   git bundle create <bundle-filename> --branches --tags
   ```
5. So from inside the test repository directory run the command, for
   example
   ```
   git bundle create ../repo-01.bundle --branches --tags
   ```
   In this example we save the bundle in the `testdata` directory. Please be
   sure to save yours there as well by using `../` before the name of the bundle.
6. Now go back to the root of the main project and be sure to commit the
   `*.bundle` file to the main project.
7. There is a function that you can use in the test to un-bundle this file
   during the test run, for example:
   
   ```go
   tmpRepo := setupARepository("repo-01")
   ```
   
   tempRepo will point to the path where the repo was extracted. Also know that
   this function will append ".bundle" to the first parameter to find the
   actual bundle.

---
[Requiring module code in a local directory]: https://golang.org/doc/modules/managing-dependencies#local_directory
[Module version numbering]: https://golang.org/doc/modules/version-numbers
