# go source files
SRC := $(shell find .. -type f -name '*.go')

all: cluster-service

cluster-service: $(SRC)
	go build -ldflags "-X main.commit=$(shell git rev-parse HEAD)"

build: cluster-service

install:
	go install -ldflags "-X main.commit=$(shell git rev-parse HEAD)"

clean:
	rm -f cluster-service

.PHONY: clean install build
