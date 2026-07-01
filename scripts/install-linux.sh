#!/usr/bin/env bash
set -euo pipefail

VERSION="${VERSION:-0.2.3}"
INSTALL_DIR="${INSTALL_DIR:-/opt/logcat}"
PORT="${PORT:-8080}"
REPO="${REPO:-jincaiw/logcat}"
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64) TARGET_ARCH="amd64" ;;
  aarch64|arm64) TARGET_ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
  cat <<EOF
Install logcat on Linux.

Usage:
  curl -fsSL https://raw.githubusercontent.com/${REPO}/v${VERSION}/scripts/install-linux.sh | sudo bash

Environment:
  VERSION=${VERSION}
  INSTALL_DIR=${INSTALL_DIR}
  PORT=${PORT}
EOF
  exit 0
fi

if [[ $EUID -ne 0 ]]; then
  echo "Please run as root, for example: sudo bash scripts/install-linux.sh" >&2
  exit 1
fi

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT
ASSET="logcat-${VERSION}-linux-${TARGET_ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/v${VERSION}/${ASSET}"

echo "Downloading $URL"
if command -v curl >/dev/null 2>&1; then
  curl -fL "$URL" -o "$TMP_DIR/$ASSET"
elif command -v wget >/dev/null 2>&1; then
  wget -O "$TMP_DIR/$ASSET" "$URL"
else
  echo "curl or wget is required" >&2
  exit 1
fi

mkdir -p "$INSTALL_DIR"
tar -xzf "$TMP_DIR/$ASSET" -C "$TMP_DIR"
cp -R "$TMP_DIR/logcat-${VERSION}-linux-${TARGET_ARCH}/." "$INSTALL_DIR/"
mkdir -p "$INSTALL_DIR/data"

cat > /etc/systemd/system/logcat.service <<EOF
[Unit]
Description=logcat syslog alert server
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
WorkingDirectory=${INSTALL_DIR}
Environment=SYSLG_ALERT_DATA_DIR=${INSTALL_DIR}/data
Environment=SYSLG_ALERT_TEMPLATES_DIR=${INSTALL_DIR}/templates
Environment=LOGCAT_OPEN_BROWSER=0
ExecStart=${INSTALL_DIR}/logcat ${PORT}
Restart=always
RestartSec=3
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable --now logcat

echo "logcat installed successfully."
echo "Open: http://<server-ip>:${PORT}"
echo "Default account: admin / admin123"
echo "Please change the password after first login."
