#!/bin/sh

set -ex

echo "Init called as: $@"

export PATH="/usr/local/bin:$PATH"

. /usr/local/bin/set-kernel-runtime.sh
. /usr/local/bin/set-sqlite-options.sh

sleep 1

# Start SQLite benchmark and export results
/bench/usr/src/sqlite/sqlite-bench --benchmarks=$BENCHMARK --cmd="$sqlite_params"

sleep 10000
