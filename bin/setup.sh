#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

# Non-installable dependencies.
dependencies=(
    "gh"
    "git"
    "go"
)
for dependency in "${dependencies[@]}"; do
    if ! command -v "${dependency}" >/dev/null 2>&1; then
        echo "Unable to find dependency: ${dependency}."
        exit 1
    fi
done

# cspell: disable

if ! command -v actionlint >/dev/null 2>&1; then
    go install github.com/rhysd/actionlint/cmd/actionlint@latest
fi

if ! command -v go-enum >/dev/null 2>&1; then
    go install github.com/abice/go-enum@latest
fi

if ! command -v gocovsh >/dev/null 2>&1; then
    go install github.com/orlangure/gocovsh@latest
fi

if ! command -v golangci-lint >/dev/null 2>&1; then
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

if ! command -v pin-github-action >/dev/null 2>&1; then
    npm install -g pin-github-action
fi
