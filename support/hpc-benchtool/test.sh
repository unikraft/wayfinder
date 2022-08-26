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

cat /results.txt

echo "Parsing results..."

mkdir -p /results

total_ft=0
total_mg=0
total_cg=0
total_is=0

# Calculate total Mop/s for each benchmark
while read result; do
  name=$(echo ${result} | awk '{ printf "%s", $1 }')
  value=$(echo ${result} | awk '{ printf "%s", $2 }')

  if [[ ${name} == *"ft"* ]]; then
    total_ft=$(echo "${total_ft} + ${value}" | bc)
  elif [[ ${name} == *"mg"* ]]; then
    total_mg=$(echo "${total_mg} + ${value}" | bc)
  elif [[ ${name} == *"cg"* ]]; then
    total_cg=$(echo "${total_cg} + ${value}" | bc)
  elif [[ ${name} == *"is"* ]]; then
    total_is=$(echo "${total_is} + ${value}" | bc)
  fi

  echo ${value} > /results/${name}
done < /results.txt

# Do the average of the Mop/s
total_ft=$(echo "${total_ft} / 4" | bc)
total_mg=$(echo "${total_mg} / 4" | bc)
total_cg=$(echo "${total_cg} / 4" | bc)
total_is=$(echo "${total_is} / 4" | bc)

# When testing omp, divide once more to the average for all thread numbers
if [[ ${results} == *".64"* ]] || [[ ${results} == *".4"* ]]; then
  total_ft=$(echo "${total_ft} / 2" | bc)
  total_mg=$(echo "${total_mg} / 2" | bc)
  total_cg=$(echo "${total_cg} / 2" | bc)
  total_is=$(echo "${total_is} / 2" | bc)
fi

echo "${total_ft}" > /results/total_ft
echo "${total_mg}" > /results/total_mg
echo "${total_cg}" > /results/total_cg
echo "${total_is}" > /results/total_is


echo "Total Mop/s for FT: ${total_ft}"
echo "Total Mop/s for MG: ${total_mg}"
echo "Total Mop/s for CG: ${total_cg}"
echo "Total Mop/s for IS: ${total_is}"

echo "Done!"
