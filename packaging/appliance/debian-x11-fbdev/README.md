# BTE Debian X11 Framebuffer Appliance

This appliance profile boots a minimal Debian system directly into BTE using:

- `systemd` to start the kiosk on `tty1`
- a small shell launcher for terminal color/session environment
- udev rules for Ledger USB devices

The result is a flashable disk image that shows the BTE terminal UI after boot,
without a login prompt on the primary display. The default graphical path uses
Xorg with the `fbdev` framebuffer driver and `nomodeset`, so startup does not
depend on Wayland, DRM card selection, or GPU acceleration.

Development images also autologin as `bte` on `tty2` for recovery. The `bte`
user is in the journal/admin groups and has passwordless sudo in this appliance
profile so recovery commands work without fighting the OS.

## Build A USB Image

Build BTE first:

```sh
go build -o bos .
```

Then build the raw disk image as root:

```sh
sudo packaging/appliance/debian-x11-fbdev/build-usb-image.sh
```

By default this creates:

```text
build/bte-debian-x11-fbdev-amd64.img
```

Write it to a USB drive:

```sh
sudo dd if=build/bte-debian-x11-fbdev-amd64.img of=/dev/sdX bs=4M status=progress conv=fsync
```

Replace `/dev/sdX` with the whole USB device, not a partition like `/dev/sdX1`.

## Test With QEMU

Install QEMU on the host:

```sh
sudo apt install qemu-system-x86
```

Build BTE and the image:

```sh
go build -o bos .
sudo packaging/appliance/debian-x11-fbdev/build-usb-image.sh
```

The build should print:

```text
Configured tty1 autologin launcher for X11 framebuffer BTE.
```

Boot the image with a graphical display:

```sh
qemu-system-x86_64 \
  -enable-kvm \
  -m 2048 \
  -smp 2 \
  -drive file=build/bte-debian-x11-fbdev-amd64.img,format=raw,if=virtio \
  -vga virtio \
  -display gtk
```

If KVM is unavailable, remove `-enable-kvm`. If GTK display support is missing,
try `-display sdl` instead.

The final `-display gtk` line is required. If the command ends at a trailing
backslash, the shell is waiting for the next line instead of starting QEMU.

The image builder creates a GPT disk image with:

- a BIOS boot partition
- an EFI system partition
- an ext4 Debian root partition
- Debian minimal base
- Linux kernel and GRUB
- Xorg, `openbox`, UTF-8 `xterm`, fonts, udev, and BTE kiosk files
- AMD/NVIDIA graphics firmware and Mesa DRM/EGL/Vulkan userspace packages
- Realtek firmware for common wired and wireless adapters
- wired DHCP via `systemd-networkd`
- `getty@tty1.service` autologin, which immediately starts graphical BTE as
  the `bte` user

Useful build overrides:

```sh
BTE_IMAGE_SIZE=10G sudo packaging/appliance/debian-x11-fbdev/build-usb-image.sh
BTE_BINARY_PATH=/path/to/bos sudo packaging/appliance/debian-x11-fbdev/build-usb-image.sh /tmp/bte.img
BTE_DEBIAN_MIRROR=http://deb.debian.org/debian sudo packaging/appliance/debian-x11-fbdev/build-usb-image.sh
BTE_NAMESERVERS="1.1.1.1 8.8.8.8" sudo packaging/appliance/debian-x11-fbdev/build-usb-image.sh
BTE_DEBUG_PASSWORD=bte sudo packaging/appliance/debian-x11-fbdev/build-usb-image.sh
```

The default image size is `8G` so the firmware, Xorg stack, and kernel headers
have enough room during installation. Use `BTE_IMAGE_SIZE` to make a smaller or
larger image.

## Image-Only Build

This profile does not include a host installer. `build-usb-image.sh` creates a disk image, mounts that image under `/tmp`, and copies files into that mounted image root:

- `./bos` to `/usr/local/bin/bte`
- `./sql/schema.sql` to `/var/lib/bte/sql/schema.sql`
- kiosk config to `/etc/bte`
- systemd unit files to `/etc/systemd/system`
- Ledger udev rules to `/etc/udev/rules.d`

Those paths are inside the mounted image, not on the host OS.

## Raw Console Fallback

The image boots into the X11 framebuffer kiosk by default. For temporary
raw-console testing, set `BTE_RAW_CONSOLE=1` in `/var/lib/bte/.bash_profile`
before the kiosk launcher block.

## Graphics Startup Notes

The image passes `nomodeset` on the kernel command line and forces Xorg to use
the framebuffer driver instead of probing DRM/KMS GPUs. If Xorg does not start,
check `/var/lib/bte/bte-session.log` from `tty2`. The launcher records the
kernel command line, loaded Nvidia modules, visible `/dev/dri/card*` devices,
and visible `/dev/fb*` devices. If Xorg hangs, the launcher times out and falls
back to the raw console.

For debugging, inspect the Xorg log:

```sh
cat /var/lib/bte/.local/share/xorg/Xorg.0.log
```

## Admin Access

The autologin launcher takes over `tty1`. `tty2` is configured as an autologin
recovery console for the `bte` user. The default development password is `bte`,
or set `BTE_DEBUG_PASSWORD` while building. The recovery user can run
`journalctl`, `systemctl`, and `sudo` without extra setup.

Keep another access path enabled while developing, such as SSH or `tty2`.

To stop the kiosk:

```sh
sudo systemctl stop getty@tty1.service
```

## Terminal Choice

The default profile uses `xterm` because it is small, widely packaged, and
avoids GPU acceleration requirements. Alacritty or another X11 terminal can be
swapped in later by changing `bte-xsession`.
