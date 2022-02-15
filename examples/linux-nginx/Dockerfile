FROM debian:jessie-slim AS nginx-filesystem

ARG NGINX_VERSION=1.15.6-alpine

RUN set -xe; \
    apt-get update; \
    apt install -y \
      curl; \
    curl -L https://github.com/google/go-containerregistry/releases/download/v0.8.0/go-containerregistry_Linux_x86_64.tar.gz -o /tmp/go-containerregistry.tar.gz; \
    tar -xf /tmp/go-containerregistry.tar.gz -C /tmp; \
    mv /tmp/crane /usr/local/bin; \
    crane export nginx:${NGINX_VERSION} /nginx.tar;

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
      openssl;

#
# Grab and prepare Lupine/Linux sources
#
ENV LUPINE_HASH=3729c7570e5f1b25e786171d0a19cef858b4923f
ENV WORKDIR="/root/workspace"
ENV MICROVM_CFG="${WORKDIR}/configs/microvm.config"
ENV NGINX_CFG="${WORKDIR}/configs/apps/nginx.config"
ENV NGINX_CONF="/root/resources/nginx.conf"

RUN set -xe; \
    git clone https://github.com/hlef/Lupine-Linux.git ${WORKDIR}; \
    cd ${WORKDIR}; \
    git checkout ${LUPINE_HASH}; \
    git submodule update --init; \
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
    ./scripts/build-with-configs.sh nopatch ${MICROVM_CFG} ${NGINX_CFG}

#
# Local tooling and filesystem
#
COPY ./kconfig-set.sh /usr/local/bin/
COPY ./image2rootfs.sh /usr/local/bin/

# Overwrite the start scripts with custom functionality, allowing for
# use with Wayfinder's test.benchtool.args:
COPY ./guest_start.sh ${WORKDIR}/scripts/guest_start.sh

# Copy from the initial build step the pre-fetched filesystem
COPY --from=nginx-filesystem /nginx.tar ${WORKDIR}/nginx.tar
