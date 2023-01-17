# Botkube Plugins

This repository shows the example Botkube [source](https://docs.botkube.io/architecture/#source) and [executor](https://docs.botkube.io/architecture/#executor) plugins.

It is intended as a template repository to start developing Botkube plugins in Go. Repository contains:

- The [`echo`](cmd/echo/main.go) executor that sends back the command that was specified,
- The [`ticker`](cmd/ticker/main.go) source that emits an event at a specified interval,

To learn more, see the [tutorial on how to use this template repository](https://docs.botkube.io/plugin/template.md).

## Requirements

- [Go](https://golang.org/doc/install) >= 1.18
- [GoReleaser](https://goreleaser.com/) >= 1.13
- [`golangci-lint`](https://golangci-lint.run/) >= 1.50

## Development

1. Clone the repository.
2. Follow the [local testing guide](https://docs.botkube.io/plugin/local-testing).
