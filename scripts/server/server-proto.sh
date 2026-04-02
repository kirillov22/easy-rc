#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

export PATH="$(go env GOPATH)/bin:$PATH"

cd "$REPO_ROOT/server"

mkdir -p generated/proto-messages

protoc \
  -I="$REPO_ROOT" \
  --go_out="$REPO_ROOT" \
  common/messages.proto

echo "Proto generated: server/generated/proto-messages/messages.pb.go"
