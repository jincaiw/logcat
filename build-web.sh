#!/bin/bash
set -e

echo "Building logcat Web Server..."

cd "$(dirname "$0")"

echo "1. Installing frontend dependencies..."
cd frontend
npm install
cd ..

echo "2. Type checking frontend..."
cd frontend
npm run typecheck
cd ..

echo "3. Building frontend..."
cd frontend
npm run build
cd ..

echo "4. Building web server binary for Linux amd64..."
GOOS=linux GOARCH=amd64 go build -o build/bin/logcat-web

echo "5. Copying templates..."
cp -r templates build/bin/templates

echo "Done! Output: build/bin/"
echo ""
echo "Files:"
echo "  - logcat-web (binary)"
echo "  - templates/ (parse_templates.json, filter_policies.json)"
echo ""
echo "Usage: ./logcat-web [port]"
echo "Default port: 8080"
