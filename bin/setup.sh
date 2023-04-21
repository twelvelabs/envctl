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

export HOMEBREW_NO_INSTALL_CLEANUP=1

# Install dependencies
count=$(jq '. | length' dependencies.json)
for ((i = 0; i < count; i++)); do
    command=$(jq -r '.['$i'].command' dependencies.json)
    install=$(jq -r '.['$i'].install' dependencies.json)

    if command -v "${command}" >/dev/null 2>&1; then
        echo "Found dependency: ${command}."
    else
        echo "Installing dependency: ${command}..."
        $install
    fi
done

# Bypass repo setup steps when in CI
if [[ "${CI:-}" == "true" ]]; then
    exit 0
fi

# Ensure $USER owns /usr/local/{bin,share}.
# Allows for running `make install` w/out sudo interruptions.
if [[ ! -w /usr/local/bin ]]; then
    set -o xtrace
    sudo chown -R "${USER}" /usr/local/bin
    set +o xtrace
fi
if [[ ! -w /usr/local/share ]]; then
    set -o xtrace
    sudo chown -R "${USER}" /usr/local/share
    set +o xtrace
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
        if [[ "${GH_PAT:-}" != "" ]]; then
            gh secret set GH_PAT --body "$GH_PAT"
        fi
        git remote set-head origin --auto
        set +o xtrace

        echo ""
        echo "Remote repo created: $(gh repo view --json url --jq .url)"
        echo ""
    fi
fi
