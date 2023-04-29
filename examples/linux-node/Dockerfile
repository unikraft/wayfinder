#
# Use a file system image prebaked.
#
FROM debian:bullseye-slim AS node-filesystem

ARG NODE_VERSION=8.17.0-slim

RUN set -xe; \
    apt-get update; \
    apt install -y --force-yes \
      curl; \
    curl -L https://github.com/google/go-containerregistry/releases/download/v0.8.0/go-containerregistry_Linux_x86_64.tar.gz -o /tmp/go-containerregistry.tar.gz; \
    tar -xf /tmp/go-containerregistry.tar.gz -C /tmp; \
    mv /tmp/crane /usr/local/bin; \
    crane export node:${NODE_VERSION} /node.tar;

#
# Make the benchmarks to use for testing the kernels
#
FROM debian:bullseye-slim AS benchmark

ENV BENCH_HASH=4a12828c6a1eed02a70c011bd080445dd319a05f

COPY ./patches/0001-Remove-references-to-getaddrinfo.patch /0001-Remove-references-to-getaddrinfo.patch

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


#
# Build the kernel
#
FROM debian:bullseye-slim

#
# Base utilities and tools
#
RUN set -xe; \
    cp /etc/apt/sources.list /tmp/sources.list; \
    sed -i 's/deb /deb-src /' /tmp/sources.list; \
    cat /tmp/sources.list >> /etc/apt/sources.list; \
    echo "deb [trusted=yes] http://dk.archive.ubuntu.com/ubuntu/ bionic main universe" >> /etc/apt/sources.list; \
    apt-get update; \
    apt-get install -y --force-yes \
      build-essential \
      git \
      apt-transport-https \
      ca-certificates \
      gcc-5 \
      curl \
      bc \
      zip \
      libssl-dev \
      libgomp1 \
      openssl; \
      update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-5 5;

#
# Grab and prepare Lupine/Linux sources
#
ENV LUPINE_HASH=3729c7570e5f1b25e786171d0a19cef858b4923f
ENV LINUX_HASH=39a8804455fb23f09157341d3ba7db6d7ae6ee76
ENV OSV_HASH=3e658b8ab32a424ae58a24d01d5b2e22da61932e
ENV WORKDIR="/root/workspace"
ENV MICROVM_CFG="${WORKDIR}/configs/microvm.config"
ENV NODE_CFG="${WORKDIR}/configs/apps/node.config"

COPY ./patches/0001-init-Increase-env-var-limit.patch /0001-init-Increase-env-var-limit.patch
COPY ./patches/0001-arch-x86-Increase-cmdline-size.patch /0001-arch-x86-Increase-cmdline-size.patch

RUN set -xe; \
    mkdir ${WORKDIR}; \
    cd ${WORKDIR}; \
    git init; \
    git remote add origin https://github.com/hlef/Lupine-Linux.git; \
    git fetch --depth 1 origin ${LUPINE_HASH}; \
    git checkout FETCH_HEAD; \
    git checkout ${LUPINE_HASH}; \
    cd ${WORKDIR}/linux; \
    git init; \
    git remote add origin https://github.com/torvalds/linux.git; \
    git fetch --depth 1 origin ${LINUX_HASH}; \
    git checkout FETCH_HEAD; \
    cd ${WORKDIR}/osv; \
    git init; \
    git remote add origin https://github.com/cloudius-systems/osv.git; \
    git fetch --depth 1 origin ${OSV_HASH}; \
    git checkout FETCH_HEAD; \
    cd ${WORKDIR}/load_entropy; \
    make;

WORKDIR ${WORKDIR}

RUN set -xe; \
    echo "CONFIG_PCI=y"               >> ${MICROVM_CFG}; \
    echo "CONFIG_VIRTIO_BLK_SCSI=y"   >> ${MICROVM_CFG}; \
    echo "CONFIG_VIRTIO_PCI_LEGACY=y" >> ${MICROVM_CFG}; \
    echo "CONFIG_VIRTIO_PCI=y"        >> ${MICROVM_CFG}; \
    echo "CONFIG_VGA_ARB_MAX_GPUS=16" >> ${MICROVM_CFG}; \
    sed -i 's/-Wall -Wundef -Wstrict-prototypes/-Wall -Wundef -Wstrict-prototypes -fno-pie/g' linux/Makefile; \
    sed -i -e "s/docker run -it -v \"\$(PWD)\/linux\":\/linux-volume --rm linuxbuild:latest//" Makefile; \
    sed -i -e "s/linux-volume/root\/workspace\/linux/" Makefile; \
    sed -i -e 's/make patch-linux || true/echo "Keep everything as it is"/g' ${WORKDIR}/scripts/build-with-configs.sh; \
    sed -i -e 's/CONFIGS=$@/CONFIGS=${@:2}/g' ${WORKDIR}/scripts/build-with-configs.sh; \
    sed -i -e 's/-j8/-j $(( $(nproc) + 1 ))/g' ${WORKDIR}/Makefile; \
    ./scripts/build-with-configs.sh patch ${MICROVM_CFG} ${NODE_CFG}

#
# Local tooling and filesystem
#
COPY ./kconfig-set.sh /usr/local/bin/
COPY ./image2rootfs.sh /usr/local/bin/

# Overwrite the start scripts with custom functionality, allowing for
# use with Wayfinder's test.benchtool.args:
COPY ./guest_start.sh ${WORKDIR}/scripts/guest_start.sh

# Copy from the initial build step the benchmark
# COPY --from=benchmark /node ${WORKDIR}/node
COPY --from=benchmark /bench.tar.gz ${WORKDIR}/bench.tar.gz
COPY --from=benchmark /nc ${WORKDIR}/nc

# Copy from the initial build step the pre-fetched filesystem
COPY --from=node-filesystem /node.tar ${WORKDIR}/node.tar
