.DEFAULT_GOAL := clean
BINARY_NAME=text_editor

fmt: 
	go fmt 

vet: fmt
	go vet

build: 
	go build -o ./target/${BINARY_NAME}-windows.exe main.go

run: build 
	./target/${BINARY_NAME}-windows 

clean: run
	go clean ./