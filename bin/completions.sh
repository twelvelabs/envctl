#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

app="envctl"
go build -o "build/$app" .

rm -rf build/completions
mkdir -p build/completions

# Generate Cobra shell completion scripts.
# See: https://github.com/spf13/cobra/blob/main/shell_completions.md
for sh in bash zsh fish; do
    "build/$app" completion "$sh" >"build/completions/$app.$sh"
done
