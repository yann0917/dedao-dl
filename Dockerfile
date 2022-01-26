FROM golang:1.17-alpine3.15 AS builder

LABEL stage=gobuilder

WORKDIR /build
# RUN adduser -u 10001 -D app-runner

ENV GOPROXY https://goproxy.cn
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# build and using UPX to compress
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -a -o dedao-dl . \
    && wget https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz \
    && tar -xvf upx-3.96-amd64_linux.tar.xz \
    && cd upx-3.96-amd64_linux \
    && chmod a+x ./upx \
    && mv ./upx /usr/local/bin/ \
    && cd ../ && rm -rf upx-3.96-amd64_linux && rm -rf upx-3.96-amd64_linux.tar.xz \
    && cd /build \
    && upx dedao-dl \
    && chmod a+x ./dedao-dl

FROM alpine:3.15

# Installs latest ffmpeg, Chromium and chinese font.
RUN echo @3.15 https://mirrors.aliyun.com/alpine/v3.15/community > /etc/apk/repositories \
    && echo @3.15 https://mirrors.aliyun.com/alpine/v3.15/main >> /etc/apk/repositories \
    && echo @edge https://mirrors.aliyun.com/alpine/edge/testing >> /etc/apk/repositories \
    && apk update \
    && apk add --no-cache ffmpeg@3.15 tzdata@3.15 chromium@3.15 \
    && apk add --no-cache --allow-untrusted harfbuzz@3.15 nss@3.15 freetype@3.15 \
    ttf-freefont@3.15 wqy-zenhei@edge \
    && cp -r -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=builder /build/dedao-dl /app/

CMD chromium-browser --headless --disable-gpu --remote-debugging-port=9222 --disable-web-security --safebrowsing-disable-auto-update --disable-sync --disable-default-apps --hide-scrollbars --metrics-recording-only --mute-audio --no-first-run --no-sandbox

ENTRYPOINT ["/app/dedao-dl"]
