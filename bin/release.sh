#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

SCRIPT_DIR="$(dirname "${BASH_SOURCE[0]}")"
source "${SCRIPT_DIR}/version.sh"

if [[ "$CURRENT_VERSION" != "$NEXT_VERSION" ]]; then
    git tag \
        --sign "$NEXT_VERSION" \
        --message "$NEXT_VERSION"
    git push origin --tags
fi
