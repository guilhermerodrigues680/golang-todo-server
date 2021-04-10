default: build

clearbin:
	rm -rf ./bin

build: cmd clearbin
	GOOS=linux GOARCH=amd64 go build -v -o bin/main-linux-amd64 ./cmd

cross: cmd clearbin
	go build -v -o bin/main ./cmd

run: cross
	./bin/main

.PHONY: default run clearbin
