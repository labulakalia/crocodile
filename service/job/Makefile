
GOPATH:=$(shell go env GOPATH)

.PHONY: run
run:
	go run main.go plugin.go

.PHONY: proto
proto:
	protoc --proto_path=../../third_party:. --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/job/job.proto
	ls proto/job/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'

.PHONY: build
build: proto

	go build -o job-srv main.go plugin.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t job-srv:latest
