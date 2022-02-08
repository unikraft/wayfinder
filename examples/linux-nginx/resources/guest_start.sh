#!/bin/sh

set -x

export PATH="/usr/local/bin:$PATH"
mount -t proc proc /proc
ulimit -n 65535
echo 1024 > /proc/sys/net/core/somaxconn

echo "APP START"

if which nginx; then
    # ip self, ip gateway
    sh guest_net.sh $1 $2
    ${@:3}
    exit
fi
