# Build image
FROM golang:1.15-alpine as build

WORKDIR /app

COPY main.go go.mod go.sum ./

RUN go build

# Application image
FROM alpine:3.12 as run

RUN addgroup -S sqs-exporter && \
	adduser -S -G sqs-exporter sqs-exporter --gecos "" --disabled-password --no-create-home

USER sqs-exporter

COPY --from=build /app/sqs-prometheus-exporter .

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

CMD ./sqs-prometheus-exporter
