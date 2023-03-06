FROM debian

# Install kernel build dependencies
RUN apt-get update -y
RUN apt-get install -y \
    bc \
    bison \
    build-essential \
    gcc-riscv64-linux-gnu \
    git \
    flex \
    libelf-dev \
    libssl-dev \
    make \
    ncurses-dev \
    tree

# Clone the Linux kernel repository (only version 6.1 is needed)
RUN git clone --depth 1 --branch v6.1 https://github.com/torvalds/linux.git
RUN git clone --depth 1 --branch 1_36_0 https://git.busybox.net/busybox

# Build the kernel
WORKDIR /linux

COPY ./patches/0001-init-Increase-env-var-limit.patch /linux/0001-init-Increase-env-var-limit.patch
RUN patch -p1 < /linux/0001-init-Increase-env-var-limit.patch

RUN make ARCH=riscv CROSS_COMPILE=riscv64-linux-gnu- defconfig
RUN make ARCH=riscv CROSS_COMPILE=riscv64-linux-gnu- -j$(nproc)

COPY kconfig-set.sh /linux/kconfig-set.sh

# Build busybox
WORKDIR /busybox

RUN CROSS_COMPILE=riscv64-linux-gnu- LDFLAGS=--static make defconfig
RUN CROSS_COMPILE=riscv64-linux-gnu- LDFLAGS=--static make -j$(nproc)

# Create a minimal rootfs
WORKDIR /rootfs
RUN CROSS_COMPILE=riscv64-linux-gnu- LDFLAGS=--static make -C /busybox install CONFIG_PREFIX=/rootfs
RUN mkdir -p /rootfs/proc /rootfs/sys /rootfs/dev /rootfs/etc/init.d
COPY ./guest_start.sh /rootfs/etc/init.d/rcS
RUN chmod +x /rootfs/etc/init.d/rcS
RUN touch /rootfs/etc/fstab
