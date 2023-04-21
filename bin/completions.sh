#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

rm -rf build/completions
mkdir -p build/completions

# Generate Cobra shell completion scripts.
# See: https://github.com/spf13/cobra/blob/main/shell_completions.md
for sh in bash zsh fish; do
    go run . completion "$sh" >"build/completions/envctl.$sh"
done
