<a name="unreleased"></a>
## [Unreleased]


<a name="0.1.0"></a>
## 0.1.0 - 2022-02-18
### Added
- Excludes field to template.json Schema.
- File for collective errors.
- Method to read in answers from template.json questions.
- Function to read the template.json file required in a template.
- Detection of 7zip in the environment.
- TODO.md
- Parsing local directories as templates.
- Error checking for appPath when extracting flag value.
- Parse local directory.
- Documentation on loadAnswers.
- loadAnswers method.
- Methdo to detect template path type.
- Methdo to detect template path type.
- Flow Diagrams.
- Local dev instructions.
- Method to detect text file via file extensions.
- Vendor dependency config.
- Template processing.
- PathExist function to stdlib.
- Method to extract a template package.
- Code coverage visuals to vscode config.
- Default config.
- Initialize a config.
- Validation of URLs domains allowed to download.
- stdlib package.
- Installation instructions to the readme.
- Algorithm.
- Debugging with Delve and VS Code Remote.
- Lead config from file.
- WIP, Download template feature.
- More tools to the dev environment.
- Errors when missing required arguments.
- Parameter parsing.
- Docker configuration.

### Changed
- Renamed from Bootup to TmplToApp.
- Where extract places files.
- Error message array name to errs.
- Download zips to a unique name in the cache.
- Upgrade to Kohirens STDLib to version 0.1.1, go.sum updated.
- Upgraded Kohirens STDLIB version to 0.1.2.
- Dev env container name.
- Read values from template.json file.
- No longer throw an error when skipping files to run through the template engine.
- Update deve environment.
- Rebrand
- Refactoring Docker
- Allow user to set the file extensions to filter in the app config.
- Set app config property tmpl path when passing in CLI local path.
- Pull in default app config using unmarshal to struct type.
- Use one app config per run.
- Removed verbosity level from app config and add tmpl.
- Move directory where go-gitter stores files.
- Updated documentation.
- Refactor flags code.
- Error handling on config.settings function.
- Refactor function to initialize a config.
- Replace hard coded value with a constant.
- Updated local development docs.
- Editor config to be compatible with GoLang code fmt.
- Made function private that are not part of the public API.
- Refactoring for better organization.
- Refactor converge template package into main.
- Upgraded to go 1.16.
- Dev default environment variables.
- Extract the stdlib out to its own repo.
- Updated the development environment.
- Turn on Go tools extension checking to stay up to date.
- Map SSH keys in dev environment.
- Rename download to template.
- Pre-install vscode-server.
- Constant TEST_TMP to be the same in all test files.
- Rename getArgs to parseArgs.
- To use a flagset in place of global flag space.
- To Go module.
- Verbosity flag name.
- editorcofig line endings to Linux.
- Replaced rsyslog with tini and tail.
- Settle on a container running in dev.

### Fixed
- Unit test.
- Erros with reading template.json file.
- VS Code Remote container setup.
- Downloaded zips and ignore empty files.
- Malformed JSON in the default config.
- Test output going to wrong directory.
- Setting exit code from main when no errors.
- Get template.parseDir working.
- Return correct error from template.parseDir.
- Test template repository name for test arguement.
- Unit test in stdlib.
- Test in stdlib.
- Broken test for downloading templates.
- Download to save archive based on URL name.
- Testing suite.
- logic bugs.
- Missed file required to run settings test.
- Settings in VS Code devcontainer config.
- Added missing constants.
- passing flags to main and running tests.

### Removed
- Unused CopyDir function.

### BREAKING CHANGE

* Renamed `tplPath` flag to `tmplPath`.
* Changed `-v` no longer set verbosity level, but is a short for
  version.


[Unreleased]: https://github.com/kohirens/tmpltoapp.git/compare/0.1.0...HEAD
