FROM alpine:3.6
ADD couchbase-exporter /usr/local/bin
CMD ["couchbase-exporter"]
