#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

# Ensure homebrew
if command -v "brew" >/dev/null 2>&1; then
    echo "Found dependency: brew."
else
    echo "Setup requires homebrew."
    echo "Please follow the instructions at https://brew.sh"
    exit 1
fi

# Ensure jq
if command -v "jq" >/dev/null 2>&1; then
    echo "Found dependency: jq."
else
    echo "Setup requires jq."
    echo "Please run: brew install jq"
    exit 1
fi

# Install dependencies
while read -r json; do
    command=$(echo "$json" | jq -r '.command')
    install=$(echo "$json" | jq -r '.install')

    if command -v "${command}" >/dev/null 2>&1; then
        echo "Found dependency: ${command}."
    else
        echo "Installing dependency: ${command}..."
        set -o xtrace
        $install
        set +o xtrace
    fi
done < <(jq --compact-output '.[]' dependencies.json)

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
