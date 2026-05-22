#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILDROOT_DIR="$ROOT_DIR/build"
KERNEL="$BUILDROOT_DIR/bzImage"
BOS_BIN="$ROOT_DIR/src/bos"

OUT_DIR="$ROOT_DIR/out"
IMAGE="$OUT_DIR/bos-usb.img"
WORK="$OUT_DIR/work"
ROOT_MNT="$WORK/root"
EFI_MNT="$WORK/efi"

IMAGE_SIZE="${IMAGE_SIZE:-1024M}"
EFI_SIZE_MIB="${EFI_SIZE_MIB:-256}"
LOOP=""

cleanup() {
  set +e
  sync
  mountpoint -q "$EFI_MNT" && umount "$EFI_MNT"
  mountpoint -q "$ROOT_MNT" && umount "$ROOT_MNT"
  [ -n "$LOOP" ] && losetup -d "$LOOP" 2>/dev/null || true
}
trap cleanup EXIT

if [ "$EUID" -ne 0 ]; then
  echo "Run as root: sudo $0"
  exit 1
fi

need() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "Missing command: $1"
    echo "Install: sudo apt install -y parted dosfstools e2fsprogs rsync util-linux systemd-boot-efi"
    exit 1
  }
}

need parted
need losetup
need mkfs.vfat
need mkfs.ext4
need rsync
need blkid
need file

if [ ! -f "$KERNEL" ]; then
  echo "Missing kernel: $KERNEL"
  exit 1
fi

if [ ! -x "$BOS_BIN" ]; then
  echo "Missing executable BOS binary: $BOS_BIN"
  echo "Build it first: ./scripts/build-local-bos.sh"
  exit 1
fi

if [ ! -d "$BUILDROOT_DIR/bin" ] || [ ! -d "$BUILDROOT_DIR/etc" ] || [ ! -d "$BUILDROOT_DIR/sbin" ]; then
  echo "The build/ directory does not look like an extracted Buildroot rootfs."
  echo "Expected build/bin, build/etc, build/sbin."
  exit 1
fi

BOOTLOADER=""
for p in \
  /usr/lib/systemd/boot/efi/systemd-bootx64.efi \
  /usr/lib/systemd/boot/efi/systemd-bootx64.efi.signed
 do
  if [ -f "$p" ]; then
    BOOTLOADER="$p"
    break
  fi
done

if [ -z "$BOOTLOADER" ]; then
  echo "Could not locate systemd-boot EFI binary."
  echo "Install: sudo apt install -y systemd-boot-efi"
  exit 1
fi

rm -rf "$WORK"
mkdir -p "$OUT_DIR" "$ROOT_MNT" "$EFI_MNT"
rm -f "$IMAGE"
truncate -s "$IMAGE_SIZE" "$IMAGE"

parted -s "$IMAGE" mklabel gpt
parted -s "$IMAGE" mkpart ESP fat32 1MiB "$((EFI_SIZE_MIB + 1))MiB"
parted -s "$IMAGE" set 1 esp on
parted -s "$IMAGE" mkpart BOS_ROOT ext4 "$((EFI_SIZE_MIB + 1))MiB" 100%

LOOP="$(losetup --find --partscan --show "$IMAGE")"
sleep 1
EFI_PART="${LOOP}p1"
ROOT_PART="${LOOP}p2"

if [ ! -b "$EFI_PART" ] || [ ! -b "$ROOT_PART" ]; then
  echo "Loop partitions missing"
  lsblk "$LOOP" || true
  exit 1
fi

mkfs.vfat -F 32 -n BOS_EFI "$EFI_PART"
mkfs.ext4 -F -L BOS_ROOT "$ROOT_PART"

mount "$ROOT_PART" "$ROOT_MNT"
mount "$EFI_PART" "$EFI_MNT"

rsync -aHAX --numeric-ids \
  --exclude='/dev/*' \
  --exclude='/proc/*' \
  --exclude='/sys/*' \
  --exclude='/tmp/*' \
  --exclude='/run/*' \
  "$BUILDROOT_DIR/" "$ROOT_MNT/"

ROOT_PARTUUID="$(blkid -s PARTUUID -o value "$ROOT_PART")"
if [ -z "$ROOT_PARTUUID" ]; then
  echo "Could not read root PARTUUID"
  exit 1
fi

install -D -m 0755 "$BOS_BIN" "$ROOT_MNT/usr/local/bin/bos"

cat > "$ROOT_MNT/etc/inittab" <<'INITTAB'
::sysinit:/bin/mount -t proc proc /proc
::sysinit:/bin/mount -t sysfs sysfs /sys
::sysinit:/bin/mount -t devtmpfs devtmpfs /dev
::sysinit:/bin/mkdir -p /dev/pts /dev/shm /run /tmp /var/log
::sysinit:/bin/mount -t devpts devpts /dev/pts -o gid=5,mode=620
::sysinit:/bin/ln -sf /dev/pts/ptmx /dev/ptmx
::sysinit:/bin/mount -t tmpfs tmpfs /run
::sysinit:/bin/mount -t tmpfs tmpfs /tmp
::sysinit:/bin/mount -t tmpfs tmpfs /var/log
::sysinit:/bin/hostname -F /etc/hostname
::sysinit:/sbin/mdev -s
::sysinit:/sbin/ifconfig lo up
::sysinit:/sbin/ifconfig eth0 up
::sysinit:/sbin/udhcpc -i eth0 -s /usr/share/udhcpc/default.script
::sysinit:/etc/init.d/rcS

console::respawn:/usr/local/bin/start-bos-session

::shutdown:/etc/init.d/rcK
::shutdown:/bin/umount -a -r
INITTAB

mkdir -p "$EFI_MNT/EFI/BOOT" "$EFI_MNT/loader/entries"
install -D -m 0644 "$BOOTLOADER" "$EFI_MNT/EFI/BOOT/BOOTX64.EFI"
install -D -m 0644 "$KERNEL" "$EFI_MNT/bzImage"

cat > "$EFI_MNT/loader/loader.conf" <<'LOADER'
default bos
timeout 3
console-mode max
editor yes
LOADER

cat > "$EFI_MNT/loader/entries/bos.conf" <<ENTRY
title BOS
linux /bzImage
options root=PARTUUID=$ROOT_PARTUUID rw rootwait console=ttyS0,115200 console=tty1 loglevel=7 ignore_loglevel
ENTRY

sync

# hard verification before unmounting
test -f "$EFI_MNT/EFI/BOOT/BOOTX64.EFI"
test -f "$EFI_MNT/bzImage"
test -f "$EFI_MNT/loader/loader.conf"
test -f "$EFI_MNT/loader/entries/bos.conf"
test -x "$ROOT_MNT/usr/local/bin/bos"
test -x "$ROOT_MNT/usr/local/bin/start-bos-session"

echo "==> ESP contents"
find "$EFI_MNT" -maxdepth 5 -type f -printf '%P\n' | sort

echo "==> BOOTX64.EFI type"
file "$EFI_MNT/EFI/BOOT/BOOTX64.EFI"

echo "==> bos.conf"
cat "$EFI_MNT/loader/entries/bos.conf"

echo
 echo "DONE: $IMAGE"
