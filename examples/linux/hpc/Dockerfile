# Make the benchmarks to use for testing the kernels
FROM debian:jessie-slim AS benchmark

COPY ./bench/make-bench.sh /make-bench.sh
COPY ./bench/run-bench-ser.sh /run-bench-ser.sh
COPY ./bench/run-bench-omp.sh /run-bench-omp.sh
COPY ./bench/0001-Remove-references-to-getaddrinfo.patch /0001-Remove-references-to-getaddrinfo.patch

RUN set -xe; \
    apt-get update; \
    apt install -y \
      curl build-essential zip git; \
    curl -L "http://aces.snu.ac.kr/download.php?p=snunpbsuite&ver=2019&email=cezar.craciunoiu@gmail.com" -o /tmp/benchmark.tar.gz; \
    tar -xf /tmp/benchmark.tar.gz -C /tmp; \
    cd /tmp/SNU_NPB_2019/NPB3.3-OMP-C; \
    sed -i 's/O3/O3 -static/g' config/make.def; \
    cp /make-bench.sh ./make-bench.sh; \
    ./make-bench.sh; \
    zip /benchmark_omp.zip ./bin/*; \
    cd /tmp/SNU_NPB_2019/NPB3.3-SER-C; \
    cp /make-bench.sh ./make-bench.sh; \
    ./make-bench.sh; \
    zip /benchmark_ser.zip ./bin/*; \
    git clone https://github.com/guzlewski/netcat.git; \
    cd netcat; \
    cp /0001-Remove-references-to-getaddrinfo.patch ./0001-Remove-references-to-getaddrinfo.patch; \
    git apply 0001-Remove-references-to-getaddrinfo.patch; \
    gcc -O0 -static -fno-pic -fno-pie netcat.c -o netcat.out -lpthread; \
    cp netcat.out /nc;

# Build test image
FROM ghcr.io/unikraft/wayfinder/linux-base

ENV REDIS_CFG="${WORKDIR}/configs/apps/redis.config"

# Redis was chosen randomly
ARG REDIS_VERSION=7.0.4
RUN crane export redis:${REDIS_VERSION} ${WORKDIR}/hpc.tar

RUN ./scripts/build-with-configs.sh patch ${MICROVM_CFG} ${REDIS_CFG}

COPY --from=benchmark /benchmark_omp.zip ${WORKDIR}/benchmark_omp.zip
COPY --from=benchmark /benchmark_ser.zip ${WORKDIR}/benchmark_ser.zip
COPY --from=benchmark /run-bench-ser.sh ${WORKDIR}/run-bench-ser.sh
COPY --from=benchmark /run-bench-omp.sh ${WORKDIR}/run-bench-omp.sh
COPY --from=benchmark /nc ${WORKDIR}/nc

COPY ./guest_start.sh ${WORKDIR}/scripts/guest_start.sh
