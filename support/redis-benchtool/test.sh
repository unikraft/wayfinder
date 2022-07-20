#!/usr/bin/env bash

# Usage: test.sh path/to/payload.html

NUM_PARALLEL_CONNS=${NUM_PARALLEL_CONNS:-50}
REQUESTS_SIZE=${REQUESTS_SIZE:-2}
NUM_TOTAL_REQUESTS=${NUM_TOTAL_REQUESTS:-100000}
DURATION=${DURATION:-10}

if [[ -z "${WAYFINDER_CORE_ID0}" ]]; then
  echo "Missing core to pin!"
  exit 1
fi

if [[ -z "${WAYFINDER_DOMAIN_IP_ADDR}" ]]; then
  echo "Missing domain IP address!"
  exit 1
fi

ping -c 1 ${WAYFINDER_DOMAIN_IP_ADDR} > /dev/null 2>&1
if [[ $? -ne 0 ]]; then
  echo "Domain IP address is not reachable!"
  exit 1
fi

echo "Starting experiment..."

set -xe

taskset -c ${WAYFINDER_CORE_ID0} \
  redis-benchmark \
    -h ${WAYFINDER_DOMAIN_IP_ADDR} \
    -c ${NUM_PARALLEL_CONNS} \
    -d ${REQUESTS_SIZE} \
    -n ${NUM_TOTAL_REQUESTS} \
    -t GET,SET,LPUSH,LPOP \
    -q \
    |& tee /results.txt

set +x

if [[ ! -f /results.txt ]]; then
  echo "No results!"
  exit 1
fi

mkdir -p /results/
echo -n "$(cat results.txt | grep -e '^SET:' | awk '{ print $2 }')" > /results/set.txt
echo -n "$(cat results.txt | grep -e '^GET:' | awk '{ print $2 }')" > /results/get.txt
echo -n "$(cat results.txt | grep -e '^LPUSH:' | awk '{ print $2 }')" > /results/lpush.txt
echo -n "$(cat results.txt | grep -e '^LPOP:' | awk '{ print $2 }')" > /results/lpop.txt

echo "Done!"
