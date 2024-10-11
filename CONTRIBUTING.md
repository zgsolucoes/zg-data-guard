# Development, Testing and Contributing

1. Make sure you have a running Docker daemon
   (Install for [MacOS](https://docs.docker.com/docker-for-mac/))
1. Use a version of Go that supports [modules](https://golang.org/cmd/go/#hdr-Modules__module_versions__and_more) (e.g. Go 1.22+)
1. Fork this repo and `git clone` somewhere to `$GOPATH/src/<your-domain>/zg-data-guard`
    * Ensure that [Go modules are enabled](https://golang.org/cmd/go/#hdr-Preliminary_module_support) (e.g. your repo path or the `GO111MODULE` environment variable are set correctly)
1. Install [golangci-lint](https://github.com/golangci/golangci-lint#install)
1. Run the linter: `make lint`
1. Confirm tests are working: `make test-verbose`
1. Create a new branch for your feature or bugfix
1. Write awesome code ...
1. Add tests for your change. Only refactoring and documentation changes require no new tests. If you are adding functionality or fixing a bug, we need tests!
1. Run the linter again: `make lint`
1. Add, commit and push your changes.
1. Submit a pull request.

Some more helpful commands:
* `make install` installs the app dependencies.
* `make test-verbose` runs tests with verbose output.
* `make html-coverage` which opens a shiny test coverage overview.
* `make build` builds the app in directory `dist/`.
* `make clean` removes directories `dist/` and `testdata/reports/`. Also removes `coverage.out`, `report-lint.html`
* `make lint` runs the linter.
* `make docs` generates the new Swagger documentation.
* `make run-with-docs` generates the new Swagger documentation and runs the app.

Check [README.md](README.md) for more information.
