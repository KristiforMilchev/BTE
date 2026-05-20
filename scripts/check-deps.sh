\
#!/usr/bin/env bash
set -euo pipefail

missing=0

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing: $1"
    missing=1
  fi
}

need parted
need losetup
need mkfs.vfat
need mkfs.ext4
need rsync
need grub-install
need blkid
need qemu-system-x86_64

if [ "$missing" -ne 0 ]; then
  echo
  echo "Install dependencies on Debian:"
  echo "  sudo apt update"
  echo "  sudo apt install -y parted dosfstools e2fsprogs rsync grub-efi-amd64-bin qemu-system-x86 ovmf util-linux"
  exit 1
fi

echo "OK: all required tools are available"
