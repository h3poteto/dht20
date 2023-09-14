.PHONY: build

GOOS=linux
GOARCH=arm
GOARM=6

install:
	go mod tidy

build: install
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build
