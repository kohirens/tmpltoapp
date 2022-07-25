# How To Build A Template JSON Manifest

Your template will need to container a `template.json` file which should list
all placeholders in your template. This is file is used at the time the template
is process with this tool. It is checked for all variables that must be
answered before process can occur. You can build one manually by using the
following format.

At minimum the `template.json` needs to contain

1. A `version` property with a  value of `0.1.0`
2. A `placeholders` object property with at least 1 template variable name
2. An optional `excludes` array property with at least 1 item to indicate a file or directory to skip processing and copy as-is.

for example:
```JSON
{
    "version": "1.0.0",
    "placeholders": {
        "appName": "a name for the application",
        "repoName": "a repository name for the application",
        "repoOrg": "your GitHub organization name"
    }
}
```

Notice the string values for each key in the `placeholders` property equates
to a question. This is because they can be used as prompts to
ask for the value when filing out the template from the CLI.

## References

* [JSON Schema](https://json-schema.org/learn/getting-started-step-by-step#intro)
