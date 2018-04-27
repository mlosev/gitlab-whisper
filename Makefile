.PHONY: build vendor

build:
	go build -i -v

vendor:
	dep ensure -v
