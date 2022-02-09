FROM debian:jessie-slim

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y curl wget git build-essential ca-certificates  gnupg lsb-release apt-transport-https sudo && \
    curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list && \
    apt-get update && \
    apt-get install -y --force-yes docker-ce && \
    mkdir -p /etc/docker/ && \
    echo '{"storage-driver": "vfs"}' >> /etc/docker/daemon.json && \
    rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/hlef/Lupine-Linux.git && \
	cd Lupine-Linux && \
	git checkout b9dc99bbd09180b0a3548583d58f9c003d4576e8 && \
	git submodule update --init;

RUN apt-get update && \
    apt-get install -y libssl-dev openssl bc && \
    rm -rf /var/lib/apt/lists/*

RUN apt-get update && \
    apt-get install -y sudo && \
    rm -rf /var/lib/apt/lists/*

RUN cd Lupine-Linux && \
    sed -i -e "s/redis nginx/nginx/" scripts/build-kernels.sh && \
    cp configs/lupine-djw-kml.config configs/lupine-djw-kml-qemu.config && \
    make patch-linux || true && \
    rm -f linux/.config || true && \
    cat configs/lupine-djw-kml-qemu.config >> linux/.config && \
    cat configs/apps/nginx.config >> linux/.config && \
    cd linux/ && \
    scripts/config --set-val CONFIG_PCI y && \
    scripts/config --set-val CONFIG_VIRTIO_BLK_SCSI y && \
    scripts/config --set-val CONFIG_VIRTIO_PCI_LEGACY y && \
    scripts/config --set-val CONFIG_VIRTIO_PCI y && \
    make oldconfig && \
    make -j8 && \
    cd ../ && \
    mkdir -p kernelbuild/lupine-djw-kml-qemu++nginx/ && \
    cp linux/vmlinux  kernelbuild/lupine-djw-kml-qemu++nginx/ && \
    INSTALL_PATH=../kernelbuild/lupine-djw-kml-qemu++nginx/ make -C linux install && \
    cd load_entropy && \
    make  && \
    cd ../../../

## Copy config files in the container
COPY nginx-nocaching.conf /
COPY nginx-caching.conf /
COPY guest_start.sh /
