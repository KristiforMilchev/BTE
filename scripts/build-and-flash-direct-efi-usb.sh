#!/usr/bin/env bash
set -euo pipefail

USB="${1:-}"

if [ -z "$USB" ]; then
  echo "Usage:"
  echo "  sudo -E $0 /dev/sdX"
  exit 1
fi

if [ "$EUID" -ne 0 ]; then
  echo "Run as root:"
  echo "  sudo -E $0 $USB"
  exit 1
fi

case "$USB" in
  /dev/sd[a-z]|/dev/nvme[0-9]n[0-9]|/dev/mmcblk[0-9]) ;;
  *)
    echo "Refusing suspicious disk path: $USB"
    echo "Use a whole disk, not a partition, for example /dev/sda"
    exit 1
    ;;
esac

PROJECT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILDROOT="$PROJECT/external/buildroot"
BOS_SRC="$PROJECT/src"
OVERLAY="$PROJECT/board/bos/rootfs-overlay"
KERNEL="$BUILDROOT/output/images/bzImage"

find_go() {
  if command -v go >/dev/null 2>&1; then
    command -v go
    return 0
  fi

  for p in \
    "/home/${SUDO_USER:-kristifor}/software/go/bin/go" \
    "/home/${SUDO_USER:-kristifor}/go/bin/go" \
    "/usr/local/go/bin/go" \
    "/usr/bin/go"
  do
    if [ -x "$p" ]; then
      echo "$p"
      return 0
    fi
  done

  return 1
}

GO_BIN="$(find_go || true)"

mkdir -p "$OVERLAY/usr/local/bin"

if [ -n "$GO_BIN" ] && [ -f "$BOS_SRC/go.mod" ]; then
  echo "==> Building BOS binary with $GO_BIN"
  cd "$BOS_SRC"
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 "$GO_BIN" build -trimpath -ldflags="-s -w" -o "$OVERLAY/usr/local/bin/bos" .
elif [ -x "$BOS_SRC/bos" ]; then
  echo "==> Reusing existing BOS binary: $BOS_SRC/bos"
  cp "$BOS_SRC/bos" "$OVERLAY/usr/local/bin/bos"
else
  echo "Could not build or find BOS binary."
  echo "Expected Go at /home/$SUDO_USER/software/go/bin/go or executable $BOS_SRC/bos"
  exit 1
fi
chmod 0755 "$OVERLAY/usr/local/bin/bos"

cat > "$OVERLAY/usr/local/bin/start-bos-session" <<'EOS'
#!/bin/sh

mount -t proc proc /proc || true
mount -t sysfs sysfs /sys || true
mount -t devtmpfs devtmpfs /dev || true

mkdir -p /dev/pts
mount -t devpts devpts /dev/pts -o gid=5,mode=620 || true
ln -sf /dev/pts/ptmx /dev/ptmx

mount -t tmpfs tmpfs /run || true

mkdir -p /run/user/0 /run/seatd
chmod 700 /run/user/0
chmod 755 /run/seatd

export XDG_RUNTIME_DIR=/run/user/0
export WLR_BACKENDS=drm
export WLR_RENDERER=pixman
export TERM=foot
export COLORTERM=truecolor
export HOME=/root
export USER=root
export LANG=en_US.UTF-8
export LC_ALL=en_US.UTF-8

echo "BOS INIT: starting seatd" > /dev/ttyS0 2>/dev/null || true
seatd -g root >> /dev/ttyS0 2>&1 &

sleep 1

echo "BOS INIT: starting cage/foot/bos on tty1" > /dev/ttyS0 2>/dev/null || true

setsid sh -c '
  exec </dev/tty1 >/dev/tty1 2>/dev/ttyS0
  cage -- foot -F -e /usr/local/bin/bos
  echo "CAGE EXITED: $?" > /dev/ttyS0
' &

while true; do
  sleep 3600
done
EOS
chmod 0755 "$OVERLAY/usr/local/bin/start-bos-session"

if [ -d "$BUILDROOT" ]; then
  echo "==> Building Buildroot"
  cd "$BUILDROOT"
  make
else
  echo "Missing Buildroot directory: $BUILDROOT"
  exit 1
fi

if [ ! -f "$KERNEL" ]; then
  echo "Missing kernel: $KERNEL"
  exit 1
fi

if ! file "$KERNEL" | grep -qi 'EFI\|Linux kernel'; then
  echo "Kernel file does not look valid: $KERNEL"
  file "$KERNEL"
  exit 1
fi

echo "==> Flashing direct EFI-stub USB to $USB"
echo "WARNING: this erases $USB"
lsblk "$USB" || true
sleep 3

umount "${USB}"* 2>/dev/null || true
wipefs -a "$USB"
partprobe "$USB" 2>/dev/null || true
sleep 1

parted -s "$USB" mklabel gpt
parted -s "$USB" mkpart ESP fat32 1MiB 512MiB
parted -s "$USB" set 1 esp on
partprobe "$USB" 2>/dev/null || true
sleep 1

PART="${USB}1"
if [[ "$USB" == /dev/nvme* || "$USB" == /dev/mmcblk* ]]; then
  PART="${USB}p1"
fi

mkfs.vfat -F32 -n BOS "$PART"

mkdir -p /mnt/bos-usb
mount "$PART" /mnt/bos-usb
rm -rf /mnt/bos-usb/*
mkdir -p /mnt/bos-usb/EFI/BOOT
cp "$KERNEL" /mnt/bos-usb/EFI/BOOT/BOOTX64.EFI
sync
umount /mnt/bos-usb

echo "DONE."
echo "Boot path: UEFI -> EFI/BOOT/BOOTX64.EFI -> Linux -> start-bos-session -> seatd -> cage -> foot -> bos"
