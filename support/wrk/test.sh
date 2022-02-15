#!/usr/bin/env bash

# Usage: test.sh path/to/payload.html

PAYLOAD_PATH=${1:-/payload.txt}
NUM_PARALLEL_CONNS=${NUM_PARALLEL_CONNS:-30}
CURL_MAX_CONNECTION_TIMEOUT=${CURL_MAX_CONNECTION_TIMEOUT:-20}
DURATION=${DURATION:-10}

if [[ -z "${WAYFINDER_CORE_ID0}" ]]; then
  echo "Missing core to pin!"
  exit 1
fi

if [[ -z "${WAYFINDER_DOMAIN_IP_ADDR}" ]]; then
  echo "Missing domain IP address!"
  exit 1
fi

# perform quick test to check domain is online
CURL_RETURN_CODE=0
CURL_OUTPUT=`curl -w httpcode=%{http_code} -m ${CURL_MAX_CONNECTION_TIMEOUT} http://${WAYFINDER_DOMAIN_IP_ADDR}:80 2> /dev/null` || CURL_RETURN_CODE=$?
if [[ ${CURL_RETURN_CODE} -ne 0 ]]; then  
  echo "curl connection failed with return code: ${CURL_RETURN_CODE}"
  exit 1
else
  HTTPCODE=$(echo "${CURL_OUTPUT}" | sed -n "s/.*\httpcode=//p")
  if [[ ${HTTPCODE} -ne 200 ]]; then
    echo "curl operation/command failed due to server return code: ${HTTPCODE}"
    exit 1
  fi
fi

echo "Starting experiment..."

set -xe

taskset -c ${WAYFINDER_CORE_ID0} \
  wrk \
    -d ${DURATION} --latency \
    -t ${NUM_PARALLEL_CONNS} \
    -c ${NUM_PARALLEL_CONNS} http://${WAYFINDER_DOMAIN_IP_ADDR}:80/${PAYLOAD_PATH} |& tee /results.txt

set +x

if [[ ! -f /results.txt ]]; then
  echo "No results!"
  exit 1
fi

mkdir -p /results/
echo "$(cat /results.txt | grep 'Requests/sec:' | awk '{ print $2 }')" > /results/throughput.txt

echo "Done!"
