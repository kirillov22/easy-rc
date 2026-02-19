#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

cd "$REPO_ROOT/client"

echo "Installing client dependencies..."
npm install

echo "Generating protobuf TypeScript stubs..."
npm run proto

echo "Building client..."
npm run build

echo "Copying dist to server/static..."
rm -rf "$REPO_ROOT/server/static"
cp -r dist "$REPO_ROOT/server/static"

echo "Client build complete: server/static/"
