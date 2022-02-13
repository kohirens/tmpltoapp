* [About](#about)
    * [Description](#description)
    * [Installation](#installation)
        * [Using Go](#using-go)
        * [Using Docker](#using-docker)
    * [Get Started](#get-started)
        * [Making a Template](#making-a-template)
        * [Using a Template](#using-a-template)

# About

Start a project from a template.

NOTICE: Name change pending... "start-from-tmpl" to make it obvious what this
tool does.


## Description

Use this tool to start a new project from a template. A template is a folder
with files, with each file processed as a Go template. Data for the template
can be passed to the template/context by passing the path to a JSON file,
where the keys will be used to fill in placeholders. Placeholders refer to the
values that need to be filled in when a file is process as a Go template.

Templates can be supplied as a local folder or URL to zip file. See an examples:

[web.go.tmpl](https://github.com/kohirens/tmpl-go-web)

## Installation

### Using Go

```
go get -u -v github.com/kohirens/tmpltoapp
```

### Using Docker

```
docker pull kohirens/tmpltoapp:latest
```

## Get Started

### Making a Template

Templates are just directories containing files, which can contain Go template
syntax that will be processed to fill in placeholders. The output will serve as
a new project.

1. Make a new directory.
2. Add folders, and if a directory should be empty, then place a file named
   "empty.dir" in it.
3. Add files, but give the file the extension you need, for example "README.md"
   1. Files can contain Golang template placeholder, so `README.md` can contain:
      ```gotemplate
      # {{ .appName }}

      ## Summary
      ```
   Note: the tool will fill in placeholders at runtime.
4. Commit the changes and push up to github.com

### Using a Template

1. Run this application with 3 parameters:
   1. a path to a template, it can also be an HTTP URL.
   2. a path to where you want to place the project (it should not exist).
   3. a path to an answer (JSON) file containing key/value pairs that will
      serve as variables. For example, an `{{ .author }}`
      placeholder would take a file that has:
      ```
      {
        "author": "your name"
      }
      ```

NOTE: There a long/short flags you can use should you need to place the values
      out of order, run the program with `-h` or `--help` for options.

### Notes About Template Processing

* You can provide a list of keys, in the `answer` file, that you want a user to
  supply when they run the template.
* Empty directories will be placed without the "empty.dir" file.

---

[Golang text/template]: https://golang.org/pkg/text/template/
