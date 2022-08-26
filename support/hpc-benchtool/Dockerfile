FROM debian:11.4 AS hpc-benchtool

RUN set -xe; \
    apt-get update; \
    apt-get install -y \
      inetutils-ping \
      bash \
      curl \
      grep \
      gawk \
      bc;

COPY ./test.sh /test.sh