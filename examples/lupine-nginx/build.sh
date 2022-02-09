#!/bin/bash

set -ex
env

export LC_ALL=C.UTF-8
export LANG=C.UTF-8

cd /Lupine-Linux/
service docker start || true

make -C linux/ oldconfig
make -j8 -C linux/
mkdir -p kernelbuild/lupine-djw-kml-qemu++nginx/
cp linux/vmlinux  kernelbuild/lupine-djw-kml-qemu++nginx/
INSTALL_PATH=../kernelbuild/lupine-djw-kml-qemu++nginx/ make -C linux install
cd load_entropy
make
cd ../

# Modify the nginx config
cat /nginx-nocaching.conf | sed "s/\$WORKER_CONNECTIONS/64/g" > /nginx.conf
export ACCESS_LOG="n"
if [[ $ACCESS_LOG == "y" ]]; then \
 export ACCESS_LOG="\/dev\/stdout"; \
else \
 export ACCESS_LOG="off"; \
fi
sed -i "s/\$ACCESS_LOG/$ACCESS_LOG/g" /nginx.conf
sed -i "s/\$KEEPALIVE_TIMEOUT/30/g" /nginx.conf

cat /nginx.conf

# build the fs
echo "echo 1024 > /proc/sys/net/core/somaxconn" >> ./scripts/guest_net.sh
sed -i -e "s/seek=20G/seek=30M/" ./scripts/image2rootfs.sh
cp /guest_start.sh ./scripts/guest_start.sh
./scripts/image2rootfs.sh nginx 1.15.6-alpine ext2

# modprobe loop
mkdir -p /mnt/nginx-tmp
mount -o loop nginx.ext2 /mnt/nginx-tmp
cp /nginx.conf /mnt/nginx-tmp/etc/nginx/nginx.conf
# Generate payload
tr -dc A-Za-z0-9 </dev/urandom | head -c 2048 > /mnt/nginx-tmp/usr/share/nginx/html/payload.txt

umount /mnt/nginx-tmp
rm -rf /mnt/nginx-tmp

cp nginx.ext2 /mnt/
cp /Lupine-Linux/kernelbuild/lupine-djw-kml-qemu++nginx/vmlinuz-4.0.0-kml+ /mnt/

