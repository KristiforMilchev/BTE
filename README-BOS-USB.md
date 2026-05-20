\
# BOS Buildroot USB Scripts

These scripts are fresh and do not depend on the old Alpine scripts.

They use the project layout from the archive:

```txt
build/              extracted Buildroot rootfs
build/bzImage       Buildroot kernel
src/                Go BOS project
src/bos             BOS binary
```

Target boot mode:

```txt
UEFI
Secure Boot OFF
GRUB EFI
Buildroot kernel
Buildroot rootfs
BOS on console
```

## Install host dependencies

```bash
sudo apt update
sudo apt install -y \
  parted dosfstools e2fsprogs rsync grub-efi-amd64-bin \
  qemu-system-x86 ovmf util-linux
```

Make sure `/sbin` and `/usr/sbin` are in your zsh path:

```bash
echo 'export PATH="$PATH:/sbin:/usr/sbin"' >> ~/.zshrc
source ~/.zshrc
```

## Build BOS binary

If `src/bos` already exists and works, skip this.

```bash
./scripts/build-local-bos.sh
```

## Create bootable UEFI image

```bash
sudo ./scripts/make-uefi-usb-image.sh
```

Output:

```txt
out/bos-usb.img
```

## Test with QEMU UEFI USB emulation

```bash
./scripts/run-uefi-qemu.sh
```

This intentionally boots the image as USB storage, not virtio.

## Flash real USB

Find your USB disk:

```bash
lsblk -o NAME,SIZE,MODEL,TYPE,MOUNTPOINTS
```

Then flash the whole disk:

```bash
sudo ./scripts/flash-usb.sh /dev/sdX
```

Example:

```bash
sudo ./scripts/flash-usb.sh /dev/sda
```

Do not pass a partition like `/dev/sda1`.

## Firmware settings

Disable Secure Boot.

Keep UEFI enabled.

Do not use shim, MOK, or signing for this version.

## Important kernel note

For bare-metal USB boot without initramfs, the Buildroot kernel must have the needed storage drivers built in, not as modules:

```txt
CONFIG_USB_XHCI_HCD=y
CONFIG_USB_STORAGE=y
CONFIG_SCSI=y
CONFIG_BLK_DEV_SD=y
CONFIG_EXT4_FS=y
CONFIG_EFI=y
CONFIG_EFI_STUB=y
```

If QEMU USB boot works but bare metal fails to mount root, rebuild the Buildroot kernel with those options built in.
