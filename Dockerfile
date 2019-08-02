FROM golang:1.12.6-stretch AS builder

WORKDIR /plugin
COPY ./go.mod ./go.sum ./main.go ./
COPY ./vendor ./vendor
COPY ./cmd ./cmd
COPY ./paseto ./paseto
COPY ./vault ./vault
RUN GOOS=linux go build -a -mod vendor -o paseto-plugin

FROM vault:1.2.0
COPY ./vault/paseto.hcl /vault/config
COPY --from=builder /plugin/paseto-plugin /etc/vault/vault_plugins
