#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
TARGET_INPUT="${1:-${VERSION:-}}"

if [[ -z "$TARGET_INPUT" || "$TARGET_INPUT" == "--help" || "$TARGET_INPUT" == "-h" ]]; then
  cat <<'EOF'
Usage:
  bash scripts/release.sh <version>

What it does:
  1. Syncs version references across the repository
  2. Generates release note skeleton if missing
  3. Runs tests and build verification
  4. Builds linux amd64/arm64 release packages
  5. Creates a local git tag vX.Y.Z
EOF
  exit 0
fi

VERSION="${TARGET_INPUT#v}"
if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Invalid version: $TARGET_INPUT" >&2
  exit 1
fi

cd "$ROOT_DIR"

bash scripts/bump-version.sh "$VERSION"

RELEASE_NOTES="docs/release-notes-v${VERSION}.md"
if [[ ! -f "$RELEASE_NOTES" ]]; then
  if [[ ! -f docs/release-notes-template.md ]]; then
    echo "Missing docs/release-notes-template.md" >&2
    exit 1
  fi
  sed "s/{{VERSION}}/${VERSION}/g" docs/release-notes-template.md > "$RELEASE_NOTES"
  echo "Created $RELEASE_NOTES"
fi

echo "Running backend tests..."
go test ./...

echo "Building release packages..."
APP_VERSION="$VERSION" TARGET_OS=linux TARGET_ARCH=amd64 bash build-web.sh
APP_VERSION="$VERSION" TARGET_OS=linux TARGET_ARCH=arm64 bash build-web.sh

TAG="v${VERSION}"
if git rev-parse "$TAG" >/dev/null 2>&1; then
  echo "Tag $TAG already exists, skipping tag creation."
else
  git tag -a "$TAG" -m "Release $TAG"
  echo "Created tag $TAG"
fi

echo "Release preparation complete."
echo "Next: review $RELEASE_NOTES, push commits, push tag, then publish GitHub Release."
