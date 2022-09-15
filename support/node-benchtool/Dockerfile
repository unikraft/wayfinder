FROM debian:11.4 AS redis-benchtool

RUN set -xe; \
    apt-get update; \
    apt-get install -y \
      inetutils-ping \
      bash \
      curl \
      grep \
      gawk;

COPY ./test.sh /test.sh