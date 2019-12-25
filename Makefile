SOURCE = $(shell find . -name *.go -type f)
bldNum = $(if $(BLD_NUM),$(BLD_NUM),9999)
version = $(if $(VERSION),$(VERSION),1.0.0)
productVersion = $(version)-$(bldNum)
ARTIFACTS = build/artifacts/exporter

build: $(SOURCE) go.mod
	mkdir -p build
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -ldflags="-s -w" -o build/couchbase-exporter

image-artifacts: build
	mkdir -p $(ARTIFACTS)
	cp build/couchbase-exporter $(ARTIFACTS)
	cp Dockerfile $(ARTIFACTS)

dist: image-artifacts
	rm -rf dist
	mkdir -p dist
	tar -C $(ARTIFACTS)/.. -czf dist/couchbase-exporter-image_$(productVersion).tgz .
