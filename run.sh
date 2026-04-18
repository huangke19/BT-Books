#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_PATH="$ROOT_DIR/bt-books"

cd "$ROOT_DIR"

if [[ -z "${HTTPS_PROXY:-}" && -z "${https_proxy:-}" ]]; then
  echo "⚠️  未设置 HTTPS_PROXY，若当前网络无法直连 Z-Library，搜索和下载可能失败。"
fi

echo "==> 编译 BT-Books"
go build -o "$BIN_PATH" .

echo "==> 启动 BT-Books TUI"
exec "$BIN_PATH"
