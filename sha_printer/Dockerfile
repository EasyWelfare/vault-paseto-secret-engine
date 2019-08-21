FROM golang:1.12.6-stretch

WORKDIR /app
COPY ./go.mod ./go.sum ./main.go ./
COPY ./vendor ./vendor
RUN CGO_ENABLED="0" GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -mod vendor -o /sha_printer && chmod +x /sha_printer
CMD ["/sha_printer"]
