#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="$ROOT_DIR/out/bos-usb.img"

if [ ! -f "$IMAGE" ]; then
  echo "Missing image: $IMAGE"
  echo "Build it first: sudo IMAGE_SIZE=1024M EFI_SIZE_MIB=256 ./scripts/make-uefi-usb-image.sh"
  exit 1
fi

OVMF_CODE=""
OVMF_VARS_SRC=""
for p in /usr/share/OVMF/OVMF_CODE.fd /usr/share/OVMF/OVMF_CODE_4M.fd; do
  [ -f "$p" ] && OVMF_CODE="$p" && break
done
for p in /usr/share/OVMF/OVMF_VARS.fd /usr/share/OVMF/OVMF_VARS_4M.fd; do
  [ -f "$p" ] && OVMF_VARS_SRC="$p" && break
done

if [ -z "$OVMF_CODE" ] || [ -z "$OVMF_VARS_SRC" ]; then
  echo "Missing OVMF firmware. Install: sudo apt install -y ovmf"
  exit 1
fi

OVMF_VARS="/tmp/BOS_OVMF_VARS.fd"
cp "$OVMF_VARS_SRC" "$OVMF_VARS"
chmod 0644 "$OVMF_VARS"

qemu-system-x86_64 \
  -enable-kvm \
  -m 4096 \
  -cpu host \
  -drive if=pflash,format=raw,readonly=on,file="$OVMF_CODE" \
  -drive if=pflash,format=raw,file="$OVMF_VARS" \
  -display gtk,gl=off \
  -device virtio-vga \
  -device usb-ehci,id=ehci \
  -device usb-kbd,bus=ehci.0 \
  -device usb-mouse,bus=ehci.0 \
  -drive file="$IMAGE",format=raw,if=ide \
  -serial stdio \
  -no-reboot
