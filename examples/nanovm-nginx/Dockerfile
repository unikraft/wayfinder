FROM debian:latest

COPY ./nginx-caching.conf /nginx-caching.conf
COPY ./nginx-nocaching.conf /nginx-nocaching.conf
COPY ./test.sh /test.sh

RUN set -xe; \
    apt-get update; \
    apt-get install -y \
      jq \
      curl \
      iperf3 \
      procps \
      bridge-utils \
      net-tools \
      qemu-kvm \
      qemu-utils \
      qemu-system-x86 \
      libarchive-tools; \
    curl https://ops.city/get.sh -sSfL | sh; \
    /root/.ops/bin/ops pkg get nginx_1.18.0; \
    cd /root/.ops/packages; \
    /root/.ops/bin/ops image create --package nginx_1.18.0 --target-cloud onprem;
