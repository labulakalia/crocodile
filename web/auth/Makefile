
GOPATH:=$(shell go env GOPATH)

.PHONY: run
run:
	go run main.go plugin.go

.PHONY: build
build:
	go build -o auth-web main.go plugin.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t auth-web:latest
