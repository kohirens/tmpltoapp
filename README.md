# go-gitter

A tool to start a project using a directory as a template.

## TOC

* [Description](#description)
* [Get Started](#get-started)

## Description

This is a tool to start a new project using another directory (recursively) as a template. The most complex part is
adding placeholders in files, which use [Golang text/template] syntax.

## Get Started

1. Make a new directory.
   1. Add folders, and if a directory should be empty, then place a file named "empty.dir" in it.
   2. Add files, but give the file the extension you need, for example "README.md"
      1. Files can contain Golang template placeholder, so `README.md` can contain:
         ```gotemplate
         {{ .appName }}
         ```
      Note the tool will fill in placeholders at runtime.
2. Run this application with 2-3 parameters:
   1. the first parameter is a path to a template, it can also be an HTTP URL.
   2. the second parameter is the path to where you want to place the project.
   3. the path to an optional YAML file with answers for common placeholders, for example your name for an
      `{{ .author }}` placeholder.

### Notes About Template Processing

* You will be asked for values for any unique placeholders found in the template files, and that do not have an answer
  in the optional YAML "answer" file.
* Empty directories will be placed without the "empty.dir" file.

--

[Golang text/template]: https://golang.org/pkg/text/template/
