FROM golang:1.23-bookworm as build

WORKDIR /app

COPY main.go go.mod go.sum ./

RUN go mod download && \
	CGO_ENABLED=0 GOOS=linux go build

FROM gcr.io/distroless/base-debian11 AS run

COPY --from=build /app/sqs-prometheus-exporter .

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT [ "./sqs-prometheus-exporter" ]
