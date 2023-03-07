#!/bin/bash

export SLUG=ghcr.io/awakari/queue-nats
export VERSION=latest
docker tag awakari/queue-nats "${SLUG}":"${VERSION}"
docker push "${SLUG}":"${VERSION}"
