FROM golang:alpine as builder

ADD . /go/src/github.com/gdgtoledo/phone-service/

WORKDIR /go/src/github.com/gdgtoledo/phone-service

RUN set -ex && \
  CGO_ENABLED=0 go build -tags netgo -v -a -ldflags '-extldflags "-static"' && \
  mv ./phone-service /usr/bin/phone-service

FROM busybox
COPY --from=builder /usr/bin/phone-service /usr/local/bin/phone-service

EXPOSE 8000

ENTRYPOINT [ "phone-service" ]