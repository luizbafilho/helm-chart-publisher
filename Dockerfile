FROM alpine:3.5

RUN apk add --no-cache ca-certificates

COPY ./helm-chart-publisher /helm-chart-publisher

CMD ["/helm-chart-publisher", "-c", "/config.yaml"]

