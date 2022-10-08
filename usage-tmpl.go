package main

var usageTmpl = `
{{define "option"}}
{{printf "  -%-11s %v" .option .info}}{{with .dv }} (default = {{.}}){{end}}
{{end}}
Usage: {{ .appName }} -[options] <args>

example: {{ .appName }} -answers "answers.json" -out-path "/tmp/new-app" "https://github.com/kohirens/tmpl-go-web"

Options:
`
