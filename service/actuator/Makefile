
GOPATH:=$(shell go env GOPATH)

.PHONY: run
run:
	go run main.go plugin.go

.PHONY: proto
proto:
	protoc --proto_path=../../third_party:. --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/actuator/actuator.proto
	ls proto/actuator/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'

.PHONY: build
build: proto
	go build -o actuator-srv main.go plugin.go

.PHONY: docker
docker: build
	docker build . -t actuator-srv:latest


.PHONY: test
test:
	go test -v ./... -cover
