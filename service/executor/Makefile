
GOPATH:=$(shell go env GOPATH)
.PHONY: run
run:
	go run main.go plugin.go

.PHONY: proto
proto:
	protoc --proto_path=../../third_party:. --proto_path=${GOPATH}/src:. --proto_path= --micro_out=../../third_party:. --go_out=. proto/executor/executor.proto

.PHONY: build
build: proto
	go build -o executor-srv main.go plugin.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t executor-srv:latest
