# How to contribute

We'd love your help!

This project is [MIT Licensed](LICENSE) and accepts contributions via GitHub pull requests. This document outlines conventions on development workflow, commit message formatting, contact points and other resources to make it easier to get your contributions accepted.

### Prerequisite

- [Go](https://go.dev), at least 1.18
- `make`

#### Develop Botkube plugins

The example plugins in this repository were developed based on the official guides for the [executor](https://docs.botkube.io/plugin/custom-executor) and [source](https://docs.botkube.io/plugin/custom-source) plugins.

To test them locally, follow the [local testing](https://docs.botkube.io/plugin/local-testing) guide.

If you update this template repository, make sure that changes are reflected in the official [Template GitHub repository](https://docs.botkube.io/plugin/template-gh-repo) documentation.

## Making A Change

- Before making any significant changes, please [open an issue](https://github.com/kubeshop/botkube-plugins-template/issues). Discussing your proposed changes ahead of time will make the contribution process smooth for everyone.

- Once we've discussed your changes, and you've got your code ready, make sure that the build steps mentioned above pass. Open your pull request against the [`main`](https://github.com/kubeshop/botkube-plugins-template/tree/main) branch.

  To learn how to do it, follow the **Contribute** section in the [Git workflow guide](https://github.com/kubeshop/botkube/tree/main/git-workflow.md).

- To avoid build failures in CI, install [`golangci-lint`](https://golangci-lint.run/usage/install/) and run:

  ```sh
  # From project root directory
  make fix-lint-issues
  ```
  This will run the `golangci-lint` tool to lint the Go code.

### Create a Pull Request

- Make sure your pull request has [good commit messages](https://chris.beams.io/posts/git-commit/):
  - Separate subject from body with a blank line
  - Limit the subject line to 50 characters
  - Capitalize the subject line
  - Do not end the subject line with a period
  - Use the imperative mood in the subject line
  - Wrap the body at 72 characters
  - Use the body to explain _what_ and _why_ instead of _how_

- Try to squash unimportant commits and rebase your changes on to the `main` branch, this will make sure we have clean log of changes.

## Support Channels

Join the Botkube-related discussion on Slack!

Create your Slack account on [Botkube](https://join.botkube.io) workspace.

To report bug or feature, use [GitHub issues](https://github.com/kubeshop/botkube-plugins-template/issues/new/choose).
