#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

ensure-dependency() {
    local dependency="${1}"
    local install_command="${2}"

    if command -v "${dependency}" >/dev/null 2>&1; then
        echo "Found dependency: ${dependency}."
    else
        echo "Installing dependency: ${dependency}..."
        $install_command
    fi
}

# Ensure homebrew.
ensure-dependency "brew" "echo 'Please follow the instructions at https://brew.sh' && exit 1"

# Ensure jq.
ensure-dependency "jq" "brew install --quiet jq"

# Ensure remaining dependencies.
count=$(jq '. | length' dependencies.json)
for ((i = 0; i < count; i++)); do
    command=$(jq -r '.['$i'].command' dependencies.json)
    install=$(jq -r '.['$i'].install' dependencies.json)
    ensure-dependency "${command}" "${install}"
done

if [[ "${CI:-}" != "true" ]]; then
    SCRIPT_DIR="$(dirname "${BASH_SOURCE[0]}")"
    "${SCRIPT_DIR}/setup-local.sh"
fi
