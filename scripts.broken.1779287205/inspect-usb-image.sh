#!/usr/bin/env bash
#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="$ROOT_DIR/out/bos-usb.img"
LOOP=""

cleanup() {
  set +e
  if mountpoint -q "$ROOT_DIR/out/inspect-root"; then sudo umount "$ROOT_DIR/out/inspect-root"; fi
  if mountpoint -q "$ROOT_DIR/out/inspect-efi"; then sudo umount "$ROOT_DIR/out/inspect-efi"; fi
  if [ -n "$LOOP" ]; then sudo losetup -d "$LOOP" 2>/dev/null || true; fi
}
trap cleanup EXIT

if [ "$EUID" -ne 0 ]; then
  echo "Run as root:"
  echo "  sudo $0"
  exit 1
fi

if [ ! -f "$IMAGE" ]; then
  echo "Missing image: $IMAGE"
  exit 1
fi

mkdir -p "$ROOT_DIR/out/inspect-root" "$ROOT_DIR/out/inspect-efi"

LOOP="$(losetup --find --partscan --show "$IMAGE")"
sleep 1

echo "Disk:"
lsblk "$LOOP"
echo

echo "Partition UUIDs:"
blkid "${LOOP}p1" "${LOOP}p2"
echo

mount "${LOOP}p1" "$ROOT_DIR/out/inspect-efi"
mount "${LOOP}p2" "$ROOT_DIR/out/inspect-root"

echo
echo "EFI files:"
find "$ROOT_DIR/out/inspect-efi" -maxdepth 5 -type f | sort

echo
echo "Root boot files:"
find "$ROOT_DIR/out/inspect-root/boot" -maxdepth 4 -type f | sort

echo
echo "inittab:"
sed -n '1,120p' "$ROOT_DIR/out/inspect-root/etc/inittab"
