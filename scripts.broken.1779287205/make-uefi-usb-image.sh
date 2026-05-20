#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

BUILDROOT_DIR="$ROOT_DIR/build"
IMAGE="$ROOT_DIR/out/bos-usb.img"
WORK="$ROOT_DIR/out/work"

KERNEL="$BUILDROOT_DIR/bzImage"

mkdir -p "$ROOT_DIR/out"
mkdir -p "$WORK"

need() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "Missing dependency: $1"
    exit 1
  }
}

need parted
need mkfs.fat
need mkfs.ext4
need losetup
need mount
need rsync
need bootctl
need blkid

if [ ! -f "$KERNEL" ]; then
  echo "Missing kernel: $KERNEL"
  exit 1
fi

if [ ! -x "$ROOT_DIR/src/bos" ]; then
  echo "Missing BOS binary: $ROOT_DIR/src/bos"
  exit 1
fi

echo "==> Cleaning output"
rm -f "$IMAGE"

echo "==> Creating image"
truncate -s 512M "$IMAGE"

echo "==> Partitioning"
parted -s "$IMAGE" mklabel gpt
parted -s "$IMAGE" mkpart ESP fat32 1MiB 64MiB
parted -s "$IMAGE" set 1 esp on
parted -s "$IMAGE" mkpart BOS_ROOT ext4 64MiB 100%

LOOP="$(sudo losetup --find --partscan --show "$IMAGE")"

EFI_PART="${LOOP}p1"
ROOT_PART="${LOOP}p2"

cleanup() {
  set +e
  sudo umount "$WORK/efi" 2>/dev/null || true
  sudo umount "$WORK/root" 2>/dev/null || true
  sudo losetup -d "$LOOP" 2>/dev/null || true
}
trap cleanup EXIT

echo "==> Formatting"
sudo mkfs.fat -F 32 -n BOS_EFI "$EFI_PART"
sudo mkfs.ext4 -F -L BOS_ROOT "$ROOT_PART"

mkdir -p "$WORK/efi"
mkdir -p "$WORK/root"

echo "==> Mounting"
sudo mount "$EFI_PART" "$WORK/efi"
sudo mount "$ROOT_PART" "$WORK/root"

echo "==> Copying rootfs"
sudo rsync -aHAX "$BUILDROOT_DIR"/ "$WORK/root"/

echo "==> Installing BOS"
sudo install -D -m 0755 \
  "$ROOT_DIR/src/bos" \
  "$WORK/root/usr/local/bin/bos"

echo "==> Writing inittab"

sudo tee "$WORK/root/etc/inittab" >/dev/null <<'INITTAB'
::sysinit:/bin/mount -t proc proc /proc
::sysinit:/bin/mount -t sysfs sysfs /sys
::sysinit:/sbin/mdev -s

tty1::respawn:/bin/sh -c 'exec </dev/tty1 >/dev/tty1 2>&1; dmesg -n 1 2>/dev/null || true; clear; printf "\n\n\n\n\n\n\n\n\n\n"; printf "                    BOS\n"; printf "               -------------\n"; sleep 1; clear; exec /usr/local/bin/bos'

::shutdown:/bin/umount -a -r
INITTAB

ROOT_PARTUUID="$(sudo blkid -s PARTUUID -o value "$ROOT_PART")"

echo "==> Installing systemd-boot"

sudo bootctl install \
  --esp-path="$WORK/efi"

sudo mkdir -p "$WORK/efi/loader/entries"

sudo cp "$KERNEL" "$WORK/efi/bzImage"

sudo tee "$WORK/efi/loader/loader.conf" >/dev/null <<'LOADER'
default bos
timeout 0
console-mode keep
editor no
LOADER

sudo tee "$WORK/efi/loader/entries/bos.conf" >/dev/null <<EOF2
title BOS
linux /bzImage
options root=PARTUUID=${ROOT_PARTUUID} rootwait rw quiet loglevel=0 vt.global_cursor_default=0 console=tty1
EOF2

sudo mkdir -p "$WORK/efi/EFI/BOOT"

if [ -f "$WORK/efi/EFI/systemd/systemd-bootx64.efi" ]; then
  sudo cp \
    "$WORK/efi/EFI/systemd/systemd-bootx64.efi" \
    "$WORK/efi/EFI/BOOT/BOOTX64.EFI"
fi

echo "==> Sync"
sync

echo
echo "Image ready:"
echo "  $IMAGE"
