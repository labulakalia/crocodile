VERSION=`git tag | tail -1`
COMMIT=`git rev-parse --short HEAD`
CFGPATH='core.toml'
BUILDDATE=`date "+%Y-%m-%d"`

BUILD_DIR=build
APP_NAME=crocodile

sources=$(wildcard *.go)

build = GOOS=$(1) GOARCH=$(2) go build -o ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2) -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" main.go 
md5 = md5sum ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2) > ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2)_checksum.txt
tar =  cp core.toml ${BUILD_DIR} && tar -cvzf ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2).tar.gz  -C ${BUILD_DIR}  $(APP_NAME)-$(1)-$(2) $(APP_NAME)-$(1)-$(2)_checksum.txt core.toml
delete = rm -rf ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2) ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2)_checksum.txt ${BUILD_DIR}/core.toml
ALL_LINUX = linux-amd64 \
	linux-386 \
	linux-arm \
	linux-arm64

ALL = $(ALL_LINUX) \
	darwin-amd64

build_linux: $(ALL_LINUX:%=build/%)

build_all: $(ALL:%=build/%)

build/%: 
	$(call build,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)))
	$(call md5,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)))
	$(call tar,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)))
	# $(call delete,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)))

clean:
	rm -rf ${BUILD_DIR}

proto: clean
	protoc --go_out=plugins=grpc:. core/proto/core.proto

build:
	go build -o crocodile -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" main.go
frontend:
	cd web && yarn && yarn run build:prod
bindata: 
	go get github.com/go-bindata/go-bindata/...
	~/go/bin/go-bindata -o=core/utils/asset/asset.go  -pkg=asset web/crocodile/... sql/... && rm -rf ./crocodile

swag:
	go get -u github.com/swaggo/swag/cmd/swag
	~/go/bin/swag init -o core/docs
vet:
	go vet main.go
runs:
	go run -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" main.go server -c ${CFGPATH}
runc:
	go run -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" main.go client -c ${CFGPATH}
version:
	go run -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" main.go version
run:
	go run -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" main.go