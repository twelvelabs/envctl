#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

# Silence cleanup message on `brew install`.
export HOMEBREW_NO_INSTALL_CLEANUP=1

# Ensure homebrew.
if command -v "brew" >/dev/null 2>&1; then
    echo "Found dependency: brew."
else
    echo "Setup requires homebrew."
    echo "Please follow the instructions at https://brew.sh"
    exit 1
fi

# Ensure jq.
if command -v "jq" >/dev/null 2>&1; then
    echo "Found dependency: jq."
else
    echo "Setup requires jq."
    echo "Please run: brew install jq"
    exit 1
fi

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

if [[ "${CI:-}" != "true" ]]; then
    SCRIPT_DIR="$(dirname "${BASH_SOURCE[0]}")"
    "${SCRIPT_DIR}/setup-local.sh"
fi
