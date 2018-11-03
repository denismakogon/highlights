#!/usr/bin/env bash

set -ex

export APP=${APP:-ffmpeg}
export DOCKER_LOCALHOST=${DOCKER_LOCALHOST:-docker.for.mac.localhost}

export STORAGE_ACCESS_KEY=${STORAGE_ACCESS_KEY:-admin}
export STORAGE_SECRET_KEY=${STORAGE_SECRET_KEY:-password}
export STORAGE_BUCKET=${STORAGE_BUCKET:-ffmpeg}
export STORAGE_REGION=${STORAGE_REGION:-us-east-1}
export S3_URL="s3://${STORAGE_ACCESS_KEY}:${STORAGE_SECRET_KEY}@${DOCKER_LOCALHOST}:9000/${STORAGE_REGION}/${STORAGE_BUCKET}"

# env vars are compliant and compatible with github.com/denismakogon/s3-pollster
fn config app ${APP} BUCKET ${STORAGE_BUCKET:-ffmpeg-demo}
fn config app ${APP} ENDPOINT "${DOCKER_LOCALHOST}:9000"
fn config app ${APP} REGION ${STORAGE_REGION:-us-phoenix-1}
fn config app ${APP} ACCESS_KEY_ID ${STORAGE_ACCESS_KEY}
fn config app ${APP} SECRET_ACCESS_KEY ${STORAGE_SECRET_KEY}

# alternative path, may not work with all login/passwords
fn config app ${APP} S3_URL ${S3_URL}

# if syslog env var was not set, ignore error
fn update app ${APP} --syslog-url=${SYSLOG_URL:-""} | true
