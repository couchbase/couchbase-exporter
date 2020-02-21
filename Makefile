SOURCE = $(shell find . -name *.go -type f)
bldNum = $(if $(BLD_NUM),$(BLD_NUM),9999)
version = $(if $(VERSION),$(VERSION),1.0.0)
productVersion = $(version)-$(bldNum)
ARTIFACTS = build/artifacts/couchbase-exporter

build: $(SOURCE) go.mod
	for platform in linux darwin ; do \
	  echo "Building $$platform binary" ; \
	  mkdir -p build/$$platform ; \
	  GOOS=$$platform GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -ldflags="-s -w" -o bin/$$platform/couchbase-exporter ; \
	done

image-artifacts: build
	mkdir -p $(ARTIFACTS)/bin/linux
	cp bin/linux/couchbase-exporter $(ARTIFACTS)/bin/linux
	cp Dockerfile* LICENSE README.md $(ARTIFACTS)

dist: image-artifacts
	rm -rf dist
	mkdir -p dist
	tar -C $(ARTIFACTS)/.. -czf dist/couchbase-exporter-image_$(productVersion).tgz .
