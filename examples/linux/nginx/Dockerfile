FROM ghcr.io/unikraft/wayfinder/linux-base

ENV NGINX_CFG="${WORKDIR}/configs/apps/nginx.config"

ARG NGINX_VERSION=1.15.6
RUN crane export nginx:${NGINX_VERSION} ${WORKDIR}/nginx.tar;

RUN ./scripts/build-with-configs.sh patch ${MICROVM_CFG} ${NGINX_CFG}

COPY ./guest_start.sh ${WORKDIR}/scripts/guest_start.sh
