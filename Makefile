PROJECT := servant

all: help

help:
	@echo "make test - run go test"
	@echo "make build - build servant"
	@echo "make release - create an tag depending on the VERSION file"
	@echo "make build-travis - compiles binaries for x64 mac/linux and creates release tar.gz files with hashsums"

build:
	go build

test:
	go test -v

release:
	git tag `cat VERSION`
	git push origin master --tags

build-dir:
	@rm -rf build && mkdir build

dist-dir:
	@rm -rf dist && mkdir dist

ci-compile: build-dir
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/darwin-amd64/servant ./
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/linux-amd64/servant ./

build-travis: ci-compile dist-dir
	$(eval FILES := $(shell ls build))
	@for f in $(FILES); do \
		(cd $(shell pwd)/build/$$f && tar -cvzf ../../dist/$$f.tar.gz *); \
		(cd $(shell pwd)/dist && shasum -a 256 $$f.tar.gz > $$f.sha256); \
		(cd $(shell pwd)/dist && md5sum $$f.tar.gz > $$f.md5); \
		echo $$f; \
	done