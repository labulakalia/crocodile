VERSION=1.1.0
COMMIT=`git rev-parse --short HEAD`
CFGPATH='core/config/core.toml'
proto:
	protoc --go_out=plugins=grpc:. core/proto/core.proto

tlskey:
	# 生成服务端私钥
	openssl ecparam -genkey -name secp384r1 -out core/tls/server.key
	# 生成自签公钥
	openssl req -new -x509 -sha256 -key core/tls/server.key -out core/tls/server.pem -days 3650

build:
	go run -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT}" core/main.go

runs:
	go run -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT}" core/main.go server -c ${CFGPATH}

runc:
	go run -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT}" core/main.go client -c ${CFGPATH}
version:
	go run -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT}" core/main.go version