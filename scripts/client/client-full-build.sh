#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

"$REPO_ROOT/scripts/client/client-install.sh"
"$REPO_ROOT/scripts/client/client-proto.sh"
"$REPO_ROOT/scripts/client/client-compile.sh"
"$REPO_ROOT/scripts/client/client-lint.sh"
"$REPO_ROOT/scripts/client/client-test.sh"

echo "Client build complete."
