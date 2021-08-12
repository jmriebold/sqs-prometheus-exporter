# SQS Prometheus Exporter

![build](https://github.com/jmriebold/sqs-prometheus-exporter/workflows/Build/badge.svg)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/jmriebold)](https://artifacthub.io/packages/search?repo=jmriebold)

A simple, lightweight Prometheus metrics exporter for [AWS's Simple Queue Service](https://aws.amazon.com/sqs/), written in Go. Potential use cases are monitoring SQS queues or scaling off SQS queues (e.g. with a [Kubernetes HPA](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)).

## Use

The SQS Prometheus Exporter can be run using the [Docker image](https://hub.docker.com/repository/docker/jmriebold/sqs-prometheus-exporter) with either the [docker-compose file](docker-compose.yml) or the [Helm chart](chart).

### Configuration

#### AWS

In order to authenticate with AWS, the SQS Prometheus Exporter will need AWS credentials either in environment variables (i.e. `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`), a creds file mounted into the container, or a role to assume via [IRSA](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html) or [kiam](https://github.com/uswitch/kiam)/[kube2iam](https://github.com/jtblin/kube2iam). For more information on authenticating with AWS, see the [official documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html).

The user or role that the service will be using to monitor the SQS queues will need the permissions contained in AWS' `AmazonSQSReadOnlyAccess` policy.

#### Application

To set the queues scraped by the Exporter, set the `SQS_QUEUE_URLS` environment variable to a comma-separated list of the SQS queue URLs. For example: `SQS_QUEUE_URLS=https://sqs.[region].amazonaws.com/[account-id]/[queue-name-1],https://sqs.[region].amazonaws.com/[account-id]/[queue-name-2]`

By default the application will scrape queue metrics every 30 seconds. To change this interval, set the `SQS_MONITOR_INTERVAL_SECONDS` environment variable.

### Metrics

The SQS Prometheus Exporter serves the following metrics:

- `sqs_approximatenumberofmessages`
- `sqs_approximatenumberofmessagesdelayed`
- `sqs_approximatenumberofmessagesnotvisible`

Each has a `queue` label, which will be populated with the queue name and metric value.

### Dependencies

SQS Prometheus Exporter uses the following Go packages:

- [aws-sdk-go](https://github.com/aws/aws-sdk-go) for monitoring the queues themselves
- [prometheus/client_golang](https://github.com/prometheus/client_golang) for metrics
- [zerolog](https://github.com/rs/zerolog) for logging
