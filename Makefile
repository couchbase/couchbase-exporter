SOURCE = $(shell find . -name *.go -type f)
bldNum = $(if $(BLD_NUM),$(BLD_NUM),9999)
version = $(if $(VERSION),$(VERSION),1.0.0)
productVersion = $(version)-$(bldNum)
ARTIFACTS = build/artifacts/couchbase-exporter
DOCKER_TAG = v1
DOCKER_USER = couchbase
.PHONY: all test clean

build: $(SOURCE) go.mod
	for platform in linux darwin ; do \
	  echo "Building $$platform binary" ; \
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

container:
	docker build -f Dockerfile -t ${DOCKER_USER}/couchbase-exporter:${DOCKER_TAG} .

config-container:
	docker build -f ./example/Dockerfile -t ${DOCKER_USER}/couchbase-exporter:${DOCKER_TAG}-config ./example/

test:
	go test ./test --timeout 60s

gen:
	go install github.com/golang/mock/mockgen@v1.6.0
	$$GOPATH/bin/mockgen -source=./pkg/util/networking.go -destination=./test/mocks/mock_client.go -package mocks

lint: 
		go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./pkg/... ./test/... ./