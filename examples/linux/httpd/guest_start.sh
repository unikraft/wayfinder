#!/bin/sh

set -ex

echo "Init called as: $@"

export PATH="/usr/local/apache2/bin:/usr/local/bin:$PATH"

. /set-kernel-runtime.sh

sleep 1

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

apachectl -X
