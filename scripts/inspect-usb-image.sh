#!/usr/bin/env bash
set -euo pipefail
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="$ROOT_DIR/out/bos-usb.img"
MNT="/mnt/bos-efi-test"
LOOP=""
cleanup(){ set +e; mountpoint -q "$MNT" && sudo umount "$MNT"; [ -n "$LOOP" ] && sudo losetup -d "$LOOP" 2>/dev/null || true; }
trap cleanup EXIT
LOOP="$(sudo losetup -Pf --show "$IMAGE")"
sudo mkdir -p "$MNT"
sudo mount "${LOOP}p1" "$MNT"
df -h "$MNT"
find "$MNT" -maxdepth 5 -type f -print | sort
file "$MNT/EFI/BOOT/BOOTX64.EFI"
echo '--- bos.conf ---'
cat "$MNT/loader/entries/bos.conf"
