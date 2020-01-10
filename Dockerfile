FROM alpine:3.6
ADD bin/linux/couchbase-exporter /usr/local/bin
ENTRYPOINT ["couchbase-exporter"]
