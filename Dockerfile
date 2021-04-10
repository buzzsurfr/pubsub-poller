FROM alpine:latest
COPY poller /poller
ENTRYPOINT ["/poller"]
