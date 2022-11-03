#!/bin/sh

set -ex

echo "Init called as: $@"

export PATH="/usr/local/bin:$PATH"

. /set-kernel-runtime.sh

sleep 1

# dirty, but not sure how to do this differently
redis-server /redis.conf
