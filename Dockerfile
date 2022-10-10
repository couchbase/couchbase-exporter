FROM scratch
ARG TARGETARCH
COPY bin/linux/couchbase-exporter-${TARGETARCH} /couchbase-exporter
ENTRYPOINT ["/couchbase-exporter"]
