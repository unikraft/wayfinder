#!/usr/bin/env bash

# Usage: test.sh

DURATION=${DURATION:-30}

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

echo "Waiting for results..."

timeout ${DURATION} curl -s ${WAYFINDER_DOMAIN_IP_ADDR}:3000 > /results.txt

if [[ ! -f /results.txt ]]; then
  echo "No results!"
  exit 1
fi

mkdir -p /results/
echo -n "$(cat /results.txt | awk '{ print $3 }')" > /results/benchmark.txt

echo "Done!"
