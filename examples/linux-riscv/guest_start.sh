#!/bin/sh

set -ex

echo "Hello RISC-V World!"

mount -t proc proc /proc

/sbin/ip addr add $WAYFINDER_DOMAIN_IP_ADDR/24 dev eth0
/sbin/ip addr add 127.0.0.1/24 dev lo
/sbin/ip link set eth0 up
/sbin/ip link set lo up
/sbin/ip route add default via $WAYFINDER_DOMAIN_IP_GW_ADDR dev eth0
