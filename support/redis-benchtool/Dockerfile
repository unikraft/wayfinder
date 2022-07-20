FROM debian:9.13 AS redis-benchtool

RUN set -xe; \
    apt-get update; \
    apt-get install -y \
      inetutils-ping \
      bash \
      curl \
      grep \
      gawk \
      redis-server;

COPY ./test.sh /test.sh