GOARCH = amd64

UNAME = $(shell uname -s)
VAULT_ADDR = http://localhost:8200
VAULT_TOKEN = root

ifndef OS
	ifeq ($(UNAME), Linux)
		OS = linux
	else ifeq ($(UNAME), Darwin)
		OS = darwin
	endif
endif

.DEFAULT_GOAL := all
PLUGIN_DIR := plugins
PLUGIN_PATH := mytest

all: fmt build start

build:
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o $(PLUGIN_DIR)/paseto main.go

start:
	vault server -dev -dev-root-token-id=root -log-level=debug -dev-plugin-dir="$(PLUGIN_DIR)"

list: 
	vault secrets list

enable:
	vault secrets enable -path=${PLUGIN_PATH}/paseto paseto

config:
	vault write ${PLUGIN_PATH}/paseto/config footer="test" ttl=10
	vault read ${PLUGIN_PATH}/paseto/config

read-token:
	vault read ${PLUGIN_PATH}/paseto/token

clean:
	rm -f ./vault/plugins/paseto

fmt:
	go fmt $$(go list ./...)

.PHONY: build clean fmt start enable
