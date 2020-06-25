FROM golang:alpine as builder

ENV CGO_ENABLED=0

RUN apk update && apk add make git

COPY . $GOPATH/src/github.com/AlbinoDrought/creamy-artifacts
WORKDIR $GOPATH/src/github.com/AlbinoDrought/creamy-artifacts

RUN make test && make install

FROM scratch

COPY --from=builder /go/bin/creamy-artifacts /go/bin/creamy-artifacts
ENTRYPOINT ["/go/bin/creamy-artifacts"]
