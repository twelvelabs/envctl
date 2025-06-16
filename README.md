# envctl

[![build](https://github.com/twelvelabs/envctl/actions/workflows/build.yml/badge.svg)](https://github.com/twelvelabs/envctl/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/twelvelabs/envctl/branch/main/graph/badge.svg)](https://codecov.io/gh/twelvelabs/envctl)

Manage project environment variables with ease. ✨

## Installation

Choose one of the following:

- Download and install the latest
  [release](https://github.com/twelvelabs/envctl/releases/latest) :octocat:

- Install with [Homebrew](https://brew.sh/) 🍺

  ```bash
  brew install twelvelabs/tap/envctl
  ```

- Install from source 💻

  ```bash
  go install github.com/twelvelabs/envctl@latest
  ```

## Usage

First, initialize your project:

```bash
envctl init
```

Which generates an `.envctl.yaml` file:

```yaml
environments:
  - name: local
    vars:
      EXAMPLE: "hello local"

  - name: prod
    vars:
      EXAMPLE: "hello prod"
```

Then exec a command in one of the named environments:

```shell
envctl exec local -- sh -c 'echo $EXAMPLE'
# => hello local

envctl exec prod -- sh -c 'echo $EXAMPLE'
# => hello prod
```

## Development

```bash
# Ensures all required dependencies are installed
# and bootstraps the project for local development.
make setup

# Run tests.
make test

# Run the app.
make run

# Show help.
make
```
