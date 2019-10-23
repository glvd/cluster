# go source files
SRC := $(shell find .. -type f -name '*.go')

all: cluster-ctrl

cluster-ctrl: $(SRC)
	go build

build: cluster-ctrl

install:
	go install

clean:
	rm -f ipfs-cluster-ctl

.PHONY: clean install build
