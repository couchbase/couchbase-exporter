FROM alpine:3.6
ADD couchbase-exporter /usr/local/bin
ADD main.go /usr/local/bin
CMD ["couchbase-exporter"]
