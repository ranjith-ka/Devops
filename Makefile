OS_NAME := $(shell uname -s | tr A-Z a-z)
TIMESTAMP ?= $$(date -u +'%Y%m%d%H%M%S')
PROTOCV := 3.11.4

# strip patch version number
PROTOCVMINOR ?= $(subst $(suffix $(PROTOCV)),,$(PROTOCV))

ifeq ($(findstring mingw64,$(OS_NAME)),mingw64)
PROTOZIP ?= protoc-$(PROTOCV)-win64.zip
PROTOLOC ?= /c/protoc
export PATH := $(PATH):/c/protoc/bin
endif

ifeq ($(findstring linux,$(OS_NAME)),linux)
PROTOZIP ?= protoc-$(PROTOCV)-linux-x86_64.zip
PROTOLOC ?= /usr/local
export PATH := $(PATH):/usr/local/bin
endif

ifeq ($(findstring darwin,$(OS_NAME)),darwin)
PROTOZIP ?= protoc-$(PROTOCV)-osx-x86_64.zip
PROTOLOC ?=  $(HOME)/.protobuf
export PATH := $(PATH):$(HOME)/.protobuf/bin
endif


default: build

install: install-protoc

install-protoc:
	@echo Installing protoc $(PROTOCV) binaries to $(PROTOLOC)
ifeq ($(findstring mingw64,$(OS_NAME)),mingw64)
	@curl -sSLO https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOCV)/$(PROTOZIP)
	@mkdir -p /c/protoc
	@unzip -o $(PROTOZIP) -d $(PROTOLOC) bin/protoc* > /dev/null
	@unzip -o $(PROTOZIP) -d $(PROTOLOC) include/* > /dev/null
	@rm -f $(PROTOZIP)
	@echo
	@echo Please update your PATH to include $(PROTOLOC)/bin
endif

ifeq ($(findstring darwin,$(OS_NAME)),darwin)
	@mkdir -p $(HOME)/.protobuf/bin
	@cd $(HOME)/.protobuf
	@curl -sSLO https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOCV)/$(PROTOZIP)
	@unzip -o $(PROTOZIP) -d $(PROTOLOC) bin/protoc > /dev/null
	@unzip -o $(PROTOZIP) -d $(PROTOLOC) include/* > /dev/null
	@rm -f $(PROTOZIP)
endif

ifeq ($(findstring linux,$(OS_NAME)),linux)
	@curl -sSLO https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOCV)/$(PROTOZIP)
	@sudo unzip -o $(PROTOZIP) -d $(PROTOLOC) bin/protoc > /dev/null
	@sudo unzip -o $(PROTOZIP) -d $(PROTOLOC) include/* > /dev/null
	@rm -f $(PROTOZIP)
endif

download:
	@echo Download go.mod dependencies
	@go mod download

install-tools: download
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

build: generate compile

generate:
	@echo Generate code
	@go generate ./...

compile:
	@echo Compile binary
	@go build .

run: build
	./Docker random

reports:
	@mkdir -p reports

bench: reports
	@echo Writing benchmark to reports/bench$(TIMESTAMP).txt
	@go test -benchmem -run=^$$ <repo> -bench . > reports/bench$(TIMESTAMP).txt

coverage: reports
	@go test -race -covermode=atomic -coverprofile=reports/coverage.out ./...
	@go tool cover -func=reports/coverage.out

dirty:
ifneq ($(DIRTY),)
	@echo modified/untracked files; echo $(DIRTY); exit 1
else
	@echo 'clean'
endif
