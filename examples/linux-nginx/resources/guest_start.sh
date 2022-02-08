#!/bin/sh

set -x

echo "Init called as: $@"

export PATH="/usr/local/bin:$PATH"
mount -t proc proc /proc
ulimit -n 65535
echo 1024 > /proc/sys/net/core/somaxconn

echo "APP START"

if which nginx; then
    # ip self, ip gateway
    sh guest_net.sh $1 $2
    # dirty, but not sure how to do this differently
    $3 $4 $5 $6 $7
    exit
fi
