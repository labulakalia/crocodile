VERSION=v1.1.0
COMMIT=`git rev-parse --short HEAD`
CFGPATH='core/config/core.toml'
BUILDDATE=`date "+%Y-%m-%d"`

BUILD_DIR=build
APP_NAME=crocodile

sources=$(wildcard *.go)

build = GOOS=$(1) GOARCH=$(2) CGO_ENABLED=1 go build -o ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2) -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" core/main.go 
md5 = md5 ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2) > ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2)_checksum.txt
tar =  tar -cvzf ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2).tar.gz  -C ${BUILD_DIR}  $(APP_NAME)-$(1)-$(2) $(APP_NAME)-$(1)-$(2)_checksum.txt
delete = rm -rf ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2) ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2)_checksum.txt
ALL_LINUX = linux-amd64 \
	linux-386 \
	linux-arm \
	linux-arm64

ALL = $(ALL_LINUX) \
	darwin-amd64 \
	windows-amd64

build_linux: $(ALL_LINUX:%=build/%)

build_all: $(ALL:%=build/%)

build/%: frontend
	$(call build,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)))
	$(call md5,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)))
	$(call tar,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)))
	$(call delete,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)))

clean:
	rm -rf ${BUILD_DIR}

proto: clean
	protoc --go_out=plugins=grpc:. core/proto/core.proto

build:
	CGO_ENABLED=1 go build -o crocodile -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" core/main.go
frontend:
	cd web && yarn run build:prod
	cd web && go-bindata -o=../core/router/api/v1/asset/asset.go  -pkg=asset ./crocodile/... && rm -rf ./crocodile

runs:
	go run -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" main.go server -c ${CFGPATH}
runc:
	go run -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" main.go client -c ${CFGPATH}
version:
	go run -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" main.go version
run:
	go run -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" main.go