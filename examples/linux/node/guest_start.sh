#!/bin/bash

set -ex

echo "Init called as: $@"

export PATH="/usr/local/bin:$PATH"

. /set-kernel-runtime.sh
. /set-node-options.sh

sleep 1

# Open the port for the client to connect to
coproc ./nc -l $WAYFINDER_DOMAIN_IP_ADDR 3000

cd ./bench

node ${NODE_OPTIONS[@]} dist/cli.js > results

cat results

tail -n 1 results > results_parsed

# Dirt format the header and send it to the client
echo -n -e "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nConnection: " "Keep-Alive" "\r\nDate: Mon, 01 Jan 1970 00:00:00 GMT GMT\r\nContent-Length: $(wc -c results_parsed)\r\n\r\n" > header
cat header results_parsed <&"${COPROC[0]}" >&"${COPROC[1]}"

# Stop the server
kill -9 $COPROC_PID

sleep 60
