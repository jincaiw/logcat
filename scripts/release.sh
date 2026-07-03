#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
TARGET_INPUT="${1:-${VERSION:-}}"
SKIP_TAG="${SKIP_TAG:-0}"

if [[ -z "$TARGET_INPUT" || "$TARGET_INPUT" == "--help" || "$TARGET_INPUT" == "-h" ]]; then
  cat <<'EOF'
Usage:
  bash scripts/release.sh <version>

What it does:
  1. Syncs version references across the repository
  2. Generates release note skeleton if missing
  3. Runs tests and build verification
  4. Builds linux amd64/arm64 release packages
  5. Creates a local git tag vX.Y.Z (unless SKIP_TAG=1)

For a one-shot publish (auto-commit + push code + tag + DockerHub + checks), use:
  bash scripts/publish-release.sh <version>
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
if [[ "$SKIP_TAG" == "1" ]]; then
  echo "Skipping tag creation because SKIP_TAG=1."
elif git rev-parse "$TAG" >/dev/null 2>&1; then
  echo "Tag $TAG already exists, skipping tag creation."
else
  git tag -a "$TAG" -m "Release $TAG"
  echo "Created tag $TAG"
fi

echo "Release preparation complete."
echo "Next steps:"
echo "  1) Review $RELEASE_NOTES"
echo "  2) Commit and push release changes"
echo "  3) Push tag: git push origin $TAG"
echo "  4) Publish GitHub Release from the pushed tag"
echo "  5) If needed, push Docker image: docker build -t qing1205/logcat:$VERSION -t qing1205/logcat:latest . && docker push qing1205/logcat:$VERSION && docker push qing1205/logcat:latest"
