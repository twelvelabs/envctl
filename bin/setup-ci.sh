#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

ensure-dependency "gitsign" "brew install --quiet sigstore/tap/gitsign"

git config --global commit.gpgsign true      # Sign all commits
git config --global tag.gpgsign true         # Sign all tags
git config --global gpg.x509.program gitsign # Use gitsign for signing
git config --global gpg.format x509          # gitsign expects x509 args

# Configure commit author as "GitHub Actions" (with correct avatar email).
# See https://github.com/orgs/community/discussions/26560
git config --global user.name "github-actions[bot]"
git config --global user.email "41898282+github-actions[bot]@users.noreply.github.com"
