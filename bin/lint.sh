#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

if stylist check; then
    echo "[lint] Pass ✅"
else
    echo "[lint] Fail 🔴"
    exit 1
fi
