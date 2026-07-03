#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
TARGET_INPUT="${1:-${VERSION:-}}"
PUSH_GIT="${PUSH_GIT:-1}"
PUSH_DOCKER="${PUSH_DOCKER:-1}"
DOCKER_REPO="${DOCKER_REPO:-qing1205/logcat}"

if [[ -z "$TARGET_INPUT" || "$TARGET_INPUT" == "--help" || "$TARGET_INPUT" == "-h" ]]; then
  cat <<'EOF'
Usage:
  bash scripts/publish-release.sh <version>

What it does:
  1. Runs the standard release flow (version sync, tests, build, local tag)
  2. Pushes main/master and the tag to origin
  3. Builds and pushes DockerHub images: :<version> and :latest

Environment:
  DOCKER_REPO=qing1205/logcat
  PUSH_GIT=0|1
  PUSH_DOCKER=0|1
EOF
  exit 0
fi

VERSION="${TARGET_INPUT#v}"
if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Invalid version: $TARGET_INPUT" >&2
  exit 1
fi

cd "$ROOT_DIR"

# Local preparation: version sync, tests, build packages, local tag.
bash scripts/release.sh "$VERSION"

TAG="v${VERSION}"

echo "Publishing release artifacts for $TAG..."

if [[ "$PUSH_GIT" != "0" ]]; then
  echo "Pushing branches and tag to origin..."
  git push origin HEAD:main HEAD:master
  git push origin "$TAG"
fi

if [[ "$PUSH_DOCKER" != "0" ]]; then
  echo "Building Docker image $DOCKER_REPO:$VERSION and $DOCKER_REPO:latest..."
  docker build -t "$DOCKER_REPO:$VERSION" -t "$DOCKER_REPO:latest" .
  echo "Pushing Docker image..."
  docker push "$DOCKER_REPO:$VERSION"
  docker push "$DOCKER_REPO:latest"
fi

echo "Release publish complete."
echo "Next: verify GitHub Release and Pages status for $TAG."
