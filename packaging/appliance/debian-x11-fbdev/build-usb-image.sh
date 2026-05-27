#!/usr/bin/env bash
set -euo pipefail

if [[ "${EUID}" -ne 0 ]]; then
  echo "Run this USB image builder as root." >&2
  exit 1
fi

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd -- "${SCRIPT_DIR}/../../.." && pwd)"

IMAGE_PATH="${1:-${REPO_ROOT}/build/bte-debian-x11-fbdev-amd64.img}"
IMAGE_SIZE="${BTE_IMAGE_SIZE:-8G}"
DEBIAN_SUITE="${BTE_DEBIAN_SUITE:-bookworm}"
DEBIAN_ARCH="${BTE_DEBIAN_ARCH:-amd64}"
DEBIAN_MIRROR="${BTE_DEBIAN_MIRROR:-http://deb.debian.org/debian}"
BINARY_PATH="${BTE_BINARY_PATH:-${REPO_ROOT}/bos}"
BTE_NAMESERVERS="${BTE_NAMESERVERS:-}"
BTE_DEBUG_PASSWORD="${BTE_DEBUG_PASSWORD:-bte}"

required_commands=(
  blkid
  chroot
  debootstrap
  findmnt
  losetup
  mkfs.ext4
  mkfs.vfat
  mount
  partprobe
  parted
  systemctl
  truncate
  umount
)

for command_name in "${required_commands[@]}"; do
  if ! command -v "${command_name}" >/dev/null 2>&1; then
    echo "Missing host command: ${command_name}" >&2
    echo "On Debian/Ubuntu, install: debootstrap parted dosfstools e2fsprogs util-linux" >&2
    exit 1
  fi
done

if [[ ! -x "${BINARY_PATH}" ]]; then
  echo "BTE binary not found or not executable: ${BINARY_PATH}" >&2
  echo "Build it first with: go build -o bos ." >&2
  exit 1
fi

BUILD_DIR="$(mktemp -d /tmp/bte-usb-image.XXXXXX)"
ROOT_MOUNT="${BUILD_DIR}/root"
LOOP_DEVICE=""

cleanup() {
  set +e
  if [[ -n "${ROOT_MOUNT}" && -d "${ROOT_MOUNT}" ]]; then
    for mount_point in \
      "${ROOT_MOUNT}/run" \
      "${ROOT_MOUNT}/sys" \
      "${ROOT_MOUNT}/proc" \
      "${ROOT_MOUNT}/dev/pts" \
      "${ROOT_MOUNT}/dev" \
      "${ROOT_MOUNT}/boot/efi" \
      "${ROOT_MOUNT}"; do
      if findmnt -R "${mount_point}" >/dev/null 2>&1; then
        umount -R "${mount_point}"
      fi
    done
  fi

  if [[ -n "${LOOP_DEVICE}" ]]; then
    losetup -d "${LOOP_DEVICE}" >/dev/null 2>&1 || true
  fi

  rm -rf "${BUILD_DIR}"
}
trap cleanup EXIT

write_image_resolv_conf() {
  local output_path="$1"
  local copied_nameserver=0

  rm -f "${output_path}"
  : >"${output_path}"

  if [[ -n "${BTE_NAMESERVERS}" ]]; then
    for nameserver in ${BTE_NAMESERVERS}; do
      echo "nameserver ${nameserver}" >>"${output_path}"
      copied_nameserver=1
    done
  else
    while read -r keyword nameserver _; do
      if [[ "${keyword}" == "nameserver" && ! "${nameserver}" =~ ^127\. && "${nameserver}" != "::1" ]]; then
        echo "nameserver ${nameserver}" >>"${output_path}"
        copied_nameserver=1
      fi
    done </etc/resolv.conf
  fi

  if [[ "${copied_nameserver}" -eq 0 ]]; then
    cat >>"${output_path}" <<EOF
nameserver 1.1.1.1
nameserver 8.8.8.8
EOF
  fi
}

image_chroot() {
  echo "Running inside image root: $*"
  chroot "${ROOT_MOUNT}" /usr/bin/env -i \
    HOME=/root \
    PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin \
    LANG=C \
    LC_ALL=C \
    BTE_DEBUG_PASSWORD="${BTE_DEBUG_PASSWORD}" \
    "$@"
}

mkdir -p "$(dirname -- "${IMAGE_PATH}")"
rm -f "${IMAGE_PATH}"
truncate -s "${IMAGE_SIZE}" "${IMAGE_PATH}"

parted -s "${IMAGE_PATH}" \
  mklabel gpt \
  mkpart bios_grub 1MiB 3MiB \
  set 1 bios_grub on \
  mkpart esp fat32 3MiB 515MiB \
  set 2 esp on \
  mkpart root ext4 515MiB 100%

LOOP_DEVICE="$(losetup --find --partscan --show "${IMAGE_PATH}")"
partprobe "${LOOP_DEVICE}"
sleep 1

EFI_PART="${LOOP_DEVICE}p2"
ROOT_PART="${LOOP_DEVICE}p3"

