FROM fnproject/go:dev as build-stage

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
ADD . /go/src/github.com/denismakogon/s3-pollster
WORKDIR /go/src/github.com/denismakogon/s3-pollster
RUN dep ensure -v
RUN go build -o pollster

FROM fnproject/go

COPY --from=build-stage /go/src/github.com/denismakogon/s3-pollster/pollster /
ENTRYPOINT ["/pollster"]
