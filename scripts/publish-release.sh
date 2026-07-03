#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
TARGET_INPUT="${1:-${VERSION:-}}"
PUSH_GIT="${PUSH_GIT:-1}"
PUSH_DOCKER="${PUSH_DOCKER:-1}"
DOCKER_REPO="${DOCKER_REPO:-qing1205/logcat}"
COMMIT_RELEASE="${COMMIT_RELEASE:-1}"

if [[ -z "$TARGET_INPUT" || "$TARGET_INPUT" == "--help" || "$TARGET_INPUT" == "-h" ]]; then
  cat <<'EOF'
Usage:
  bash scripts/publish-release.sh <version>

What it does:
  1. Runs the standard release flow without creating a tag
  2. Commits the release changes
  3. Creates and pushes the release tag
  4. Builds and pushes DockerHub images: :<version> and :latest

Environment:
  DOCKER_REPO=qing1205/logcat
  PUSH_GIT=0|1
  PUSH_DOCKER=0|1
  COMMIT_RELEASE=0|1
EOF
  exit 0
fi

VERSION="${TARGET_INPUT#v}"
if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Invalid version: $TARGET_INPUT" >&2
  exit 1
fi

cd "$ROOT_DIR"

if ! command -v git >/dev/null 2>&1; then
  echo "git is required" >&2
  exit 1
fi
if [[ "$PUSH_DOCKER" != "0" ]] && ! command -v docker >/dev/null 2>&1; then
  echo "docker is required when PUSH_DOCKER=1" >&2
  exit 1
fi
if command -v gh >/dev/null 2>&1; then
  gh auth status -h github.com >/dev/null 2>&1 || echo "Warning: gh auth status failed; GitHub release verification may be limited."
fi
if [[ "$PUSH_DOCKER" != "0" ]]; then
  docker info >/dev/null 2>&1 || { echo "Docker daemon is not available" >&2; exit 1; }
fi

# Local preparation: version sync, tests, build packages, release notes, no tag yet.
SKIP_TAG=1 bash scripts/release.sh "$VERSION"
TAG="v${VERSION}"

if [[ "$COMMIT_RELEASE" != "0" ]]; then
  if git status --porcelain | grep -q .; then
    git add -A
    git commit -m "chore(release): ${TAG}"
    echo "Committed release changes for $TAG"
  else
    echo "No release changes to commit."
  fi
fi

if git rev-parse "$TAG" >/dev/null 2>&1; then
  echo "Tag $TAG already exists locally. Delete it before re-running this command." >&2
  exit 1
fi

git tag -a "$TAG" -m "Release $TAG"
echo "Created tag $TAG"

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