mkfs.vfat -F 32 -n BTE_EFI "${EFI_PART}"
mkfs.ext4 -F -L BTE_ROOT "${ROOT_PART}"

mkdir -p "${ROOT_MOUNT}"
mount "${ROOT_PART}" "${ROOT_MOUNT}"
mkdir -p "${ROOT_MOUNT}/boot/efi"
mount "${EFI_PART}" "${ROOT_MOUNT}/boot/efi"

debootstrap \
  --arch="${DEBIAN_ARCH}" \
  --variant=minbase \
  --include=bash,ca-certificates,systemd-sysv,usrmerge \
  "${DEBIAN_SUITE}" \
  "${ROOT_MOUNT}" \
  "${DEBIAN_MIRROR}"

mount --bind /dev "${ROOT_MOUNT}/dev"
mount --bind /dev/pts "${ROOT_MOUNT}/dev/pts"
mount -t proc proc "${ROOT_MOUNT}/proc"
mount -t sysfs sys "${ROOT_MOUNT}/sys"
mount -t tmpfs tmpfs "${ROOT_MOUNT}/run"

write_image_resolv_conf "${ROOT_MOUNT}/etc/resolv.conf"

image_chroot /bin/bash -c '
set -euo pipefail
export DEBIAN_FRONTEND=noninteractive

apt_get() {
  apt-get \
    -o Acquire::Retries=5 \
    -o Acquire::http::Timeout=30 \
    -o Acquire::https::Timeout=30 \
    "$@"
}

sed -i -E "s/^(deb .* )main$/\1main contrib non-free non-free-firmware/" /etc/apt/sources.list
apt_get update
apt_get install -y --no-install-recommends \
  fonts-dejavu-core \
  fonts-jetbrains-mono \
  firmware-amd-graphics \
  firmware-linux-free \
  firmware-misc-nonfree \
  firmware-nvidia-gsp \
  firmware-realtek \
  grub-common \
  grub-efi-amd64-bin \
  grub-pc-bin \
  grub2-common \
  initramfs-tools \
  libdrm-amdgpu1 \
  libegl-mesa0 \
  libgbm1 \
  libpam-systemd \
  libvulkan1 \
  libgl1-mesa-dri \
  linux-image-amd64 \
  linux-headers-amd64 \
  locales \
  mesa-vulkan-drivers \
  openbox \
  sudo \
  systemd-resolved \
  udev \
  x11-xserver-utils \
  xinit \
  xserver-xorg-core \
  xserver-xorg-video-fbdev \
  xserver-xorg-input-all \
  xserver-xorg-legacy \
  xserver-xorg-video-vesa \
  xserver-xorg-video-all \
  xterm \
  zstd
'

ROOT_UUID="$(blkid -s UUID -o value "${ROOT_PART}")"
EFI_UUID="$(blkid -s UUID -o value "${EFI_PART}")"

cat >"${ROOT_MOUNT}/etc/fstab" <<EOF
UUID=${ROOT_UUID} / ext4 defaults,noatime 0 1
UUID=${EFI_UUID} /boot/efi vfat umask=0077 0 1
EOF

cat >"${ROOT_MOUNT}/etc/hostname" <<EOF
bte
EOF

cat >"${ROOT_MOUNT}/etc/hosts" <<EOF
127.0.0.1 localhost
127.0.1.1 bte

::1 localhost ip6-localhost ip6-loopback
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
EOF

mkdir -p "${ROOT_MOUNT}/etc/systemd/network"
cat >"${ROOT_MOUNT}/etc/systemd/network/20-wired.network" <<EOF
[Match]
Name=en* eth*

[Network]
DHCP=yes
EOF

mkdir -p \
  "${ROOT_MOUNT}/etc/bte" \
  "${ROOT_MOUNT}/etc/X11" \
  "${ROOT_MOUNT}/etc/X11/xorg.conf.d" \
  "${ROOT_MOUNT}/etc/sudoers.d" \
  "${ROOT_MOUNT}/etc/systemd/system/getty@tty1.service.d" \
  "${ROOT_MOUNT}/etc/systemd/system/getty@tty2.service.d" \
  "${ROOT_MOUNT}/etc/systemd/system" \
  "${ROOT_MOUNT}/etc/udev/rules.d" \
  "${ROOT_MOUNT}/usr/local/bin" \
  "${ROOT_MOUNT}/var/lib/bte/data" \
  "${ROOT_MOUNT}/var/lib/bte/sql"

install -m 0755 "${BINARY_PATH}" "${ROOT_MOUNT}/usr/local/bin/bte"
install -m 0755 "${SCRIPT_DIR}/bte-launch" "${ROOT_MOUNT}/usr/local/bin/bte-launch"
install -m 0755 "${SCRIPT_DIR}/bte-xsession" "${ROOT_MOUNT}/usr/local/bin/bte-xsession"
install -m 0755 "${SCRIPT_DIR}/bte-console-colors" "${ROOT_MOUNT}/usr/local/bin/bte-console-colors"
install -m 0644 "${SCRIPT_DIR}/bte-bash-profile" "${ROOT_MOUNT}/var/lib/bte/.bash_profile"
install -m 0644 "${SCRIPT_DIR}/ledger.rules" "${ROOT_MOUNT}/etc/udev/rules.d/20-ledger.rules"
install -m 0644 "${REPO_ROOT}/sql/schema.sql" "${ROOT_MOUNT}/var/lib/bte/sql/schema.sql"

