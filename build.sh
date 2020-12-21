#!/bin/sh

mkdir "Releases"

# 【darwin/amd64】
echo "start build darwin/amd64 >>>"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags '-w -s' -o ./Releases/dedao-darwin-amd64 main.go

# 【windows/amd64】
echo "start build windows/amd64 >>>"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags '-w -s' -o ./Releases/dedao-windows-amd64.exe main.go

# 【linux/amd64】
echo "start build linux/amd64 >>>"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o ./Releases/dedao-linux-amd64 main.go

echo "All build success!!!"
