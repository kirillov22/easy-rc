#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

mkdir -p "$REPO_ROOT/bin"
cd "$REPO_ROOT/server"
go build -tags debug -o "$REPO_ROOT/bin/server" .

echo "Dev build complete: bin/server"
