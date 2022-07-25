# How To Build A Template JSON Manifest

This JSON will contain a list of all placeholders in your template.

At minimum the `template.json` needs to contain
1. version with a  value of `0.1.0`
2. A `placeholders` property with at least 1 template variable name

for example:
```JSON
{
    "version": "0.1.0",
    "placeholders": {
        "appName": "a name for the application",
        "repoName": "a repository name for the application",
        "repoOrg": "your GitHub organization name"
    }
}
```

Notice the string values for the `placeholder` for each key equest to a question
to as for that variables value. This is because they are used to ask for the
value when filing out the template from the CLI.

## References

* [JSON Schema](https://json-schema.org/learn/getting-started-step-by-step#intro)
