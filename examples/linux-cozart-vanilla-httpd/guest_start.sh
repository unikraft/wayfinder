#!/bin/sh

set -ex

# Usage: guest_start.sh $IP_ADDR $IP_GW

echo "Init called as: $@"

export PATH="/usr/local/bin:$PATH"

mount -t proc proc /proc
ulimit -n 65535
echo 1024 > /proc/sys/net/core/somaxconn

# ip self, ip gateway
./busybox-x86_64 ip addr add $WAYFINDER_DOMAIN_IP_ADDR/24 dev eth0
./busybox-x86_64 ip addr add 127.0.0.1/24 dev lo
./busybox-x86_64 ip link set eth0 up
./busybox-x86_64 ip link set lo up
./busybox-x86_64 ip route add default via $WAYFINDER_DOMAIN_IP_GW_ADDR dev eth0

tr -dc A-Za-z0-9 < /dev/urandom | head -c $PAYLOAD_SIZE > /usr/local/apache2/htdocs/index.html

cat << EOF >> /usr/local/apache2/conf/httpd.conf
ContentDigest ${CONTENT_DIGEST}
EnableMMAP ${ENABLE_MMAP}
EnableSendfile ${ENABLE_SEND_FILE}
ExtendedStatus ${EXTENDED_STATUS}
FlushMaxPipeLined ${FLUSH_MAX_PIPELINED}
FlushMaxThreshold ${FLUSH_MAX_THRESHOLD}
KeepAlive ${KEEP_ALIVE}
KeepAliveTimeout ${KEEP_ALIVE_TIMEOUT}
LimitRequestBody ${LIMIT_REQUEST_BODY}
LimitRequestFields ${LIMIT_REQUEST_FIELDS}
LimitRequestFieldSize ${LIMIT_REQUEST_FILED_SIZE}
LimitRequestLine ${LIMIT_REQUEST_LINE}
LogLevel ${LOG_LEVEL}
MaxKeepAliveRequests ${MAX_KEEP_ALIVE_REQUESTS}
MaxRangeOverlaps ${MAX_RANGE_OVERLAPS}
MaxRangeReversals ${MAX_RANGE_REVERSALS}
MaxRanges ${MAX_RANGES}
Mutex ${MUTEX}
ReadBufferSize ${READ_BUFFER_SIZE}
ServerTokens ${SERVER_TOKENS}
TimeOut ${TIME_OUT}
EOF

apachectl -k start
sleep 100000
