#!/bin/bash

# power-monitor.sh $url

url=$1

curl -s --connect-timeout 1 $url | jq -r ".power"
