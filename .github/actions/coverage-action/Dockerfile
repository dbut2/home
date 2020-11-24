FROM golang:1.14.4-buster

ENV GO111MODULE=on
RUN go get github.com/ory/go-acc

WORKDIR /coverage
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY templates/ templates/
COPY fonts/ fonts/
COPY pkg/ pkg/
COPY cmd/main/main.go cmd/main/main.go
COPY cmd/main/config.go cmd/main/config.go

WORKDIR /
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
