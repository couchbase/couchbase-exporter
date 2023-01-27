SOURCE = $(shell find . -name *.go -type f)
bldNum = $(if $(BLD_NUM),$(BLD_NUM),9999)
version = $(if $(VERSION),$(VERSION),1.0.0)
productVersion = $(version)-$(bldNum)
# The git revision, infinitely more useful than an arbitrary build number.
REVISION := $(shell git rev-parse HEAD)

ARTIFACTS = build/artifacts/couchbase-exporter
DOCKER_TAG = v1
DOCKER_USER = couchbase
.PHONY: all test clean

GOPATH := $(shell go env GOPATH)
GOBIN := $(if $(GOPATH),$(GOPATH)/bin,$(HOME)/go/bin)
GOLINT_VERSION := v1.42.1
# These are propagated into each binary so we can tell for sure the exact build
# that a binary came from.
LDFLAGS = \
  -s -w \
  -X github.com/couchbase/couchbase-exporter/pkg/version.Version=$(version) \
  -X github.com/couchbase/couchbase-exporter/pkg/version.BuildNumber=$(bldNum) \
  -X github.com/couchbase/couchbase-exporter/pkg/revision.gitRevision=$(REVISION)

build: $(SOURCE) go.mod
	for platform in linux darwin ; do \
		for arch in amd64 arm64 ; do \
	  	echo "Building $$platform $$arch binary " ; \
	  	GOOS=$$platform GOARCH=$$arch CGO_ENABLED=0 GO111MODULE=on go build -ldflags="$(LDFLAGS)" -o bin/$$platform/couchbase-exporter-$$arch ; \
	  done \
	done

image-artifacts: build
	mkdir -p $(ARTIFACTS)/bin/linux
	cp bin/linux/couchbase-exporter-* $(ARTIFACTS)/bin/linux/
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
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLINT_VERSION)
	$(GOBIN)/golangci-lint run ./pkg/... ./test/...
