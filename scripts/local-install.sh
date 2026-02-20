#!/usr/bin/env bash
# Build the provider + SDK and install into an SST v3 project for local development.
#
# Usage:
#   ./scripts/local-install.sh /path/to/your-sst-project
#
# Prerequisites:
#   - Go, Node.js, Bun, and Pulumi CLI installed
#   - The SST project has been initialized (`.sst/platform/` exists)
#     Run `npx sst deploy --stage dev` once first if it doesn't.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
VERSION="1.0.0-alpha.0"
PROVIDER_NAME="pulumi-resource-runpod"

if [ $# -lt 1 ]; then
    echo "Usage: $0 <path-to-sst-project>"
    echo ""
    echo "Examples:"
    echo "  $0 ../my-sst-app"
    echo "  $0 /Users/me/projects/infra"
    exit 1
fi

SST_PROJECT="$(cd "$1" && pwd)"
SST_PLATFORM="$SST_PROJECT/.sst/platform"

if [ ! -d "$SST_PLATFORM" ]; then
    echo "Error: $SST_PLATFORM does not exist."
    echo "Run 'npx sst deploy --stage dev' in your SST project first to initialize it."
    exit 1
fi

echo "==> Building provider binary..."
cd "$REPO_ROOT"
go build \
    -ldflags "-X github.com/runpod/pulumi-runpod/provider.Version=$VERSION" \
    -o "bin/$PROVIDER_NAME" \
    ./provider/cmd/pulumi-resource-runpod/

echo "==> Installing Pulumi plugin..."
pulumi plugin install resource runpod "$VERSION" --file "bin/$PROVIDER_NAME"
cp "bin/$PROVIDER_NAME" "$(go env GOPATH)/bin/$PROVIDER_NAME"

echo "==> Building Node.js SDK..."
cd "$REPO_ROOT/sdk/nodejs"
npm install --silent
npm run build

echo "==> Packing SDK tarball..."
TARBALL=$(npm pack --silent 2>/dev/null)
TARBALL_PATH="$REPO_ROOT/sdk/nodejs/$TARBALL"

echo "==> Installing SDK into SST project..."
cd "$SST_PLATFORM"
bun add "$TARBALL_PATH"

echo ""
echo "Done! Provider and SDK installed into: $SST_PROJECT"
echo ""
echo "In your sst.config.ts:"
echo ""
echo '  async run() {'
echo '    const runpod = await import("@runpod/pulumi");'
echo '    const provider = new runpod.Provider("runpod", {'
echo '      apiKey: process.env.RUNPOD_API_KEY,'
echo '    });'
echo '    // ... create resources with { provider }'
echo '  }'
