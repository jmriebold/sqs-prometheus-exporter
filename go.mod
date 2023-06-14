module sqs-prometheus-exporter

go 1.16

require (
	github.com/aws/aws-sdk-go-v2/config v1.18.26
	github.com/aws/aws-sdk-go-v2/service/sqs v1.23.1
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/prometheus/client_golang v1.15.1
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.10.1 // indirect
	github.com/rs/zerolog v1.29.1
	golang.org/x/sys v0.9.0 // indirect
)
