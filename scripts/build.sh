#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

echo "Generating protobuf..."
"$REPO_ROOT/scripts/proto.sh"

echo "Building client..."
"$REPO_ROOT/scripts/client.sh"

cd "$REPO_ROOT/server"

echo "Running go vet..."
go vet ./...

echo "Running tests..."
go test ./...

echo "Building production binary..."
mkdir -p "$REPO_ROOT/server/bin"
go build -o "$REPO_ROOT/server/bin" .

echo "Production build complete: ./server/bin/server"
