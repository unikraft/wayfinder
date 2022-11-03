FROM debian:jessie-slim AS base

ENV LUPINE_HASH=3729c7570e5f1b25e786171d0a19cef858b4923f
ENV LINUX_HASH=39a8804455fb23f09157341d3ba7db6d7ae6ee76
ENV OSV_HASH=3e658b8ab32a424ae58a24d01d5b2e22da61932e
ENV WORKDIR="/root/workspace"
ENV MICROVM_CFG="${WORKDIR}/configs/microvm.config"

# Base utilities and tools
RUN cp /etc/apt/sources.list /tmp/sources.list
RUN sed -i 's/deb /deb-src /' /tmp/sources.list
RUN cat /tmp/sources.list >> /etc/apt/sources.list

RUN apt-get update
RUN apt-get install -y \
    build-essential \
    git \
    apt-transport-https \
    ca-certificates \
    curl \
    bc \
    libssl-dev \
    libgomp1 \
    openssl \
    zip

# Crane
RUN curl -L https://github.com/google/go-containerregistry/releases/download/v0.8.0/go-containerregistry_Linux_x86_64.tar.gz -o /tmp/go-containerregistry.tar.gz
RUN tar -xf /tmp/go-containerregistry.tar.gz -C /tmp
RUN mv /tmp/crane /usr/local/bin

# Lupine/Linux
WORKDIR ${WORKDIR}

RUN git init
RUN git remote add origin https://github.com/hlef/Lupine-Linux.git
RUN git fetch --depth 1 origin ${LUPINE_HASH}
RUN git checkout FETCH_HEAD
RUN git checkout ${LUPINE_HASH}

WORKDIR ${WORKDIR}/linux

RUN git init
RUN git remote add origin https://github.com/torvalds/linux.git
RUN git fetch --depth 1 origin ${LINUX_HASH}
RUN git checkout FETCH_HEAD
RUN git checkout ${LINUX_HASH}

WORKDIR ${WORKDIR}/osv

RUN git init
RUN git remote add origin https://github.com/cloudius-systems/osv.git
RUN git fetch --depth 1 origin ${OSV_HASH}
RUN git checkout FETCH_HEAD

WORKDIR ${WORKDIR}/load_entropy

RUN make

# Copy and apply patches
WORKDIR ${WORKDIR}/linux

COPY ../patches /patches
RUN git apply /patches/0001-init-Increase-env-var-limit.patch
RUN git apply /patches/0001-arch-x86-Increase-cmdline-size.patch

# Build base kernel
WORKDIR ${WORKDIR}

RUN echo "CONFIG_PCI=y"               >> ${MICROVM_CFG}
RUN echo "CONFIG_VIRTIO_BLK_SCSI=y"   >> ${MICROVM_CFG}
RUN echo "CONFIG_VIRTIO_PCI_LEGACY=y" >> ${MICROVM_CFG}
RUN echo "CONFIG_VIRTIO_PCI=y"        >> ${MICROVM_CFG}
RUN echo "CONFIG_VGA_ARB_MAX_GPUS=16" >> ${MICROVM_CFG}
RUN sed -i -e "s/docker run -it -v \"\$(PWD)\/linux\":\/linux-volume --rm linuxbuild:latest//" Makefile
RUN sed -i -e "s/linux-volume/root\/workspace\/linux/" Makefile
RUN sed -i -e 's/make patch-linux || true/echo "Keep everything as it is"/g' ${WORKDIR}/scripts/build-with-configs.sh
RUN sed -i -e 's/CONFIGS=$@/CONFIGS=${@:2}/g' ${WORKDIR}/scripts/build-with-configs.sh
RUN sed -i -e 's/-j8/-j $(( $(nproc) + 1 ))/g' ${WORKDIR}/Makefile
RUN ./scripts/build-with-configs.sh patch ${MICROVM_CFG}

# Copy scripts
COPY ../scripts/image2rootfs.sh /usr/local/bin/image2rootfs.sh
COPY ../scripts/kconfig-set.sh /usr/local/bin/kconfig-set.sh
COPY ../scripts/set-kernel-build.sh /usr/local/bin/set-kernel-build.sh
COPY ../scripts/set-kernel-runtime.sh ${WORKDIR}/scripts/set-kernel-runtime.sh
