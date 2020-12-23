# go-gitter

Start a new project using a repos' directory structure as a manifest/blueprint.

## TOC

* [Description](#description)
* [Get Started](#get-started)

## Description

Copies one directory structure to another then process any files in the
directoreis as Go templates.

## Get Started

### Making a Template

Directories/repos are called templates, because they serve as a template for a
new project.

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

1. Run this application with 2-3 parameters:
   1. the first parameter is a path to a template, it can also be an HTTP URL.
   2. the second parameter is the path to where you want to place the project.
   3. you can pass in the local path to an asnswer file with the flags
      `-a, --answers` the path to an optional YAML file with answers for common
      placeholders, for example your name for an
      `{{ .author }}` placeholder.

### Notes About Template Processing

* You will be asked for values for any unique placeholders found in the template
  files, and that do not have an answer in the optional YAML "answer" file.
* Empty directories will be placed without the "empty.dir" file.

--

[Golang text/template]: https://golang.org/pkg/text/template/
