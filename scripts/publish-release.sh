#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
TARGET_INPUT="${1:-${VERSION:-}}"
PUSH_GIT="${PUSH_GIT:-1}"
PUSH_DOCKER="${PUSH_DOCKER:-1}"
DOCKER_REPO="${DOCKER_REPO:-qing1205/logcat}"
COMMIT_RELEASE="${COMMIT_RELEASE:-1}"
VERIFY_RELEASE="${VERIFY_RELEASE:-1}"
DRY_RUN="${DRY_RUN:-0}"
PAGES_URL="${PAGES_URL:-https://logcat.mujizi.com/}"

resolve_github_repo() {
  local remote
  remote="$(git -C "$ROOT_DIR" remote get-url origin 2>/dev/null || true)"
  if [[ "$remote" =~ github\.com[:/]([^/]+)/([^/.]+)(\.git)?$ ]]; then
    echo "${BASH_REMATCH[1]}/${BASH_REMATCH[2]}"
  fi
}

GH_REPO="${GH_REPO:-$(resolve_github_repo)}"

usage() {
  cat <<'EOF'
Usage:
  bash scripts/publish-release.sh [--dry-run] <version>

What it does:
  1. Runs the standard release flow without creating a tag
  2. Commits the release changes
  3. Creates and pushes the release tag
  4. Builds and pushes DockerHub images: :<version> and :latest
  5. Runs release checks via scripts/release-check.sh

Environment:
  DOCKER_REPO=qing1205/logcat
  GH_REPO=owner/repo
  PAGES_URL=https://logcat.mujizi.com/
  PUSH_GIT=0|1
  PUSH_DOCKER=0|1
  COMMIT_RELEASE=0|1
  VERIFY_RELEASE=0|1
  DRY_RUN=0|1
EOF
}

if [[ $# -gt 0 ]]; then
  case "$1" in
    --help|-h)
      usage
      exit 0
      ;;
    --dry-run)
      DRY_RUN=1
      shift
      TARGET_INPUT="${1:-${VERSION:-}}"
      ;;
  esac
fi

if [[ -z "$TARGET_INPUT" ]]; then
  usage
  exit 1
fi

VERSION="${TARGET_INPUT#v}"
if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Invalid version: $TARGET_INPUT" >&2
  exit 1
fi

TAG="v${VERSION}"

if [[ "$DRY_RUN" == "1" ]]; then
  echo "Dry run mode"
  echo "Version: $VERSION"
  echo "Tag: $TAG"
  echo "Repo: ${GH_REPO:-unknown}"
  echo "Docker repo: $DOCKER_REPO"
  echo "Would run: SKIP_TAG=1 bash scripts/release.sh \"$VERSION\""
  [[ "$COMMIT_RELEASE" != "0" ]] && echo "Would commit release changes"
  echo "Would create tag: $TAG"
  [[ "$PUSH_GIT" != "0" ]] && echo "Would push main/master and tag to origin"
  [[ "$PUSH_DOCKER" != "0" ]] && echo "Would build and push Docker images: $DOCKER_REPO:$VERSION and :latest"
  [[ "$VERIFY_RELEASE" != "0" ]] && echo "Would run scripts/release-check.sh $VERSION"
  exit 0
fi

cd "$ROOT_DIR"

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || { echo "$1 is required" >&2; exit 1; }
}

need_cmd git
if [[ "$PUSH_DOCKER" != "0" ]]; then
  need_cmd docker
  docker info >/dev/null 2>&1 || { echo "Docker daemon is not available" >&2; exit 1; }
fi
if [[ "$VERIFY_RELEASE" != "0" ]]; then
  need_cmd python3
fi

if command -v gh >/dev/null 2>&1; then
  gh auth status -h github.com >/dev/null 2>&1 || echo "Warning: gh auth status failed; release verification may be limited."
fi

SKIP_TAG=1 bash scripts/release.sh "$VERSION"

if [[ "$COMMIT_RELEASE" != "0" ]]; then
  if git status --porcelain | grep -q .; then
    git add -A
    git commit -m "chore(release): ${TAG}"
    echo "Committed release changes for $TAG"
  else
    echo "No release changes to commit."
  fi
else
  if git status --porcelain | grep -q .; then
    echo "COMMIT_RELEASE=0 but the worktree is dirty. Commit the release changes before tagging." >&2
    exit 1
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

if [[ "$VERIFY_RELEASE" != "0" ]]; then
  echo "Running release checks..."
  GH_REPO="$GH_REPO" PAGES_URL="$PAGES_URL" DOCKER_REPO="$DOCKER_REPO" bash scripts/release-check.sh "$VERSION"
fi

if [[ "$GH_REPO" == */* ]]; then
  GITHUB_RELEASE_URL="https://github.com/${GH_REPO}/releases/tag/${TAG}"
else
  GITHUB_RELEASE_URL="unknown"
fi
DOCKER_TAG_URL="https://hub.docker.com/r/${DOCKER_REPO}/tags?name=${VERSION}"

echo "Release publish complete."
echo "Final URLs:"
echo "  GitHub Release: ${GITHUB_RELEASE_URL}"
echo "  Pages: ${PAGES_URL}"
echo "  DockerHub tags: ${DOCKER_TAG_URL}"
echo "Next: if GitHub Release pages are still building, wait a moment and re-check the links for $TAG."
