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

brew bundle install --no-lock

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

# Ensure local repo
if ! git rev-parse --is-inside-work-tree &>/dev/null; then
    if gum confirm "Create local git repo?"; then
        set -o xtrace
        git init
        git add .
        git commit -m "feat: initial commit"
        set +o xtrace
    fi
fi

# Ensure remote repo
if ! gh repo view --json url &>/dev/null; then
    if gum confirm "Create remote git repo?"; then
        set -o xtrace
        gh repo create
        sleep 1
        git remote set-head origin --auto
        set +o xtrace

        echo ""
        echo "Remote repo created: $(gh repo view --json url --jq .url)"
        echo ""
    fi
fi
