#!/bin/sh
set -e

TARGET_DIR="$1"

# BOS init
cat > "$TARGET_DIR/etc/inittab" <<'INITTAB'
::sysinit:/bin/mount -t proc proc /proc
::sysinit:/bin/mount -t sysfs sysfs /sys
::sysinit:/bin/mount -t devtmpfs devtmpfs /dev
::sysinit:/bin/mkdir -p /dev/pts /dev/shm /run /tmp /var/log /run/xdg
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

console::respawn:/usr/local/bin/start-bos-session

::shutdown:/etc/init.d/rcK
::shutdown:/bin/umount -a -r
INITTAB
