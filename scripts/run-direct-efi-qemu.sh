#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILDROOT_DIR="$ROOT_DIR/external/buildroot"
KERNEL="$BUILDROOT_DIR/output/images/bzImage"
ESP="/tmp/bos-esp"

if [ ! -f "$KERNEL" ]; then
  KERNEL="$ROOT_DIR/build/bzImage"
fi

if [ ! -f "$KERNEL" ]; then
  echo "Missing kernel bzImage. Build Buildroot first."
  exit 1
fi

rm -rf "$ESP"
mkdir -p "$ESP/EFI/BOOT"
cp "$KERNEL" "$ESP/EFI/BOOT/BOOTX64.EFI"

qemu-system-x86_64 \
  -enable-kvm \
  -m 2048 \
  -cpu host \
  -bios /usr/share/OVMF/OVMF_CODE.fd \
  -display gtk,gl=off \
  -device virtio-vga \
  -device usb-ehci,id=ehci \
  -device usb-kbd,bus=ehci.0 \
  -device usb-mouse,bus=ehci.0 \
  -drive file=fat:rw:"$ESP",format=raw,media=disk \
  -serial stdio
