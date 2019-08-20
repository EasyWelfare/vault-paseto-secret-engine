FROM golang:1.12.6-stretch AS builder

WORKDIR /plugin
COPY ./go.mod ./go.sum ./main.go ./
COPY ./vendor ./vendor
COPY ./paseto ./paseto
COPY ./vault ./vault
RUN CGO_ENABLED="0" GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -mod vendor -o /paseto .

FROM vault:1.2.0
RUN mkdir -p /etc/vault/vault_plugins
COPY --from=builder /paseto /etc/vault/vault_plugins/
CMD ["vault", "server", "-config=/etc/vault/local.hcl", "-log-level=debug"]
