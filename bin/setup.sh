#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

# Bootstrap dependencies.
dependencies=(
    "brew"
    "go"
    "npm"
)
for dependency in "${dependencies[@]}"; do
    if ! command -v "${dependency}" >/dev/null 2>&1; then
        echo "Unable to find dependency: ${dependency}."
        exit 1
    fi
done

# cspell: disable

brew bundle install

if ! command -v go-enum >/dev/null 2>&1; then
    go install github.com/abice/go-enum@latest
fi

if ! command -v gocovsh >/dev/null 2>&1; then
    go install github.com/orlangure/gocovsh@latest
fi

if ! command -v cspell >/dev/null 2>&1; then
    npm install -g cspell
fi

if ! command -v pin-github-action >/dev/null 2>&1; then
    npm install -g pin-github-action
fi
