FROM golang:1.24-alpine AS builder
RUN apk update && apk add --no-cache ca-certificates

ENV CGO_ENABLED=0 GO111MODULE=on GOOS=linux

WORKDIR /

COPY go.* ./
RUN   --mount=type=cache,target=/go/pkg/mod \
    go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o main ./cmd/api/main.go

FROM scratch

ENV HTTP_PORT=4004

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /main .
EXPOSE $HTTP_PORT

CMD ["/main"]