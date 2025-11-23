FROM golang:1.25-trixie AS build

WORKDIR /app

COPY main.go go.mod go.sum ./

RUN go mod download && \
	CGO_ENABLED=0 GOOS=linux go build

RUN useradd -u 10001 scratchuser


FROM scratch AS run

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /app/sqs-prometheus-exporter .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

USER scratchuser
ENTRYPOINT [ "./sqs-prometheus-exporter" ]
