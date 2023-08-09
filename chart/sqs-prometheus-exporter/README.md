# SQS Prometheus Exporter

A simple, lightweight Prometheus metrics exporter for [AWS's Simple Queue Service](https://aws.amazon.com/sqs/), written in Go. Potential use cases are monitoring SQS queues or scaling off SQS queues (e.g. with a [Kubernetes HPA](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)).

## Installation

```bash
helm repo add jmriebold https://jmriebold.github.io/charts
helm install release-name --set sqs.region=[region-name] \
  --set sqs.queueUrls[0]=https://sqs.[region].amazonaws.com/[account-id]/[queue-name-1] \
  jmriebold/sqs-prometheus-exporter
```

### Configuration

#### AWS

In order to authenticate with AWS, the SQS Prometheus Exporter will need AWS credentials either in environment variables (i.e. `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`), a creds file mounted into the container, or a role to assume via [IRSA](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html) or [kiam](https://github.com/uswitch/kiam)/[kube2iam](https://github.com/jtblin/kube2iam). For more information on authenticating with AWS, see the [official documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html).

The user or role that the service will be using to monitor the SQS queues will need the permissions contained in AWS' `AmazonSQSReadOnlyAccess` policy.

#### Application

| Name                         | Description                | Default       |
|------------------------------|----------------------------|---------------|
| `sqs.region`                 | AWS region name            | `"us-west-2"` |
| `sqs.queueUrls`              | List of AWS SQS queue URLs | `[]`          |
| `sqs.monitorIntervalSeconds` | Interval for polling SQS   | `30`          |

### Metrics

The SQS Prometheus Exporter serves the following metrics:

- `sqs_approximatenumberofmessages`
- `sqs_approximatenumberofmessagesdelayed`
- `sqs_approximatenumberofmessagesnotvisible`

Each has a `queue` label, which will be populated with the queue name and metric value.

### Examples

#### IRSA

To use IRSA for granting SQS Prometheus Exporter access to SQS:

```bash
helm install release-name \
  --set sqs.region=[region-name] \
  --set sqs.queueUrls[0]=https://sqs.[region].amazonaws.com/[account-id]/[queue-name-1] \
  --set serviceAccount.create=true \
  --set serviceAccount.annotations."eks\.amazonaws\.com/role-arn"=[role-arn] \
  jmriebold/sqs-prometheus-exporter
```

#### kiam/kube2iam

To use Kiam/Kube2Iam for granting SQS Prometheus Exporter access to SQS:

```bash
helm install release-name \
  --set sqs.region=[region-name] \
  --set sqs.queueUrls[0]=https://sqs.[region].amazonaws.com/[account-id]/[queue-name-1] \
  --set podAnnotations."iam\.amazonaws\.com/role"=[role-name] \
  jmriebold/sqs-prometheus-exporter
```

#### Environment Variables

To use AWS environment variables for granting SQS Prometheus Exporter access to SQS:

```bash
helm install release-name \
  --set sqs.region=[region-name] \
  --set sqs.queueUrls[0]=https://sqs.[region].amazonaws.com/[account-id]/[queue-name-1] \
  --set extraEnv.AWS_ACCESS_KEY_ID=[aws-access-key] \
  --set extraEnv.AWS_SECRET_ACCESS_KEY=[aws-secret-access-key] \
  jmriebold/sqs-prometheus-exporter
```
