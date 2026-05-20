\
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

IMAGE_SIZE="${IMAGE_SIZE:-512M}"
EFI_SIZE_MIB="${EFI_SIZE_MIB:-64}"

LOOP=""

cleanup() {
  set +e
  if mountpoint -q "$EFI_MNT"; then sudo umount "$EFI_MNT"; fi
  if mountpoint -q "$ROOT_MNT"; then sudo umount "$ROOT_MNT"; fi
  if [ -n "$LOOP" ]; then sudo losetup -d "$LOOP" 2>/dev/null || true; fi
}
trap cleanup EXIT

if [ "$EUID" -ne 0 ]; then
  echo "Run as root:"
  echo "  sudo $0"
  exit 1
fi

need() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "Missing command: $1"
    echo "Install:"
    echo "  sudo apt install -y parted dosfstools e2fsprogs rsync grub-efi-amd64-bin util-linux"
    exit 1
  }
}

need parted
need losetup
need mkfs.vfat
need mkfs.ext4
need rsync
need grub-install
need blkid

if [ ! -f "$KERNEL" ]; then
  echo "Missing kernel: $KERNEL"
  exit 1
fi

if [ ! -x "$BOS_BIN" ]; then
  echo "Missing executable BOS binary: $BOS_BIN"
  echo "Build it first:"
  echo "  ./scripts/build-local-bos.sh"
  exit 1
fi

if [ ! -d "$BUILDROOT_DIR/bin" ] || [ ! -d "$BUILDROOT_DIR/etc" ] || [ ! -d "$BUILDROOT_DIR/sbin" ]; then
  echo "The build/ directory does not look like an extracted Buildroot rootfs."
  echo "Expected at least build/bin, build/etc, build/sbin."
  exit 1
fi

echo "==> Cleaning output"
rm -rf "$WORK"
mkdir -p "$OUT_DIR" "$ROOT_MNT" "$EFI_MNT"

echo "==> Creating disk image: $IMAGE"
rm -f "$IMAGE"
truncate -s "$IMAGE_SIZE" "$IMAGE"

echo "==> Partitioning GPT image"
parted -s "$IMAGE" mklabel gpt
parted -s "$IMAGE" mkpart ESP fat32 1MiB "$((EFI_SIZE_MIB + 1))MiB"
parted -s "$IMAGE" set 1 esp on
parted -s "$IMAGE" mkpart BOS_ROOT ext4 "$((EFI_SIZE_MIB + 1))MiB" 100%

LOOP="$(losetup --find --partscan --show "$IMAGE")"
sleep 1

EFI_PART="${LOOP}p1"
ROOT_PART="${LOOP}p2"

if [ ! -b "$EFI_PART" ] || [ ! -b "$ROOT_PART" ]; then
  echo "Loop partitions were not created correctly."
  lsblk "$LOOP" || true
  exit 1
fi

echo "==> Formatting partitions"
mkfs.vfat -F 32 -n BOS_EFI "$EFI_PART"
mkfs.ext4 -F -L BOS_ROOT "$ROOT_PART"

echo "==> Mounting partitions"
mount "$ROOT_PART" "$ROOT_MNT"
mkdir -p "$ROOT_MNT/boot/efi"
mount "$EFI_PART" "$EFI_MNT"

echo "==> Copying Buildroot rootfs"
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

echo "==> Overriding inittab for BOS appliance"
cat > "$ROOT_MNT/etc/inittab" <<'EOF'
::sysinit:/bin/mount -t proc proc /proc
::sysinit:/bin/mount -t sysfs sysfs /sys
::sysinit:/bin/mount -t devtmpfs devtmpfs /dev
::sysinit:/bin/mkdir -p /dev/pts /dev/shm /run /tmp /var/log
::sysinit:/bin/mount -t devpts devpts /dev/pts
::sysinit:/bin/mount -t tmpfs tmpfs /run
::sysinit:/bin/mount -t tmpfs tmpfs /tmp
::sysinit:/bin/mount -t tmpfs tmpfs /var/log
::sysinit:/bin/hostname -F /etc/hostname
::sysinit:/sbin/mdev -s
::sysinit:/sbin/ifconfig lo up
::sysinit:/sbin/ifconfig eth0 up
::sysinit:/sbin/udhcpc -i eth0 -s /usr/share/udhcpc/default.script
::sysinit:/etc/init.d/rcS

console::respawn:/usr/local/bin/bos

::shutdown:/etc/init.d/rcK
::shutdown:/bin/umount -a -r
EOF

echo "==> Installing BOS binary"
install -D -m 0755 "$BOS_BIN" "$ROOT_MNT/usr/local/bin/bos"

echo "==> Installing direct UKI bootloader, no GRUB"

ROOT_PARTUUID="$(blkid -s PARTUUID -o value "$ROOT_PART")"
if [ -z "$ROOT_PARTUUID" ]; then
  echo "Could not read root PARTUUID"
  exit 1
fi

STUB=""
for candidate in \
  /usr/lib/systemd/boot/efi/linuxx64.efi.stub \
  /lib/systemd/boot/efi/linuxx64.efi.stub
do
  if [ -f "$candidate" ]; then
    STUB="$candidate"
    break
  fi
done

if [ -z "$STUB" ]; then
  echo "Missing Linux EFI stub."
  echo "Install package providing linuxx64.efi.stub."
  exit 1
fi

mkdir -p "$EFI_MNT/EFI/BOOT"
mkdir -p "$WORK/uki"

cat > "$WORK/uki/cmdline" <<EOF
root=PARTUUID=${ROOT_PARTUUID} rootwait rw quiet loglevel=0 vt.global_cursor_default=0 console=tty1
EOF

cat > "$WORK/uki/os-release" <<'EOF'
ID=bos
NAME="BOS Appliance"
PRETTY_NAME="BOS Appliance"
EOF

objcopy \
  --add-section .osrel="$WORK/uki/os-release" --change-section-vma .osrel=0x20000 \
  --add-section .cmdline="$WORK/uki/cmdline" --change-section-vma .cmdline=0x30000 \
  --add-section .linux="$KERNEL" --change-section-vma .linux=0x2000000 \
  "$STUB" "$EFI_MNT/EFI/BOOT/BOOTX64.EFI"

echo "==> Final sync"
sync

echo
echo "DONE:"
echo "  $IMAGE"
echo
echo "Test with:"
echo "  ./scripts/run-uefi-qemu.sh"
echo
echo "Flash with:"
echo "  sudo ./scripts/flash-usb.sh /dev/sdX"
