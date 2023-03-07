FROM golang:1.20.1-alpine3.17 AS builder
WORKDIR /go/src/queue-nats
COPY . .
RUN \
    apk add protoc protobuf-dev make git && \
    make build

FROM alpine:3.17.0
COPY --from=builder /go/src/queue-nats/queue-nats /bin/queue-nats
ENTRYPOINT ["/bin/queue-nats"]
