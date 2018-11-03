#!/usr/bin/env bash

set -xe

export DOCKER_LOCALHOST=docker.for.mac.localhost

export STORAGE_ACCESS_KEY=${STORAGE_ACCESS_KEY:-admin}
export STORAGE_SECRET_KEY=${STORAGE_SECRET_KEY:-password}
export STORAGE_BUCKET=${STORAGE_BUCKET:-ffmpeg}
export STORAGE_REGION=${STORAGE_REGION:-us-east-1}
export S3_URL="s3://${STORAGE_ACCESS_KEY}:${STORAGE_SECRET_KEY}@${DOCKER_LOCALHOST}:9000/${STORAGE_REGION}/${STORAGE_BUCKET}"

pushd ../scripts
    ./start.sh
    ./test.sh
popd

fn --verbose deploy --all app.yml --local --build-arg S3_URL=${S3_URL}

pushd ../scripts
    ./configure.sh
popd
