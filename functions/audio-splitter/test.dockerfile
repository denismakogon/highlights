FROM denismakogon/ffmpeg-debian:golang as build-stage

ARG S3_URL

ADD . /go/src/func/
WORKDIR /go/src/func/

RUN S3_URL=${S3_URL} go test -v ./... && go build -o func

FROM denismakogon/ffmpeg-debian:runtime
WORKDIR /function
COPY --from=build-stage /go/src/func/func /function/
ENTRYPOINT ["./func"]
