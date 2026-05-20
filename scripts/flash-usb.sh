\
#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="$ROOT_DIR/out/bos-usb.img"
DISK="${1:-}"

if [ "$EUID" -ne 0 ]; then
  echo "Run as root:"
  echo "  sudo $0 /dev/sdX"
  exit 1
fi

if [ -z "$DISK" ]; then
  echo "Usage:"
  echo "  sudo $0 /dev/sdX"
  exit 1
fi

if [ ! -b "$DISK" ]; then
  echo "Not a block device: $DISK"
  exit 1
fi

if [ ! -f "$IMAGE" ]; then
  echo "Missing image: $IMAGE"
  echo "Build it first:"
  echo "  sudo ./scripts/make-uefi-usb-image.sh"
  exit 1
fi

case "$DISK" in
  /dev/sd[a-z]|/dev/nvme[0-9]n[0-9]|/dev/mmcblk[0-9]) ;;
  *)
    echo "Refusing suspicious disk path: $DISK"
    echo "Use a whole disk, not a partition, for example /dev/sda"
    exit 1
    ;;
esac

echo "About to overwrite:"
lsblk "$DISK"
echo
echo "Image:"
ls -lh "$IMAGE"
echo
read -r -p "Type YES to overwrite $DISK: " confirm
if [ "$confirm" != "YES" ]; then
  echo "Aborted."
  exit 1
fi

umount "${DISK}"* 2>/dev/null || true

dd if="$IMAGE" of="$DISK" bs=16M status=progress conv=fsync
sync

echo
echo "Flashed $IMAGE to $DISK"
echo "Now eject:"
echo "  sudo eject $DISK"
