FROM unikraft/kraft:staging AS kraft

ARG UK_KRAFT_GITHUB_TOKEN=
ENV UK_KRAFTRC=/usr/src/unikraft/.kraftrc
ENV UK_CACHEDIR=/usr/src/unikraft/.kraftcache
ENV UK_WORKDIR=/usr/src/unikraft 

LABEL org.opencontainers.image.source="https://github.com/unikraft/wayfinder"
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
      libarchive-tools; \
    echo "" > $KRAFTRC; \
    kraft list add https://github.com/unikraft/unikraft.git; \
    kraft list add https://github.com/lancs-net/lib-lwip.git; \
    kraft list add https://github.com/unikraft/lib-pthread-embedded.git; \
    kraft list add https://github.com/unikraft/lib-newlib.git; \
    kraft list add https://github.com/unikraft/lib-nginx.git; \
    kraft list add https://github.com/unikraft/app-nginx.git; \
    kraft list update; \
    kraft list pull unikraft@staging; \
    kraft list pull lib/pthread-embedded@staging; \
    kraft list pull lib/newlib@staging; \
    kraft list pull lib/lwip@staging; \
    kraft list pull lib/nginx@staging; \
    kraft list pull --no-deps app/nginx; \
    cd /usr/src/unikraft/apps/nginx; \
    kraft configure \
      -u unikraft@staging \
      -u pthread-embedded@staging \
      -u newlib@staging \
      -u lwip@staging \
      -u nginx@staging; \
    kraft build --no-progress --fast; \
    unset UK_KRAFT_GITHUB_TOKEN
