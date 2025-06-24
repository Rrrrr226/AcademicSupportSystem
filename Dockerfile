FROM golang:latest AS builder

COPY . /build

WORKDIR /build

RUN set -ex \
    && GO111MODULE=auto CGO_ENABLED=0 go build -ldflags "-s -w -extldflags '-static' -X 'HelpStudent/core/version.SysVersion=$(git show -s --format=%h)'" -o App

FROM alpine:latest

WORKDIR /Serve

COPY --from=builder /build/App ./App

RUN  echo 'https://dl-cdn.alpinelinux.org/alpine/latest-stable/main' > /etc/apk/repositories \
    && echo 'https://dl-cdn.alpinelinux.org/alpine/latest-stable/community' >>/etc/apk/repositories \
    && apk update && apk add tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

ENTRYPOINT [ "/Serve/App", "server" ]