\
#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SRC_DIR="$ROOT_DIR/src"

if [ ! -f "$SRC_DIR/go.mod" ]; then
  echo "Missing Go module: $SRC_DIR/go.mod"
  exit 1
fi

cd "$SRC_DIR"

# Build a static Linux binary. This is preferred for an appliance rootfs.
# If your Ledger/HID dependency requires CGO, install musl-tools and kernel headers first.
if command -v musl-gcc >/dev/null 2>&1; then
  echo "==> Building static BOS with musl/cgo"
  CC=musl-gcc CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -mod=vendor -trimpath \
    -ldflags="-linkmode external -extldflags '-static' -s -w" \
    -o "$SRC_DIR/bos" .
else
  echo "==> musl-gcc not found; building static pure-Go BOS"
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -mod=vendor -trimpath -ldflags="-s -w" \
    -o "$SRC_DIR/bos" .
fi

file "$SRC_DIR/bos"
ldd "$SRC_DIR/bos" || true
