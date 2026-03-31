#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

mkdir -p "$REPO_ROOT/server/bin"
cd "$REPO_ROOT/server"
go build -tags debug -o "$REPO_ROOT/server/bin" .

echo "Dev build complete: server/bin/server"
