# Contributing

Here's some ways to help:

* Pick an un-checked item from [Status](https://github.com/osteele/liquid#status). Let me know you want to work on it – I have ideas for some of these.
* Pick an item from [Other Differences](https://github.com/osteele/liquid#other-differences).
* Search the source for FIXME and TODO.
* Improve the [code coverage](https://coveralls.io/github/osteele/liquid?branch=master). Once it reaches 90%, we can submit a PR to [Awesome Go](https://github.com/avelino/awesome-go/)!

Review the [pull request template](https://github.com/osteele/liquid/blob/master/.github/PULL_REQUEST_TEMPLATE.md) before you get too far along on coding.

A note on lint: `nolint: gocyclo` has been used to disable cyclomatic complexity checks on generated functions, hand-written parsers, and some of the generic interpreter functions. IMO this check isn't appropriate for those classes of functions. This isn't a license to disable cyclomatic complexity checks or lint in general.

## Cookbook

### Set up your machine

Fork and clone the repo.

[Install go](https://golang.org/doc/install#install). On macOS running Homebrew, `brew install go` is easier than the linked instructions.

Install package dependencies and development tools:

* `make install-dev-tools`
* `go get -t ./...`

### Test and Lint

```bash
go test ./...
make lint
```

### Preview the Documentation

```bash
godoc -http=:6060
open http://localhost:6060/pkg/github.com/osteele/liquid/
```

### Work on the Expression Parser and Lexer

To work on the lexer, install Ragel. On macOS: `brew install ragel`.

Do this after editing `scanner.rl` or `expressions.y`:

```bash
go generate ./...
```

Test just the scanner:

```bash
cd expression
ragel -Z scanner.rl && go test
```