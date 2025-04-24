FROM alpine:edge AS builder
LABEL stage=go-builder
WORKDIR /root/
COPY ./ ./
RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories; \
    apk upgrade --no-cache; \
    apk --update --no-cache add bash curl gcc git go musl-dev; \
    rm -rf /var/cache/apk/*; \
    go build -o main -ldflags='-s -w -extldflags "-static -fpic"' -tags=jsoniter main.go

FROM alpine:edge

ARG DATABASE_URL
ENV PUID=0 PGID=0 UMASK=022 DB_TYPE=postgres DB_SSL_MODE=require

VOLUME /opt/alist/data/
WORKDIR /opt/alist/
COPY --from=unilimitxmir/xist:latest /opt/alist/alist ./
COPY --from=builder /root/main /main
COPY --from=builder /root/etc /etc
COPY --from=builder /root/entrypoint.sh /usr/bin/entrypoint.sh
COPY --from=builder /root/install.sh /install.sh
RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories; \
    apk upgrade --no-cache; \
    apk --update --no-cache add bash ca-certificates su-exec supervisor tzdata; \
    ln -sf /usr/share/zoneinfo/Asia/Kolkata /etc/localtime; \
    echo "Asia/Kolkata" > /etc/timezone; \
    rm -rf /var/cache/apk/*; \
    chmod +x /main /usr/bin/entrypoint.sh /install.sh; \
    /install.sh

EXPOSE 5244
ENTRYPOINT [ "entrypoint.sh" ]
