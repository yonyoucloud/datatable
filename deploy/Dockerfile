FROM alpine:latest

LABEL maintainer="bishh@yonyou.com"

ARG TARGETARCH

RUN apk upgrade --no-cache && \
    apk add --no-cache --update tzdata && \
    rm -rf /tmp/* /var/cache/apk/*

COPY ${TARGETARCH}/datatable /datatable/datatable
COPY web /datatable/web
COPY etc/config.yaml /datatable/etc/config.yaml

RUN \
mkdir -p /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && \
ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
echo "Asia/Shanghai" > /etc/timezone

WORKDIR /

CMD /bin/sh -c "cd /datatable && ./datatable"