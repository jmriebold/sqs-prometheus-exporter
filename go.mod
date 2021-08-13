module sqs-prometheus-exporter

go 1.16

require (
	github.com/aws/aws-sdk-go-v2/config v1.6.0
	github.com/aws/aws-sdk-go-v2/service/sqs v1.7.1
	github.com/prometheus/client_golang v1.11.0
	github.com/rs/zerolog v1.23.0
)
