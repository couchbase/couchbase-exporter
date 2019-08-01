FROM alpine:3.6
ADD main /usr/local/bin
ADD main.go /usr/local/bin
CMD ["main"]
