FROM golang:alpine as builder

ADD . /go/src/github.com/gdgtoledo/phone-users/

WORKDIR /go/src/github.com/gdgtoledo/phone-users

RUN set -ex && \
  CGO_ENABLED=0 go build -tags netgo -v -a -ldflags '-extldflags "-static"' && \
  mv ./phone-users /usr/bin/phone-users

FROM busybox
COPY --from=builder /usr/bin/phone-users /usr/local/bin/phone-users

EXPOSE 8000

ENTRYPOINT [ "phone-users" ]