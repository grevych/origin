ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -tags gobox_dev -o run-app -v ./cmd/origin

FROM debian:bookworm

COPY --from=builder /usr/src/app/trace.yaml /run/config/gobox/
COPY --from=builder /usr/src/app/run-app /usr/local/bin/
RUN mkdir -p /run/config/gobox
RUN touch /run/config/gobox/config.yaml
CMD ["run-app"]
