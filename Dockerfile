FROM alpine:3.6
ADD couchbase_exporter /usr/local/bin
ADD main.go /usr/local/bin
CMD ["couchbase_exporter"]
