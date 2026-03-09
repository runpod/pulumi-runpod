#!/usr/bin/env bash
# Run go mod tidy across all Go modules in the repo.
# Run this after any dependency change to keep all modules in sync.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

MODULES=(
  "$REPO_ROOT"
  "$REPO_ROOT/provider"
  "$REPO_ROOT/sdk/go/runpod"
  "$REPO_ROOT/examples/go"
)

for mod in "${MODULES[@]}"; do
  echo "==> go mod tidy in $mod"
  (cd "$mod" && go mod tidy)
done

echo ""
echo "All modules tidied."
