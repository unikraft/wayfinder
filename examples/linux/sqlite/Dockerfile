# Make the benchmarks to use for testing the kernels
FROM debian:jessie-slim AS benchmark

ENV BENCH_HASH=b3f6b4c709e38093c9b89598a9ff5fe3f60df6f5

COPY ./bench/serve_file.c /usr/src/sqlite/serve_file.c
COPY ./bench/0001-Add-compatibility-with-external-scrapers.patch /usr/src/sqlite/0001-Add-compatibility-with-external-scrapers.patch
COPY ./bench/0002-Add-option-to-execute-custom-commands.patch /usr/src/sqlite/0002-Add-option-to-execute-custom-commands.patch

RUN set -xe; \
    apt-get update; \
    apt install -y \
      curl build-essential git wget unzip; \
    cd /usr/src/sqlite; \
    git init; \
    git remote add origin https://github.com/ukontainer/sqlite-bench.git; \
    git fetch --depth 1 origin ${SQLITE_BENCH_HASH}; \
    git checkout FETCH_HEAD; \
    rm sqlite3.c sqlite3.h sqlite3ext.h; \
    wget --no-check-certificate https://www.sqlite.org/2019/sqlite-amalgamation-3300100.zip; \
    unzip sqlite-amalgamation-3300100.zip; \
    cp ./sqlite-amalgamation-3300100/sqlite3.c sqlite3.c; \
    cp ./sqlite-amalgamation-3300100/sqlite3.h sqlite3.h; \
    cp ./sqlite-amalgamation-3300100/sqlite3ext.h sqlite3ext.h; \
    rm ./sqlite-amalgamation-3300100.zip; \
    rm -rf ./sqlite-amalgamation-3300100 ./.git; \
    git apply ./0001-Add-compatibility-with-external-scrapers.patch; \
    git apply ./0002-Add-option-to-execute-custom-commands.patch; \
    make; \
    cd ..; \
    tar -cvf /usr/src/bench.tar /usr/src/sqlite;

# Build test image
FROM ghcr.io/unikraft/wayfinder/linux-base

# TODO: This file doesn't exist
ENV SQLITE_CFG="${WORKDIR}/configs/apps/sqlite.config"

# Nginx was selected randomly
ARG NGINX_VERSION=1.15.6
RUN crane export nginx:${NGINX_VERSION} ${WORKDIR}/sqlite.tar;

RUN ./scripts/build-with-configs.sh patch ${MICROVM_CFG} ${SQLITE_CFG}

COPY --from=benchmark /usr/src/bench.tar ${WORKDIR}/bench.tar

COPY ./set-sqlite-options.sh ${WORKDIR}/scripts/set-sqlite-options.sh
COPY ./guest_start.sh ${WORKDIR}/scripts/guest_start.sh
