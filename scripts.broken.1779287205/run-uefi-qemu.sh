#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="$ROOT_DIR/out/bos-usb.img"

if [ ! -f "$IMAGE" ]; then
  echo "Missing image: $IMAGE"
  echo "Build it first:"
  echo "  sudo ./scripts/make-uefi-usb-image.sh"
  exit 1
fi

OVMF_CODE=""
for p in \
  /usr/share/OVMF/OVMF_CODE.fd \
  /usr/share/ovmf/OVMF.fd \
  /usr/share/qemu/OVMF.fd
do
  if [ -f "$p" ]; then
    OVMF_CODE="$p"
    break
  fi
done

if [ -z "$OVMF_CODE" ]; then
  echo "Missing OVMF firmware. Install:"
  echo "  sudo apt install ovmf"
  exit 1
fi

# Use usb-storage here because the real target is a USB stick.
# This catches more USB-root boot problems than virtio does.
qemu-system-x86_64 \
  -machine q35,accel=kvm:tcg \
  -m 1024 \
  -vga std \
  -drive if=pflash,format=raw,readonly=on,file="$OVMF_CODE" \
  -drive id=bosdisk,if=none,format=raw,file="$IMAGE" \
  -device qemu-xhci,id=xhci \
  -device usb-host,vendorid=0x2c97,productid=0x4000 \

  -device usb-storage,drive=bosdisk,bootindex=1 \
  -display gtk \
  -serial stdio \
  -no-reboot
