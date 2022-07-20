#!/usr/bin/env bash

# Usage: test.sh

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

retries=3

while [[ $retries -gt 0 ]]; do
  taskset -c ${WAYFINDER_CORE_ID0} \
    redis-benchmark \
      -h ${WAYFINDER_DOMAIN_IP_ADDR} \
      -c ${NUM_PARALLEL_CONNS} \
      -d ${REQUESTS_SIZE} \
      -n ${NUM_TOTAL_REQUESTS} \
      -t GET,SET,LPUSH,LPOP \
      -q \
      |& tee /results.txt

  if [[ $? -eq 0 && $(cat /results.txt | grep "Connection refused") == "" ]]; then
    break
  fi
  retries=$((retries - 1))
  sleep 1
done

if [[ $retries -eq 0 ]]; then
  echo "Failed to run redis-benchmark!"
  exit 1
fi

set +x

if [[ ! -f /results.txt ]]; then
  echo "No results!"
  exit 1
fi

mkdir -p /results/
echo -n "$(cat /results.txt | grep -e '^SET:'   | awk -F'[ \n]' '{ print $4 }' | awk -F'[.]' '{ print $1 }')" > /results/set.txt
echo -n "$(cat /results.txt | grep -e '^GET:'   | awk -F'[ \n]' '{ print $4 }' | awk -F'[.]' '{ print $1 }')" > /results/get.txt
echo -n "$(cat /results.txt | grep -e '^LPUSH:' | awk -F'[ \n]' '{ print $4 }' | awk -F'[.]' '{ print $1 }')" > /results/lpush.txt
echo -n "$(cat /results.txt | grep -e '^LPOP:'  | awk -F'[ \n]' '{ print $4 }' | awk -F'[.]' '{ print $1 }')" > /results/lpop.txt

echo "Done!"
