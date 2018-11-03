#!/usr/bin/env bash

export STORAGE_ACCESS_KEY=${STORAGE_ACCESS_KEY:-admin}
export STORAGE_SECRET_KEY=${STORAGE_SECRET_KEY:-password}

docker rm -f minio1 || true
docker run -d -v `pwd`/.data:/data -p 9000:9000  --rm --name minio1 \
    -e "MINIO_ACCESS_KEY=${STORAGE_ACCESS_KEY}" \
    -e "MINIO_SECRET_KEY=${STORAGE_SECRET_KEY}" minio/minio  server /data
