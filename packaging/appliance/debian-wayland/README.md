# BTE Debian Wayland Kiosk

This appliance profile boots a minimal Debian system directly into BTE using:

- `cage` as the single-application Wayland compositor
- `foot` as the terminal emulator
- `systemd` to start the kiosk on `tty1`
- a small shell launcher for terminal color/session environment
- udev rules for Ledger USB devices

The result is a flashable disk image that shows the BTE terminal UI after boot, without a login prompt on the primary display.

## Build A USB Image

Build BTE first:

```sh
go build -o bos .
```

Then build the raw disk image as root:

```sh
sudo packaging/appliance/debian-wayland/build-usb-image.sh
```

By default this creates:

```text
build/bte-debian-wayland-amd64.img
```

Write it to a USB drive:

```sh
sudo dd if=build/bte-debian-wayland-amd64.img of=/dev/sdX bs=4M status=progress conv=fsync
```

Replace `/dev/sdX` with the whole USB device, not a partition like `/dev/sdX1`.

The image builder creates a GPT disk image with:

- a BIOS boot partition
- an EFI system partition
- an ext4 Debian root partition
- Debian minimal base
- Linux kernel and GRUB
- `cage`, `foot`, fonts, udev, and BTE kiosk files
- wired DHCP via `systemd-networkd`

Useful build overrides:

```sh
BTE_IMAGE_SIZE=6G sudo packaging/appliance/debian-wayland/build-usb-image.sh
BTE_BINARY_PATH=/path/to/bos sudo packaging/appliance/debian-wayland/build-usb-image.sh /tmp/bte.img
BTE_DEBIAN_MIRROR=http://deb.debian.org/debian sudo packaging/appliance/debian-wayland/build-usb-image.sh
BTE_NAMESERVERS="1.1.1.1 8.8.8.8" sudo packaging/appliance/debian-wayland/build-usb-image.sh
```

## Image-Only Build

This profile does not include a host installer. `build-usb-image.sh` creates a disk image, mounts that image under `/tmp`, and copies files into that mounted image root:

- `./bos` to `/usr/local/bin/bte`
- `./sql/schema.sql` to `/var/lib/bte/sql/schema.sql`
- kiosk config to `/etc/bte`
- systemd unit files to `/etc/systemd/system`
- Ledger udev rules to `/etc/udev/rules.d`

Those paths are inside the mounted image, not on the host OS.

## Admin Access

The kiosk disables the normal login prompt on `tty1` only. Keep another access path enabled while developing:

- SSH, or
- another TTY such as `tty2`

To stop the kiosk:

```sh
sudo systemctl stop bte-kiosk.service
```

To re-enable a login prompt on `tty1`:

```sh
sudo systemctl enable --now getty@tty1.service
sudo systemctl disable --now bte-kiosk.service
```

## Terminal Choice

This profile uses `foot` because it is small, fast, Wayland-native, and packaged in Debian. Alacritty can be swapped in later by changing `bte-kiosk.service` to launch `alacritty` inside `cage`.
