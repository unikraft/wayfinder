#!/bin/sh

set -ex

# Usage: guest_start.sh $IP_ADDR $IP_GW

echo "Init called as: $@"

export PATH="/usr/local/bin:$PATH"

mount -t proc proc /proc
ulimit -n 65535
echo 1024 > /proc/sys/net/core/somaxconn

# ip self, ip gateway
./busybox-x86_64 ip addr add $WAYFINDER_DOMAIN_IP_ADDR/24 dev eth0
./busybox-x86_64 ip addr add 127.0.0.1/24 dev lo
./busybox-x86_64 ip link set eth0 up
./busybox-x86_64 ip link set lo up
./busybox-x86_64 ip route add default via $WAYFINDER_DOMAIN_IP_GW_ADDR dev eth0

tr -dc A-Za-z0-9 < /dev/urandom | head -c $PAYLOAD_SIZE > /data/www/index.html
tr -dc A-Za-z0-9 < /dev/urandom | head -c $PAYLOAD_SIZE > /usr/share/nginx/html/index.html

# Modify values inside the configuration (for NGINX)
if [ "$OPEN_FILE_CACHE" = "nocaching" ]; then
cat <<EOF >/etc/nginx/nginx.conf
worker_processes auto;
error_log stderr;
pid /tmp/nginx.pid;
user daemon daemon;
daemon off;
master_process off;
worker_rlimit_nofile $WORKER_RLIMIT_NOFILE;

events {
  worker_connections $WORKER_CONNECTIONS;
}

http {
  include       mime.types;
  default_type  application/octet-stream;
  log_format    main  '\$remote_addr - \$remote_user [\$time_local] "\$request" '
                      '\$status \$body_bytes_sent "\$http_referer" '
                      '"\$http_user_agent" "\$http_x_forwarded_for"';
  server_tokens $SERVER_TOKENS;

  sendfile $SENDFILE;
  tcp_nopush $TCP_NOPUSH;
  keepalive_timeout $KEEPALIVE_TIMEOUT;

  access_log    /dev/null;

  server {
    listen       80;
    listen       [::]:80;
    server_name  localhost;

    location / {
      root   /data/www;
      index  index.html;
    }
  }
}
EOF
else
cat <<EOF >/etc/nginx/nginx.conf
worker_processes auto;
error_log stderr;
pid /tmp/nginx.pid;
user daemon daemon;
daemon off;
master_process off;
worker_rlimit_nofile $WORKER_RLIMIT_NOFILE;

events {
  worker_connections $WORKER_CONNECTIONS;
}

http {
  include       mime.types;
  default_type  application/octet-stream;

  # caching
  open_file_cache max=200000 inactive=20s;
  open_file_cache_valid 30s;
  open_file_cache_min_uses 2;
  open_file_cache_errors on;

  log_format    main  '\$remote_addr - \$remote_user [\$time_local] "\$request" '
                      '\$status \$body_bytes_sent "\$http_referer" '
                      '"\$http_user_agent" "\$http_x_forwarded_for"';
  server_tokens $SERVER_TOKENS;

  sendfile $SENDFILE;
  tcp_nopush $TCP_NOPUSH;
  keepalive_timeout $KEEPALIVE_TIMEOUT;

  access_log    /dev/null;

  server {
    listen       80;
    listen       [::]:80;
    server_name  localhost;

    location / {
      root   /data/www;
      index  index.html;
    }
  }
}
EOF
fi

# dirty, but not sure how to do this differently
nginx -c /etc/nginx/nginx.conf
