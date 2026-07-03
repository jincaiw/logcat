#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
VERSION_INPUT="${1:-}"

if [[ -z "$VERSION_INPUT" || "$VERSION_INPUT" == "--help" || "$VERSION_INPUT" == "-h" ]]; then
  cat <<'EOF'
Usage:
  bash scripts/bump-version.sh <version>

This script updates the repository version references from the current version
stored in VERSION to the target version.
EOF
  exit 0
fi

VERSION="${VERSION_INPUT#v}"
if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Invalid version: $VERSION_INPUT" >&2
  exit 1
fi

python3 - "$ROOT_DIR" "$VERSION" <<'PY'
import re
import sys
from pathlib import Path

root = Path(sys.argv[1])
version = sys.argv[2]
old_version = (root / 'VERSION').read_text().strip()

(root / 'VERSION').write_text(version + '\n')

text_files = [
    root / 'README.md',
    root / 'README-zh.md',
    root / 'docs/index.md',
    root / 'docs/index.html',
    root / 'docs/installation.md',
    root / 'docs/installation.html',
    root / 'docker-compose.yml',
    root / 'build-web.sh',
    root / 'scripts/install-linux.sh',
]

for path in text_files:
    text = path.read_text()
    text = text.replace(f'v{old_version}', f'v{version}')
    text = text.replace(old_version, version)
    path.write_text(text)

# Go runtime version
path = root / 'pkg/constants/constants.go'
text = path.read_text()
text = re.sub(r'(AppVersion\s*=\s*")[^"]+(")', lambda m: f'{m.group(1)}{version}{m.group(2)}', text, count=1)
path.write_text(text)

# Frontend version marker
path = root / 'frontend/src/version.ts'
text = path.read_text()
text = re.sub(r"(APP_VERSION\s*=\s*')[^']+(')", lambda m: f"{m.group(1)}{version}{m.group(2)}", text, count=1)
path.write_text(text)

# Frontend package metadata
path = root / 'frontend/package.json'
text = path.read_text()
text = re.sub(r'("version"\s*:\s*")[^"]+(")', lambda m: f'{m.group(1)}{version}{m.group(2)}', text, count=1)
path.write_text(text)

path = root / 'frontend/package-lock.json'
text = path.read_text()
text = re.sub(r'("version"\s*:\s*")[^"]+(")', lambda m: f'{m.group(1)}{version}{m.group(2)}', text, count=2)
path.write_text(text)

PY

echo "Version bumped to ${VERSION}."
