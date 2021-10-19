SHELL := /bin/bash
BASEDIR := $(shell pwd)
SYNC_CMD := synctl
PACK_CMD := packer
DEFAULT_VERSION := v1.0.0
version ?= $(DEFAULT_VERSION)
build_time := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
build_commit := $(shell git rev-parse --short HEAD)
packer_cmd_path := github.com/shipengqi/keel-pkg/app/packer/cmd
sync_cmd_path := github.com/shipengqi/keel-pkg/app/packer/cmd

.PHONY: all show setup build build-sync build-pack test benchmark clean help
all: help


.ONESHELL:
build: build-sync build-pack
	@echo "build $(SYNC_CMD) and $(PACK_CMD) completed."

# -tags="containers_image_openpgp" https://github.com/containers/image/issues/268
build-sync:
	@echo "building command: $(SYNC_CMD) ..."
	@CGO_ENABLED=0 go build -tags="containers_image_openpgp" -mod=mod -o synctl -ldflags \
	"-X '$(sync_cmd_path).Version=$(version)' -X '$(sync_cmd_path).BuildTime=$(build_time)' -X '$(sync_cmd_path).GitCommit=$(build_commit)'" \
	./app/synchronizer/

build-pack:
	@echo "building command: $(PACK_CMD) ..."
	@CGO_ENABLED=0 go build -mod=mod -o packer -ldflags \
	"-X '$(packer_cmd_path).Version=$(version)' -X '$(packer_cmd_path).BuildTime=$(build_time)' -X '$(packer_cmd_path).GitCommit=$(build_commit)'" \
	./app/packer/

test:
	@echo "go test ..."

benchmark:
	@echo "go test benchmark ..."

clean:
	@echo "cleaning ..."
	@rm -rf ./$(SYNC_CMD)
	@rm -rf ./$(PACK_CMD)

.ONESHELL:
show:
	@echo "current directory: $(BASEDIR)"
	@echo CGO_ENABLED=0 go build -tags="containers_image_openpgp" -mod=mod -o synctl -ldflags \
          	"-X '$(sync_cmd_path).Version=$(version)' -X '$(sync_cmd_path).BuildTime=$(build_time)' -X '$(sync_cmd_path).GitCommit=$(build_commit)'" \
          	./app/synchronizer/
	@echo CGO_ENABLED=0 go build -mod=mod -o packer -ldflags \
          	"-X '$(packer_cmd_path).Version=$(version)' -X '$(packer_cmd_path).BuildTime=$(build_time)' -X '$(packer_cmd_path).GitCommit=$(build_commit)'" \
          	./app/packer/

help:
	@echo "make              - help"
	@echo "make build        - build synctl and packer"
	@echo "     version        the version of commands, default is 'v1.0.0'. e.g. 'make build version=v1.1.2'"
	@echo "make build-sync   - build synctl"
	@echo "     version        the version of synctl command, default is 'v1.0.0'. e.g. 'make build-sync version=v1.1.2'"
	@echo "make build-pack   - build packer"
	@echo "     version        the version of packer command, default is 'v1.0.0'. e.g. 'make build-pack version=v1.1.2'"
	@echo "make clean        - remove binary file and prune image"
