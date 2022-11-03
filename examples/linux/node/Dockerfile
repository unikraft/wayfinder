# Make the benchmarks to use for testing the kernels
FROM debian:bullseye-slim AS benchmark

ENV BENCH_HASH=4a12828c6a1eed02a70c011bd080445dd319a05f

COPY ./bench/0001-Remove-references-to-getaddrinfo.patch /0001-Remove-references-to-getaddrinfo.patch

RUN set -xe; \
    apt-get update; \
    apt install -y \
      curl build-essential git npm; \
    mkdir /tmp/node-bench; \
    cd /tmp/node-bench; \
    git init; \
    git remote add origin https://github.com/v8/web-tooling-benchmark.git; \
    git fetch --depth 1 origin ${BENCH_HASH}; \
    git checkout FETCH_HEAD; \
    npm install; \
    rm -rf .git; \
    tar -czf /tmp/bench.tar.gz .; \
    cp /tmp/bench.tar.gz /bench.tar.gz; \
    cd /tmp; \
    git clone https://github.com/guzlewski/netcat.git; \
    cd netcat; \
    cp /0001-Remove-references-to-getaddrinfo.patch ./0001-Remove-references-to-getaddrinfo.patch; \
    git apply 0001-Remove-references-to-getaddrinfo.patch; \
    gcc -O0 -static -fno-pic -fno-pie netcat.c -o netcat.out -lpthread; \
    cp netcat.out /nc;

# Build test image
FROM ghcr.io/unikraft/wayfinder/linux-base

ENV NODE_CFG="${WORKDIR}/configs/apps/node.config"

ARG NODE_VERSION=8.17.0-slim
RUN crane export node:${NODE_VERSION} ${WORKDIR}/node.tar

RUN ./scripts/build-with-configs.sh patch ${MICROVM_CFG} ${NODE_CFG}

# COPY --from=benchmark /node ${WORKDIR}/node
COPY --from=benchmark /bench.tar.gz ${WORKDIR}/bench.tar.gz
COPY --from=benchmark /nc ${WORKDIR}/nc

COPY ./set-node-options.sh ${WORKDIR}/scripts/set-node-options.sh
COPY ./guest_start.sh ${WORKDIR}/scripts/guest_start.sh
