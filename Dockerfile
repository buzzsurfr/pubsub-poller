FROM golang:1.16-alpine AS builder
COPY . /app
WORKDIR /app
ENV CGO_ENABLED=0
RUN go build -o /main

FROM alpine:latest
LABEL org.opencontainers.image.source=https://github.com/buzzsurfr/pubsub-poller
COPY --from=builder /main /poller
ENTRYPOINT ["/poller"]