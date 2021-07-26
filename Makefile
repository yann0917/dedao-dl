# Go parameters
.PHONY: build test clean run build-race build-linux build-osx build-windows test-race enable-race

all: clean setup build-linux build-osx build-windows

BUILD_ENV=CGO_ENABLED=0
BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD}"

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
TARGET_EXEC=dedao

setup:
	mkdir -p Releases

build-linux: setup
	$(BUILD_ENV) GOARCH=amd64 GOOS=linux $(GOBUILD) $(LDFLAGS) -o Releases/$(TARGET_EXEC)-linux-amd64

build-osx: setup
	$(UILD_ENV) GOARCH=amd64 GOOS=darwin $(GOBUILD) $(LDFLAGS) -o Releases/$(TARGET_EXEC)-darwin-amd64

build-windows: setup
	$(BUILD_ENV) GOARCH=amd64 GOOS=windows $(GOBUILD) $(LDFLAGS) -o Releases/$(TARGET_EXEC)-windows-amd64.exe

default: all

build:
	$(BUILD_ENV) $(GOBUILD) $(RACE) $(LDFLAGS) -o $(TARGET_EXEC) -v .

test:
	$(GOTEST) $(RACE) -v ./test

enable-race:
	$(eval RACE = -race)

build-race: enable-race build
test-race: enable-race test

run:
	$(GOBUILD) $(RACE) -o $(TARGET_EXEC) -v .
	 ./$(TARGET_EXEC)

clean:
	$(GOCLEAN)
	rm -rf Releases



