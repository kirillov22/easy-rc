#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

cd "$REPO_ROOT/server"
go build -tags debug -o server .

echo "Dev build complete: server/server"
