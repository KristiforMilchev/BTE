#!/usr/bin/env bash
set -euo pipefail

PROJECT="/home/kristifor/Projects/ledger-test"
BOS_SRC="$PROJECT/src"
OUT_DIR="$PROJECT/out"
IMAGE="$OUT_DIR/bos-debian-appliance.img"

IMAGE_SIZE="${IMAGE_SIZE:-4G}"
DEBIAN_RELEASE="${DEBIAN_RELEASE:-bookworm}"
MIRROR="${MIRROR:-http://deb.debian.org/debian}"

WORK="$OUT_DIR/debian-work"
ROOT_MNT="$WORK/root"
EFI_MNT="$WORK/efi"

LOOP=""

cleanup() {
  set +e
  sync
  mountpoint -q "$EFI_MNT" && umount "$EFI_MNT"
  mountpoint -q "$ROOT_MNT/dev/pts" && umount "$ROOT_MNT/dev/pts"
  mountpoint -q "$ROOT_MNT/dev" && umount "$ROOT_MNT/dev"
  mountpoint -q "$ROOT_MNT/proc" && umount "$ROOT_MNT/proc"
  mountpoint -q "$ROOT_MNT/sys" && umount "$ROOT_MNT/sys"
  mountpoint -q "$ROOT_MNT" && umount "$ROOT_MNT"
  [ -n "$LOOP" ] && losetup -d "$LOOP" 2>/dev/null || true
}
trap cleanup EXIT

if [ "$EUID" -ne 0 ]; then
  echo "Run as root:"
  echo "  sudo $0"
  exit 1
fi

need() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "Missing: $1"
    echo "Install: sudo apt install -y mmdebstrap parted dosfstools e2fsprogs util-linux rsync golang-go"
    exit 1
  }
}

need mmdebstrap
need parted
need losetup
need mkfs.vfat
need mkfs.ext4
need rsync

GO_BIN="$(command -v go || true)"
if [ -z "$GO_BIN" ]; then
  for p in \
    /home/kristifor/software/go/bin/go \
    /usr/local/go/bin/go \
    /usr/bin/go
  do
    [ -x "$p" ] && GO_BIN="$p" && break
  done
fi

if [ -z "$GO_BIN" ]; then
  echo "Go not found"
  exit 1
fi

echo "==> Building BOS binary"

mkdir -p "$OUT_DIR"
cd "$BOS_SRC"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  "$GO_BIN" build -trimpath -ldflags="-s -w" \
  -o "$OUT_DIR/bos" .

chmod 0755 "$OUT_DIR/bos"

echo "==> Creating image: $IMAGE"

rm -rf "$WORK"
mkdir -p "$ROOT_MNT" "$EFI_MNT"
rm -f "$IMAGE"

truncate -s "$IMAGE_SIZE" "$IMAGE"

parted -s "$IMAGE" mklabel gpt
parted -s "$IMAGE" mkpart ESP fat32 1MiB 513MiB
parted -s "$IMAGE" set 1 esp on
parted -s "$IMAGE" mkpart BOS_ROOT ext4 513MiB 100%

LOOP="$(losetup --find --partscan --show "$IMAGE")"
sleep 1

EFI_PART="${LOOP}p1"
ROOT_PART="${LOOP}p2"

mkfs.vfat -F32 -n BOS_EFI "$EFI_PART"
mkfs.ext4 -F -L BOS_ROOT "$ROOT_PART"

mount "$ROOT_PART" "$ROOT_MNT"
mkdir -p "$ROOT_MNT/boot/efi"
mount "$EFI_PART" "$EFI_MNT"

ROOT_PARTUUID="$(blkid -s PARTUUID -o value "$ROOT_PART")"

echo "==> Installing minimal Debian rootfs"

DEBROOT="$WORK/debroot"
rm -rf "$DEBROOT"
mkdir -p "$DEBROOT"

mmdebstrap \
  --variant=minbase \
  --components=main,contrib,non-free,non-free-firmware \
  --include=systemd-sysv,linux-image-amd64,initramfs-tools,systemd-boot,locales,dbus-user-session,libpam-systemd,udev,kmod,seatd,cage,foot,fonts-dejavu-core,firmware-amd-graphics,firmware-linux-free,ca-certificates,sudo,bash,login,passwd \
  "$DEBIAN_RELEASE" \
  "$DEBROOT" \
  "$MIRROR"

rsync -aHAX --numeric-ids "$DEBROOT/" "$ROOT_MNT/"
rm -rf "$DEBROOT"

echo "==> Configuring appliance rootfs"

mount -t proc proc "$ROOT_MNT/proc"
mount -t sysfs sysfs "$ROOT_MNT/sys"
mount --bind /dev "$ROOT_MNT/dev"
mount -t devpts devpts "$ROOT_MNT/dev/pts"

