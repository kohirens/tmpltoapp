package main

var usageTmpl = `
{{define "option"}}
{{printf "  -%-11s %v" .option .info}}{{with .dv }} (default = {{.}}){{end}}
{{end}}
Usage: {{ .appName }} -[options] <args>

example: {{ .appName }} -answers "answers.json" -out-path "/tmp/new-app" "https://github.com/kohirens/tmpl-go-web"

Options:
`

var usageManifest = `
Generate a template.json in the {{.appName}} schema format containing all the specified templates placeholders.
This is meant for template designers to reduce error in adding placeholders to your file manually.
Placeholders are also known ad "Actions"--data evaluations-- in Go.

Usage: {{.appName}} <path>

example: {{.appName}} ./
`
