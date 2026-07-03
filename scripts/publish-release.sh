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
detect_github_repo() {
  local remote
  remote="$(git -C "$ROOT_DIR" remote get-url origin 2>/dev/null || true)"
  if [[ "$remote" =~ github\.com[:/]([^/]+)/([^/.]+)(\.git)?$ ]]; then
    echo "${BASH_REMATCH[1]}/${BASH_REMATCH[2]}"
  fi
}

GH_REPO="${GH_REPO:-$(detect_github_repo)}"
PAGES_URL="${PAGES_URL:-https://logcat.mujizi.com/}"

usage() {
  cat <<'EOF'
Usage:
  bash scripts/publish-release.sh [--dry-run] <version>

What it does:
  1. Runs the standard release flow without creating a tag
  2. Commits the release changes
  3. Creates and pushes the release tag
  4. Builds and pushes DockerHub images: :<version> and :latest
  5. Verifies GitHub Release / Pages / DockerHub when enabled

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
  [[ "$VERIFY_RELEASE" != "0" ]] && echo "Would verify GitHub Release / Pages / DockerHub"
  exit 0
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
if [[ "$VERIFY_RELEASE" != "0" ]] && ! command -v python3 >/dev/null 2>&1; then
  echo "python3 is required when VERIFY_RELEASE=1" >&2
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

retry() {
  local attempts="$1"
  local delay="$2"
  shift 2
  local n=1
  while true; do
    if "$@"; then
      return 0
    fi
    if [[ "$n" -ge "$attempts" ]]; then
      return 1
    fi
    sleep "$delay"
    n=$((n + 1))
  done
}

verify_github_release_once() {
  gh api "repos/${GH_REPO}/releases/tags/${TAG}" --jq '{tag_name:.tag_name,name:.name,html_url:.html_url,published_at:.published_at}' >/tmp/logcat-release-check.json
}

verify_pages_once() {
  python3 - "$PAGES_URL" "$VERSION" <<'PY'
import sys
import urllib.request

url = sys.argv[1]
version = sys.argv[2]
with urllib.request.urlopen(url, timeout=20) as r:
    body = r.read().decode('utf-8', 'ignore')
if version not in body:
    raise SystemExit(f'Pages check failed: version {version} not found at {url}')
print(f'Pages check passed: {url}')
PY
}

verify_dockerhub_once() {
  python3 - "$DOCKER_REPO" "$VERSION" <<'PY'
import json
import sys
import urllib.request

repo = sys.argv[1]
version = sys.argv[2]
for tag in [version, 'latest']:
    url = f'https://hub.docker.com/v2/repositories/{repo}/tags/{tag}'
    with urllib.request.urlopen(url, timeout=20) as r:
        data = json.load(r)
    if not data.get('name'):
        raise SystemExit(f'DockerHub check failed: missing tag {tag}')
print(f'DockerHub check passed: {repo}:{version} and :latest')
PY
}

verify_github_release() {
  if [[ -z "$GH_REPO" || "$GH_REPO" != */* ]]; then
    echo "Skip GitHub Release check: GH_REPO is not set or invalid."
    return 0
  fi
  if ! command -v gh >/dev/null 2>&1; then
    echo "Skip GitHub Release check: gh not installed."
    return 0
  fi
  if ! gh auth status -h github.com >/dev/null 2>&1; then
    echo "Skip GitHub Release check: gh not authenticated."
    return 0
  fi
  retry 12 10 verify_github_release_once
  echo "GitHub Release check passed: ${TAG}"
  cat /tmp/logcat-release-check.json
}

verify_pages() {
  retry 12 10 verify_pages_once
}

verify_dockerhub() {
  retry 12 10 verify_dockerhub_once
}

if [[ "$VERIFY_RELEASE" != "0" ]]; then
  echo "Verifying release outputs..."
  verify_github_release
  verify_pages
  verify_dockerhub
fi

echo "Release publish complete."
echo "Next: if GitHub Release pages are still building, wait a moment and re-check the links for $TAG."
