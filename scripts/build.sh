#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

"$REPO_ROOT/scripts/proto.sh"
"$REPO_ROOT/scripts/client/client-full-build.sh"
"$REPO_ROOT/scripts/server/server-full-build.sh"

echo "Full build complete."
