default: build

clearbin:
	rm -rf ./bin

build: cmd clearbin
	GOOS=linux GOARCH=amd64 go build -v -o bin/main-linux-amd64 ./cmd

cross: cmd clearbin
	go build -v -o bin/main ./cmd

run: cross
	./bin/main

protobuffer:
	protoc --go_out=./transport/grpc --go-grpc_out=./transport/grpc ./protobuffer-files/*.proto

.PHONY: default run clearbin protobuffer
