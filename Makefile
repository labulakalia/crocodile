VERSION=v1.1.0
COMMIT=`git rev-parse --short HEAD`
CFGPATH='core/config/core.toml'
BUILDDATE=`date "+%Y-%m-%d"`
proto:
	protoc --go_out=plugins=grpc:. core/proto/core.proto

build:
	go build -o crocodile -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" core/main.go

runs:
	go run -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" core/main.go server -c ${CFGPATH}

runc:
	go run -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" core/main.go client -c ${CFGPATH}
version:
	go run -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" core/main.go version