#!/bin/sh
export PATH="/usr/local/bin:$PATH"
mkdir -p /trusted
mount -t proc proc /proc
ulimit -n 65535
echo "These are the params $@"
./busybox-x86_64 ip addr add $1 dev eth0
./busybox-x86_64 ip addr add 127.0.0.1/24 dev lo
./busybox-x86_64 ip link set eth0 up
./busybox-x86_64 ip link set lo up
./busybox-x86_64 ip route add default via $2 dev eth0

echo "APP START"

if which nginx; then
#    sh guest_net.sh
    cp `which nginx` /trusted
    if echo $@ | grep trusted - > /dev/null; then
        echo ========KML=========
        /trusted/libc.so /trusted/nginx -g 'daemon off;error_log stderr debug;'
    else
        echo ========NOKML=========
        $@ -g 'daemon off;error_log stderr debug;'
    fi
    exit
fi
