FROM unikraft/kraft:staging AS kraft

ARG UK_KRAFT_GITHUB_TOKEN=
ENV UK_KRAFTRC=/usr/src/unikraft/.kraftrc
ENV UK_CACHEDIR=/usr/src/unikraft/.kraftcache
ENV UK_WORKDIR=/usr/src/unikraft 

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
    kraft list add https://github.com/unikraft/lib-lwip.git; \
    kraft list add https://github.com/unikraft/lib-pthread-embedded.git; \
    kraft list add https://github.com/unikraft/lib-newlib.git; \
    kraft list add https://github.com/unikraft/lib-nginx.git; \
    kraft list add https://github.com/unikraft/lib-tinyalloc.git; \
    kraft list add https://github.com/unikraft/lib-mimalloc.git; \
    kraft list add https://github.com/unikraft/lib-tlsf.git; \
    kraft list add https://github.com/unikraft/app-nginx.git; \
    kraft list update; \
    kraft list pull unikraft@0.7.0; \
    kraft list pull lib/pthread-embedded@0.7.0; \
    kraft list pull lib/newlib@0.7.0; \
    kraft list pull lib/lwip@0.7.0; \
    kraft list pull lib/tinyalloc@0.7.0; \
    kraft list pull lib/mimalloc@0.7.0; \
    kraft list pull lib/tlsf@0.7.0; \
    kraft list pull lib/nginx@0.6; \
    kraft list pull --no-deps app/nginx; \
    cd /usr/src/unikraft/apps/nginx; \
    kraft configure \
      -u unikraft@0.7.0 \
      -u pthread-embedded@0.7.0 \
      -u newlib@0.7.0 \
      -u lwip@0.7.0 \
      -u nginx@0.6; \
    kraft build --no-progress --fast; \
    unset UK_KRAFT_GITHUB_TOKEN
