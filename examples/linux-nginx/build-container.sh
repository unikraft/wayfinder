#!/bin/bash

LUPINE_HASH=3729c7570e5f1b25e786171d0a19cef858b4923f
SUDO=sudo

set -x

# gain sudo privileges
$SUDO -v

SCRIPTS_PATH=$(dirname "$0")
SCRIPTS_PATH=$(cd "$SCRIPTS_PATH" && pwd)

NGINX_CONF="${SCRIPTS_PATH}/resources/nginx.conf"
# modified scripts with a few tweaks for our use-case
GUEST_START="${SCRIPTS_PATH}/resources/guest_start.sh"
GUEST_NET="${SCRIPTS_PATH}/resources/guest_net.sh"
TMP_FOLDER="${SCRIPTS_PATH}/generated-data"
mkdir -p $TMP_FOLDER

# =============================================================================
# Lupine clone
# =============================================================================

LUPINE_DIR="${SCRIPTS_PATH}/generated-data/Lupine-Linux"

if [ ! -d "${LUPINE_DIR}" ]; then
  pushd $TMP_FOLDER
  git clone https://github.com/hlef/Lupine-Linux.git
  git checkout b9dc99bbd09180b0a3548583d58f9c003d4576e8

  pushd $LUPINE_DIR
  git checkout $LUPINE_HASH
  git submodule update --init

  pushd ./load_entropy
  make
  popd # ./load_entropy

  popd # $LUPINE_DIR

  popd # $TMP_FOLDER
fi

# =============================================================================
# FS images
# Note: we have to do this here because it requires docker and docker in docker
# is a pain.
# =============================================================================

if [ ! -f "${TMP_FOLDER}/nginx.ext2" ]; then
  pushd ${LUPINE_DIR}
  cp ${GUEST_START} ./scripts/guest_start.sh
  cp ${GUEST_NET} ./scripts/guest_net.sh

  # reduce default size of images, 20G is way to much
  sed -i -e "s/seek=20G/seek=30M/" ./scripts/image2rootfs.sh

  # build FS image from Alpine
  $SUDO ./scripts/image2rootfs.sh nginx 1.15.6-alpine ext2
  mv ${LUPINE_DIR}/nginx.ext2 ${TMP_FOLDER}/nginx.ext2

  # Cleanup a bit
  git checkout ./scripts/image2rootfs.sh
  git checkout ./scripts/guest_net.sh
  git checkout ./scripts/guest_start.sh
  popd

  # set nginx configuration in the FS image
  $SUDO modprobe loop
  $SUDO mkdir -p /mnt/nginx-tmp
  $SUDO mount -o loop ${TMP_FOLDER}/nginx.ext2 /mnt/nginx-tmp
  $SUDO cp ${NGINX_CONF} /mnt/nginx-tmp/etc/nginx/nginx.conf
  $SUDO umount /mnt/nginx-tmp
  $SUDO rm -rf /mnt/nginx-tmp
fi

# =============================================================================
# Docker image
# =============================================================================

docker build -t hlefeuvre/linux-nginx -f linux-nginx.dockerfile .
