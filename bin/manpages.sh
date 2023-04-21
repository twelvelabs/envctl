#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

rm -rf build/manpages
mkdir -p build/manpages

# Generate Cobra man pages.
# See: https://github.com/spf13/cobra/blob/main/doc/man_docs.md
go run . man | gzip -c -9 >build/manpages/envctl.1.gz
