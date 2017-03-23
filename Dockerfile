FROM alpine:3.5

RUN apk add --no-cache ca-certificates

COPY ./dist/helm-chart-publisher_linux-amd64 /helm-chart-publisher

CMD ["/helm-chart-publisher", "-c", "/config.yaml"]

