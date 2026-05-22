#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SRC_DIR="$ROOT_DIR/src"
OUT="$SRC_DIR/bos"

find_go() {
  if command -v go >/dev/null 2>&1; then
    command -v go
    return 0
  fi

  for p in \
    "$HOME/software/go/bin/go" \
    "/home/${SUDO_USER:-}/software/go/bin/go" \
    "/usr/local/go/bin/go" \
    "/usr/bin/go"
  do
    if [ -x "$p" ]; then
      echo "$p"
      return 0
    fi
  done

  return 1
}

if [ ! -f "$SRC_DIR/go.mod" ]; then
  echo "Missing Go module: $SRC_DIR/go.mod"
  exit 1
fi

GO_BIN="$(find_go || true)"
if [ -z "$GO_BIN" ]; then
  echo "Could not find go. Run this script as your normal user, or install Go."
  exit 1
fi

cd "$SRC_DIR"

if command -v musl-gcc >/dev/null 2>&1; then
  echo "==> Building static BOS with musl/cgo"
  CC=musl-gcc CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    "$GO_BIN" build -trimpath \
    -ldflags="-linkmode external -extldflags '-static' -s -w" \
    -o "$OUT" .
else
  echo "==> Building static BOS with pure Go"
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    "$GO_BIN" build -trimpath -ldflags="-s -w" \
    -o "$OUT" .
fi

chmod 0755 "$OUT"
file "$OUT"
ldd "$OUT" || true
