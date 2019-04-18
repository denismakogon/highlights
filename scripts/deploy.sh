#!/usr/bin/env bash

pushd ../functions/audio-splitter && fn --verbose deploy --app ffmpeg --local --build-arg S3_URL=${S3_URL}; popd
pushd ../functions/audio-processor && fn --verbose deploy --app ffmpeg --local --build-arg S3_URL=${S3_URL}; popd
pushd ../functions/amplification-threshold && fn --verbose deploy --app ffmpeg --local --build-arg S3_URL=${S3_URL}; popd
pushd ../functions/timecode-statistics && fn --verbose deploy --app ffmpeg --local --build-arg S3_URL=${S3_URL}; popd
pushd ../functions/csv-converter && fn --verbose deploy --app ffmpeg --local --build-arg S3_URL=${S3_URL}; popd
pushd ../functions/peak-frequency-plotter && fn --verbose deploy --app ffmpeg --local --build-arg S3_URL=${S3_URL}; popd
