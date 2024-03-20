.DEFAULT_GOAL := clean
BINARY_NAME=text_editor
GOGC=100

fmt: 
	go fmt 

vet: fmt
	go vet

build: 
	GOARCH=amd64 GOOS=windows CGO_ENABLED=1 go build -o ./target/${BINARY_NAME}-windows.exe main.go

run: build 
	GOGC=${GOGC} GODEBUG=gctrace=1 ./target/${BINARY_NAME}-windows 

clean: run
	go clean ./