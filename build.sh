#!/bin/sh

set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
RELEASE_DIR="$ROOT_DIR/Releases"
WEBUI_DIR="$ROOT_DIR/webui"

mkdir -p "$RELEASE_DIR"

echo "start build webui >>>"
cd "$WEBUI_DIR"
if [ ! -d node_modules ]; then
  npm ci
fi
npm run build

cd "$ROOT_DIR"

# 【darwin/amd64】
echo "start build darwin/amd64 >>>"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags '-w -s' -o "$RELEASE_DIR/dedao-darwin-amd64" main.go

# 【windows/amd64】
echo "start build windows/amd64 >>>"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags '-w -s' -o "$RELEASE_DIR/dedao-windows-amd64.exe" main.go

# 【linux/amd64】
echo "start build linux/amd64 >>>"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o "$RELEASE_DIR/dedao-linux-amd64" main.go

echo "All build success!!!"
