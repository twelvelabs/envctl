#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

app="envctl"
go build -o "build/$app" .

rm -rf build/manpages
mkdir -p build/manpages

# Generate Cobra man pages.
# See: https://github.com/spf13/cobra/blob/main/doc/man_docs.md
"build/$app" man | gzip -c -9 >build/manpages/envctl.1.gz
