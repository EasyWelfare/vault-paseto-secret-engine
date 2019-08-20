GOARCH = amd64

UNAME = $(shell uname -s)

ifndef OS
	ifeq ($(UNAME), Linux)
		OS = linux
	else ifeq ($(UNAME), Darwin)
		OS = darwin
	endif
endif

.DEFAULT_GOAL := all
PLUGIN_DIR := plugins

all: fmt build start

build:
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o "$(PLUGIN_DIR)"/paseto main.go

start:
	vault server -dev -dev-root-token-id=root -log-level=debug -dev-plugin-dir="$(PLUGIN_DIR)"

enable:
	vault secrets enable -path=test paseto

config:
	vault write test/paseto/config footer="" ttl=10
	vault read test/paseto/config

read-token:
	vault read test/paseto/token

clean:
	rm -f ./vault/plugins/paseto

fmt:
	go fmt $$(go list ./...)

.PHONY: build clean fmt start enable
