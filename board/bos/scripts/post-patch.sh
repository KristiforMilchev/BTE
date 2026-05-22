#!/bin/sh
set -e

for f in $(find "$BUILD_DIR/package/ncurses" -name '*.mk'); do
    sed -i '/terminfo\/y/d' "$f" || true
done

find "$BUILD_DIR" -path '*ncurses*.mk' -exec sed -i '/terminfo\/y/d' {} \; || true
