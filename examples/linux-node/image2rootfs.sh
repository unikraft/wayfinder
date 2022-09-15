#!/usr/bin/env bash
# Usage: ./image2rootfs.sh input_archive output_filesystem

set -xe

_die() {
  echo "$*" 1>&2 ; exit 1;
}

WORKDIR=$(pwd)
APP="${1%.*}"
FS="${2:-ext3}"
MOUNTDIR=${MOUNTDIR:-$(mktemp -d)}

# Generate the filesystem
dd if=/dev/zero of=${APP}.${FS} bs=1 count=0 seek=1024M
yes | mkfs."${FS}" "${APP}.${FS}"
mount "${APP}.${FS}" ${MOUNTDIR}
tar -xvf ${1} -C ${MOUNTDIR}

# Install devices
mknod -m 666 ${MOUNTDIR}/dev/null    c 1 3
mknod -m 666 ${MOUNTDIR}/dev/zero    c 1 5
mknod -m 666 ${MOUNTDIR}/dev/ptmx    c 5 2
mknod -m 666 ${MOUNTDIR}/dev/tty     c 5 0
mknod -m 444 ${MOUNTDIR}/dev/random  c 1 8
mknod -m 444 ${MOUNTDIR}/dev/urandom c 1 9
mknod -m 660 ${MOUNTDIR}/dev/mem     c 1 1

# Unmount and finish
umount ${MOUNTDIR} || true
rmdir ${MOUNTDIR} || true
