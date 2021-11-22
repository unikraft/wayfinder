FROM debian:9.13 AS wrk

RUN set -xe; \
    apt-get update; \
    apt-get install -y \
      build-essential \
      unzip \
      curl \
      git; \
    git clone https://github.com/wg/wrk.git /wrk; \
    cd /wrk; \
    make

FROM debian:9.13 AS test

RUN set -xe; \
    apt-get update; \
    apt-get install -y \
      curl

COPY --from=wrk /wrk/wrk /usr/local/bin/wrk
COPY ./test.sh /test.sh
