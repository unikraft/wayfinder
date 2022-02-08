#!/bin/bash

# verbose output
set -x

IMAGES=images/
BASEIP=172.190.0
NETIF=tux0
LOG=/tmp/sanitycheck-wrk.txt
touch $LOG

echo "creating bridge"
brctl addbr $NETIF || true
ifconfig $NETIF ${BASEIP}.1
killall -9 qemu-system-x86
pkill -9 qemu-system-x86

function benchmark_nginx_server {
  wrk -t 5 -d10s -c 5 http://${1}/index.html | tee -a ${2}
}

function cleanup {
  # kill all children (evil)
  ifconfig $NETIF down
  brctl delbr $NETIF
  rm /root/nginx.ext2.disposible
  killall -9 qemu-system-x86
  pkill -9 qemu-system-x86
  pkill -P $$
}

trap "cleanup" EXIT

cp /root/linux-nginx.ext2 /root/nginx.ext2.disposible

qemu-guest -k /root/linux-nginx-qemu.kernel \
  -d /root/nginx.ext2.disposible \
  -a "root=/dev/vda rw console=ttyS0 init=/guest_start.sh ${BASEIP}.1 ${BASEIP}.2 nginx" \
  -m 1024 -b ${NETIF} -x

# make sure that the server has properly started
sleep 3

# benchmark
benchmark_nginx_server ${BASEIP}.2 $LOG
#curl http://${BASEIP}.2/index.html --noproxy ${BASEIP}.2 --output -

# stop server
killall -9 qemu-system-x86
pkill -9 qemu-system-x86
rm /root/nginx.ext2.disposible
