#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

cd "$REPO_ROOT/client"

echo "Generating protobuf TypeScript stubs..."
npm run proto
