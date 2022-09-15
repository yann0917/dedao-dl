FROM golang:1.19-alpine3.16 AS builder

LABEL stage=gobuilder

WORKDIR /build
# RUN adduser -u 10001 -D app-runner

ENV GOPROXY https://goproxy.cn
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETPLATFORM

# build and using UPX to compress
RUN if [ "$TARGETPLATFORM" = "linux/amd64" ]; then ARCHITECTURE=amd64; elif [ "$TARGETPLATFORM" = "linux/arm/v7" ]; then ARCHITECTURE=arm; elif [ "$TARGETPLATFORM" = "linux/arm64" ]; then ARCHITECTURE=arm64; else ARCHITECTURE=amd64; fi \
    && CGO_ENABLED=0 GOARCH=${ARCHITECTURE} GOOS=linux go build -ldflags="-s -w" -a -o dedao-dl . \
    && wget https://github.com/upx/upx/releases/download/v3.96/upx-3.96-${ARCHITECTURE}_linux.tar.xz \
    && tar -xvf upx-3.96-${ARCHITECTURE}_linux.tar.xz \
    && cd upx-3.96-${ARCHITECTURE}_linux \
    && chmod a+x ./upx \
    && mv ./upx /usr/local/bin/ \
    && cd ../ && rm -rf upx-3.96-${ARCHITECTURE}_linux && rm -rf upx-3.96-${ARCHITECTURE}_linux.tar.xz \
    && cd /build \
    && upx dedao-dl \
    && chmod a+x ./dedao-dl

FROM alpine:3.16

# Installs latest ffmpeg, Chromium and chinese font.
RUN echo @3.16 https://mirrors.aliyun.com/alpine/v3.16/community > /etc/apk/repositories \
    && echo @3.16 https://mirrors.aliyun.com/alpine/v3.16/main >> /etc/apk/repositories \
    && echo @edge https://mirrors.aliyun.com/alpine/edge/testing >> /etc/apk/repositories \
    && apk update \
    && apk add --no-cache ffmpeg@3.16 tzdata@3.16 chromium@3.16 \
    && apk add --no-cache --allow-untrusted harfbuzz@3.16 nss@3.16 freetype@3.16 \
    ttf-freefont@3.16 wqy-zenhei@edge \
    && cp -r -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=builder /build/dedao-dl /app/

CMD chromium-browser --headless --disable-gpu --remote-debugging-port=9222 --disable-web-security --safebrowsing-disable-auto-update --disable-sync --disable-default-apps --hide-scrollbars --metrics-recording-only --mute-audio --no-first-run --no-sandbox

ENTRYPOINT ["/app/dedao-dl"]