cat >"${ROOT_MOUNT}/etc/sudoers.d/90-bte-debug" <<'EOF'
bte ALL=(ALL) NOPASSWD:ALL
EOF
chmod 0440 "${ROOT_MOUNT}/etc/sudoers.d/90-bte-debug"

cat >"${ROOT_MOUNT}/etc/X11/Xwrapper.config" <<'EOF'
allowed_users=console
needs_root_rights=auto
EOF

cat >"${ROOT_MOUNT}/etc/X11/xorg.conf" <<'EOF'
Section "Device"
    Identifier "BTE Framebuffer"
    Driver "fbdev"
    Option "fbdev" "/dev/fb0"
EndSection

Section "Monitor"
    Identifier "BTE Monitor"
EndSection

Section "Screen"
    Identifier "BTE Screen"
    Device "BTE Framebuffer"
    Monitor "BTE Monitor"
    DefaultDepth 24
EndSection

Section "ServerFlags"
    Option "AutoAddGPU" "false"
EndSection
EOF

cat >"${ROOT_MOUNT}/etc/systemd/system/getty@tty1.service.d/autologin.conf" <<'EOF'
[Unit]
After=systemd-modules-load.service

[Service]
ExecStart=
ExecStart=-/sbin/agetty --autologin bte --noclear %I linux
EOF

cat >"${ROOT_MOUNT}/etc/systemd/system/getty@tty2.service.d/autologin.conf" <<'EOF'
[Service]
ExecStart=
ExecStart=-/sbin/agetty --autologin bte --noclear %I linux
EOF

image_chroot /bin/bash -c '
set -euo pipefail
for group_name in adm plugdev input video render systemd-journal; do
  if ! getent group "${group_name}" >/dev/null; then
    groupadd --system "${group_name}"
  fi
done

if ! id bte >/dev/null 2>&1; then
  useradd --create-home --home-dir /var/lib/bte --shell /bin/bash --groups adm,plugdev,input,video,render,systemd-journal bte
else
  usermod --append --groups adm,plugdev,input,video,render,systemd-journal bte
fi

echo "bte:${BTE_DEBUG_PASSWORD}" | chpasswd
chown -R bte:bte /var/lib/bte
sed -i "s/^# *en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/" /etc/locale.gen
locale-gen
cat >/etc/default/locale <<EOF
LANG=en_US.UTF-8
LANGUAGE=en_US:en
LC_ALL=en_US.UTF-8
EOF
'

systemctl --root="${ROOT_MOUNT}" enable systemd-networkd.service
systemctl --root="${ROOT_MOUNT}" disable systemd-networkd-wait-online.service || true
systemctl --root="${ROOT_MOUNT}" enable systemd-resolved.service
systemctl --root="${ROOT_MOUNT}" enable getty@tty1.service
systemctl --root="${ROOT_MOUNT}" enable getty@tty2.service
echo "Configured tty1 autologin launcher for X11 framebuffer BTE."
echo "Configured tty2 autologin recovery console for user bte."

mkdir -p "${ROOT_MOUNT}/boot/grub"
cat >"${ROOT_MOUNT}/etc/default/grub" <<'EOF'
GRUB_DEFAULT=0
GRUB_TIMEOUT=1
GRUB_DISTRIBUTOR="BTE"
GRUB_CMDLINE_LINUX_DEFAULT="quiet nomodeset"
GRUB_CMDLINE_LINUX=""
EOF

image_chroot update-initramfs -u
image_chroot update-grub
image_chroot grub-install --target=x86_64-efi --efi-directory=/boot/efi --bootloader-id=BTE --removable --recheck
image_chroot grub-install --target=i386-pc --recheck "${LOOP_DEVICE}"

image_chroot apt-get clean
rm -rf "${ROOT_MOUNT}/var/lib/apt/lists/"*
rm -f "${ROOT_MOUNT}/etc/machine-id"
touch "${ROOT_MOUNT}/etc/machine-id"
rm -f "${ROOT_MOUNT}/etc/resolv.conf"
ln -s /run/systemd/resolve/stub-resolv.conf "${ROOT_MOUNT}/etc/resolv.conf"

sync
if [[ -n "${SUDO_UID:-}" && -n "${SUDO_GID:-}" ]]; then
  chown "${SUDO_UID}:${SUDO_GID}" "${IMAGE_PATH}"
fi
echo "Created flashable USB image: ${IMAGE_PATH}"
echo "Write it to USB with: sudo dd if=${IMAGE_PATH} of=/dev/sdX bs=4M status=progress conv=fsync"
