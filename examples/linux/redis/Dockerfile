FROM ghcr.io/unikraft/wayfinder/linux-base

ENV REDIS_CFG="${WORKDIR}/configs/apps/redis.config"

ARG REDIS_VERSION=7.0.4
RUN crane export redis:${REDIS_VERSION} ${WORKDIR}/redis.tar

RUN ./scripts/build-with-configs.sh patch ${MICROVM_CFG} ${REDIS_CFG}

COPY ./redis-template.conf ${WORKDIR}/redis-template.conf
COPY ./guest_start.sh ${WORKDIR}/scripts/guest_start.sh
