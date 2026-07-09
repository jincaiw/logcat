#!/bin/bash
set -e

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
DEFAULT_VERSION="$(cat "$ROOT_DIR/VERSION" 2>/dev/null || true)"
DEFAULT_VERSION="${DEFAULT_VERSION:-0.2.10}"
APP_VERSION="${APP_VERSION:-$DEFAULT_VERSION}"
TARGET_OS="${TARGET_OS:-linux}"
TARGET_ARCH="${TARGET_ARCH:-amd64}"
OUTPUT_DIR="build/logcat-${APP_VERSION}-${TARGET_OS}-${TARGET_ARCH}"

if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
  echo "Usage: APP_VERSION=${DEFAULT_VERSION} TARGET_OS=linux TARGET_ARCH=amd64 bash build-web.sh"
  exit 0
fi

echo "Building logcat ${APP_VERSION} Web Server for ${TARGET_OS}/${TARGET_ARCH}..."

cd "$ROOT_DIR"

rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

echo "1. Installing frontend dependencies..."
cd frontend
npm ci

echo "2. Building frontend..."
npm run build
cd ..

echo "3. Building web server binary..."
BIN_NAME="logcat"
if [[ "$TARGET_OS" == "windows" ]]; then
  BIN_NAME="logcat.exe"
fi
GOOS="$TARGET_OS" GOARCH="$TARGET_ARCH" CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "$OUTPUT_DIR/$BIN_NAME" .

echo "4. Copying templates and docs..."
cp -r templates "$OUTPUT_DIR/templates"
cp README.md "$OUTPUT_DIR/README.md"

cat > "$OUTPUT_DIR/start.sh" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
export SYSLG_ALERT_DATA_DIR="${SYSLG_ALERT_DATA_DIR:-$DIR/data}"
export SYSLG_ALERT_TEMPLATES_DIR="${SYSLG_ALERT_TEMPLATES_DIR:-$DIR/templates}"
exec "$DIR/logcat" "${1:-8080}"
EOF
chmod +x "$OUTPUT_DIR/start.sh"

TARBALL="build/logcat-${APP_VERSION}-${TARGET_OS}-${TARGET_ARCH}.tar.gz"
tar -C build -czf "$TARBALL" "logcat-${APP_VERSION}-${TARGET_OS}-${TARGET_ARCH}"

echo "Done!"
echo "  Directory: $OUTPUT_DIR"
echo "  Archive:   $TARBALL"
echo ""
echo "Usage:"
echo "  tar -xzf $TARBALL"
echo "  cd logcat-${APP_VERSION}-${TARGET_OS}-${TARGET_ARCH}"
echo "  ./start.sh 8080"
