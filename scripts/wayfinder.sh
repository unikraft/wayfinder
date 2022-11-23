#!/bin/bash

set -xe

POSTGRES_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' wayfinder-postgres)
REDIS_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' wayfinder-redis)

POSTGRES_HOST="${POSTGRES_IP}" REDIS_ADDR="${REDIS_IP}:6379" ./dist/wayfinderd -c ./config/wayfinderd.yaml "$@"
