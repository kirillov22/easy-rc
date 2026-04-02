#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

"$REPO_ROOT/scripts/server/server-proto.sh"
"$REPO_ROOT/scripts/server/server-compile.sh"
"$REPO_ROOT/scripts/server/server-lint.sh"
"$REPO_ROOT/scripts/server/server-test.sh"

echo "Server build complete."
