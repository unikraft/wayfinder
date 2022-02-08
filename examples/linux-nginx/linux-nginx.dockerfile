FROM debian:jessie-slim

# =============================================================================
# Setup, dependencies
# =============================================================================

ENV LUPINE_HASH=3729c7570e5f1b25e786171d0a19cef858b4923f
ENV NGINX_CONF="/root/resources/nginx.conf"
ENV GUEST_START="/root/resources/guest_start.sh"
ENV BUILD_FOLDER="/root/generated-data"

WORKDIR /root

# add deb-src entires for apt build-dep
RUN cp /etc/apt/sources.list /tmp/sources.list
RUN sed -i 's/deb /deb-src /' /tmp/sources.list
RUN cat /tmp/sources.list >> /etc/apt/sources.list

# dependencies
RUN apt-get update
RUN apt install -y build-essential git apt-transport-https \
                   ca-certificates curl qemu qemu-system-x86

RUN apt install -y build-essential bc libssl-dev openssl

# copy needed resources
RUN mkdir -p /root/resources
COPY ./resources/* /root/resources/
RUN mkdir -p $BUILD_FOLDER

# =============================================================================
# Linux kernel build
# =============================================================================

ENV LUPINE_DIR="/root/generated-data/Lupine-Linux"
ENV MICROVM_CFG="${LUPINE_DIR}/configs/microvm.config"
ENV NGINX_CFG="${LUPINE_DIR}/configs/apps/nginx.config"

WORKDIR $BUILD_FOLDER
# clone Lupine, we use their scripts as they're convenient even though we're
# actually building microvm and not lupine per se
# RUN git clone https://github.com/hlef/Lupine-Linux.git
COPY ./generated-data/Lupine-Linux Lupine-Linux

WORKDIR $LUPINE_DIR
RUN git checkout $LUPINE_HASH
RUN git submodule update --init

# build qemu/kvm version
RUN echo "CONFIG_PCI=y"               >> $MICROVM_CFG
RUN echo "CONFIG_VIRTIO_BLK_SCSI=y"   >> $MICROVM_CFG
RUN echo "CONFIG_VIRTIO_PCI_LEGACY=y" >> $MICROVM_CFG
RUN echo "CONFIG_VIRTIO_PCI=y"        >> $MICROVM_CFG
RUN echo "CONFIG_VGA_ARB_MAX_GPUS=16" >> $MICROVM_CFG

# no need to build in docker
RUN sed -i -e "s/docker run -it -v \"\$(PWD)\/linux\":\/linux-volume --rm linuxbuild:latest//" Makefile
RUN sed -i -e "s/linux-volume/root\/generated-data\/Lupine-Linux\/linux/" Makefile

# only build for microvm
RUN ./scripts/build-with-configs.sh nopatch $MICROVM_CFG $NGINX_CFG

# =============================================================================
# Cleanup
# =============================================================================

WORKDIR /root

RUN cp $LUPINE_DIR/kernelbuild/microvm++nginx/vmlinuz-4.0.0 \
                   /root/linux-nginx-qemu.kernel
COPY ./generated-data/nginx.ext2 /root/linux-nginx.ext2

# useful to have it here for debugging
RUN apt install -y wget socat uuid-runtime bridge-utils net-tools psmisc vim
RUN wget https://raw.githubusercontent.com/unikraft/kraft/6217d48668cbdf0847c7864bc6368a6adb94f6a6/scripts/qemu-guest
RUN chmod a+x /root/qemu-guest
RUN wget https://github.com/unikraft/eurosys21-artifacts/raw/master/tools/wrk
RUN chmod u+x wrk
COPY ./resources/sanity_check.sh /root/
COPY ./resources/kconfig-set.sh /root/
