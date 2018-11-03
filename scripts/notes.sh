#!/usr/bin/env bash

ffmpeg -loglevel panic -y -i http://mirrors.standaloneinstaller.com/video-sample/metaxas-keller-Bell.m4v -acodec copy -f wav - > sox-test.wav


awk '$2 > 0.4 { print $1 }' < sox.dat | awk '(2 < NR)'


ffmpeg -loglevel panic -y -i http://mirrors.standaloneinstaller.com/video-sample/metaxas-keller-Bell.m4v -acodec copy -f wav - | sox -t raw -r 48000 -b 16 -e signed - new-sox.dat


sort -k2,2nr -k1,1 new-sox.dat

awk 'NR > 2' new-sox.dat | awk '$2 > 0.6 { print $1 "," $2}'

echo -e '{"threshold_value":0.499985,"stat_uuid":"af503ec9-08a8-4aff-a3a9-c9c3b0d8475f"}' | fn invoke ffmpeg csv-converter

echo -e '{"threshold_value":0.499985,"csv_file_url":"http://docker.for.mac.localhost:9000/ffmpeg/af503ec9-08a8-4aff-a3a9-c9c3b0d8475f.csv?X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Credential=admin%2F20181101%2Fus-east-1%2Fs3%2Faws4_request\u0026X-Amz-Date=20181101T084722Z\u0026X-Amz-Expires=3600\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Signature=60de48f6828baa2f83f8528ead22b4346044913202d2bc19020a0a6043334a57"}'