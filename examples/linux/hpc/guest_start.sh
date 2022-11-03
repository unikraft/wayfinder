#!/bin/bash

set -ex

echo "Init called as: $@"

export PATH="/usr/local/bin:$PATH"

. /set-kernel-runtime.sh

sleep 1

# Open the port for the client to connect to
coproc ./nc -l $WAYFINDER_DOMAIN_IP_ADDR 3000

# Switch to what type of benchmarks you want to run
if [ "$BENCHMARK" == "serial" ]; then
  cd ./bench-ser
else
  cd ./bench-omp
fi

./run-bench.sh

cd result

# concatenate results in a single file
if [ "$BENCHMARK" == "serial" ]; then
  for benchmark in ft mg cg is
  do
      for class in S W A B
      do
          cat $benchmark.$class.out | grep Mop/s | \
          awk -v benchmark=$benchmark -v class=$class \
          '{ print benchmark "." class " " $4 }' >> concatenated_results
      done
  done
else
  for benchmark in ft mg cg is
  do
      for class in S W A B
      do
          for num_thread in 4 16
          do
              cat $benchmark.$class.$num_thread.out | grep Mop/s/thread | tr -s " " | \
              awk -v benchmark=$benchmark -v class=$class -v num_thread=$num_thread \
              '{ print benchmark "." class "." num_thread " " $3 }' >> concatenated_results
          done
      done
  done
fi

# Dirt format the header and send it to the client
echo -n -e "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nConnection: " "Keep-Alive" "\r\nDate: Mon, 01 Jan 1970 00:00:00 GMT GMT\r\nContent-Length: $(wc -c concatenated_results)\r\n\r\n" > header
cat header concatenated_results <&"${COPROC[0]}" >&"${COPROC[1]}"

# Stop the server
kill -9 $COPROC_PID