cat > "$ROOT_MNT/etc/fstab" <<EOF
PARTUUID=$ROOT_PARTUUID / ext4 defaults,noatime 0 1
EOF

echo "bos" > "$ROOT_MNT/etc/hostname"

cat > "$ROOT_MNT/etc/locale.gen" <<'EOF'
en_US.UTF-8 UTF-8
EOF

chroot "$ROOT_MNT" locale-gen
chroot "$ROOT_MNT" update-locale LANG=en_US.UTF-8 LC_ALL=en_US.UTF-8

install -D -m 0755 "$OUT_DIR/bos" "$ROOT_MNT/usr/local/bin/bos"

chroot "$ROOT_MNT" useradd -m -s /bin/bash bos || true
chroot "$ROOT_MNT" passwd -d bos || true

chroot "$ROOT_MNT" usermod -aG video,input,render bos || true
chroot "$ROOT_MNT" getent group seat >/dev/null 2>&1 && chroot "$ROOT_MNT" usermod -aG seat bos || true

mkdir -p "$ROOT_MNT/etc/systemd/system/getty@tty1.service.d"

cat > "$ROOT_MNT/etc/systemd/system/getty@tty1.service.d/override.conf" <<'EOF'
[Service]
ExecStart=
ExecStart=-/sbin/agetty --autologin bos --noclear %I $TERM
EOF

cat > "$ROOT_MNT/home/bos/.bash_profile" <<'EOF'
if [ "$(tty)" = "/dev/tty1" ] && [ -z "${WAYLAND_DISPLAY:-}" ]; then
  export LANG=en_US.UTF-8
  export LC_ALL=en_US.UTF-8
  export XDG_RUNTIME_DIR="/run/user/$(id -u)"
  export WLR_BACKENDS=drm
  export TERM=foot
  exec cage -- foot -F -e /usr/local/bin/bos
fi
EOF

chroot "$ROOT_MNT" chown bos:bos /home/bos/.bash_profile
chroot "$ROOT_MNT" chmod 0644 /home/bos/.bash_profile

chroot "$ROOT_MNT" systemctl enable getty@tty1.service
chroot "$ROOT_MNT" systemctl enable seatd || true

echo "==> Installing systemd-boot fallback loader"

mkdir -p "$EFI_MNT/EFI/BOOT"
mkdir -p "$EFI_MNT/loader/entries"

BOOTLOADER="$ROOT_MNT/usr/lib/systemd/boot/efi/systemd-bootx64.efi"

if [ ! -f "$BOOTLOADER" ]; then
  echo "Missing systemd-boot EFI binary: $BOOTLOADER"
  exit 1
fi

install -D -m 0644 "$BOOTLOADER" "$EFI_MNT/EFI/BOOT/BOOTX64.EFI"

KERNEL="$(ls "$ROOT_MNT"/boot/vmlinuz-* | head -n1)"
INITRD="$(ls "$ROOT_MNT"/boot/initrd.img-* | head -n1)"

install -D -m 0644 "$KERNEL" "$EFI_MNT/vmlinuz"
install -D -m 0644 "$INITRD" "$EFI_MNT/initrd.img"

cat > "$EFI_MNT/loader/loader.conf" <<'EOF'
default bos
timeout 0
console-mode max
editor no
EOF

cat > "$EFI_MNT/loader/entries/bos.conf" <<EOF
title BOS
linux /vmlinuz
initrd /initrd.img
options root=PARTUUID=$ROOT_PARTUUID rw rootwait quiet loglevel=3
EOF

echo "==> Verifying image"

find "$EFI_MNT" -maxdepth 4 -type f -print
test -f "$EFI_MNT/EFI/BOOT/BOOTX64.EFI"
test -f "$EFI_MNT/vmlinuz"
test -f "$EFI_MNT/initrd.img"
test -f "$EFI_MNT/loader/entries/bos.conf"
test -x "$ROOT_MNT/usr/local/bin/bos"

sync

echo
echo "DONE:"
echo "  $IMAGE"
echo
echo "Test with:"
echo "  qemu-system-x86_64 \\"
echo "    -enable-kvm \\"
echo "    -m 4096 \\"
echo "    -cpu host \\"
echo "    -drive if=pflash,format=raw,readonly=on,file=/usr/share/OVMF/OVMF_CODE.fd \\"
echo "    -drive if=pflash,format=raw,file=/tmp/BOS_OVMF_VARS.fd \\"
echo "    -display gtk,gl=off \\"
echo "    -device virtio-vga \\"
echo "    -device usb-ehci,id=ehci \\"
echo "    -device usb-kbd,bus=ehci.0 \\"
echo "    -device usb-mouse,bus=ehci.0 \\"
echo "    -drive file=$IMAGE,format=raw,if=ide"
echo
echo "Flash with:"
echo "  sudo dd if=$IMAGE of=/dev/sdX bs=4M status=progress oflag=sync"

