#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
TARGET_INPUT=""
MODE="all"
GH_REPO="${GH_REPO:-}"
PAGES_URL="${PAGES_URL:-https://logcat.mujizi.com/}"
DOCKER_REPO="${DOCKER_REPO:-qing1205/logcat}"
CHECK_GITHUB="${CHECK_GITHUB:-1}"
CHECK_PAGES="${CHECK_PAGES:-1}"
CHECK_DOCKER="${CHECK_DOCKER:-1}"
ATTEMPTS="${ATTEMPTS:-12}"
DELAY_SECONDS="${DELAY_SECONDS:-10}"

usage() {
  cat <<'EOF'
Usage:
  bash scripts/release-check.sh [--github-only|--pages-only|--docker-only|--all] <version>

What it does:
  1. Verifies GitHub Release exists for v<version>
  2. Verifies Pages contains the version content
  3. Verifies DockerHub has :<version> and :latest
  4. Prints final URLs for the release artifacts

Environment:
  GH_REPO=owner/repo
  PAGES_URL=https://logcat.mujizi.com/
  DOCKER_REPO=qing1205/logcat
  CHECK_GITHUB=0|1
  CHECK_PAGES=0|1
  CHECK_DOCKER=0|1
  ATTEMPTS=12
  DELAY_SECONDS=10
EOF
}

detect_github_repo() {
  local remote
  remote="$(git -C "$ROOT_DIR" remote get-url origin 2>/dev/null || true)"
  if [[ "$remote" =~ github\.com[:/]([^/]+)/([^/.]+)(\.git)?$ ]]; then
    echo "${BASH_REMATCH[1]}/${BASH_REMATCH[2]}"
  fi
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --help|-h)
      usage
      exit 0
      ;;
    --github-only)
      MODE="github"
      shift
      ;;
    --pages-only)
      MODE="pages"
      shift
      ;;
    --docker-only)
      MODE="docker"
      shift
      ;;
    --all)
      MODE="all"
      shift
      ;;
    --)
      shift
      break
      ;;
    *)
      TARGET_INPUT="$1"
      shift
      break
      ;;
  esac
done

if [[ -z "$TARGET_INPUT" && $# -gt 0 ]]; then
  TARGET_INPUT="$1"
fi

if [[ -z "$TARGET_INPUT" || "$TARGET_INPUT" == "--help" || "$TARGET_INPUT" == "-h" ]]; then
  usage
  exit 0
fi

VERSION="${TARGET_INPUT#v}"
if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Invalid version: $TARGET_INPUT" >&2
  exit 1
fi

TAG="v${VERSION}"
GH_REPO="${GH_REPO:-$(detect_github_repo)}"
GITHUB_RELEASE_URL="https://github.com/${GH_REPO}/releases/tag/${TAG}"
DOCKER_TAG_URL="https://hub.docker.com/r/${DOCKER_REPO}/tags?name=${VERSION}"

declare -a WANT_GITHUB WANT_PAGES WANT_DOCKER
WANT_GITHUB=()
WANT_PAGES=()
WANT_DOCKER=()

case "$MODE" in
  github) CHECK_PAGES=0; CHECK_DOCKER=0 ;;
  pages) CHECK_GITHUB=0; CHECK_DOCKER=0 ;;
  docker) CHECK_GITHUB=0; CHECK_PAGES=0 ;;
esac

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

check_github_release() {
  if [[ "$CHECK_GITHUB" == "0" ]]; then
    echo "Skip GitHub Release check."
    return 0
  fi
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
  retry "$ATTEMPTS" "$DELAY_SECONDS" verify_github_release_once
  echo "GitHub Release check passed: ${TAG}"
  cat /tmp/logcat-release-check.json
}

check_pages() {
  if [[ "$CHECK_PAGES" == "0" ]]; then
    echo "Skip Pages check."
    return 0
  fi
  retry "$ATTEMPTS" "$DELAY_SECONDS" verify_pages_once
}

check_dockerhub() {
  if [[ "$CHECK_DOCKER" == "0" ]]; then
    echo "Skip DockerHub check."
    return 0
  fi
  retry "$ATTEMPTS" "$DELAY_SECONDS" verify_dockerhub_once
}

echo "Running release checks for ${TAG}..."
check_github_release
check_pages
check_dockerhub

echo "Release checks complete for ${TAG}."
echo "URLs:"
echo "  GitHub Release: ${GITHUB_RELEASE_URL}"
echo "  Pages: ${PAGES_URL}"
echo "  DockerHub tags: ${DOCKER_TAG_URL}"
