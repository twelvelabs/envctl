#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

# Ensure $USER owns /usr/local/{bin,share}.
# Allows for running `make install` w/out sudo interruptions.
if [[ ! -w /usr/local/bin ]]; then
    echo "Changing owner of /usr/local/bin to $USER."
    sudo chown -R "${USER}" /usr/local/bin
fi
if [[ ! -w /usr/local/share ]]; then
    echo "Changing owner of /usr/local/share to $USER."
    sudo chown -R "${USER}" /usr/local/share
fi

# Ensure local repo.
if ! git rev-parse --is-inside-work-tree &>/dev/null; then
    if gum confirm "Create local git repo?"; then
        echo "Creating local repo."
        git init
        git add .
        git commit -m "feat: initial commit"
    fi
fi

# Ensure local repo hooks.
if [[ -d .git ]]; then
    echo "Updating .git/hooks."
    mkdir -p .git/hooks
    rm -f .git/hooks/*.sample
    cp -f bin/githooks/* .git/hooks/
    chmod +x .git/hooks/*
fi

# Ensure remote repo.
if ! gh repo view --json url &>/dev/null; then
    if gum confirm "Create remote git repo?"; then
        echo "Creating remote repo."
        gh repo create
        sleep 1
        echo "Remote repo created: $(gh repo view --json url --jq .url)"

        if [[ "${GH_PAT:-}" == "" ]]; then
            echo "Enter a GitHub PAT with the 'repo' scope:"
            echo "  - This is required to publish to the homebrew tap repo in CI."
            echo "  - To create a new one, go to https://github.com/settings/tokens/new"
            GH_PAT=$(gum input --password)
        fi
        echo "Setting GH_PAT repo secret."
        gh secret set GH_PAT --body "$GH_PAT"

        echo "Setting 'remotes/origin/HEAD'."
        git remote set-head origin --auto
    fi
fi
