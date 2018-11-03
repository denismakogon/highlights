FROM denismakogon/ffmpeg-debian:golang as build-stage

ARG S3_URL

RUN apt-get update \
  && apt-get install -qy sox libsox-fmt-all

RUN ffmpeg -loglevel panic -y -i http://mirrors.standaloneinstaller.com/video-sample/metaxas-keller-Bell.m4v -vn -f wav - > /tmp/sox.wav

ADD . /go/src/func
WORKDIR /go/src/func/
RUN S3_URL=${S3_URL} go test -v ./... && go build -o func

FROM denismakogon/ffmpeg-debian:runtime

RUN apt-get update \
  && apt-get install -qy sox libsox-fmt-all
WORKDIR /function
COPY --from=build-stage /go/src/func/func /function/

ENTRYPOINT ["./func"]
