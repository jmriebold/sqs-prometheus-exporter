# SQS Prometheus Exporter

![build](https://github.com/jmriebold/sqs-prometheus-exporter/workflows/Build/badge.svg)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/sqs-prometheus-exporter)](https://artifacthub.io/packages/helm/sqs-prometheus-exporter/sqs-prometheus-exporter)

A simple, lightweight Prometheus metrics exporter for [AWS' Simple Queue Service](https://aws.amazon.com/sqs/), written in Go. Potential use cases are monitoring SQS queues or scaling off SQS queues (e.g. with a [Kubernetes HPA](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)).

# Use

The SQS Prometheus Exporter can be run using the [Docker image](https://github.com/jmriebold/sqs-prometheus-exporter/pkgs/container/sqs-prometheus-exporter) with either the [docker-compose file](docker-compose.yml) or the [Helm chart](chart). To use the Helm chart, execute `helm install sqs-prometheus-exporter oci://ghcr.io/jmriebold/charts/sqs-prometheus-exporter`.

# Configuration

## AWS

In order to authenticate with AWS, the SQS Prometheus Exporter will need AWS credentials either in environment variables (i.e. `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`), a creds file mounted into the container, or a role to assume via [IRSA](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html) or [kiam](https://github.com/uswitch/kiam)/[kube2iam](https://github.com/jtblin/kube2iam). For more information on authenticating with AWS, see the [official documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html).

The user or role that the service will be using to monitor the SQS queues will need the permissions contained in AWS' `AmazonSQSReadOnlyAccess` policy.

## Application

To set the queues scraped by the Exporter, set the `SQS_QUEUE_URLS` environment variable to a comma-separated list of the SQS queue URLs. For example: `SQS_QUEUE_URLS=https://sqs.[region].amazonaws.com/[account-id]/[queue-name-1],https://sqs.[region].amazonaws.com/[account-id]/[queue-name-2]`

By default the application will scrape queue metrics every 30 seconds. To change this interval, set the `SQS_MONITOR_INTERVAL_SECONDS` environment variable.

# Examples

Examples of how to run SQS Prometheus Exporter locally (such as via [systemd](examples/systemd)) can be found under the [examples](examples) directory. See the [chart README](chart/sqs-prometheus-exporter/README.md) for Helm-specific examples.

# Metrics

The SQS Prometheus Exporter serves the following metrics:

- `sqs_approximatenumberofmessages`
- `sqs_approximatenumberofmessagesdelayed`
- `sqs_approximatenumberofmessagesnotvisible`

Each has a `queue` label, which will be populated with the queue name and metric value.

# Dependencies

SQS Prometheus Exporter uses the following Go packages:

- [aws-sdk-go](https://github.com/aws/aws-sdk-go) for monitoring the queues themselves
- [prometheus/client_golang](https://github.com/prometheus/client_golang) for metrics
- [zerolog](https://github.com/rs/zerolog) for logging

## Dev Dependencies

- [golangci](https://github.com/golangci/golangci-lint) code quality tools | `go tool golangci-lint run`
- [skaffold](https://skaffold.dev) local Kubernetes development tool

# Local Dev

## Requirements

* Kubernetes cluster (ideally [Minikube](https://minikube.sigs.k8s.io), however managed or other types of clusters work as well)
* [kubectl](https://kubernetes.io/docs/reference/kubectl)
* [helm](https://helm.sh/)
* [skaffold](https://skaffold.dev)

## Setup

1. Install local dev tools
2. Create a Kubernetes cluster with your tool of choice and set kubectl context accordingly
3. Execute `skaffold dev`

SQS Prometheus Exporter should now be running alongside LocalStack and the Kube-Prometheus-Stack, where you can test your changes.
