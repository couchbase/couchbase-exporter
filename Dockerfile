FROM scratch
COPY bin/linux/couchbase-exporter /
ENTRYPOINT ["/couchbase-exporter"]
