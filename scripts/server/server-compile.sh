#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

cd "$REPO_ROOT/server"

echo "Building production binary..."
mkdir -p "$REPO_ROOT/server/bin"
go build -o "$REPO_ROOT/server/bin/easy-rc" .

echo "Production build complete: server/bin/easy-rc"
