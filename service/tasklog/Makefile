
GOPATH:=$(shell go env GOPATH)

.PHONY: run
run:
	go run main.go plugin.go

.PHONY: proto
proto:
	protoc --proto_path=../../third_party:. --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/tasklog/tasklog.proto

.PHONY: build
build: proto

	go build -o tasklog-srv main.go plugin.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t tasklog-srv:latest
