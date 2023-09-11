# Botkube Plugins

This repository shows the example Botkube [source](https://docs.botkube.io/architecture/#source) and [executor](https://docs.botkube.io/architecture/#executor) plugins.

It is intended as a template repository to start developing Botkube plugins in Go. Repository contains:

- The [`echo`](cmd/echo/main.go) executor that sends back the command that was specified.
- The [`ticker`](cmd/ticker/main.go) source that emits an event each time the configured time duration elapses.
- The [`forwarder`](cmd/forwarder/main.go) source that echos the message sent as an incoming webhook request.
- The [`msg`](cmd/msg/main.go) executor that showcases the interactive message capabilities.
- The release [GitHub Action](https://github.com/features/actions) jobs:
	- that creates [GitHub release](.github/workflows/release.yml) with plugin binaries and index file each time a new tag is pushed.
		- See: https://github.com/kubeshop/botkube-plugins-template/releases/latest
	- that updates [GitHub Pages](.github/workflows/pages-release.yml) with plugin binaries and index file each time a new tag is pushed.
		- See: https://kubeshop.github.io/botkube-plugins-template/

To learn more, see the [tutorial on how to use this template repository](https://docs.botkube.io/plugin/quick-start).

## Requirements

- [Go](https://golang.org/doc/install) >= 1.18
- [GoReleaser](https://goreleaser.com/) >= 1.13
- [`golangci-lint`](https://golangci-lint.run/) >= 1.50

## Development

1. Clone the repository.
2. Follow the [local testing guide](https://docs.botkube.io/plugin/local-testing).
