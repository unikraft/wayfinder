FROM debian:jessie-slim AS sqlite-filesystem

ARG NGINX_VERSION=1.15.6

RUN set -xe; \
    apt-get update; \
    apt install -y \
      curl; \
    curl -L https://github.com/google/go-containerregistry/releases/download/v0.8.0/go-containerregistry_Linux_x86_64.tar.gz -o /tmp/go-containerregistry.tar.gz; \
    tar -xf /tmp/go-containerregistry.tar.gz -C /tmp; \
    mv /tmp/crane /usr/local/bin; \
    crane export nginx:${NGINX_VERSION} /sqlite.tar;

FROM debian:jessie-slim

#
# Base utilities and tools
#
RUN set -xe; \
    cp /etc/apt/sources.list /tmp/sources.list; \
    sed -i 's/deb /deb-src /' /tmp/sources.list; \
    cat /tmp/sources.list >> /etc/apt/sources.list; \
    apt-get update; \
    apt install -y \
      build-essential \
      git \
      apt-transport-https \
      ca-certificates \
      curl \
      bc \
      libssl-dev \
      wget \
      unzip \
      openssl;


ENV SQLITE_BENCH_HASH=b3f6b4c709e38093c9b89598a9ff5fe3f60df6f5
ENV WORKDIR="/root/workspace"

COPY ./sqlite-bench/serve_file.c /usr/src/sqlite/serve_file.c
COPY ./sqlite-bench/0001-Add-compatibility-with-external-scrapers.patch /usr/src/sqlite/0001-Add-compatibility-with-external-scrapers.patch
COPY ./sqlite-bench/0002-Add-option-to-execute-custom-commands.patch /usr/src/sqlite/0002-Add-option-to-execute-custom-commands.patch

#
# Bechtool build
#
RUN set -xe; \
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
    tar -cvf /usr/src/sqlite_bench.tar /usr/src/sqlite; \
    mkdir -p ${WORKDIR}; \
    cp /usr/src/sqlite_bench.tar ${WORKDIR}/sqlite_bench.tar

#
# Grab and prepare Lupine/Linux sources
#
ENV LUPINE_HASH=3729c7570e5f1b25e786171d0a19cef858b4923f
ENV LINUX_HASH=39a8804455fb23f09157341d3ba7db6d7ae6ee76
ENV OSV_HASH=3e658b8ab32a424ae58a24d01d5b2e22da61932e
ENV WORKDIR="/root/workspace"
ENV MICROVM_CFG="${WORKDIR}/configs/microvm.config"
ENV SQLITE_CFG="${WORKDIR}/configs/sqlite.config"

COPY ./patches/0001-init-Increase-env-var-limit.patch /0001-init-Increase-env-var-limit.patch
COPY ./patches/0001-arch-x86-Increase-cmdline-size.patch /0001-arch-x86-Increase-cmdline-size.patch

RUN set -xe; \
    mkdir -p ${WORKDIR}; \
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
    sed -i -e "s/docker run -it -v \"\$(PWD)\/linux\":\/linux-volume --rm linuxbuild:latest//" Makefile; \
    sed -i -e "s/linux-volume/root\/workspace\/linux/" Makefile; \
    sed -i -e 's/make patch-linux || true/echo "Keep everything as it is"/g' ${WORKDIR}/scripts/build-with-configs.sh; \
    sed -i -e 's/CONFIGS=$@/CONFIGS=${@:2}/g' ${WORKDIR}/scripts/build-with-configs.sh; \
    sed -i -e 's/-j8/-j $(( $(nproc) + 1 ))/g' ${WORKDIR}/Makefile; \
    touch ${SQLITE_CFG}; \
    ./scripts/build-with-configs.sh patch ${MICROVM_CFG} ${SQLITE_CFG}

#
# Local tooling and filesystem
#
COPY ./kconfig-set.sh /usr/local/bin/
COPY ./image2rootfs.sh /usr/local/bin/

# Overwrite the start scripts with custom functionality, allowing for
# use with Wayfinder's test.benchtool.args:
COPY ./guest_start.sh ${WORKDIR}/scripts/guest_start.sh


# Copy from the initial build step the pre-fetched filesystem
COPY --from=sqlite-filesystem /sqlite.tar ${WORKDIR}/sqlite.tar
