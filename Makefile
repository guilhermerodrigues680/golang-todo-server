build: cmd
	rm -rf bin
	GOOS=linux GOARCH=amd64 go build -v -o bin/main-linux-amd64 ./cmd

cross: cmd
	rm -rf bin
	go build -v -o bin/main ./cmd
