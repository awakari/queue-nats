# Contents

1. [Overview](#1-overview)<br/>
   1.1. [Purpose](#11-purpose)<br/>
   1.2. [Definitions](#12-definitions)<br/>
   &nbsp;&nbsp;&nbsp;1.2.1. [Message](#121-message)<br/>
   &nbsp;&nbsp;&nbsp;1.2.2. [Queue](#122-queue)<br/>
2. [Configuration](#2-configuration)<br/>
3. [Deployment](#3-deployment)<br/>
   3.1. [Prerequisites](#31-prerequisites)<br/>
   3.2. [Bare](#32-bare)<br/>
   3.3. [Docker](#33-docker)<br/>
   3.4. [K8s](#34-k8s)<br/>
   &nbsp;&nbsp;&nbsp;3.4.1. [Helm](#341-helm) <br/>
4. [Usage](#4-usage)<br/>
   4.1. [Create/Update Queue](#41-createupdate-queue)<br/>
   4.2. [Submit Message](#42-submit-message)<br/>
   4.3. [Poll Messages](#43-poll-messages)<br/>
5. [Design](#5-design)<br/>
   5.1. [Requirements](#51-requirements)<br/>
   5.2. [Approach](#52-approach)<br/>
   &nbsp;&nbsp;&nbsp;5.2.1. [Data Schema](#521-data-schema)<br/>
   5.3. [Limitations](#53-limitations)<br/>
6. [Contributing](#6-contributing)<br/>
   6.1. [Versioning](#61-versioning)<br/>
   6.2. [Issue Reporting](#62-issue-reporting)<br/>
   6.3. [Building](#63-building)<br/>
   6.4. [Testing](#64-testing)<br/>
   &nbsp;&nbsp;&nbsp;6.4.1. [Functional](#641-functional)<br/>
   &nbsp;&nbsp;&nbsp;6.4.2. [Performance](#642-performance)<br/>
   6.5. [Releasing](#65-releasing)<br/>

# 1. Overview

Nats JetStream implementation for the Awakari Queue service. 

## 1.1. Purpose

To be used both:
* internally to transfer messages between the Awakari system components
* externally as a resolved message destination

## 1.2. Definitions

### 1.2.1. Message

Awakari Queue service works with messages in [CloudEvents](https://cloudevents.io/) format.

### 1.2.2. Queue

The service allows to create a named queues to submit the messages to. 
In addition to the name the queue has a non-negative length limit.

# 2. Configuration

The service is configurable using the environment variables:

| Variable      | Example value      | Description                                                     |
|---------------|--------------------|-----------------------------------------------------------------|
| API_PORT      | `8080`             | gRPC API port                                                   |
| LOG_LEVEL     | `-4`               | [Logging level](https://pkg.go.dev/golang.org/x/exp/slog#Level) |
| NATS_URI      | `nats://nats:4222` | NATS service URI                                                |
| NATS_USER     | `awakari`          | NATS user name                                                  |
| NATS_PASSWORD | `*******`          | NATS password                                                   |

# 3. Deployment

## 3.1. Prerequisites

NATS in the JetStream mode is required to be available.

```shell
helm install nats bitnami/nats \
  --set jetstream.enabled=true \
  --set persistence.enabled=true,resourceType="statefulset" \
  --set auth.user=awakari \
  --set auth.password=awakari \
  --set replicaCount=2
```

## 3.2. Bare

Preconditions:
1. Build "queue-nats" executive using ```make build```

Then run the command:
```shell
API_PORT=8080 \
NATS_URI=nats://localhost:4222 \
./queue-nats
```

## 3.3. Docker

```shell
make run
```

## 3.4. K8s

### 3.4.1. Helm

Create a helm package from the sources:
```shell
helm package helm/queue-nats/
```

Install the helm chart:
```shell
helm install queue-nats ./queue-nats-<CHART_VERSION>.tgz
```

where
* `<CHART_VERSION>` is the helm chart version

# 4. Usage

## 4.1. Create/Update Queue

Example command:
```shell
grpcurl \
  -plaintext \
  -proto api/grpc/service.proto \
  -d @ \
  localhost:8080 \
  queue.Service/SetQueue
```
Example payload:
```json
{
  "name": "q4",
  "limit": 1
}
```

## 4.2. Submit Message

Example command:
```shell
grpcurl \
  -plaintext \
  -proto api/grpc/service.proto \
  -d @ \
  localhost:8080 \
  queue.Service/SubmitMessage
```
Example payload:
```json
{
    "queue": "q4",
    "msg": {
      "id": "3426d090-1b8a-4a09-ac9c-41f2de24d5ac",
      "type": "example.type",
      "source": "example/uri",
      "spec_version": "1.0",
      "attributes": {
        "subject": {
          "ce_string": "Obi-Wan Kenobi"
        },
        "time": {
          "ce_timestamp": "1985-04-12T23:20:50.52Z"
        }
      },
      "text_data": "I felt a great disturbance in the force"
    }
}
```

## 4.3. Poll Messages

Example command:
```shell
grpcurl \
  -plaintext \
  -proto api/grpc/service.proto \
  -d @ \
  localhost:8080 \
  queue.Service/Poll
```
Example payload:
```json
{
  "queue": "q4",
  "limit": 1000
}
```

# 5. Design

## 5.1. Requirements

| #     | Summary                                | Description                                                                                                 |
|-------|----------------------------------------|-------------------------------------------------------------------------------------------------------------|
| REQ-1 | TODO                                   | TODO                                                                                                        |

## 5.2. Approach

TODO

## 5.3. Limitations

| #     | Summary | Description |
|-------|---------|-------------|
| LIM-1 | TODO    | TODO        |

# 6. Contributing

## 6.1. Versioning

The service uses the [semantic versioning](http://semver.org/).
The single source of the version info is the git tag:
```shell
git describe --tags --abbrev=0
```

## 6.2. Issue Reporting

TODO

## 6.3. Building

```shell
make build
```
Generates the sources from proto files, compiles and creates the `queue-nats` executable.

## 6.4. Testing

### 6.4.1. Functional

```shell
make test
```

### 6.4.2. Performance

TODO

## 6.5. Releasing

To release a new version (e.g. `1.2.3`) it's enough to put a git tag:
```shell
git tag -v1.2.3
git push --tags
```

The corresponding CI job is started to build a docker image and push it with the specified tag (+latest).
