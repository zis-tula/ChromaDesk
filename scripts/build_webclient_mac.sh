#!/bin/bash
set -e

echo "---------------------------------------------------------------"
echo "Build WebClient (Vue 3 + Vite)"
echo "---------------------------------------------------------------"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
WEBCLIENT_DIR="$SCRIPT_DIR/../WebClient"

cd "$WEBCLIENT_DIR"

echo "[*] Installing dependencies..."
npm install

echo "[*] Building..."
npm run build

echo "[*] Copying remote.html and assets to dist..."
cp -r js dist/
cp remote.html dist/
cp favicon.ico dist/ 2>/dev/null || true
[ -d images ] && cp -r images dist/

echo "[*] WebClient build complete: $WEBCLIENT_DIR/dist"
