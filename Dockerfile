FROM alpine:3.6
ADD build/couchbase-exporter /usr/local/bin
CMD ["couchbase-exporter"]
