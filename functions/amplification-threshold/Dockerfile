FROM denismakogon/ffmpeg-debian:golang as build-stage

RUN apt-get update \
  && apt-get install -qy sox libsox-fmt-all

RUN ffmpeg -loglevel panic -y -i http://mirrors.standaloneinstaller.com/video-sample/metaxas-keller-Bell.m4v -vn -f wav - | sox -t wav - /tmp/sox.dat

ADD . /go/src/func
WORKDIR /go/src/func/
# for the test purposes
RUN go test -v ./... && go build -o func

FROM debian:stretch-slim
WORKDIR /function
COPY --from=build-stage /go/src/func/func /function/

ENTRYPOINT ["./func"]
