# SQS Prometheus Exporter

A simple, lightweight Prometheus metrics exporter for AWS's Simple Queue Service, written in Go.

## Use

The SQS Prometheus Exporter can be run using the Docker image with either the docker-compose file or the Helm chart.

### Configuration

#### Authentication

In order to authenticate with AWS, the SQS Prometheus Exporter will need AWS credentials either in environment variables (i.e. `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, and `AWS_REGION`), or in the form of a creds file mounted into the container. For more information on authenticating with AWS, see the [official documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html).

#### Application

To set the queues scraped by the Exporter, set the `SQS_QUEUE_URLS` environment variable to a comma-separated list of the SQS queue URLs. For example: `SQS_QUEUE_URLS=https://sqs.[region].amazonaws.com/[account-id]/[queue-name-1],https://sqs.[region].amazonaws.com/[account-id]/[queue-name-2]`

By default the application will scrape queue metrics every 30 seconds. To change this interval, set the `SQS_MONITOR_INTERVAL_SECONDS` environment variable.

### Metrics Exported

The SQS Prometheus Exporter exports the following metrics:

* `sqs_approximatenumberofmessages`
* `sqs_approximatenumberofmessagesdelayed`
* `sqs_approximatenumberofmessagesnotvisible`

Each has a `queue` label, which will be populated with the name of each queue.
