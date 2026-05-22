#!/usr/bin/env bash
set -euo pipefail
USB="${1:-}"
if [ -z "$USB" ]; then
  echo "Usage: sudo $0 /dev/sdX"
  exit 1
fi
PART="${USB}1"
if [[ "$USB" == /dev/nvme* || "$USB" == /dev/mmcblk* ]]; then
  PART="${USB}p1"
fi
mkdir -p /mnt/bos-usb-inspect
mount "$PART" /mnt/bos-usb-inspect
find /mnt/bos-usb-inspect -maxdepth 5 -type f -print
file /mnt/bos-usb-inspect/EFI/BOOT/BOOTX64.EFI || true
umount /mnt/bos-usb-inspect
