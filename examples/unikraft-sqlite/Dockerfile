FROM unikraft/kraft:staging AS kraft

ARG UK_KRAFT_GITHUB_TOKEN=
ENV UK_KRAFTRC=/usr/src/unikraft/.kraftrc
ENV UK_CACHEDIR=/usr/src/unikraft/.kraftcache
ENV UK_WORKDIR=/usr/src/unikraft 
ENV SQLITE_BENCH_HASH=b3f6b4c709e38093c9b89598a9ff5fe3f60df6f5

COPY ./sqlite-bench/Makefile.uk /usr/src/unikraft/apps/sqlite/Makefile.uk
COPY ./sqlite-bench/serve_file.c /usr/src/unikraft/apps/sqlite/serve_file.c
COPY ./sqlite-bench/0001-Add-compatibility-with-external-scrapers.patch /usr/src/sqlite-bench/0001-Add-compatibility-with-external-scrapers.patch

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
    kraft list add https://github.com/unikraft/lib-sqlite.git; \
    kraft list add https://github.com/unikraft/lib-tinyalloc.git; \
    kraft list add https://github.com/unikraft/lib-mimalloc.git; \
    kraft list add https://github.com/unikraft/lib-tlsf.git; \
    kraft list add https://github.com/unikraft/app-sqlite.git; \
    kraft list update; \
    kraft list pull unikraft@0.9.0; \
    kraft list pull lib/pthread-embedded@0.9.0; \
    kraft list pull lib/newlib@0.9.0; \
    kraft list pull lib/lwip@0.9.0; \
    kraft list pull lib/tinyalloc@0.9.0; \
    kraft list pull lib/mimalloc@0.9.0; \
    kraft list pull lib/tlsf@0.9.0; \
    kraft list pull lib/sqlite@0.9.0; \
    kraft list pull --no-deps app/sqlite@0.9.0; \
    cd /usr/src/sqlite-bench; \
    git init; \
    git remote add origin https://github.com/ukontainer/sqlite-bench.git; \
    git fetch --depth 1 origin ${SQLITE_BENCH_HASH}; \
    git checkout FETCH_HEAD; \
    git apply 0001-Add-compatibility-with-external-scrapers.patch; \
    cd /usr/src/unikraft/apps/sqlite; \
    cp /usr/src/sqlite-bench/bench.h /usr/src/sqlite-bench/benchmark.c \
      /usr/src/sqlite-bench/histogram.c /usr/src/sqlite-bench/main.c \
      /usr/src/sqlite-bench/random.c /usr/src/sqlite-bench/raw.c \
      /usr/src/sqlite-bench/sqlite3.h /usr/src/sqlite-bench/sqlite3ext.h \
      /usr/src/sqlite-bench/util.c /usr/src/unikraft/apps/sqlite; \
    sed -i 's/UK_FEATURE_RXQ_INTR_AVAILABLE/UK_NETDEV_F_RXQ_INTR/g' /usr/src/unikraft/libs/lwip/uknetdev.c; \
    sed -i 's/libraries:/libraries:\n  lwip:\n    version: 0.9.0/g' /usr/src/unikraft/apps/sqlite/kraft.yaml; \
    kraft -v configure -F \
      -u unikraft@0.9.0 \
      -u pthread-embedded@0.9.0 \
      -u newlib@0.9.0 \
      -u lwip@0.9.0 \
      -u lib/sqlite@0.9.0 \
      -u app/sqlite@0.9.0 \
      --arch x86_64 \
      --plat kvm \
      --yes LIBVFSCORE_ROOTFS_INITRD \
      --yes LIBCPIO \
      --yes LIBUKNETDEV \
      --yes VIRTIO_PCI \
      --no  LIB9PFS \
      --no  LIBUK9P \
      --no LIBSQLITE_MAIN_FUNCTION; \
    set -xe; \
    kraft build --no-fetch --no-progress --fast; \
    unset UK_KRAFT_GITHUB_TOKEN
