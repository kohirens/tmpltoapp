* [About](#about)
    * [Description](#description)
    * [Installation](#installation)
        * [Using Go](#using-go)
        * [Using Docker](#using-docker)
    * [Get Started](#get-started)
        * [Making a Template](#making-a-template)
        * [Using a Template](#using-a-template)

# About

Start an app from a template.


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
go get github.com/kohirens/tmpltoapp
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
   ".empty" in it. It does need any text in it.
3. Add your files, the extension does not matter as it will be processed as a Go template. For example imagine "README.md" with the content.
   1. Files can contain a GoLang template placeholder, so `README.md` can contain:
      ```gotemplate
      # {{ .appName }}

      ## Summary
      ```
   Note: the placeholder `{{ .appName }}` will be replaced with the apps name at runtime.
4. Commit the changes and push up to your repo.

### Using a Template

Run this application with 3 parameters:
1. a path to a template, a URL to a zip or local folder.
2. a path to where you want to place the project (it should not exist).
3. a path to an answer (JSON) file containing key/value pairs that will
   serve as variables. The name is not imp For example, an `{{ .author }}`
   placeholder would take a file that has:
   ```
   {
      "author": "your name here",
      "appName": "awesomeAppName"
   }
   ```
NOTICE: Name change pending... "start-from-tmpl" to make it obvious what this
tool does.

4. make a `template.json` file which acts as a manifest, it needs the following:
   ```
   {
      "version": "0.2.0",
      "variables": {
        "appName": "Name of the app",
        "author": "Code author name"
      },
      "excludes": [
         ".gitignore"
      ]
   }
   ```
NOTE: There are command line flags should you need to place the arguments
      out of order. Run the program with `-h` or `--help` for options.

### Notes About Template Processing

* All variables are treated as strings.
* If any variables are in the `template.json` that are supplied by an answer JSON, then processing will halt and ask for them. 
* Empty directories will be placed without the ".empty" file.
* Files listed in the `excludes` list are output to the final app directory without template processing.
* Template are processed with the Go lib [Golang text/template].
---

[Golang text/template]: https://golang.org/pkg/text/template/
