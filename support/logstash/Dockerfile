FROM docker.elastic.co/logstash/logstash:6.2.4

USER root

RUN set -xe; \
    rm -f /usr/share/logstash/pipeline/logstash.conf; \
    bin/logstash-plugin install --version 5.0.3 logstash-output-influxdb; \
    bin/logstash-plugin install logstash-filter-metricize; \
    yum install gettext -y

USER logstash

COPY pipeline/ /usr/share/logstash/pipeline-templates/
COPY config/ /usr/share/logstash/config/
COPY init/ /opt/docker-init/

ENV INPUT_PORT=12345
ENV OUTPUT_HOST=influxdb
ENV OUTPUT_PORT=8086
ENV OUTPUT_USER=
ENV OUTPUT_PASS=

ENTRYPOINT [ "/opt/docker-init/entrypoint" ]

