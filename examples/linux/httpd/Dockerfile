FROM ghcr.io/unikraft/wayfinder/linux-base

ENV APACHE_CFG="${WORKDIR}/configs/apps/httpd.config"

ARG APACHE_VERSION=2.4
RUN crane export httpd:${APACHE_VERSION} ${WORKDIR}/httpd.tar;

RUN ./scripts/build-with-configs.sh patch ${MICROVM_CFG} ${APACHE_CFG};

COPY ./guest_start.sh ${WORKDIR}/scripts/guest_start.sh
