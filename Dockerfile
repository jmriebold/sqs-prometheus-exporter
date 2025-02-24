FROM golang:1.24-bookworm AS build

WORKDIR /app

COPY main.go go.mod go.sum ./

RUN go mod download && \
	CGO_ENABLED=0 GOOS=linux go build

FROM alpine:3.21.3 AS run

COPY --from=build /app/sqs-prometheus-exporter .

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

ENTRYPOINT [ "./sqs-prometheus-exporter" ]
