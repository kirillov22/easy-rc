#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

echo "Copying dist to server/static..."
rm -rf "$REPO_ROOT/server/static"
cp -r "$REPO_ROOT/client/dist" "$REPO_ROOT/server/static"

echo "Client deploy complete: server/static/"
